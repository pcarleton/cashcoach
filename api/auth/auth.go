package auth

import (
	"crypto/rsa"
	"encoding/json"
	"encoding/gob"
	"encoding/hex"
	"fmt"
  "time"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
)

func init() {
	// Gob encoding for gorilla/sessions
	gob.Register(&Profile{})
}

// TODO: probably put this in storage
type Profile struct {
	Email, DisplayName, ImageURL string
}

const (
	googleEndPoint = "https://www.googleapis.com/oauth2/v1/certs"
	sessionName          = "default_session"
	emailKey             = "email"
)

var googleCerts map[string]*rsa.PublicKey

type SessionHandler struct {
  store sessions.Store
}

func (s *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request, profile *Profile) error {
	session, err := s.store.New(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[emailKey] = profile
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

// profileFromSession retreives the Google+ profile from the default session.
// Returns nil if the profile cannot be retreived (e.g. user is logged out).
func (s *SessionHandler) ProfileFromSession(r *http.Request) *Profile {
	session, err := s.store.Get(r, sessionName)

	if err != nil {
		return nil
	}
	profile, ok := session.Values[emailKey].(*Profile)
	if !ok {
		return nil
	}
	return profile
}

func CreateSessionHandler(cookieSecret string) (*SessionHandler, error) {
	secret, err := hex.DecodeString(cookieSecret)
	if err != nil {
		return nil, err
	}
	sessionStore := sessions.NewCookieStore(secret)
	sessionStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int((time.Hour * 3).Seconds()),
	}

  return &SessionHandler{sessionStore}, nil
}


func getGoogleCerts() (map[string]*rsa.PublicKey, error) {
	if googleCerts == nil {
		resp, err := http.Get(googleEndPoint)
		if err != nil {
			return nil, err
		}

		rawCerts := make(map[string]string)

		respBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(respBody, &rawCerts)

		if err != nil {
			return nil, err
		}

		googleCerts = make(map[string]*rsa.PublicKey)

		for k, v := range rawCerts {
			rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(v))

			if err != nil {
				return nil, err
			}
			googleCerts[k] = rsaKey
		}

	}
	return googleCerts, nil
}

func VerifyToken(tokenString string, certs map[string]*rsa.PublicKey) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Only accept RSA signed JWT's
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		keyID := token.Header["kid"].(string)
		rsaKey := certs[keyID]
		return rsaKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token: %v", token)
	}

	return token, nil
}

func VerifyGoogleJwt(tokenString string) (*jwt.Token, error) {
  // TODO: Invalidate these certs periodically because they go stale
	certs, err := getGoogleCerts()
	if err != nil {
		return nil, err
	}

	return VerifyToken(tokenString, certs)
}

func ProfileFromJwt(token *jwt.Token) (*Profile, error) {
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

