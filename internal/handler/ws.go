package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// https://vishalrana9915.medium.com/understanding-websockets-in-depth-6eb07ab298b3

// client sends http get with upgrade header
// server responds with 101 switching protocols after validating upgrade header
// a handshake is done between client and server
// after handshake, client and server can send messages to each other
// client headers - connection: upgrade, upgrade: websocket, sec-websocket-key: base64 encoded key, sec-websocket-version: 13

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients   = make(map[*websocket.Conn]string)
	clientMux sync.RWMutex
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	clientMux.Lock()
	clients[conn] = username
	clientMux.Unlock()

	broadcast(Message{
		Username: "system",
		Message:  username + " has joined the chat",
	})

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			break
		}

		msg.Username = username
		broadcast(msg)
	}

	clientMux.Lock()
	delete(clients, conn)
	clientMux.Unlock()

	broadcast(Message{
		Username: "system",
		Message:  username + " has left the chat",
	})

}

func broadcast(msg Message) {
	clientMux.RLock()
	defer clientMux.RUnlock()

	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Failed to write message to client: %v", err)
			client.Close()
		}
	}
}
