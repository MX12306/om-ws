package main

import (
	"log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			go log.Print("Register")
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			go log.Print("Unregister")
		case m := <-h.broadcast:
			go log.Print("Broadcast")
			for c := range h.connections {
				select {
				case c.sendBinary <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
