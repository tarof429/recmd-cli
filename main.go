package main

import (
	"fmt"

	recmd "github.com/tarof429/recmd"
)

func main() {

	cmd := recmd.NewCommand("ls /bin", "List files")

	sc := recmd.ScheduleCommand(cmd, recmd.RunShellScriptCommand)

	fmt.Println(sc.Stdout)

	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

}
