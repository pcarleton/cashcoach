package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pcarleton/cashcoach/api/auth"
	"github.com/pcarleton/cashcoach/api/storage"
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

type authorizedHandler func(*auth.Profile, http.ResponseWriter, *http.Request) *appError

func handleAuth(handler authorizedHandler) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		profile := config.Sessions.ProfileFromSession(r)
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

func meHandler(profile *auth.Profile, w http.ResponseWriter, r *http.Request) *appError {
	return respondJson(w, profile)
}

func transactionsHandler(profile *auth.Profile, w http.ResponseWriter, r *http.Request) *appError {
	person, err := config.Get(profile.DisplayName)

	if err != nil {
		return appErrorf(err, "Error loading bank info from database.")
	}

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	referenceTime := "2006-01-02"

	transactions, err := config.Plaid.Transactions(
		person.Accounts[0].Token, lastMonth.Format(referenceTime), now.Format(referenceTime))

	if err != nil {
		return appErrorf(err, "Error getting transactions")
	}

	return respondJson(w, transactions)
}

type JwtRequest struct {
	IDToken string `json:"idtoken"`
}

func unmarshal(v interface{}, r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(reqBody, v)

	if err != nil {
		return err
	}

	return nil
}

func jwtHandler(w http.ResponseWriter, r *http.Request) *appError {
	req := new(JwtRequest)

	err := unmarshal(req, r)

	if err != nil {
		return appErrorf(err, "bad request")
	}

	token, err := auth.VerifyGoogleJwt(req.IDToken)

	if err != nil {
		return appErrorf(err, "bad jwt")
	}

	profile, err := auth.ProfileFromJwt(token)
	if err != nil {
		return appErrorf(err, "couldn't create profile")
	}

	err = config.Sessions.CreateSession(w, r, profile)
	if err != nil {
		return appErrorf(err, "couldn't create session")
	}

	alreadyExists, err := config.Create(profile.Email)

	if err != nil {
		return appErrorf(err, "problem checking database")
	}

	var message string
	if alreadyExists {
		message = "sign-in success"
	} else {
		message = "new person created"
	}

	return respondJson(w, message)
}

func accountsHandler(profile *auth.Profile, w http.ResponseWriter, r *http.Request) *appError {
	result, err := config.Get(profile.Email)

	if err != nil {
		return appErrorf(err, "couldn't find %s", profile.Email)
	}

	return respondJson(w, result)
}

type AddAccountRequest struct {
	Name        string
	PublicToken string
}

func addAccount(profile *auth.Profile, w http.ResponseWriter, r *http.Request) *appError {
	person, err := config.Get(profile.Email)

	if err != nil {
		return appErrorf(err, "couldn't find %s", profile.Email)
	}

	req := AddAccountRequest{}
	err = unmarshal(&req, r)

	if err != nil {
		return appErrorf(err, "bad request")
	}

	resp, err := config.Plaid.Exchange(req.PublicToken)

	if err != nil {
		return appErrorf(err, "problem exchanging public token")
	}

	acct := storage.Account{Name: req.Name, Token: resp.AccessToken}

	person.Accounts = append(person.Accounts, acct)

	err = config.Update(person)

	if err != nil {
		return appErrorf(err, "problem saving")
	}

	return respondJson(w, "saved new account")
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

	http.Handle("/api/me", appHandler(handleAuth(meHandler)))
	http.Handle("/api/transactions", appHandler(handleAuth(transactionsHandler)))
	http.Handle("/api/jwt", appHandler(jwtHandler))
	http.Handle("/api/accounts", appHandler(handleAuth(accountsHandler)))
	http.Handle("/api/accounts/add", appHandler(handleAuth(addAccount)))

	log.Println("Serving...")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
