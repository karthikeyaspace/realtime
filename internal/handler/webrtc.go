package handler

import "net/http"

// /webrtc

func WebRTCHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WEB RTC HANDLER"))
}
