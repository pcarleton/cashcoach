package lib

import (
	"github.com/spf13/viper"
	"github.com/pcarleton/cashcoach/api/plaid"
)

type Account struct{
  Name string
  Token string
}

func GetAccounts() ([]Account, error) {
    accounts := make([]Account, 0, 10)
    err := viper.UnmarshalKey("accounts", &accounts)
    if err != nil {
      return nil, err
    }
    return accounts, nil
}

func GetAccount(name string) (*Account, error) {
  accts, err := GetAccounts()
  if err != nil {
    return nil, err
  }
  for _, acct := range(accts) {
    if acct.Name == name {
      return &acct, nil
    }
  }
  return nil, nil
}

func GetClient() plaid.Client {
  // TODO: Memoize?
	return plaid.NewClient(
		viper.GetString("client_id"),
		viper.GetString("client_secret"),
		plaid.DevURL)
}
