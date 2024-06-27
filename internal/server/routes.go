package server

import "net/http"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(pongMessage))
}
