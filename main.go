package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"franquel.in/bajidspotifyserver/config"
	"franquel.in/bajidspotifyserver/gcp"
)

func main() {
	fmt.Println("Starting Bajid server!")

	gcpProjectId := config.RequireEnvVar("GCP_PROJECT_ID")
	spotifyClientID := config.RequireEnvVar("SPOTIFY_CLIENT_ID")
	authorizedRedirectURI := config.RequireEnvVar("SPOTIFY_REDIRECT_URI")
	loginRedirectURL := spotifyRedirectURL(spotifyClientID, authorizedRedirectURI)
	log.Printf("Login redirect URL: %s\n", loginRedirectURL)

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
	http.Handle("/login", http.RedirectHandler(loginRedirectURL, 303))

	port := config.RequireEnvVar("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Bajid listening on port %s", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	log.Fatal(err)
}

const spotifyScopes = "user-read-private user-read-email"

func spotifyRedirectURL(clientID string, redirectURI string) string {
	encScopes := url.QueryEscape(spotifyScopes)
	encRedirectURI := url.QueryEscape(redirectURI)

	return "https://accounts.spotify.com/authorize" +
		"?response_type=code" +
		"&client_id=" + clientID +
		"&scope=" + encScopes +
		"&redirect_uri=" + encRedirectURI
}
