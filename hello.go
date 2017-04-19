package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"errors"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "session-name"
	stateKey = "session-id"
	oauthCallbackUrl = "http://localhost:5000/oauth2callback"
)

var MissingConfigError = errors.New("Missing client secrets config path in environment")
var BadConfigError = errors.New("Unparseable config data")

var envVars = map[string]string {
	"CLIENT_SECRETS_PATH": "Path to OAuth client id and secrets json file",
	//"COOKIE_SECRET": "Hex encoded 32 byte secret for encrypting cookies",
}

var conf *oauth2.Config

// TODO: Read the secret from environment (use securecookie.GenerateRandomKey())
var store = sessions.NewCookieStore([]byte("secret"))

// Credentials which stores google ids.
type Credentials struct {
	Web struct {
		ClientId string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"web"`
}


func exitWithError(err error) {
	fmt.Printf("Exiting due to error: %v\n", err)
	os.Exit(1)
}

func checkEnv() {
	missing := make([]string, 0)
	for envKey, desc := range envVars {
		_, ok := os.LookupEnv(envKey)
		if !ok {
			missing = append(missing, fmt.Sprintf("%s: %s", envKey, desc))
		}
	}
	if len(missing) > 0 {
		fmt.Println("Missing environment variables:")
		for _, v := range missing {
			fmt.Println(v)
		}
		os.Exit(1)
	}
}



func init() {
	checkEnv()
	credsPath := os.Getenv("CLIENT_SECRETS_PATH")

	var cred Credentials
	file, err := ioutil.ReadFile(credsPath)
	if err != nil {
		exitWithError(MissingConfigError)
	}
	err = json.Unmarshal(file, &cred)

	if err != nil {
		exitWithError(BadConfigError)
	}

	conf = &oauth2.Config{
		ClientID:     cred.Web.ClientId,
		ClientSecret: cred.Web.ClientSecret,
		RedirectURL:  oauthCallbackUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func logHandler(msg string) (func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request from %v\n", r.RemoteAddr)
		w.Write([]byte(msg))
	}
}



func main() {
	http.HandleFunc("/", logHandler("hello") )
	http.HandleFunc("/oauth2callback", logHandler("oauth2callback") )
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":5000", nil))
}