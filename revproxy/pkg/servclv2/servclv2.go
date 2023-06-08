package servcl2

import (
	"fmt"
	"net/http"

)


func HandleMessage(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	log.Printf("Got connection: %s", r.Proto)
	// Send a message back to the client
	w.Write([]byte("Hello"))
}

func HandlePusher(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	log.Printf("Got connection: %s", r.Proto)

	// Handle 2nd request, must be before push to prevent recursive calls.
	// Don't worry - Go protect us from recursive push by panicking.
	if r.URL.Path == "/2nd" {
		log.Println("Handling 2nd")
		w.Write([]byte("Hello Again!"))
		return
	}

	// Handle 1st request
	log.Println("Handling 1st")

	// Server push must be before response body is being written.
	// In order to check if the connection supports push, we should use
	// a type-assertion on the response writer.
	// If the connection does not support server push, or that the push
	// fails we just ignore it - server pushes are only here to improve
	// the performance for HTTP/2 clients.
	pusher, ok := w.(http.Pusher)
	if !ok {
		log.Println("Can't push to client")
	} else {
		err := pusher.Push("/2nd", nil)
		if err != nil {
			log.Printf("Failed push: %v", err)
		}
	}

	// Send response body
	w.Write([]byte("Hello"))
}


