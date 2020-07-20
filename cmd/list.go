/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	recmd "github.com/tarof429/recmd"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Long:  `List commands.`,
	Run: func(cmd *cobra.Command, args []string) {

		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to obtain home directory path %v\n", err)
		}

		readCmds, err := recmd.ReadCmdHistoryFile(homeDir)

		// layout := "Mon 01 02 2006 15:04:05"

		output := fmt.Sprintf("%.15s\t\t%-40s\t%-50s\n", "COMMAND HASH", "COMMAND STRING", "COMMAND COMMENT")
		//fmt.Printf("%.15s\t\t%.30s\t%.30s\n", "COMMAND HASH", "COMMAND STRING", "COMMAND COMMENT")

		for _, c := range readCmds {
			// Maybe these are nice to have
			// creationTime := c.Creationtime.Format(layout)
			// modTime := c.Creationtime.Format(layout)
			cmdHash := c.CmdHash[0:15]

			var cmdString string
			var comment string

			if len(c.CmdString) > 40 {
				cmdString = c.CmdString[0:40] + "..."
			} else {
				cmdString = c.CmdString
			}

			if len(c.Comment) > 50 {
				comment = c.Comment[0:50] + "..."
			} else {
				comment = c.Comment
			}

			output = fmt.Sprintf(output+"%.15s\t\t%-40s\t%-50s\n", cmdHash, cmdString, comment)
			//fmt.Printf("%.15s\t\t%.30s\t\t%.30s\n", c.CmdHash, c.CmdString, c.Comment)
		}

		fmt.Print(output)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
