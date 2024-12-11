package websocket

type Hub struct {
	clients    map[*socketClient]bool // Registered Clients
	broadcast  chan []byte            // Inbound messages from clients
	register   chan *socketClient     // Register requests from clients
	unregister chan *socketClient     // Unregister requests from clients
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *socketClient),
		unregister: make(chan *socketClient),
		clients:    make(map[*socketClient]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
