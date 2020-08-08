package main

import (
	"fmt"
	"log"
	"net/http"

	"franquel.in/bajidspotifyserver/config"
)

func main() {
	fmt.Println("Starting Bajid server!")

	http.HandleFunc("/", handler)

	port := config.RequireEnvVar("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Bajid listening on port %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	log.Fatal(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("received a request")
	_, _ = fmt.Fprintf(w, "OK")
}
