// Copyright © 2017 Paul Carleton
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pcarleton/cashcoach/cash/lib"
	"github.com/pcarleton/sheets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sheetsCmd represents the sheets command
var sheetsCmd = &cobra.Command{
	Use:   "sheets",
	Short: "Interact with Google Sheets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create Google sheets from TSV",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := lib.GetSheetsClient()
		fname := lib.StringFlagOrDie(cmd, "file")
		title := lib.StringFlagOrDie(cmd, "title")

		reader, err := os.Open(fname)
		if err != nil {
			log.Fatalf("Unable to open file: %v", err)
		}

		r, err := client.CreateSpreadsheetFromTsv(title, reader)
		if err != nil {
			log.Fatalf("Unable to import file: %v", err)
		}
		log.Print("Created spreadsheet.")

		log.Print("Sharing...")
		email := viper.GetString("email")

		err = r.Share(email)
		if err != nil {
			log.Fatalf("Unable to share file: %v", err)
		}

		log.Printf("Complete! View at: %s\n", r.Url())
	},
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import a TSV to an existing Google sheet",
	Run: func(cmd *cobra.Command, args []string) {
		fname := lib.StringFlagOrDie(cmd, "file")
		ssId := lib.StringFlagOrDie(cmd, "spreadsheet")
		sheetName := lib.StringFlagOrDie(cmd, "name")

		client := lib.GetSheetsClient()

		reader, err := os.Open(fname)
		if err != nil {
			log.Fatalf("Unable to open file: %v", err)
		}
		data := sheets.TsvToArr(reader)

		ss, err := client.GetSpreadsheet(ssId)
		if err != nil {
			log.Fatalf("Unable to find spreadsheet: %v", err)
		}

		sheet := ss.GetSheet(sheetName)

		if sheet == nil {
			sheet, err = ss.AddSheet(sheetName)

			if err != nil {
				log.Fatalf("Unable to add sheet: %v", err)
			}
		}

		err = sheet.Update(data)
		if err != nil {
			log.Fatalf("Unable to add data to sheet: %v", err)
		}

		log.Printf("Complete! View at: %s\n", ss.Url())
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete file from Drive",
	Run: func(cmd *cobra.Command, args []string) {
		client := lib.GetSheetsClient()
		fileId := args[0]

		err := client.Delete(fileId)
		if err != nil {
			log.Fatalf("Unable to delete file: %v", err)
		}

		log.Printf("Deleted %s.\n", fileId)
	},
}

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "share the specified file with the user specified in config",
	Run: func(cmd *cobra.Command, args []string) {
		client := lib.GetSheetsClient()
		fileID := args[0]

		email := viper.GetString("email")

		err := client.ShareFile(fileID, email)
		if err != nil {
			log.Fatalf("Unable to share file: %v", err)
		}

		log.Printf("Successfully shared %s with %s\n", fileID, email)
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list Google sheets",
	Run: func(cmd *cobra.Command, args []string) {
		client := lib.GetSheetsClient()
		files, err := client.ListFiles(flagQuery)
		if err != nil {
			log.Fatalf("Unable to retrieve data from drive. %v", err)
		}

		ids, err := cmd.Flags().GetBool("ids")
		mime, err := cmd.Flags().GetBool("mime")

		if err != nil {
			log.Fatalf("Unable to parse flags. %v", err)
		}

		headers := []string{"Name"}
		if ids {
			headers = append(headers, "ID")
		} else {
			headers = append(headers, "Link")
		}

		if mime {
			headers = append(headers, "Mime Type")
		}

		// TODO: Better table printing
		fmt.Println(strings.Join(headers, "\t"))

		if len(files) > 0 {
			for _, i := range files {
				pieces := []string{i.Name}
				if ids {
					pieces = append(pieces, i.Id)
				} else {
					link := fmt.Sprintf(lib.LinkTmpl, i.Id)
					pieces = append(pieces, link)
				}

				if mime {
					pieces = append(pieces, i.MimeType)
				}
				fmt.Println(strings.Join(pieces, "\t"))
			}
		} else {
			fmt.Println("No files found.")
		}
	},
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull a sheet from Google Sheets",
	Run: func(cmd *cobra.Command, args []string) {
		ssId := lib.StringFlagOrDie(cmd, "spreadsheet")
		sheetName := lib.StringFlagOrDie(cmd, "name")

		client := lib.GetSheetsClient()

		ss, err := client.GetSpreadsheetWithData(ssId)
		if err != nil {
			log.Fatalf("Unable to find spreadsheet: %v", err)
		}

		sheet := ss.GetSheet(sheetName)
		if sheet == nil {
			log.Fatalf("Unable to find sheet: %s", sheetName)
		}

		matrix, err := sheet.GetContents()
		if err != nil {
			log.Fatalf("Unable to fetch sheet contents: %v", err)
		}

		for _, row := range matrix {
			fmt.Println(strings.Join(row, "\t"))
		}
	},
}

var flagQuery string

func init() {
	RootCmd.AddCommand(sheetsCmd)

	sheetsCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&flagQuery, "query", "q", "", "Query to pass to Files.list")
	listCmd.Flags().BoolP("ids", "i", false, "Only list ID's")
	listCmd.Flags().BoolP("mime", "m", false, "Print the file mime types")

	sheetsCmd.AddCommand(shareCmd)

	shareCmd.Flags().String("email", "", "Email to share file with")
	viper.BindPFlag("email", shareCmd.Flags().Lookup("email"))

	sheetsCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "", "The file to read data from, if not set use STDIN")
	createCmd.Flags().StringP("title", "t", "", "The title to give the spreadsheet")

	sheetsCmd.AddCommand(deleteCmd)
	sheetsCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("file", "f", "", "The file to read data from, if not set use STDIN")
	importCmd.Flags().StringP("spreadsheet", "s", "", "The ID of the spreadsheet to import to")
	importCmd.Flags().StringP("name", "n", "", "The name of the sheet to import to")

	sheetsCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringP("spreadsheet", "s", "", "The ID of the spreadsheet to pull to")
	pullCmd.Flags().StringP("name", "n", "", "The name of the sheet to pull to")
}
