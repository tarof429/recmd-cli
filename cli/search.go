package cli

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

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command. It takes one parameter, the command description.
// Multiple results may be returned; the output format is JSON (for the moment). Any command with
// a matching string is returned; this makes it possible, for example, to show all commmands
// that have 'docker' as part of its description.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a command by its comment",
	Long:  `Search for a command by its comment`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error: any part of the comment must be specified")
			os.Exit(1)
		}

		InitTool()

		description := strings.Trim(args[0], "")

		ret, err := SearchCmd(description)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
			return
		}

		if len(ret) > 0 {
			data, _ := json.MarshalIndent(ret, "", "\t")
			fmt.Println(string(data))
		}
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
