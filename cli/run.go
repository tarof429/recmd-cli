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
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// backgroundFlag is a flag that determines whether to run a command in the background
	backgroundFlag bool
)

// The runCmd represents the run command. It takes one parameter, the command hash.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command",
	Long:  `Run a command`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Error: the hash must be specified")
			os.Exit(1)
		}

		InitTool()

		commandHash := strings.Trim(args[0], "")

		// backgroundFlag is a flag that determines whether to run a command in the background
		ret, err := RunCmd(commandHash, backgroundFlag)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
			return
		}

		// backgroundFlag is a flag that determines whether to run a command in the background
		if backgroundFlag == false {
			fmt.Println(ret.Coutput)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVarP(&backgroundFlag, "b", "b", false, "Run command in the background")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
