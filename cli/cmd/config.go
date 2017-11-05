// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/viper"
)

type Account struct{
  Name string
  Token string
}


// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interacts with the saved config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")

    fmt.Printf("Client ID: %s\n", viper.GetString("client_id"))
    // TODO: hide secrets on printing by default
    fmt.Printf("Client secret: %s\n", viper.GetString("client_secret"))

    accounts := make([]Account, 0, 10)
    err := viper.UnmarshalKey("accounts", &accounts)
    if err != nil {
      panic(err)
    }
    fmt.Println("Accounts")
    fmt.Println("Name\tToken")

    for _, acct := range(accounts) {
      fmt.Printf("%s\t%s\n", acct.Name, acct.Token)
    }





	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
