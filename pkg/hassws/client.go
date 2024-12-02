package hassws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrAuthInvalid            = errors.New("invalid access token")
	ErrAuthUnexpectedResponse = errors.New("unexpected response during authentication")
)

type Client struct {
	cfg  ClientConfig
	conn *websocket.Conn

	msgID atomic.Int64

	// Each request has a unique ID and any response will have the same ID. To
	// provide a synchronous API, we store a channel for each request and stream
	// the response there.
	responses map[int]chan []byte
	mutex     sync.RWMutex
}

type ClientConfig struct {
	Host  string
	Token string
}

func NewWebsocketAPI(config ClientConfig) *Client {
	return &Client{
		cfg:       config,
		responses: make(map[int]chan []byte),
	}
}

func (c *Client) authenticate() error {
	// Read auth_required message from server
	var authRequired AuthChallenge
	if err := c.read(&authRequired); err != nil {
		return err
	}

	log.Println("Authenticating")

	// Send auth message with access token
	if err := c.send(AuthRequest{
		Type:        "auth",
		AccessToken: c.cfg.Token,
	}); err != nil {
		return err
	}

	// Read auth_ok or auth_invalid response
	var authResponse AuthResponse
	if err := c.read(&authResponse); err != nil {
		return err
	}

	slog.Debug("Received auth response", "msg", authResponse)

	if authResponse.Type == "auth_invalid" {
		return fmt.Errorf("%w: %s", ErrAuthInvalid, authResponse.Message)
	}

	if authResponse.Type != "auth_ok" {
		return fmt.Errorf("%w: %s", ErrAuthUnexpectedResponse, authResponse.Message)
	}

	log.Println("Authenticated")

	return nil
}

func (c *Client) Close() error {
	return c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
}

func (c *Client) shutdown() error {
	return c.conn.Close()
}

func (c *Client) Connect() error {
	log.Println("Connecting to", c.cfg.Host)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/api/websocket", c.cfg.Host), nil)
	if err != nil {
		return err
	}

	log.Println("Connected")

	c.conn = conn

	if err := c.authenticate(); err != nil {
		return err
	}

	// Start listening for messages
	go c.listen()

	return nil
}

// Listen for messages from the websocket and dispatch to listener channels.
func (c *Client) listen() {
	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Received close message, bye")
				c.shutdown()
				return
			}

			panic(fmt.Errorf("error reading from websocket: %w", err))
		}

		slog.Debug("Received message", "msg", string(b))

		// Get message ID
		var msg CommandMessage
		if err := json.Unmarshal(b, &msg); err != nil {
			slog.Error("Error unmarshalling message", "error", err)
			continue
		}

		c.mutex.RLock()
		ch, ok := c.responses[msg.ID]
		c.mutex.RUnlock()

		if !ok {
			slog.Warn("No listeners for message", "id", msg.ID)
			continue
		}

		go func(ch chan []byte) {
			defer func() {
				if r := recover(); r != nil {
					c.removeMessageResponseListener(msg.ID)
					slog.Debug("Removed listener due to panic sending to channel", "id", msg.ID)
				}
			}()

			ch <- b
		}(ch)
	}
}

// Generate a unique ID for each message.
func (c *Client) nextMsgID() int {
	return int(c.msgID.Add(1))
}

// Read a message from the websocket and unmarshal it into the target.
func (c *Client) read(target any) error {
	_, b, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}

	slog.Debug("Received message", "msg", string(b))

	return json.Unmarshal(b, target)
}

// Send a message to the websocket.
func (c *Client) send(msg any) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	slog.Debug("Writing message", "msg", string(b))

	return c.conn.WriteMessage(websocket.TextMessage, b)
}

// Add a listener channel for a response to a specific sent message.
func (c *Client) addMessageResponseListener(msgID int) (ch chan []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ch = make(chan []byte)
	c.responses[msgID] = ch

	return ch
}

func (c *Client) removeMessageResponseListener(msgID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.responses, msgID)
}

// Send a message to the websocket and return a channel to listen for responses.
func (c *Client) sendMessageStreamResponses(msgBytes []byte) (ch chan []byte, err error) {
	var msg jsonMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return nil, err
	}

	msgID := c.nextMsgID()
	msg["id"] = msgID

	ch = c.addMessageResponseListener(msgID)

	if err := c.send(msg); err != nil {
		return nil, err
	}

	return ch, nil
}

// Send a message to the websocket and wait for a response.
func (c *Client) sendMessageWaitResponse(msgBytes []byte) (response []byte, err error) {
	ch, err := c.sendMessageStreamResponses(msgBytes)
	if err != nil {
		return nil, err
	}

	// Close channel after receiving first response
	defer func() {
		close(ch)
	}()

	return c.readMesssageFromChannel(ch)
}

// Read a message from a listener channel.
func (c *Client) readMesssageFromChannel(ch chan []byte) (response []byte, err error) {
	select {
	case res := <-ch:
		return res, nil
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// Subscribe to home assistant events.
func (c *Client) SubscribeEvents(eventType string, handler func(EventMessage)) error {
	msg := subscribeEventsRequest{
		ID:        c.nextMsgID(),
		Type:      MessageTypeSubscribeEvents,
		EventType: eventType,
	}

	reqBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ch, err := c.sendMessageStreamResponses(reqBytes)
	if err != nil {
		return err
	}

	// First message contains the initial response about the subscription
	resBytes, err := c.readMesssageFromChannel(ch)
	if err != nil {
		close(ch)
		return err
	}

	var res subscribeEventsResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		close(ch)
		return err
	}

	if !res.Success {
		close(ch)
		return fmt.Errorf("unexpected response: %s", resBytes)
	}

	// Create Goroutine to event messages and dispatch to handler
	go func(ch chan []byte) {
		for b := range ch {
			var msg EventMessage
			if err := json.Unmarshal(b, &msg); err != nil {
				slog.Error("Error unmarshalling event message", "error", err)
				continue
			}

			handler(msg)
		}
	}(ch)

	return nil
}

func (c *Client) CallService(msg CallServiceRequest) (CallServiceResponse, error) {
	reqBytes, err := json.Marshal(msg)
	if err != nil {
		return CallServiceResponse{}, err
	}

	resBytes, err := c.sendMessageWaitResponse(reqBytes)
	if err != nil {
		return CallServiceResponse{}, err
	}

	var resp CallServiceResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return CallServiceResponse{}, err
	}

	return resp, nil
}
