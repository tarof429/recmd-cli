package cli

// Common functions

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

// DisplayQueue lists the queue of commands in a table format
func DisplayQueue(ret []Command) {
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)

	defer w.Flush()

	show := func(a, b, c, d interface{}) {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", a, b, c, d)
	}

	show("HASH", "COMMAND", "DESCRIPTION", "STATUS")

	for _, c := range ret {
		cmdHash := c.CmdHash[0:15]

		var cmdString string
		var Description string

		if len(c.CmdString) > 40 {
			cmdString = c.CmdString[0:40] + "..."
		} else {
			cmdString = c.CmdString
		}

		if len(c.Description) > 50 {
			Description = c.Description[0:50] + "..."
		} else {
			Description = c.Description
		}

		Status := c.Status

		show(cmdHash, cmdString, Description, Status)
	}
}

// Display lists the given list of commands in a table format
func Display(ret []Command) {
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)

	defer w.Flush()

	show := func(a, b, c, d interface{}) {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", a, b, c, d)
	}

	show("HASH", "COMMAND", "DESCRIPTION", "DURATION")

	for _, c := range ret {
		cmdHash := c.CmdHash[0:15]

		var cmdString string
		var Description string

		if len(c.CmdString) > 40 {
			cmdString = c.CmdString[0:40] + "..."
		} else {
			cmdString = c.CmdString
		}

		if len(c.Description) > 50 {
			Description = c.Description[0:50] + "..."
		} else {
			Description = c.Description
		}

		var durationString string

		if int(c.Duration.Minutes()) > 0 {
			minutes := strconv.FormatFloat(c.Duration.Minutes(), 'f', 0, 64)
			seconds := strconv.FormatFloat(c.Duration.Seconds(), 'f', 0, 64)
			durationString = minutes + " minute(s) " + seconds + " second(s)"
		} else {

			seconds := strconv.FormatFloat(c.Duration.Seconds(), 'f', 0, 64)
			if seconds == "-0" {
				durationString = "-"
			} else {
				durationString = seconds + " second(s)"
			}

		}

		show(cmdHash, cmdString, Description, durationString)

	}
}
