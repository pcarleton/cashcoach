package lib

import (
  "os"
  "log"

	"github.com/spf13/viper"
  "github.com/pcarleton/sheets"
)

const (
  LinkTmpl = "https://docs.google.com/spreadsheets/d/%s"
  SheetLinkTmpl = "https://docs.google.com/spreadsheets/d/%s/edit#gid=%s"
)

func GetSheetsClient() *sheets.Client {
	robotCredsPath := viper.GetString("robot_creds")

  r, err := os.Open(robotCredsPath)
  if err != nil {
    log.Fatalf("Unable to read credentials: %s", err)
  }

  client, err := sheets.NewServiceAccountClient(r)

  if err != nil {
    log.Fatalf("Unable to create sheets client: %s", err)
  }

  return client
}

