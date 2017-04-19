package main

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"net/url"
)

var config *Config

func init() {
	v, err := parseConfig()
	if err != nil {
		panic(fmt.Errorf("Error parsing config file (%s): %s", v.ConfigFileUsed(), err))
	}
	config, err = makeConfig(v)
	if err != nil {
		panic(fmt.Errorf("Error setting up environment: %s", err))
	}
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


func meHandler(w http.ResponseWriter, r *http.Request) *appError {
	profile := profileFromSession(r)

	var js []byte
	var err error
	if profile == nil {
		js, err = json.Marshal(map[string]string{"error": "Unauthenticated"})
	} else {
		js, err = json.Marshal(map[string]string{"name": profile.DisplayName})
	}

	if err != nil {
		return appErrorf(err, "Problem marshaling JS")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func entryPointHandler(w http.ResponseWriter, r *http.Request) *appError {
	http.ServeFile(w, r, "app/app.html")
	return nil
}

func main() {
	http.Handle("/", appHandler(logHandler("<a href='login?redirect=restricted'>Login</a>")))
	http.Handle("/app", restricted(entryPointHandler))
	http.Handle("/oauth2callback", appHandler(oauthCallbackHandler))
	http.Handle("/login", appHandler(loginHandler))
	http.Handle("/logout", appHandler(logoutHandler))
	http.Handle("/me", appHandler(meHandler))

	log.Fatal(http.ListenAndServe(":5000", nil))
}