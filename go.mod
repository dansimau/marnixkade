module github.com/dansimau/marnixkade

go 1.22.10

require github.com/davecgh/go-spew v1.1.1 // indirect

require github.com/dansimau/hal v0.0.0-20241204133341-b5e97b88b9fb

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/dansimau/hal => ../hal
