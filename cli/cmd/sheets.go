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

	"github.com/spf13/cobra"
)

// sheetsCmd represents the sheets command
var sheetsCmd = &cobra.Command{
	Use:   "sheets",
	Short: "Interact with Google Sheets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sheets called")
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list Google sheets",
	Run: func(cmd *cobra.Command, args []string) {
    lib.ListSpreadsheets(flagQuery)
	},
}

var flagQuery string

func init() {
	RootCmd.AddCommand(sheetsCmd)

	sheetsCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&flagQuery, "query", "q", "", "Query to pass to Files.list")
}
