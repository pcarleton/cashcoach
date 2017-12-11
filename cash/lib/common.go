package lib

import (
  "fmt"
  "log"
  "time"
  "encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/pcarleton/cashcoach/api/plaid"
)

const (
  DateFmt = plaid.DateFmt
)

type Account struct{
  Name string
  Token string
  Nicknames map[string]string
}


func DateOrDie(date string) time.Time {
  result, err := time.Parse(DateFmt, date)
  if err != nil {
    log.Fatalf("Unable to parse date %s: %v", date, err)
  }
  return result
}

func StringFlagOrDie(cmd *cobra.Command, flag string) string {
  result, err := cmd.Flags().GetString(flag)
  if err != nil {
    log.Fatalf("Unable to parse flag %s: %v", flag, err)
  }
  return result
}

func IntFlagOrDie(cmd *cobra.Command, flag string) int {
  result, err := cmd.Flags().GetInt(flag)
  if err != nil {
    log.Fatalf("Unable to parse flag %s: %v", flag, err)
  }
  return result
}



func (a *Account) NickMap(accts []plaid.Account) map[string]string {
  nickMap := make(map[string]string)

  for _, acct := range(accts) {
    nickMap[acct.ID] = a.Nicknames[acct.Mask]
  }

  return nickMap
}

func Masks(accts []plaid.Account) []string {
  masks := make([]string, len(accts))

  for i, acct := range(accts) {
    masks[i] = acct.Mask
  }

  return masks
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
  Start time.Time
  End time.Time
}

func NDaysAgo(n int) time.Time {
  today := time.Now()
  return today.Add(time.Duration(n * -24) * time.Hour )
}

func LastNDays(n int) Interval {
  return Interval{
    Start: NDaysAgo(n),
    End: time.Now(),
  }
}

func OutputJson(val interface{}) error {
  valj, err := json.Marshal(val)
  if err != nil {
    return err
  }
  fmt.Println(string(valj))
  return nil
}

