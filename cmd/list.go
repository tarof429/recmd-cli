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
	"time"

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

		// Mon Jan 2 15:04:05 MST 2006
		layout := "2006-01-02"

		for _, c := range readCmds {
			creationTime, _ := time.Parse(layout, c.Creationtime.String())
			fmt.Printf("%s\t%s\t%s\t%v\t%v\n", c.CmdHash, c.CmdString, c.Comment, creationTime, c.Modificationtime)
		}
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
