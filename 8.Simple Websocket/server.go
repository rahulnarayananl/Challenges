package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		// Wait for a message from the client
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Print received message
		log.Printf("Received message: %s\n", msg)

		// Echo the message back to the client
		err = conn.WriteMessage(messageType, []byte("Server's Message "+string(msg)))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}

func main() {
	// Setup WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the server
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
