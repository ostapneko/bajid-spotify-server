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

const cookieExpiration = 7 * 24 * time.Hour // 1 week
const bajidSpotifyTokenKey = "bajid-spotify-token"
const playerPath = "/player"
const loginPath = "/login"

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

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/favicon.ico")
	})

	http.HandleFunc(loginPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	})

	http.HandleFunc(playerPath, func(w http.ResponseWriter, r *http.Request) {
		if checkBajidSpotifyCookie(w, r) {
			http.ServeFile(w, r, "public/player.html")
		}
	})

	http.HandleFunc("/auth/spotify/login", handleOauthLogin(oauthConf))

	http.HandleFunc("/auth/spotify/callback", handleOauthCallback(oauthConf))

	http.Handle("/css/", http.FileServer(http.Dir("./public")))

	http.Handle("/js/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/", handleWelcome)

	port := config.RequireEnvVar("PORT")

	log.Printf("Bajid listening on port %s", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	log.Fatal(err)
}

func handleWelcome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if checkBajidSpotifyCookie(w, r) {
		http.Redirect(w, r, playerPath, http.StatusTemporaryRedirect)
	}
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
			http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
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

		expiration := time.Now().Add(cookieExpiration)

		cookie := &http.Cookie{
			Name:    bajidSpotifyTokenKey,
			Value:   token.AccessToken,
			Expires: expiration,
			Path:    "/",
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, playerPath, http.StatusTemporaryRedirect)
	}
}

// checkBajidSpotifyCookie returns true if there is a valid bajidSpotifyToken and rediretcts
func checkBajidSpotifyCookie(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie(bajidSpotifyTokenKey)

	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return false
	}

	return true
}
