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
  "os"
  "log"
  "encoding/json"
  "time"

	"github.com/spf13/cobra"

	"github.com/pcarleton/cashcoach/cash/lib"
	"github.com/pcarleton/cashcoach/cash/lib/ledger"
	"github.com/pcarleton/cashcoach/api/plaid"
)

var ledgerImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a TSV into ledger format",
	Run: func(cmd *cobra.Command, args []string) {
    // TODO: read from stdin if no file specified

    filename := lib.StringFlagOrDie(cmd, "file")

    reader, err := os.Open(filename)
    if err != nil {
      log.Fatalf("Unable to open file: %v", err)
    }

    var trans []plaid.Transaction

    decoder := json.NewDecoder(reader)

    err = decoder.Decode(&trans)
    if err != nil {
      log.Fatalf("Unable to read json", err)
    }

    for _, t := range trans {
      lTrans := makeLTrans(&t, "taylor")
      fmt.Println(lTrans.String())
      fmt.Println("")
    }
	},
}

func splitTrans(t *plaid.Transaction, acct1, acct2 string) ledger.Transaction {
  changes := []ledger.Change{
    {ledger.Expense(t.Category...), t.Amount},
    {ledger.Liability(t.AccountID, acct1), -1*t.Amount/2},
    {Account: ledger.Liability(t.AccountID, acct2)},
  }

  date, err := time.Parse(plaid.DateFmt, t.Date)
  if err != nil {
    log.Fatalf("Unable to parse date: %s", err)
  }

  return ledger.Transaction{
    Date: date,
    Description: t.Name,
    Changes: changes,
  }
}


func makeLTrans(t *plaid.Transaction, acct string) ledger.Transaction {
  changes := []ledger.Change{
    {ledger.Expense(t.Category...), t.Amount},
    {Account: ledger.Liability(t.AccountID, acct)},
  }

  date, err := time.Parse(plaid.DateFmt, t.Date)
  if err != nil {
    log.Fatalf("Unable to parse date: %s", err)
  }

  return ledger.Transaction{
    Date: date,
    Description: t.Name,
    Changes: changes,
  }
}

// ledgerCmd represents the ledger command
var ledgerCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Interact with the ledger format",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(ledgerCmd)

  ledgerCmd.AddCommand(ledgerImportCmd)
  ledgerImportCmd.Flags().StringP("file", "f", "", "File to read transaction data from")
}
