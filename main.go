package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"franquel.in/bajidspotifyserver/config"
	"franquel.in/bajidspotifyserver/gcp"
	"franquel.in/bajidspotifyserver/spotify"
)

func main() {
	fmt.Println("Starting Bajid server!")

	gcpProjectId := config.RequireEnvVar("GCP_PROJECT_ID")
	spotifyClientID := config.RequireEnvVar("SPOTIFY_CLIENT_ID")
	authorizedRedirectURI := config.RequireEnvVar("SPOTIFY_REDIRECT_URI")

	sm, err := gcp.NewSecretManager(gcpProjectId)

	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	spotifyClientSecret, err := sm.GetSecret(ctx, "SPOTIFY_CLIENT_SECRET")

	oauthConf := spotify.NewOauthConf(spotifyClientID, spotifyClientSecret, authorizedRedirectURI)

	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/auth/spotify/login", handleOauthLogin(oauthConf))
	http.HandleFunc("/auth/spotify/callback", handleOauthCallback(oauthConf))

	port := config.RequireEnvVar("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Bajid listening on port %s", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	log.Fatal(err)
}

func handleOauthLogin(oauthConf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := spotify.GenState()
		cookie := spotify.MkCookie(state)
		http.SetCookie(w, cookie)
		u := oauthConf.AuthCodeURL(state)
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func handleOauthCallback(oauthConf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState, _ := r.Cookie("oauthstate")

		if r.FormValue("state") != oauthState.Value {
			log.Println("invalid oauth google state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		ctx := context.Background()
		token, err := oauthConf.Exchange(ctx, code)

		if err != nil {
			log.Printf("error while exchanging token: %s\n", err)
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)

		cookie := &http.Cookie{
			Name:    "bajid-spotify-token",
			Value:   token.AccessToken,
			Expires: expiration,
		}

		http.SetCookie(w, cookie)

		w.Write([]byte("OK"))
	}
}
