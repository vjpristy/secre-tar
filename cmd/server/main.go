package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/vjpristy/secre-tar/internal/config"
	"github.com/vjpristy/secre-tar/internal/network"
)

var (
	clients    = make(map[*network.Connection]bool)
	clientsMux sync.Mutex
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := network.HandleConnections(w, r)
		if err != nil {
			log.Println("Error handling connection:", err)
			return
		}
		defer conn.Close()

		clientsMux.Lock()
		clients[conn] = true
		clientsMux.Unlock()

		for {
			message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				clientsMux.Lock()
				delete(clients, conn)
				clientsMux.Unlock()
				break
			}
			log.Printf("Received message: %s", message)
			broadcastMessage(message, conn)
		}
	})

	log.Printf("Starting server on %s\n", cfg.ServerAddress)
	err = http.ListenAndServe(cfg.ServerAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func broadcastMessage(message []byte, sender *network.Connection) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for client := range clients {
		if client != sender {
			err := client.WriteMessage(message)
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
