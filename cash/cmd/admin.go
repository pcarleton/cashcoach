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

	"github.com/spf13/cobra"

	"github.com/pcarleton/cashcoach/cash/lib"
)

var publicTokenCmd = &cobra.Command{
	Use:   "public_token",
	Short: "Get a plaid public token for an account",
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

		resp, err := client.CreatePublicToken(acct.Token)

		if err != nil {
			log.Fatal(err)
		}

  	fmt.Println(resp.PublicToken)
	},
}

// plaidCmd represents the plaid command
var plaidCmd = &cobra.Command{
	Use:   "plaid",
	Short: "Interact with the plaid API",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(plaidCmd)

  plaidCmd.AddCommand(publicTokenCmd)
}
