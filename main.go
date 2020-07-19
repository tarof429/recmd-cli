package main

import (
	"fmt"

	recmd "github.com/tarof429/recmd"
)

func main() {

	cmd := recmd.Command{}
	cmd.CmdString = "foo"
	fmt.Println(cmd.CmdString)
}
