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
	//"encoding/json"
	"fmt"
  "strconv"
	"log"
	"os"
	"time"
  "bufio"
  "strings"

	"github.com/spf13/cobra"

	"github.com/pcarleton/cashcoach/api/plaid"
	"github.com/pcarleton/cashcoach/cash/lib"
	"github.com/pcarleton/cashcoach/cash/lib/ledger"
)

var ledgerImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a TSV into ledger format",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: read from stdin if no file specified

		filename := lib.StringFlagOrDie(cmd, "file")

    ttrans, err :=  readTsv(filename)

		if err != nil {
			log.Fatalf("Unable to load transactions: %v", err)
		}

    fmt.Printf("%+v", ttrans)
    
    

		//reader, err := os.Open(filename)



		//var trans []plaid.Transaction

		//decoder := json.NewDecoder(reader)

		//err = decoder.Decode(&trans)
		//if err != nil {
		//	log.Fatalf("Unable to read json", err)
		//}

		//for _, t := range trans {
		//	lTrans := makeLTrans(&t, "taylor")
		//	fmt.Println(lTrans.String())
		//	fmt.Println("")
		//}
	},
}

type TableTrans struct {
  Account string // Nick name, human readable
  Date time.Time
  Description string
  Category string
  Label string
  Amount float64
}

func readTsv(fileName string) ([]TableTrans, error) {
		reader, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}

    scanner := bufio.NewScanner(reader)
    scanner.Scan()
    headers := strings.Split(scanner.Text(), "\t")
    fmt.Println(headers)

    ttrans := make([]TableTrans, 1)

    for scanner.Scan() {
      line := scanner.Text()
      pieces := strings.Split(line, "\t")
      trans := make(map[string]string)

      for idx, header := range headers {
        trans[header] = pieces[idx]
      }

      date, err := time.Parse("2006-01-02", trans["date"])
      if err != nil {
        return nil, fmt.Errorf("Invalid date %s in line {%s}", trans["date"], line)
      }

      amount, err := strconv.ParseFloat(trans["amount"], 64)
      if err != nil {
        return nil, fmt.Errorf("Invalid amount %s in line {%s}", trans["amount"], line)
      }
      ttrans = append(ttrans, TableTrans{
        Account: trans["account"],
        Date: date,
        Description: trans["description"],
        Category: trans["category"],
        Amount: amount,
      })
    }

    return ttrans, nil
}

func splitTrans(t *plaid.Transaction, acct1, acct2 string) ledger.Transaction {
	changes := []ledger.Change{
		{ledger.Expense(t.Category...), t.Amount},
		{ledger.Liability(t.AccountID, acct1), -1 * t.Amount / 2},
		{Account: ledger.Liability(t.AccountID, acct2)},
	}

	date, err := time.Parse(plaid.DateFmt, t.Date)
	if err != nil {
		log.Fatalf("Unable to parse date: %s", err)
	}

	return ledger.Transaction{
		Date:        date,
		Description: t.Name,
		Changes:     changes,
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
		Date:        date,
		Description: t.Name,
		Changes:     changes,
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
