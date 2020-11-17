package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk through the commands",
	Long:  `"Walk through the commands.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("This operation will walk through all the commands.")
		fmt.Println("The command will only be run if you answer 'y'.")

		fmt.Print("Do you want to continue? (y/N) ")

		var resp string
		fmt.Scanln(&resp)

		if strings.Trim(resp, "") != "y" {
			return
		}

		InitTool()

		ret, err := List()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
			return
		}

		total := len(ret)

		for index, cmd := range ret {

			fmt.Println("*****************************")
			fmt.Printf("Walking through command %x/%x\n", index+1, total)
			fmt.Println("*****************************")
			fmt.Println()

			ret, err := ShowCmd(cmd.CmdHash)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
				return
			}
			fmt.Println(ret)
			fmt.Println()

			fmt.Print("Do you want to run this command? (y/N) ")

			// Reset the response
			resp = ""
			fmt.Scanln(&resp)

			if strings.Trim(resp, "") != "y" {
				continue
			}
			sc, err := RunCmd(cmd.CmdHash, false)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
				return
			}

			fmt.Println(sc.Coutput)

			// Wait for the user to hit 'enter' before showing the next command
			fmt.Print("Type any key to continue ")
			resp = ""
			fmt.Scanln(&resp)
		}
	},
}

func init() {
	rootCmd.AddCommand(walkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
