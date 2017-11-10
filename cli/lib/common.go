package lib

import (
  "time"

	"github.com/spf13/viper"
	"github.com/pcarleton/cashcoach/api/plaid"
)

const (
  DateFmt = "2006-01-02"
)

type Account struct{
  Name string
  Token string
  Nicknames map[string]string
}

func (a *Account) NickMap(accts []plaid.Account) map[string]string {
  nickMap := make(map[string]string)

  for _, acct := range(accts) {
    nickMap[acct.ID] = a.Nicknames[acct.Mask]
  }

  return nickMap
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

type Interval struct {
  Start string
  End string
}


func LastNDays(n int) Interval {
  today := time.Now()
  nDaysAgo := today.Add(time.Duration(n * -24) * time.Hour )

  return Interval{
    Start: nDaysAgo.Format(DateFmt),
    End: today.Format(DateFmt),
  }
}

func TodayStr() string {
  return time.Now().Format(DateFmt)
}



