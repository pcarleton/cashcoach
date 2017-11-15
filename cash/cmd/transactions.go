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
	"github.com/spf13/cobra"
)

func pickInterval() lib.Interval {
	if flagLastN != 0 {
		return lib.LastNDays(flagLastN)
	}

	start := flagStart
	end := flagEnd
	if end == "" && start != "" {
		// Pick today as end
		end = lib.TodayStr()
	}

	if end == "" && start == "" {
		return lib.LastNDays(30)
	}

	return lib.Interval{
		Start: start,
		End:   end,
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
		interval := pickInterval()

		resp, err := client.Transactions(acct.Token, interval.Start, interval.End)

		if err != nil {
			log.Fatal(err)
		}

    delimiter, err := cmd.Flags().GetString("delimiter")
		if err != nil {
			log.Fatal(err)
		}


		log.Printf("Showing transactions for %s to %s", interval.Start, interval.End)

		headers := []string{
			"account",
			"date",
			"description",
			"category",
			"amount",
		}

    nickMap := acct.NickMap(resp.Accounts)

    log.Printf("Nicks: %+v", nickMap)
    log.Printf("Names: %+v", acct.Nicknames)

		fmt.Println(strings.Join(headers, "\t"))
		// TODO: Allow JSON output via flags
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

var flagStart string
var flagEnd string
var flagLastN int

func init() {
	RootCmd.AddCommand(transactionsCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	transactionsCmd.PersistentFlags().StringVarP(&flagStart, "start", "s", "", "Start date to find transactions")
	transactionsCmd.PersistentFlags().StringVarP(&flagEnd, "end", "e", "", "End date to find transactions (inclusive)")
	transactionsCmd.PersistentFlags().IntVarP(&flagLastN, "lastN", "l", 0, "Fecth transactions for the last N days")
	transactionsCmd.Flags().StringP("delimiter", "d", "\t", "Delimiter to use for printing")
}
