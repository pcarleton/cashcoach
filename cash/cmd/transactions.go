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
  "time"
	"os"
	"strings"

	"github.com/pcarleton/cashcoach/cash/lib"
	"github.com/pcarleton/cashcoach/api/plaid"
	"github.com/spf13/cobra"
)



func pickInterval(cmd *cobra.Command) lib.Interval {
  lastN := lib.IntFlagOrDie(cmd, "lastN")
	if lastN != 0 {
		return lib.LastNDays(lastN)
	}

  var startT time.Time
  var endT time.Time

  start := lib.StringFlagOrDie(cmd, "start")
  end := lib.StringFlagOrDie(cmd, "end")

  if start != "" {
    startT = lib.DateOrDie(start)
  } else {
    startT = lib.NDaysAgo(30)
  }

  if end != "" {
    if start == "" {
      log.Fatalf("Must specify --start if you specify an --end")
    }
    endT = lib.DateOrDie(end)
  } else {
    endT = time.Now()
  }

	return lib.Interval{
		Start: startT,
		End:   endT,
	}
}

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Fetch transactions for a particular account",
	Args:  cobra.MinimumNArgs(1),
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

		client := lib.GetClient()
		interval := pickInterval(cmd)

		resp, err := client.Transactions(acct.Token, interval.Start, interval.End)

		if err != nil {
			log.Fatal(err)
		}

    delimiter, err := cmd.Flags().GetString("delimiter")
		if err != nil {
			log.Fatal(err)
		}


    masks := lib.Masks(resp.Accounts)
    log.Printf("Accounts ending in: %s", strings.Join(masks, ", "))
		log.Printf("%s to %s", interval.Start, interval.End)

    nickMap := acct.NickMap(resp.Accounts)

    jsonOut, err := cmd.Flags().GetBool("json")
		if err != nil {
			log.Fatal(err)
		}

    if jsonOut {
      // Might regret messing with the data like this later... 
      newTrans := make([]plaid.Transaction, len(resp.Transactions))
      for i, t := range resp.Transactions {
				t.AccountID = nickMap[t.AccountID]
        newTrans[i] = t
      }

      lib.OutputJson(newTrans)
      return
    }


		headers := []string{
			"account",
			"date",
			"description",
			"category",
			"amount",
		}


		fmt.Println(strings.Join(headers, "\t"))
		for _, trans := range resp.Transactions {
			pieces := []string{
				nickMap[trans.AccountID],
				trans.Date,
				trans.Name,
				strings.Join(trans.Category, ":"),
				fmt.Sprintf("%.2f", trans.Amount),
			}

			fmt.Println(strings.Join(pieces, delimiter))
		}
	},
}


func init() {
	RootCmd.AddCommand(transactionsCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	transactionsCmd.PersistentFlags().StringP("start", "s", "", "Start date to find transactions (like 2006-01-03)")
	transactionsCmd.PersistentFlags().StringP("end", "e", "", "End date to find transactions (inclusive)")
	transactionsCmd.PersistentFlags().IntP("lastN", "l", 0, "Fecth transactions for the last N days")
	transactionsCmd.Flags().StringP("delimiter", "d", "\t", "Delimiter to use for printing")
	transactionsCmd.Flags().BoolP("json", "j", false, "When true, output transaction data as JSON")
}
