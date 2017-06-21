package main

import (
	"encoding/hex"
	"time"

	"github.com/gorilla/sessions"
	"github.com/pcarleton/cashcoach/plaid"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/plus/v1"
	"gopkg.in/mgo.v2"
)

func parseConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")           // call multiple times to add many search paths
	v.AddConfigPath("$HOME/.cashcoach") // call multiple times to add many search paths
	v.AddConfigPath(".")                // optionally look for config in the working directory

	return v, v.ReadInConfig()
}

type Config struct {
	OAuthConfig  *oauth2.Config
	SessionStore sessions.Store
	Storage
	Plaid plaid.Client
}

func getFakeStorage(v *viper.Viper) Storage {
	return &FakeStorage{[]string{v.GetString("plaid.access_token")}}
}

func getMongoStorage(v *viper.Viper) (Storage, error) {
	session, err := mgo.Dial(v.GetString("mongo.host"))

	if err != nil {
		return nil, err
	}
	storage := &MongoStorage{session}

	return storage, nil
}

func makeConfig(v *viper.Viper) (*Config, error) {
	// TODO: Consider passing this in
	scopes := []string{plus.UserinfoProfileScope, plus.UserinfoEmailScope}

	oauthConf := &oauth2.Config{
		ClientID:     v.GetString("google.client_id"),
		ClientSecret: v.GetString("google.client_secret"),
		RedirectURL:  v.GetString("oauth_callback_url"),
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	secret, err := hex.DecodeString(v.GetString("cookie_secret"))
	if err != nil {
		return nil, err
	}
	sessionStore := sessions.NewCookieStore(secret)
	sessionStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int((time.Hour * 3).Seconds()),
	}

	plaidClient := plaid.NewClient(
		v.GetString("plaid.client_id"),
		v.GetString("plaid.client_secret"),
		plaid.DevURL)

	storage, err := getMongoStorage(v)

	if err != nil {
		return nil, err
	}

	return &Config{oauthConf, sessionStore, storage, plaidClient}, nil
}
