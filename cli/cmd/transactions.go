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
  "github.com/pcarleton/cashcoach/cli/lib"
)

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Fetch transactions for a particular account",
  Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

    // TODO: Allow specifying multiple accounts?
    acct, err := lib.GetAccount(args[0])

    if err != nil {
      panic(err)
    }

    if acct == nil {
      fmt.Println("No account with that name.")
      os.Exit(1)
    }

    fmt.Printf("Found %s\n", acct.Name)

    client := lib.GetClient()

    // TODO: Specify date range via flags
    resp, err := client.Transactions(acct.Token, "2017-10-01", "2017-10-07")

    if err != nil {
      log.Fatal(err)
    }

    headers := []string{
      "account",
      "date",
      "description",
      "category",
      "amount",
    }

    fmt.Println(strings.Join(headers,"\t"))
    // TODO: Allow JSON output via flags
    for _, trans := range(resp.Transactions) {
      pieces := []string{
        trans.AccountID,
        trans.Date,
        trans.Name,
        strings.Join(trans.Category, ":"),
        fmt.Sprintf("%f", trans.Amount),
      }

      fmt.Println(strings.Join(pieces, "\t"))
    }
  },
}

func init() {
	RootCmd.AddCommand(transactionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transactionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transactionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
