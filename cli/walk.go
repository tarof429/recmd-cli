package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/manifoldco/promptui"
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

		prompt := promptui.Prompt{
			Label:     "Do you want to continue?",
			IsConfirm: true,
		}

		_, err := prompt.Run()

		if err != nil {
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

			// Ask the user whether to run the command.
			// If the user types 'y', then run the comand. If the user types 'N', then
			// move on to the next command. If the user inputs anything else (such as Ctrl-C) then abort.
			prompt := promptui.Prompt{
				Label:     "Do you want to run this command",
				IsConfirm: true,
			}

			ret, err = prompt.Run()

			if ret != "y" {
				continue
			}

			sc, err := RunCmd(cmd.CmdHash, false)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
				return
			}

			fmt.Println(sc.Coutput)
			time.Sleep(time.Second)

			// // Similar to the above prompt. If the user types Ctrl-C, then abort. Otherwise, continue to the next command.
			// prompt = promptui.Prompt{
			// 	Label:     "Type any key to continue",
			// 	IsConfirm: false,
			// }

			// ret, err = prompt.Run()

			// if ret == "" {
			// 	continue
			// } else if err != nil {
			// 	return
			// }

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
