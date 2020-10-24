package spotify

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const authURL = "https://accounts.spotify.com/authorize"
const tokenURL = "https://accounts.spotify.com/api/token"
const cookieExpiration = 1 * 24 * time.Hour // 1 day

func NewOauthConf(clientID string, clientSecret string, redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURI,
		Scopes:      []string{"user-read-private", "user-read-email", "streaming", "user-modify-playback-state"},
	}
}

func GenState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func MkCookie(state string) *http.Cookie {
	expiration := time.Now().Add(cookieExpiration)
	return &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
}
