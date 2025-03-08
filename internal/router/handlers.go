package router

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Log a message on the server side
	fmt.Println("You hit the endpoint:", r.URL.Path)

	// Send a response to the client
	w.Write([]byte("I'm Mister Derp!"))
}
