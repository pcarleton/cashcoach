package main

import (
	"github.com/spf13/viper"
	"log"
	"fmt"
	"os"

	"github.com/pcarleton/cashcoach/plaid"
)

func main() {
	viper.BindEnv("client_id")
	viper.BindEnv("client_secret")
	viper.BindEnv("access_token")


	client := plaid.NewClient(
		viper.GetString("client_id"),
		viper.GetString("secret"),
		plaid.DevURL)

	accessToken := viper.GetString("access_token")

	listTransactions(accessToken, client)
}

func listTransactions(accessToken string, client plaid.Client) {
	//updateToken(accessToken, client)
	resp, err := client.Transactions(accessToken,
		"2017-04-01", "2017-04-07")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Resp: %+v", resp)
}

func updateToken(accessToken string, client plaid.Client) {
	resp, err := client.UpdateAccessToken(accessToken)

	if err != nil {
		log.Fatalf("Error upgrading token: %v", err)
	}

	fmt.Printf("New token: %v", resp.AccessToken)

	os.Setenv("access_token", resp.AccessToken)
}
