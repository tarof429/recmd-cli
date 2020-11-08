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

	"github.com/spf13/cobra"
)

// queueCmd lists commands that are queued for execution
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Get the queue of commands",
	Long:  `"Get the queue of commands."`,
	Run: func(cmd *cobra.Command, args []string) {

		InitTool()

		ret, err := Queue()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
			return
		}

		DisplayQueue(ret)
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
