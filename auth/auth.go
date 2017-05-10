package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

const (
	googleEndPoint = "https://www.googleapis.com/oauth2/v1/certs"
)

var googleCerts map[string]*rsa.PublicKey

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
	certs, err := getGoogleCerts()
	if err != nil {
		return nil, err
	}

	return VerifyToken(tokenString, certs)
}
