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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	recmd "github.com/tarof429/recmd-cli/recmd"
)

var command string
var message string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a command",
	Long:  `Add a command`,

	Run: func(cmd *cobra.Command, args []string) {
		if command == "" || message == "" {
			cmd.Usage()
			os.Exit(1)
		}

		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to obtain home directory path %v\n", err)
		}

		testCmd := recmd.NewCommand(command, message)

		ret := recmd.WriteCmdHistoryFile(homeDir, testCmd)

		if ret == false {
			fmt.Fprintf(os.Stderr, "An error occurred while writing the history file\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&command, "command", "c", "", "Command line")
	addCmd.Flags().StringVarP(&message, "description", "d", "", "Description")
	// addCmd.AddCommand(&comment, "message", "m", 1, "Message about the command")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
