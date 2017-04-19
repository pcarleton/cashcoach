package main

import (
	"net/url"
	"errors"
	"net/http"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	"log"
	"google.golang.org/api/plus/v1"
	"encoding/gob"
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
	ID, DisplayName, ImageURL string
}

func loginHandler(w http.ResponseWriter, r *http.Request) *appError {
	sessionID := uuid.NewV4().String()

	oauthFlowSession, err := cookieStore.New(r, sessionID)
	if err != nil {
		return appErrorf(err, "could not create oauth session: %v", err)
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		return appErrorf(err, "invalid redirect URL: %v", err)
	}
	oauthFlowSession.Values[oauthFlowRedirectKey] = redirectURL

	if err := oauthFlowSession.Save(r, w); err != nil {
		return appErrorf(err, "could not save session: %v", err)
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	authUrl := conf.AuthCodeURL(sessionID, oauth2.ApprovalForce,
		oauth2.AccessTypeOnline)
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
	oauthFlowSession, err := cookieStore.Get(r, r.FormValue("state"))
	if err != nil {
		return appErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
	// Validate this callback request came from the app.
	if !ok {
		return appErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	tok, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return appErrorf(err, "could not get auth token: %v", err)
	}

	session, err := cookieStore.New(r, sessionName)
	if err != nil {
		return appErrorf(err, "could not get default session: %v", err)
	}

	person, err := fetchProfile(context.Background(), tok)

	if err != nil {
		return appErrorf(err, "could not fetch userinfo")
	}

	log.Printf("Person: %v\n", person)
	session.Values[oauthTokenSessionKey] = tok
	session.Values[emailKey] = stripProfile(person)
	if err := session.Save(r, w); err != nil {
		return appErrorf(err, "could not save session: %v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}


// fetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func fetchProfile(ctx context.Context, tok *oauth2.Token) (*plus.Person, error) {
	client := oauth2.NewClient(ctx, conf.TokenSource(ctx, tok))
	plusService, err := plus.New(client)
	if err != nil {
		return nil, err
	}
	return plusService.People.Get("me").Do()
}

// logoutHandler clears the default session.
func logoutHandler(w http.ResponseWriter, r *http.Request) *appError {
	session, err := cookieStore.New(r, sessionName)
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
	session, err := cookieStore.Get(r, sessionName)
	if err != nil {
		return nil
	}
	tok, ok := session.Values[oauthTokenSessionKey].(*oauth2.Token)
	if !ok || !tok.Valid() {
		return nil
	}
	profile, ok := session.Values[emailKey].(*Profile)
	if !ok {
		return nil
	}
	return profile
}


// stripProfile returns a subset of a plus.Person.
func stripProfile(p *plus.Person) *Profile {
	return &Profile{
		ID:          p.Id,
		DisplayName: p.DisplayName,
		ImageURL:    p.Image.Url,
	}
}