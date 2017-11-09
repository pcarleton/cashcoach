package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/drive/v3"
)

const (
	sheetsScope = "https://www.googleapis.com/auth/spreadsheets"
	sheetMimeType = "application/vnd.google-apps.spreadsheet"
)

type Service struct {
	Sheets *sheets.Service
	Drive *drive.Service
}


func getClient() *http.Client {
	robotCredsPath := viper.GetString("robot_creds")
	b, err := ioutil.ReadFile(robotCredsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.JWTConfigFromJSON(b, sheetsScope, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse robot creds to config: %v", err)
	}
	ctx := context.Background()
	client := config.Client(ctx)
	
	return client

}

func GetService() Service {
	client:= getClient()
	sheetsSrv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	driveSrv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	return Service{sheetsSrv, driveSrv}
}

func ListSpreadsheets() {
	srv := GetService()
	r, err := srv.Drive.Files.List().PageSize(10).
			Fields("nextPageToken, files(id, name, mimeType)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from drive. %v", err)
	}

  if err != nil {
    log.Fatalf("Unable to retrieve files: %v", err)
  }

  fmt.Println("Files:")
  if len(r.Files) > 0 {
    for _, i := range r.Files {
      fmt.Printf("%s (%s) %s\n", i.Name, i.Id, i.MimeType)
    }
  } else {
    fmt.Println("No files found.")
  }
}

