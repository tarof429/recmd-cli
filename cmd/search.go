/*
Copyright © 2020 Taro Fukunaga <tarof429@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	recmd "github.com/tarof429/recmd"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a command",
	Long:  `Search for a command`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Must specify the command")
			os.Exit(1)
		}

		value := args[0]

		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to obtain home directory path %v\n", err)
		}

		ret, err := recmd.SelectCmd(homeDir, value)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to find command %v\n", err)
		}

		data, _ := json.MarshalIndent(ret, "", "\t")
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
