package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pcarleton/cashcoach/auth"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var config *Config

func logHandler(msg string) func(w http.ResponseWriter, r *http.Request) *appError {
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

type MeMsg struct {
	Name string `json:"name"`
}

type ErrorMsg struct {
	Error string `json:"error"`
}

type authorizedHandler func(*Profile, http.ResponseWriter, *http.Request) *appError

func handleAuth(handler authorizedHandler) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		profile := profileFromSession(r)
		if profile == nil {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return nil
		}
		return handler(profile, w, r)
	}
}

func respondJson(w http.ResponseWriter, v interface{}) *appError {
	js, err := json.Marshal(v)
	if err != nil {
		return appErrorf(err, "Problem marshaling JS")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func meHandler(profile *Profile, w http.ResponseWriter, r *http.Request) *appError {
	return respondJson(w, profile)
}

type Storage interface {
	GetPerson(string) (*Person, error)
}

type FakeStorage struct {
	Tokens []string
}

func (f *FakeStorage) GetPerson(email string) (*Person, error) {
	p := &Person{
		Email: email,
		Banks: []Bank{Bank{"bank1", f.Tokens[0]}},
	}
	return p, nil
}

type MongoStorage struct {
	Session *mgo.Session
}

func (s *MongoStorage) GetPerson(email string) (*Person, error) {
	c := s.Session.DB("test").C("people")

	result := new(Person)
	err := c.Find(bson.M{"name": email}).One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func transactionsHandler(profile *Profile, w http.ResponseWriter, r *http.Request) *appError {
	person, err := config.GetPerson(profile.DisplayName)

	if err != nil {
		return appErrorf(err, "Error loading bank info from database.")
	}

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	referenceTime := "2006-01-02"

	transactions, err := config.Plaid.Transactions(
		person.Banks[0].Token, lastMonth.Format(referenceTime), now.Format(referenceTime))

	if err != nil {
		return appErrorf(err, "Error getting transactions")
	}

	return respondJson(w, transactions)
}

func entryPointHandler(w http.ResponseWriter, r *http.Request) *appError {
	http.ServeFile(w, r, "app/app.html")
	return nil
}

type JwtRequest struct {
	IDToken string `json:"idtoken"`
}

func jwtHandler(w http.ResponseWriter, r *http.Request) *appError {
	req := new(JwtRequest)

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return appErrorf(err, "bad req")
	}

	err = json.Unmarshal(reqBody, &req)

	token, err := auth.VerifyGoogleJwt(req.IDToken)

	if err != nil {
		return appErrorf(err, "bad jwt")
	}

	profile, err := profileFromJwt(token)
	if err != nil {
		return appErrorf(err, "couldn't create profile")
	}

	err = createSession(w, r, profile)
	if err != nil {
		return appErrorf(err, "couldn't create session")
	}

	return respondJson(w, token.Claims)
}

type Bank struct {
	Name  string
	Token string
}

type Person struct {
	Email string
	Banks []Bank
}

func dummyMongoHandler(w http.ResponseWriter, r *http.Request) *appError {
	result, err := config.GetPerson("Ale")

	if err != nil {
		return appErrorf(err, "couldn't find")
	}

	return respondJson(w, result)
}

func main() {
	v, err := parseConfig()
	if err != nil {
		panic(fmt.Errorf("Error parsing config file (%s): %s", v.ConfigFileUsed(), err))
	}
	config, err = makeConfig(v)
	if err != nil {
		panic(fmt.Errorf("Error setting up environment: %s", err))
	}

	http.Handle("/", appHandler(logHandler("<a href='login?redirect=app'>Login</a>")))
	http.Handle("/app", restricted(entryPointHandler))
	http.Handle("/oauth2callback", appHandler(oauthCallbackHandler))
	http.Handle("/login", appHandler(loginHandler))
	http.Handle("/logout", appHandler(logoutHandler))
	// TODO: Delete all the other handlers?
	http.Handle("/api/me", appHandler(handleAuth(meHandler)))
	http.Handle("/api/transactions", appHandler(handleAuth(transactionsHandler)))
	http.Handle("/api/jwt", appHandler(jwtHandler))
	http.Handle("/api/mongo", appHandler(dummyMongoHandler))

	log.Fatal(http.ListenAndServe(":5001", nil))
}
