package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"franquel.in/bajidspotifyserver/config"
	"franquel.in/bajidspotifyserver/gcp"
)

func main() {
	fmt.Println("Starting Bajid server!")

	gcpProjectId := config.RequireEnvVar("GCP_PROJECT_ID")

	sm, err := gcp.NewSecretManager(gcpProjectId)

	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	_, err = sm.GetSecret(ctx, "SPOTIFY_CLIENT_SECRET")

	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))

	port := config.RequireEnvVar("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Bajid listening on port %s", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	log.Fatal(err)
}
