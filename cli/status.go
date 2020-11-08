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
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command. It takes one parameter, the command hash.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status",
	Long:  `Get the status`,
	Run: func(cmd *cobra.Command, args []string) {

		InitTool()

		_, err := StatusCmd()

		if err != nil {
			fmt.Println("Not running")
			return
		} else {
			fmt.Println("Server is up!")
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
