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


	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/pcarleton/cashcoach/cli/lib"
)

// sheetsCmd represents the sheets command
var sheetsCmd = &cobra.Command{
	Use:   "sheets",
	Short: "Interact with Google Sheets",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("sheets called")
	},
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import TSV to Google sheets",
	Run: func(cmd *cobra.Command, args []string) {
    srv := lib.GetService()
    fname := args[0]

    reader, err := os.Open(fname)
    if err != nil {
      log.Fatalf("Unable to open file: %v", err)
    }

		err = srv.ImportSpreadsheet(fname, reader)
    if err != nil {
      log.Fatalf("Unable to import file: %v", err)
    }

	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete file from Drive",
	Run: func(cmd *cobra.Command, args []string) {
    srv := lib.GetService()
    fileId := args[0]

    err := srv.Delete(fileId)
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
    srv := lib.GetService()
    fileID := args[0]

    email := viper.GetString("email")

    err := srv.ShareFile(fileID, email)
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
    srv := lib.GetService()
		files, err := srv.ListSpreadsheets(flagQuery)
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

	sheetsCmd.AddCommand(importCmd)
	sheetsCmd.AddCommand(deleteCmd)
}
