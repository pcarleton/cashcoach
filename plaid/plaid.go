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

	//updateToken(accessToken, client)
	resp, err := client.Transactions(accessToken,
		"2017-04-01", "2017-04-07")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Resp: %+v", resp)
}

func updateToken(accessToken string, client Client) {
	resp, err := client.UpdateAccessToken(accessToken)

	if err != nil {
		log.Fatalf("Error upgrading token: %v", err)
	}

	fmt.Printf("New token: %v", resp.AccessToken)

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
	AccessToken string`json:"access_token_v1"`
}

type UpdateAccessTokenResponse struct {
	AccessToken string`json:"access_token"`
	RequestID string`json:"request_id"`
}

func (c *Client) UpdateAccessToken(accessToken string) (UpdateAccessTokenResponse, error) {
	endpoint := "/item/access_token/update_version"

	req := UpdateAccessTokenRequest{c.clientID, c.secret, accessToken}

	resp := UpdateAccessTokenResponse{}
	err := c.post(endpoint, req, resp)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Client) post(endpoint string, req interface{}, resp interface{}) error {

	jsonText, err := json.Marshal(req)
	if err != nil {
		return err
	}

	target := c.baseURL + endpoint

	log.Printf("Requesting %s", target)

	log.Printf("Request: %s", jsonText)

	postResp, err := c.client.Post(target, "application/json",
		bytes.NewReader(jsonText))

	if err != nil {
		return err
	}


	respBody, err := ioutil.ReadAll(postResp.Body)

	log.Printf("Response: %s", respBody)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBody, resp); err != nil {
		return err
	}

	log.Printf("Unmarshalled: %+v", resp)

	return nil

}



type TransactionRequest struct {
	ClientID string`json:"client_id"`
	Secret string`json:"secret"`
	AccessToken string`json:"access_token"`
	StartDate string`json:"start_date"`
	EndDate string `json:"end_date"`
}

type Balance struct {
	Available 	float64`json:"available"`
	Current	float64`json:"current"`
	Limit float64`json:"limit"`
}

type Account struct {
	ID string`json:"account_id"`
	ItemID string`json:"item_id"`
	InstitutionID string`json:"institution_id"`
	Balances Balance`json:"balances"`
	Name string `json:"name"`
	Mask string `json:"mask"`
	OfficialName string`json:"official_name"`
	Type string`json:"type"`
	Subtype string`json:"subtype"`

}

type TransactionResponse struct {
	Accounts []Account`json:"accounts"`
	Transactions []Transaction`json:"transactions"`
	Item Item`json:"item"`
	RequestID string`json:"request_id"`
	TotalTransactions int32`json:"total_transactions"`
}

type ApiError struct {
	Type string`json:"error_type"`
	Code string`json:"error_code"`
	Message string`json:"error_message"`
	DisplayMessage string`json:"display_message"`
}

type Item struct {
	AvailableProducts []string`json:"available_products"`
	BilledProducts []string`json:"billed_products"`
	Error ApiError`json:"error"`
	InstitutionID string`json:"institution_id"`
	ItemID string`json:"item_id"`
}


type Transaction struct {
	ID string`json:"transaction_id"`
	AccountID string`json:"account_id"`
	Category []string`json:"category"`
	CategoryID string`json:"category_id"`
	Type string`json:"transaction_type"`
	Amount float64`json:"amount"`
	Date string`json:"date"`
	Pending bool`json:"pending"`
	AccountOwner string`json:"account_owner"`
	Name string `json:"string"`
}


func (c *Client) Transactions(accessToken, startDate, endDate string) (TransactionResponse, error) {
	endpoint := "/transactions/get"

	request := TransactionRequest{
		ClientID: c.clientID,
		Secret: c.secret,
		AccessToken: accessToken,
		StartDate: startDate,
		EndDate: endDate,
	}

	resp := TransactionResponse{}
	err := c.post(endpoint, request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}