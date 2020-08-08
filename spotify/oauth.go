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

func NewOauthConf(clientID string, clientSecret string, redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURI,
		Scopes:      []string{"user-read-private", "user-read-email"},
	}
}

func GenState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func MkCookie(state string) *http.Cookie {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	return &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
}
