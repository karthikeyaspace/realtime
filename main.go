package main

import (
	"log"
	"net/http"

	"github.com/karthikeyaspace/realtime/internal/config"
	"github.com/karthikeyaspace/realtime/internal/handler"
	"github.com/karthikeyaspace/realtime/internal/middleware"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Start() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /sse", handler.SSEHandler)
	router.HandleFunc("GET /ws", handler.WebSocketHandler)
	router.HandleFunc("GET /webrtc", handler.WebRTCHandler)

	middleware := middleware.Cors(middleware.Logger(router))

	server := &http.Server{

		Addr:    s.addr,
		Handler: middleware,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server:", err)
	}

	return nil
}

func main() {
	port := config.NewConfig().Port
	server := NewAPIServer(port)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	log.Printf("Starting API server at http://localhost:%s", port)
}
