package plaid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DevURL     = "https://development.plaid.com"
	SandboxURL = "https://sandbox.plaid.com"
	DateFmt    = "2006-01-02"
)

type Client struct {
	baseURL  string
	client   *http.Client
	clientID string
	secret   string
}

func NewClient(clientID, secret, baseURL string) Client {
	return Client{
		baseURL:  baseURL,
		client:   new(http.Client),
		clientID: clientID,
		secret:   secret,
	}
}

type UpdateAccessTokenRequest struct {
	ClientID    string `json:"client_id"`
	Secret      string `json:"secret"`
	AccessToken string `json:"access_token_v1"`
}

type UpdateAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	RequestID   string `json:"request_id"`
}

func (c *Client) UpdateAccessToken(accessToken string) (UpdateAccessTokenResponse, error) {
	endpoint := "/item/access_token/update_version"

	req := UpdateAccessTokenRequest{c.clientID, c.secret, accessToken}

	resp := UpdateAccessTokenResponse{}
	err := c.post(endpoint, req, &resp)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

type ApiError struct {
	Response *ErrorResponse
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%+v", e.Response)
}

func (c *Client) post(endpoint string, req interface{}, resp interface{}) error {

	jsonText, err := json.Marshal(req)
	if err != nil {
		return err
	}

	target := c.baseURL + endpoint

	postResp, err := c.client.Post(target, "application/json",
		bytes.NewReader(jsonText))

	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(postResp.Body)

	if err != nil {
		return err
	}

	if postResp.StatusCode != 200 {
		errResp := ErrorResponse{}
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return err
		}

		return ApiError{&errResp}
	}

	if err := json.Unmarshal(respBody, resp); err != nil {
		return err
	}

	return nil
}

type TransactionRequest struct {
	ClientID    string `json:"client_id"`
	Secret      string `json:"secret"`
	AccessToken string `json:"access_token"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type Balance struct {
	Available float64 `json:"available"`
	Current   float64 `json:"current"`
	Limit     float64 `json:"limit"`
}

type Account struct {
	ID            string  `json:"account_id"`
	ItemID        string  `json:"item_id"`
	InstitutionID string  `json:"institution_id"`
	Balances      Balance `json:"balances"`
	Name          string  `json:"name"`
	Mask          string  `json:"mask"`
	OfficialName  string  `json:"official_name"`
	Type          string  `json:"type"`
	Subtype       string  `json:"subtype"`
}

type TransactionResponse struct {
	Accounts          []Account     `json:"accounts"`
	Transactions      []Transaction `json:"transactions"`
	Item              Item          `json:"item"`
	RequestID         string        `json:"request_id"`
	TotalTransactions int32         `json:"total_transactions"`
}

type ErrorResponse struct {
	Type           string `json:"error_type"`
	Code           string `json:"error_code"`
	Message        string `json:"error_message"`
	DisplayMessage string `json:"display_message"`
}

type Item struct {
	AvailableProducts []string      `json:"available_products"`
	BilledProducts    []string      `json:"billed_products"`
	Error             ErrorResponse `json:"error"`
	InstitutionID     string        `json:"institution_id"`
	ItemID            string        `json:"item_id"`
}

type Transaction struct {
	ID           string   `json:"transaction_id"`
	AccountID    string   `json:"account_id"`
	Category     []string `json:"category"`
	CategoryID   string   `json:"category_id"`
	Type         string   `json:"transaction_type"`
	Amount       float64  `json:"amount"`
	Date         string   `json:"date"`
	Pending      bool     `json:"pending"`
	AccountOwner string   `json:"account_owner"`
	Name         string   `json:"name"`
}

type PublicTokenRequest struct {
	ClientID    string `json:"client_id"`
	Secret      string `json:"secret"`
	AccessToken string `json:"access_token"`
}

type PublicTokenResponse struct {
	RequestID         string        `json:"request_id"`
	PublicToken string        `json:"public_token"`
}

func (c *Client) Transactions(accessToken string, startDate, endDate time.Time) (TransactionResponse, error) {
	endpoint := "/transactions/get"

	request := TransactionRequest{
		ClientID:    c.clientID,
		Secret:      c.secret,
		AccessToken: accessToken,
		StartDate:   startDate.Format(DateFmt),
		EndDate:     endDate.Format(DateFmt),
	}

	resp := TransactionResponse{}
	err := c.post(endpoint, request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}


func (c *Client) CreatePublicToken(accessToken string) (PublicTokenResponse, error) {
	endpoint := "/item/public_token/create"

	request := PublicTokenRequest{
		ClientID:    c.clientID,
		Secret:      c.secret,
		AccessToken: accessToken,
	}

	resp := PublicTokenResponse{}
	err := c.post(endpoint, request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type ExchangeRequest struct {
	ClientID    string `json:"client_id"`
	Secret      string `json:"secret"`
	PublicToken string `json:"public_token"`
}

type ExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ItemID      string `json:"item_id"`
	ErrorResponse
}

func (c *Client) Exchange(publicToken string) (ExchangeResponse, error) {
	endpoint := "/item/public_token/exchange"

	request := ExchangeRequest{
		ClientID:    c.clientID,
		Secret:      c.secret,
		PublicToken: publicToken,
	}

	resp := ExchangeResponse{}
	err := c.post(endpoint, request, &resp)
	if err != nil {
		return resp, err
	}

	if resp.Message != "" {
		return resp, fmt.Errorf(resp.Message)
	}

	return resp, nil
}
