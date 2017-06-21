package main

import (
	"encoding/gob"
	"net/http"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

const (
	sessionName          = "default_session"
	emailKey             = "email"
	oauthFlowRedirectKey = "oauth_redirect"
	oauthTokenSessionKey = "oauth_token"
)

func init() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	gob.Register(&Profile{})
}

type Profile struct {
	Email, DisplayName, ImageURL string
}

func createSession(w http.ResponseWriter, r *http.Request, profile *Profile) error {
	session, err := config.SessionStore.New(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[emailKey] = profile
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func profileFromJwt(token *jwt.Token) (*Profile, error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("Invalid JWT claims: %+v", token.Claims)
	}

	email, ok := claims["email"].(string)
	name, ok := claims["name"].(string)
	image, ok := claims["picture"].(string)

	if !ok {
		return nil, fmt.Errorf("Invalid JWT claims: %+v", token.Claims)
	}

	return &Profile{
		Email:       email,
		DisplayName: name,
		ImageURL:    image,
	}, nil

}

// profileFromSession retreives the Google+ profile from the default session.
// Returns nil if the profile cannot be retreived (e.g. user is logged out).
func profileFromSession(r *http.Request) *Profile {
	session, err := config.SessionStore.Get(r, sessionName)
	if err != nil {
		return nil
	}
	profile, ok := session.Values[emailKey].(*Profile)
	if !ok {
		return nil
	}
	return profile
}
