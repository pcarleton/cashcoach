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

    ids, err := cmd.Flags().GetBool("ids")

		if err != nil {
			log.Fatalf("Unable to retrieve data from drive. %v", err)
		}
		fmt.Println("Files:")
		if len(files) > 0 {
			for _, i := range files {
        if ids {
				  fmt.Printf("%s\t%s\n", i.Name, i.Id)
        } else {
          link := fmt.Sprintf(lib.LinkTmpl, i.Id)
				  fmt.Printf("%s\t%s\n", i.Name, link)
        }
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

	sheetsCmd.AddCommand(importCmd)
	sheetsCmd.AddCommand(shareCmd)

  shareCmd.Flags().String("email", "", "Email to share file with")
  viper.BindPFlag("email", shareCmd.Flags().Lookup("email"))
}
