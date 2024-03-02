package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients   = make(map[*websocket.Conn]struct{})
	clientsMu sync.Mutex
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket connection requested")

	// Check the origin of the request to allow connections only from specific origins
	// You may want to change this to your specific needs
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = struct{}{}
	clientsMu.Unlock()

	log.Println("WebSocket connection established")

	for {
		// Read message from the browser
		log.Println("Waiting to read message from the browser...")
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			return
		}

		// Print the received message for debugging
		log.Println("Received message:", string(p))

		// Broadcast message to all clients
		broadcastMessage(string(p))
	}
}

func broadcastMessage(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		log.Println("Broadcasting message to client...")
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error broadcasting message to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling index request...")
	http.ServeFile(w, r, "index.html")
}

func main() {
	log.Println("Starting WebSocket server on port 8080")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}