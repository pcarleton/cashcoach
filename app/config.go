package main

import (
	"encoding/hex"

	"github.com/gorilla/sessions"
	"github.com/pcarleton/hello/plaid"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/plus/v1"
)

func parseConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")           // call multiple times to add many search paths
	v.AddConfigPath("$HOME/.cashcoach") // call multiple times to add many search paths
	v.AddConfigPath(".")                // optionally look for config in the working directory

	// TODO: Store this somewhere
	v.BindEnv("access_token")
	return v, v.ReadInConfig()
}

type Config struct {
	OAuthConfig  *oauth2.Config
	SessionStore sessions.Store
	Storage
	Plaid plaid.Client
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

	plaidClient := plaid.NewClient(
		v.GetString("plaid.client_id"),
		v.GetString("plaid.client_secret"),
		plaid.DevURL)

	storage := FakeStorage{[]string{v.GetString("access_token")}}

	return &Config{oauthConf, sessionStore, &storage, plaidClient}, nil
}
