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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command",
	Long:  `Run a command`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Error: either the command or hash must be specified")
			os.Exit(1)
		}

		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to obtain home directory path %v\n", err)
			return
		}

		value := args[0]

		selectedCmd, cerr := recmd.SelectCmd(homeDir, value)

		if cerr != nil {
			fmt.Fprintf(os.Stderr, "Unable to read history file: %s\n", err)
			return
		}

		sc := recmd.ScheduleCommand(selectedCmd, recmd.RunShellScriptCommandWithSpinner)

		if sc.ExitStatus != 0 {
			fmt.Println(sc.Stderr)
		} else {
			fmt.Println(sc.Stdout)
		}

		ret := recmd.UpdateCommandDuration(homeDir, selectedCmd, sc.Duration)

		if ret != true {
			fmt.Fprintf(os.Stderr, "Error while updating command")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
