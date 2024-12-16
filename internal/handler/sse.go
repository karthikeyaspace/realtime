// sse.go
package handler

import (
	"fmt"
	"net/http"
	"time"
)

func SSEHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			fmt.Println("Connection closed from client")
			return
		case t := <-ticker.C:
			fmt.Fprintf(w, "data: %s\n\n", t.Format(time.RFC3339))
			flusher.Flush()
		}
	}
}
