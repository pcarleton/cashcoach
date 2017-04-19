package main

import (
	"golang.org/x/oauth2"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/api/plus/v1"
	"golang.org/x/oauth2/google"
	"encoding/hex"
)


func parseConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")  // call multiple times to add many search paths
	v.AddConfigPath("$HOME/.cashcoach")  // call multiple times to add many search paths
	v.AddConfigPath(".")               // optionally look for config in the working directory
	return v, v.ReadInConfig()
}

type Config struct {
	OAuthConfig *oauth2.Config
	SessionStore sessions.Store
}

func makeConfig(v *viper.Viper) (*Config, error) {
	// TODO: Consider passing this in
	scopes := []string{plus.UserinfoProfileScope}

	oauthConf := &oauth2.Config{
		ClientID:     v.GetString("client_id"),
		ClientSecret: v.GetString("client_secret"),
		RedirectURL:  v.GetString("oauth_callback_url"),
		Scopes: scopes,
		Endpoint: google.Endpoint,
	}

	secret, err := hex.DecodeString(v.GetString("cookie_secret"))
	if err != nil {
		return nil, err
	}
	sessionStore := sessions.NewCookieStore(secret)

	return &Config{oauthConf, sessionStore}, nil
}