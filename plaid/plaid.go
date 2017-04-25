package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	devURL = "https://development.plaid.com"
)

func main() {
	viper.BindEnv("client_id")
	viper.BindEnv("client_secret")
	viper.BindEnv("access_token")

	client := Client{
		baseURL:devURL,
		client: new(http.Client),
		clientID: viper.GetString("client_id"),
		secret: viper.GetString("client_secret"),
	}

	accessToken := viper.GetString("access_token")
	client.Transactions(accessToken,
		"2017-04-01", "2017-04-07")




}

func updateToken(accessToken string, client Client) {
	resp, err := client.UpdateAccessToken(accessToken)

	if err != nil {
		log.Fatalf("Error upgrading token: %v", err)
	}

	fmt.Printf("New token: %s", resp.AccessToken)

	os.Setenv("access_token", resp.AccessToken)
}


type Client struct {
	baseURL string
	client *http.Client
	clientID string
	secret string
}

//type ClientInfo struct {
//	ClientID string`json:"client_id"`
//	Secret string`json:"secret"`
//}
//
//type AccountToken struct {
//	AccessToken string`json:"access_token"`
//}
//
//type StringDateRange struct {
//	StartDate string`json:"start_date"`
//	EndDate string `json:"end_date"`
//}


type UpdateAccessTokenRequest struct {
	ClientID string`json:"client_id"`
	Secret string`json:"secret"`
	AccessToken string`json:"access_token"`
}

type UpdateAccessTokenResponse struct {
	AccessToken string`json:"access_token"`
	RequestID string`json:"request_id"`
}

func (c *Client) UpdateAccessToken(accessToken string) (UpdateAccessTokenResponse, error) {
	endpoint := "/item/access_token/update_version"

	req := UpdateAccessTokenRequest{c.clientID, c.secret, accessToken}

	resp := UpdateAccessTokenResponse{}
	rawResp, err := c.post(endpoint, req)

	if err != nil {
		return resp, err
	}

	if err := json.Unmarshal(rawResp, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Client) post(endpoint string, v interface{}) ([]byte, error) {

	jsonText, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	target := c.baseURL + endpoint

	resp, err := c.client.Post(target, "application/json",
		bytes.NewReader(jsonText))

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil

}



type TransactionRequest struct {
	ClientID string`json:"client_id"`
	Secret string`json:"secret"`
	AccessToken string`json:"access_token"`
	StartDate string`json:"start_date"`
	EndDate string `json:"end_date"`
}


type Transaction struct {
	ID string`json:"transaction_id"`
	AccountID string`json:"account_id"`
	Category []string`json:"category"`
	CategoryID string`json:"category_id"`
	Type string`json:"transaction_type"`
	Amount string`json:"amount"`
	Date string`json:"date"`
	Pending bool`json:"pending"`
	AccountOwner string`json:"account_owner"`
	Name string `json:"string"`
}



func (c *Client) Transactions(accessToken, startDate, endDate string) error {
	endpoint := "/transactions/get"

	jsonText, err := json.Marshal(TransactionRequest{
		ClientID: c.clientID,
		Secret: c.secret,
		AccessToken: accessToken,
		StartDate: startDate,
		EndDate: endDate,
	})

	if err != nil {
		return err
	}

	reader := bytes.NewReader(jsonText)

	transUrl := c.baseURL + endpoint
	resp, err := c.client.Post(transUrl, "application/json", reader)

	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Got: %s", respBody)
	return nil
}