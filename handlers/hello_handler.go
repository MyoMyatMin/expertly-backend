package handlers

import (
	"net/http"
)

// write a handler returning a simple string
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
