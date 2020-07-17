package recmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadCmdHistoryFile(t *testing.T) {

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

	readCmds := ReadCmdHistoryFile(".")

	if !reflect.DeepEqual(cmds, readCmds) {
		t.Error("The command histories are not equal")
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
}

func TestWriteCmdHistoryFile(t *testing.T) {

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

	WriteCmdHistoryFile(".", cmd3)

	readCmds := ReadCmdHistoryFile(".")

	// Add cmd3 to our slice for comparison
	updatedCmds := append(cmds, cmd3)

	if !reflect.DeepEqual(updatedCmds, readCmds) {
		t.Error("The command histories are not equal")
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
}

func TestWriteSameCommands(t *testing.T) {

	cmd := NewCommand("ls", "list files")
	cmd2 := NewCommand("ls", "Show disk usage")

	if WriteCmdHistoryFile(".", cmd) == false {
		t.Error("Unable to write " + cmd.CmdString)
	}

	if WriteCmdHistoryFile(".", cmd2) == true {
		t.Error("Accidentally wrote " + cmd2.CmdString)
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

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

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
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

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
}

func TestRunMockCommand(t *testing.T) {

	cmd := NewCommand("ls", "List files")

	sc := ScheduleCommand(cmd, RunMockCommand)

	// if sc.ExitStatus != 0 {
	// 	t.Error("The exit status of the command was not 0")
	// }
	data, _ := json.MarshalIndent(sc, "", "\t")

	fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
}

func TestRunCommand(t *testing.T) {

	cmd := NewCommand("ls /", "List files")

	sc := ScheduleCommand(cmd, RunCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}
	data, _ := json.MarshalIndent(sc, "", "\t")

	fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

}

func TestRunCommandInvalid(t *testing.T) {

	cmd := NewCommand("lslsls", "List files")

	sc := ScheduleCommand(cmd, RunCommand)

	if sc.ExitStatus == 0 {
		t.Error("The exit status of the command was 0")
	}
	data, _ := json.MarshalIndent(sc, "", "\t")

	fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

}

func TestRunCommandMultiple(t *testing.T) {

	cmd := NewCommand("cd /; ls; cd /home; ls", "List files")

	sc := ScheduleCommand(cmd, RunCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	data, _ := json.MarshalIndent(sc, "", "\t")

	fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

}

func TestRunCommandLongRunning(t *testing.T) {

	cmd := NewCommand("sleep 1", "Take a brief nap")

	sc := ScheduleCommand(cmd, RunCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	data, _ := json.MarshalIndent(sc, "", "\t")

	fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

}

func TestRunByCommandString(t *testing.T) {

	cmd := NewCommand("uname", "Show my name")

	result := WriteCmdHistoryFile(".", cmd)

	if result != true {
		t.Error("Unable to write history file")
	}

	ret := SelectCmd(".", "commandString", "uname")

	sc := ScheduleCommand(ret, RunCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	data, _ := json.MarshalIndent(sc, "", "\t")

	fmt.Println(string(data))

	if sc.ExitStatus != 0 {
		t.Error("Exit status was not 0")
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
}

// This test attempts to run a command by a matching hash. For this test, it will fail on purpose due
// to an invalid hash. In real life, users or the UI will most likely select commands by their hash since
// it should be a unique identifier.
func TestRunByCommandHash(t *testing.T) {

	cmd := NewCommand("uname", "Show my name")

	result := WriteCmdHistoryFile(".", cmd)

	if result != true {
		t.Error("Unable to write history file")
	}

	// Attempt to select a command whose hash is 'uname'. This should never succeed.
	ret := SelectCmd(".", "commandHash", "uname")

	// Show the contents of the command. The fields should be empty.
	data, _ := json.MarshalIndent(ret, "", "\t")
	fmt.Println(string(data))

	if ret.CmdHash != "" {
		t.Error("Accidentally did not find an empty command")
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

}

func TestDeleteCommand(t *testing.T) {
	cmd := Command{"abc", "cp", "comment a"}
	cmd2 := Command{"def", "mv", "comment b"}
	cmd3 := Command{"ghk", "sleep", "comment c"}

	cmds := []Command{cmd, cmd2, cmd3}

	// Minitest 1: Test whether we can remove a command by its hash
	ret := DeleteCmd(cmds, "commandHash", "abc")

	data, _ := json.MarshalIndent(ret, "", "\t")
	fmt.Println(string(data))

	for _, cmd := range ret {
		if cmd.CmdHash == "abc" {
			t.Error("Command was not deleted as expected")
		}
	}

	// Minitest 2: Test whether we can remove a command by its name
	cmds = append(ret, cmd)

	ret = DeleteCmd(cmds, "commandString", "mv")

	data, _ = json.MarshalIndent(ret, "", "\t")
	fmt.Println(string(data))

	for _, cmd := range ret {
		if cmd.CmdHash == "def" {
			t.Error("Command was not deleted as expected")
		}
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})
}

func TestOverwriteCmdHistoryFile(t *testing.T) {

	cmd := NewCommand("cp", "comment a")
	cmd2 := NewCommand("mv", "comment b")
	cmd3 := NewCommand("sleep", "comment c")

	if WriteCmdHistoryFile(".", cmd) == false {
		t.Error("Unable to write " + cmd.CmdString)
	}

	time.Sleep(1 * time.Second)
	if WriteCmdHistoryFile(".", cmd2) == false {
		t.Error("Unable to write " + cmd2.CmdString)
	}

	if WriteCmdHistoryFile(".", cmd3) == false {
		t.Error("Unable to write " + cmd3.CmdString)
	}

	cmds := ReadCmdHistoryFile(".")

	ret := DeleteCmd(cmds, "commandString", "sleep")

	// This change will remove the sleep command from the history file
	result := OverwriteCmdHistoryFile(".", ret)

	if result != true {
		t.Error("Unable to overwrite history file")
	}

	cmds = ReadCmdHistoryFile(".")

	updatedData, _ := json.MarshalIndent(cmds, "", "\t")

	fmt.Println(string(updatedData))

	for _, cmd := range ret {
		if cmd.CmdString == "sleep" {
			t.Error("Command was not deleted as expected")
		}
	}

	t.Cleanup(func() {
		os.Remove(historyFile)
	})

}
