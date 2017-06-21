package main

import (
	"github.com/pcarleton/cashcoach/api/auth"
	"github.com/pcarleton/cashcoach/api/plaid"
	"github.com/pcarleton/cashcoach/api/storage"
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
	Sessions     *auth.SessionHandler
	storage.Storage
	Plaid plaid.Client
}

func getFakeStorage(v *viper.Viper) storage.Storage {
	return &storage.FakeStorage{[]string{v.GetString("plaid.access_token")}}
}

func getMongoStorage(v *viper.Viper) (storage.Storage, error) {
	session, err := mgo.Dial(v.GetString("mongo.host"))

	if err != nil {
		return nil, err
	}
	storage := &storage.MongoStorage{session}

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

	plaidClient := plaid.NewClient(
		v.GetString("plaid.client_id"),
		v.GetString("plaid.client_secret"),
		plaid.DevURL)

	sessionHandler, err := auth.CreateSessionHandler(v.GetString("cookie_secret"))

	if err != nil {
		return nil, err
	}

	storage, err := getMongoStorage(v)

	if err != nil {
		return nil, err
	}

	return &Config{oauthConf, sessionHandler, storage, plaidClient}, nil
}
