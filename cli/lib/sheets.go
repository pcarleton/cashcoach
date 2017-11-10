package lib

import (
	"io/ioutil"
	"net/http"
	"log"
  "io"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/drive/v3"
)

const (
	sheetsScope = "https://www.googleapis.com/auth/spreadsheets"
	sheetMimeType = "application/vnd.google-apps.spreadsheet"

  LinkTmpl = "https://docs.google.com/spreadsheets/d/%s"
  SheetLinkTmpl = "https://docs.google.com/spreadsheets/d/%s/edit#gid=%s"
)

type Service struct {
	Sheets *sheets.Service
	Drive *drive.Service
}




func (srv *Service) ListSpreadsheets(query string) ([]*drive.File, error) {
	r, err := srv.Drive.Files.List().PageSize(10).
			Q(query).
			Fields("nextPageToken, files(id, name)").Do()

	if err != nil {
    return nil, err
	}

  return r.Files, nil
}

func (srv *Service) ImportSpreadsheet(ssName string, data io.Reader) error {
  newFile := drive.File{Name: ssName, MimeType: sheetMimeType}

  _, err := srv.Drive.Files.Create(&newFile).Media(data).Do()

  // TODO: Return spreadsheet or something
  return err
}

func (srv *Service) ShareFile(fileID, email string) error {
  perm := drive.Permission{
    EmailAddress: email,
    Role: "owner",
    Type: "user",
  }

  req := srv.Drive.Permissions.Create(fileID, &perm).TransferOwnership(true)


  _, err := req.Do()
  return err
}

func getClient() *http.Client {
	robotCredsPath := viper.GetString("robot_creds")
	b, err := ioutil.ReadFile(robotCredsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

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

