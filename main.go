package main

import "net/http"

type Client struct {
	broker *Broker
	send   chan []byte
}
type Broker struct {
	// Inbound messages from clients.
	broadcast chan []byte

	// Registered clients.
	clients map[*Client]bool

	// Register clients.
	register chan *Client
}

func NewClient() *Client {
	return &Client{}
}

func NewBroker() *Broker {
	return &Broker{
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
		register:  make(chan *Client),
	}
}

func (b *Broker) run() {

	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
		case message := <-b.broadcast:
			for client := range b.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(b.clients, client)
				}

			}
		}
	}

}

func main() {
	broker := NewBroker()
	go broker.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(broker, w, r)
	})
	http.ListenAndServe("1111", nil)
}

func serveWs(broker *Broker, w http.ResponseWriter, r *http.Request) {}
func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}
