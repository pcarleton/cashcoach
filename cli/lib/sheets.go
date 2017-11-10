package lib

import (
  "fmt"
	"io/ioutil"
	"net/http"
	"log"
  "bufio"
  "os"
  "strings"

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
  Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Service struct {
	Sheets *sheets.Service
	Drive *drive.Service
}

func (srv *Service) ListSpreadsheets(query string) ([]*drive.File, error) {
	r, err := srv.Drive.Files.List().PageSize(10).
			Q(query).
			Fields("nextPageToken, files(id, name, mimeType)").Do()

	if err != nil {
    return nil, err
	}

  return r.Files, nil
}

func (srv *Service) ImportSpreadsheet(ssName string, data [][]interface{}) (*sheets.Spreadsheet, error) {
  ss := &sheets.Spreadsheet{
    Properties: &sheets.SpreadsheetProperties{Title: ssName},
  }
  ss2, err := srv.Sheets.Spreadsheets.Create(ss).Do()
  if err != nil {
    return nil, err
  }


  bottomLeft := CellPos{len(data), len(data[0])}

  aRange := ARange(CellPos{}, bottomLeft)

  vRange := &sheets.ValueRange{
    Range: aRange,
    Values: data,
  }

  req := srv.Sheets.Spreadsheets.Values.Update(ss2.SpreadsheetId, aRange, vRange)

  req.ValueInputOption("USER_ENTERED")

  _, err = req.Do()

  if err != nil {
    return nil, err
  }

  // TODO: Return spreadsheet or something
  return ss2,  err
}

func (srv *Service) ShareFile(fileID, email string) error {
  perm := drive.Permission{
    EmailAddress: email,
    Role: "writer",
    Type: "user",
  }

  req := srv.Drive.Permissions.Create(fileID, &perm).SendNotificationEmail(false)


  _, err := req.Do()
  return err
}

func (srv *Service) Delete(fileID string) error {
  req := srv.Drive.Files.Delete(fileID)
  err := req.Do()
  return err
}

func aRangeLetter(idx int) string {
  secondLet := idx % len(Alphabet)

  if idx > len(Alphabet) {
    firstLet := idx / len(Alphabet)
    return fmt.Sprintf("%s%s",
    Alphabet[firstLet:firstLet + 1], Alphabet[secondLet:secondLet + 1])
  }

  return fmt.Sprintf("%s", Alphabet[secondLet:secondLet+1])
}

type CellPos struct {
  Row int
  Col int
}

func (c CellPos) ANotation() string {
  return fmt.Sprintf("%s%d", aRangeLetter(c.Col), c.Row + 1)
}

func ARange(start, end CellPos) string {
  return fmt.Sprintf("%s:%s", start.ANotation(), end.ANotation())
}

func TsvToArr(fname string) ([][]interface{}, error) {
    reader, err := os.Open(fname)
    if err != nil {
      return nil, err
    }
    delimiter := "\t"

    scanner := bufio.NewScanner(reader)

    data := make([][]interface{}, 0)

    for scanner.Scan() {
      pieces := strings.Split(scanner.Text(), delimiter)
      data = append(data, strToInterface(pieces))
    }

    return data, nil
}

func strToInterface(strs []string) []interface{} {
      arr := make([]interface{}, len(strs))

      for i, s := range(strs) {
        arr[i] = s
      }
      return arr
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

