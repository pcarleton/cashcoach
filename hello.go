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
	"google.golang.org/api/plus/v1"
	"encoding/hex"
	"net/url"
)

const (
	oauthCallbackUrl = "http://localhost:5000/oauth2callback"
	secretsPathEnv = "CLIENT_SECRETS_PATH"
	cookieSecretEnv = "COOKIE_SECRET"
)

var MissingConfigError = errors.New("Missing client secrets config path in environment")
var BadConfigError = errors.New("Unparseable config data")

var envVars = map[string]string {
	secretsPathEnv: "Path to OAuth client id and secrets json file",
	cookieSecretEnv: "Hex encoded 32 byte secret for encrypting cookies",
}

var conf *oauth2.Config

// TODO: Read the secret from environment (use securecookie.GenerateRandomKey())
var cookieStore sessions.Store


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

	// Set up OAuth parameters
	credsPath := os.Getenv(secretsPathEnv)

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
			plus.UserinfoProfileScope,
		},
		Endpoint: google.Endpoint,
	}

	// Set up Cookie store with secret
	encoded := os.Getenv(cookieSecretEnv)
	secret, err := hex.DecodeString(encoded)
	if err != nil {
		exitWithError(err)
	}
	cookieStore = sessions.NewCookieStore(secret)
}

func logHandler(msg string) (func(w http.ResponseWriter, r *http.Request) *appError) {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		log.Printf("request from %v\n", r.RemoteAddr)
		w.Write([]byte(msg))
		return nil
	}
}


// http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

type restricted func(http.ResponseWriter, *http.Request) *appError

func makeLoginUrl(path string) string {
	vals := url.Values{}
	vals.Set("redirect", path)

	return "/login?" + vals.Encode()
}

func (fn restricted) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	profile := profileFromSession(r)
	if profile == nil {
		loginUrl, err := validateRedirectURL(makeLoginUrl(r.URL.Path))

		if err != nil {
			log.Printf("Error redirecting to login %#v", err)
			http.Error(w, "Problem redirecting to login", 500)
			return
		}

		http.Redirect(w, r, loginUrl, 302)
	}

	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}


func main() {
	http.Handle("/", appHandler(logHandler("<a href='login?redirect=restricted'>Login</a>")))
	http.Handle("/restricted", restricted(logHandler("hello")))
	http.Handle("/oauth2callback", appHandler(oauthCallbackHandler))
	http.Handle("/login", appHandler(loginHandler))
	http.Handle("/logout", appHandler(logoutHandler))
	log.Fatal(http.ListenAndServe(":5000", nil))
}