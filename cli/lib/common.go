package lib

import (
	"github.com/spf13/viper"
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
