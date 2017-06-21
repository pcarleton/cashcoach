package main

import (
	"github.com/spf13/viper"
	"log"
	"fmt"
	"os"

	"github.com/pcarleton/cashcoach/api/plaid"
)

func main() {
	viper.BindEnv("CLIENT_ID")
	viper.BindEnv("CLIENT_SECRET")
	viper.BindEnv("ACCESS_TOKEN")


	client := plaid.NewClient(
		viper.GetString("client_id"),
		viper.GetString("client_secret"),
		plaid.DevURL)

	accessToken := viper.GetString("access_token")

	listTransactions(accessToken, client)
}

func listTransactions(accessToken string, client plaid.Client) {
	accessToken = updateToken(accessToken, client)
	resp, err := client.Transactions(accessToken,
		"2017-04-01", "2017-04-07")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Resp: %+v", resp)
}

func updateToken(accessToken string, client plaid.Client) string {
	resp, err := client.UpdateAccessToken(accessToken)

	if err != nil {
		log.Fatalf("Error upgrading token: %v", err)
	}

	fmt.Printf("New token: %v", resp.AccessToken)

	os.Setenv("access_token", resp.AccessToken)
  return resp.AccessToken
}
