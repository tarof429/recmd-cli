package recmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestReadCmdHistoryFile(t *testing.T) {

	os.Remove(historyFile)

	cmd := Command{"abc", "ls", "list files"}
	cmd2 := Command{"def", "df", "Show disk usage"}

	cmds := []Command{cmd, cmd2}

	// Convert the struct to JSON
	data, err := json.MarshalIndent(cmds, "", "\t")

	if err != nil {
		t.Error(err)
	}

	mode := int(0644)

	ioutil.WriteFile(historyFile, data, os.FileMode(mode))

	readCmds := ReadCmdHistoryFile()

	if !reflect.DeepEqual(cmds, readCmds) {
		t.Error("The command histories are not equal")
	}
}

func TestWriteCmdHistoryFile(t *testing.T) {

	os.Remove(historyFile)

	cmd := NewCommand("ls", "list files")
	cmd2 := NewCommand("df", "Show disk usage")

	cmds := []Command{cmd, cmd2}

	// Convert the struct to JSON
	data, err := json.MarshalIndent(cmds, "", "\t")

	if err != nil {
		t.Error(err)
	}

	mode := int(0644)

	// Write two commands to cmd_history.json
	ioutil.WriteFile(historyFile, data, os.FileMode(mode))

	// Define a new command
	cmd3 := NewCommand("top", "Show active processes")

	WriteCmdHistoryFile(cmd3)

	readCmds := ReadCmdHistoryFile()

	// Add cmd3 to our slice for comparison
	updatedCmds := append(cmds, cmd3)

	if !reflect.DeepEqual(updatedCmds, readCmds) {
		t.Error("The command histories are not equal")
	}

}

func TestWriteSameCommands(t *testing.T) {

	os.Remove(historyFile)

	cmd := NewCommand("ls", "list files")
	cmd2 := NewCommand("ls", "Show disk usage")

	if WriteCmdHistoryFile(cmd) == false {
		t.Error("Unable to write " + cmd.CmdString)
	}

	if WriteCmdHistoryFile(cmd2) == true {
		t.Error("Accidentally wrote " + cmd2.CmdString)
	}

}

func TestNewCommand(t *testing.T) {

	cmd := NewCommand("df /usr", "Find disk usage")

	if cmd.CmdString != "df /usr" {
		t.Error(cmd.CmdString + " was not expected")
	}

	cmd2 := NewCommand("free", "Find memory usage")

	if cmd2.CmdString != "free" {
		t.Error(cmd.CmdString + " was not expected")
	}
}

func TestMultipleNewCommand(t *testing.T) {

	cmd := NewCommand("df /usr fdfdfsaasf fsfadf", "Find disk usage")
	cmd2 := NewCommand("df /usr fdfdfsaasf fsfadf", "Find disk usage")

	// Test whether the hashes are the same. They should be because the command line is the same.
	if cmd.CmdHash != cmd2.CmdHash {

		// fmt.Println(cmd.CmdHash)
		// fmt.Println(cmd2.CmdHash)

		t.Error("The hashes for the two commands were not the same")
	}
}
