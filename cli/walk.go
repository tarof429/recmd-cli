package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk through the commands",
	Long:  `"Walk through the commands.`,
	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Printf("\nThis operation will walk through all the commands.")
		// fmt.Printf("\nThe command will only be run if you answer 'y'.\n")

		// prompt := promptui.Prompt{
		// 	Label:     "Continue",
		// 	IsConfirm: true,
		// }

		// _, err := prompt.Run()

		// if err != nil {
		// 	return
		// }

		InitTool()

		ret, err := List()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
			return
		}

		total := len(ret)

		for index, cmd := range ret {

			pterm.Bold.Printf(pterm.LightBlue("Command "))
			header := "Command " + strconv.Itoa(index+1) + "/" + strconv.Itoa(total)

			pterm.Bold.Printf(pterm.LightBlue(header) + "\n\n")

			ret, err := ShowCmd(cmd.CmdHash)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Command failed: please run './recmd start' and try again.\n")
				return
			}

			var CustomTable = pterm.TablePrinter{
				HeaderStyle:    &pterm.Style{pterm.FgLightBlue, pterm.Bold},
				Style:          &pterm.Style{pterm.FgWhite},
				Separator:      " | ",
				SeparatorStyle: &pterm.ThemeDefault.TableSeparatorStyle,
			}

			CustomTable.WithHasHeader().WithData(pterm.TableData{
				{"Field", "Value"},
				{"Hash", cmd.CmdHash},
				{"Description", cmd.Description},
				{"Command", cmd.CmdString},
			}).Render()

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

			pterm.Println()
			pterm.FgLightYellow.Println(sc.Coutput)

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
