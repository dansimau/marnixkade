package haltest

import (
	"testing"

	"github.com/dansimau/home-automation/pkg/hal"
	"github.com/dansimau/home-automation/pkg/hassws"
	"gotest.tools/v3/assert"
)

func NewClientServer(t *testing.T) (*hal.Connection, *hassws.Server, func()) {
	// Create test server
	server, err := hassws.NewServer()
	assert.NilError(t, err)

	// Create client and connection
	client := hassws.NewWebsocketAPI(hassws.ClientConfig{
		Host:  server.ListenAddress(),
		Token: "test-token",
	})
	conn := hal.NewConnection(client)

	// Create test entity and register it
	entity := hal.NewEntity("test.entity")
	conn.RegisterEntities(entity)

	// Start connection
	err = conn.Start()
	assert.NilError(t, err)

	return conn, server, func() {
		conn.HomeAssistant().Close()
		server.Close()
	}
}
