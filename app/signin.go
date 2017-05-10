package main

import (
	"encoding/gob"
	"errors"
	"net/http"
	"net/url"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/plus/v1"
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

func loginHandler(w http.ResponseWriter, r *http.Request) *appError {
	sessionID := uuid.NewV4().String()

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		return appErrorf(err, "invalid redirect URL: %v", err)
	}

	profile := profileFromSession(r)

	if profile != nil {
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}

	oauthFlowSession, err := config.SessionStore.New(r, sessionID)
	if err != nil {
		return appErrorf(err, "could not create oauth session: %v", err)
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	oauthFlowSession.Values[oauthFlowRedirectKey] = redirectURL

	if err := oauthFlowSession.Save(r, w); err != nil {
		return appErrorf(err, "could not save session: %v", err)
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	authUrl := config.OAuthConfig.AuthCodeURL(sessionID, oauth2.AccessTypeOnline, oauth2.ApprovalForce)
	http.Redirect(w, r, authUrl, http.StatusFound)

	return nil
}

// validateRedirectURL checks that the URL provided is valid.
// If the URL is missing, redirect the user to the application's root.
// The URL must not be absolute (i.e., the URL must refer to a path within this
// application).
func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL must be absolute")
	}
	return path, nil
}

// oauthCallbackHandler completes the OAuth flow, retreives the user's profile
// information and stores it in a session.
func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) *appError {
	oauthFlowSession, err := config.SessionStore.Get(r, r.FormValue("state"))
	if err != nil {
		return appErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
	// Validate this callback request came from the app.
	if !ok {
		return appErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	tok, err := config.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return appErrorf(err, "could not get auth token: %v", err)
	}

	person, err := fetchProfile(context.Background(), tok)

	if err != nil {
		return appErrorf(err, "could not fetch userinfo")
	}

	profile := stripProfile(person)

	err = createSession(w, r, profile)
	if err != nil {
		appErrorf(err, "could not create session: %v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
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

// fetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func fetchProfile(ctx context.Context, tok *oauth2.Token) (*plus.Person, error) {
	client := oauth2.NewClient(ctx, config.OAuthConfig.TokenSource(ctx, tok))
	plusService, err := plus.New(client)
	if err != nil {
		return nil, err
	}
	return plusService.People.Get("me").Do()
}

// logoutHandler clears the default session.
func logoutHandler(w http.ResponseWriter, r *http.Request) *appError {
	session, err := config.SessionStore.New(r, sessionName)
	if err != nil {
		return appErrorf(err, "could not get default session: %v", err)
	}
	session.Options.MaxAge = -1 // Clear session.
	if err := session.Save(r, w); err != nil {
		return appErrorf(err, "could not save session: %v", err)
	}
	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}

// profileFromSession retreives the Google+ profile from the default session.
// Returns nil if the profile cannot be retreived (e.g. user is logged out).
func profileFromSession(r *http.Request) *Profile {
	session, err := config.SessionStore.Get(r, sessionName)
	if err != nil {
		return nil
	}
	// TODO: Should I check something here?
	//tok, ok := session.Values[oauthTokenSessionKey].(*oauth2.Token)
	//if !ok || !tok.Valid() {
	//	return nil
	//}
	profile, ok := session.Values[emailKey].(*Profile)
	if !ok {
		return nil
	}
	return profile
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

// stripProfile returns a subset of a plus.Person.
func stripProfile(p *plus.Person) *Profile {
	return &Profile{
		Email:       p.Emails[0].Value,
		DisplayName: p.DisplayName,
		ImageURL:    p.Image.Url,
	}
}
