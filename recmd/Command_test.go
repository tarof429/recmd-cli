package recmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

const testdataDir = "testdata"
const testHistoryFile = testdataDir + "/.cmd_history.json"

func TestMain(m *testing.M) {
	fmt.Println("Running tests...")

	os.Remove(testHistoryFile)

	err := os.RemoveAll(testdataDir)

	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(testdataDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	// err = os.RemoveAll(testdataDir)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	os.Exit(status)
}

func TestReadCmdHistoryFile(t *testing.T) {

	cmd := Command{"abc", "ls", "list files", time.Now(), time.Now()}
	cmd2 := Command{"def", "df", "Show disk usage", time.Now(), time.Now()}

	cmds := []Command{cmd, cmd2}

	// Convert the struct to JSON
	data, err := json.MarshalIndent(cmds, "", "\t")

	if err != nil {
		t.Error(err)
	}

	mode := int(0644)

	ioutil.WriteFile(testHistoryFile, data, os.FileMode(mode))

	readCmds, err := ReadCmdHistoryFile(testdataDir)

	if err != nil {
		t.Error("an error occured while reading command history")
	}

	// Check to make sure that the commands we wrote are the ones we created by comparing hashes
	for i := 0; i < len(cmds); i++ {
		if cmds[i].CmdHash != readCmds[i].CmdHash {
			t.Error("The command histories are not equal")
		}

	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}

func TestCreateCmdHistoryFiel(t *testing.T) {
	if CreateCmdHistoryFile(testdataDir) == false {
		t.Fail()
	}

}

func TestWriteCmdHistoryFile(t *testing.T) {

	cmd := NewCommand("find / -type f -perm 0777 -print -exec chmod 644 {} \\;", "Find all 777 permission files and use chmod command to set permissions to 644")
	cmd2 := NewCommand("df --print-type --total --human-readable /home /dev/sda6", "Display Total Information of Partitions in Human Readable Terms")

	cmds := []Command{cmd, cmd2}

	// Convert the struct to JSON
	data, err := json.MarshalIndent(cmds, "", "\t")

	if err != nil {
		t.Error(err)
	}

	mode := int(0644)

	// Write two commands to cmd_history.json
	ioutil.WriteFile(testHistoryFile, data, os.FileMode(mode))

	// Define a new command
	cmd3 := NewCommand("rsync -avz rpmpkgs/ root@192.168.0.101:/home/", "Copy a Directory from Local Server to a Remote Server")

	WriteCmdHistoryFile(testdataDir, cmd3)

	readCmds, err := ReadCmdHistoryFile(testdataDir)

	found1, found2, found3 := false, false, false

	for _, c := range readCmds {
		if c.CmdString == "find / -type f -perm 0777 -print -exec chmod 644 {} \\;" {
			found1 = true
		} else if c.CmdString == "df --print-type --total --human-readable /home /dev/sda6" {
			found2 = true
		} else if c.CmdString == "rsync -avz rpmpkgs/ root@192.168.0.101:/home/" {
			found3 = true
		}
	}

	if found1 == false {
		t.Error("The ls command was missing from history")
	}

	if found2 == false {
		t.Error("The df command was missing from history")
	}

	if found3 == false {
		t.Error("The top command was missing from history")
	}

	// An attempt to display the commands in a table
	// fmt.Printf("COMMAND ID\tCOMMAND\tDESCRIPTION\n")

	// for _, c := range readCmds {

	// 	fmt.Printf("%s\t%s\t%s\t%v\t%v\n", c.CmdHash, c.CmdString, c.Comment, c.Creationtime, c.Modificationtime)

	// }

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}

func TestWriteSameCommands(t *testing.T) {

	cmd := NewCommand("ls", "list files")
	cmd2 := NewCommand("ls", "Show disk usage")

	if WriteCmdHistoryFile(testdataDir, cmd) == false {
		t.Error("Unable to write " + cmd.CmdString)
	}

	if WriteCmdHistoryFile(testdataDir, cmd2) == true {
		t.Error("Accidentally wrote " + cmd2.CmdString)
	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
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
		os.Remove(testHistoryFile)
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
		os.Remove(testHistoryFile)
	})
}

func TestRunMockCommand(t *testing.T) {

	cmd := NewCommand("ls", "List files")

	sc := ScheduleCommand(cmd, RunMockCommand)

	if sc.ExitStatus != 99 {
		t.Error("The exit status of the command was not 0")
	}
	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	time.Sleep(1 * time.Second)

	sc = ScheduleCommand(cmd, RunMockCommand)

	if sc.ExitStatus != 99 {
		t.Error("The exit status of the command was not 0")
	}

	// data, _ = json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}

func TestRunShellScriptCommand(t *testing.T) {

	cmd := NewCommand("ls /", "List files")

	sc := ScheduleCommand(cmd, RunShellScriptCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}
	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	// t.Cleanup(func() {
	// 	os.Remove(testHistoryFile)
	// })

}

func TestRunShellScriptCommandInvalid(t *testing.T) {

	cmd := NewCommand("lslsls", "List files")

	sc := ScheduleCommand(cmd, RunShellScriptCommand)

	if sc.ExitStatus == 0 {
		t.Error("The exit status of the command was 0")
	}
	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})

}

func TestRunShellScriptCommandMultiple(t *testing.T) {

	cmd := NewCommand("cd /; ls; cd /home; ls", "List files")

	sc := ScheduleCommand(cmd, RunShellScriptCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})

}

func TestRunShellScriptCommandLongRunning(t *testing.T) {

	cmd := NewCommand("sleep 1", "Take a brief nap")

	sc := ScheduleCommand(cmd, RunShellScriptCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})

}

func TestRunByCommandString(t *testing.T) {

	cmd := NewCommand("uname", "Show my name")

	result := WriteCmdHistoryFile(testdataDir, cmd)

	if result != true {
		t.Error("Unable to write history file")
	}

	ret, err := SelectCmd(testdataDir, "commandString", "uname")

	if err != nil {
		t.Error("Unable to read history file")
	}

	sc := ScheduleCommand(ret, RunShellScriptCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	// data, _ := json.MarshalIndent(sc, "", "\t")
	// fmt.Println(string(data))

	if sc.ExitStatus != 0 {
		t.Error("Exit status was not 0")
	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}

// This test attempts to run a command by a matching hash. For this test, it will fail on purpose due
// to an invalid hash. In real life, users or the UI will most likely select commands by their hash since
// it should be a unique identifier.
func TestRunByCommandHash(t *testing.T) {

	cmd := NewCommand("uname", "Show my name")

	result := WriteCmdHistoryFile(testdataDir, cmd)

	if result != true {
		t.Error("Unable to write history file")
	}

	// Attempt to select a command whose hash is 'uname'. This should never succeed.
	ret, err := SelectCmd(testdataDir, "commandHash", "uname")

	if err != nil {
		t.Error("Unable to read history file")
	}

	// Show the contents of the command. The fields should be empty.
	// data, _ := json.MarshalIndent(ret, "", "\t")
	// fmt.Println(string(data))

	if ret.CmdHash != "" {
		t.Error("Accidentally did not find an empty command")
	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})

}

func TestDeleteCommandUsingCommandHash(t *testing.T) {
	cmd := Command{"abc", "cp", "comment a", time.Now(), time.Now()}
	cmd2 := Command{"def", "mv", "comment b", time.Now(), time.Now()}
	cmd3 := Command{"ghk", "sleep", "comment c", time.Now(), time.Now()}

	WriteCmdHistoryFile(testdataDir, cmd)
	WriteCmdHistoryFile(testdataDir, cmd2)
	WriteCmdHistoryFile(testdataDir, cmd3)

	foundIndex := DeleteCmd(testdataDir, "def", "commandHash")

	if foundIndex == -1 {
		t.Error("Command was not deleted")
	}

	ret, err := ReadCmdHistoryFile(testdataDir)

	if err != nil {
		t.Error("Received an error")
	}

	found := false
	for _, cmd := range ret {
		if cmd.CmdHash == "def" {
			found = true

		}
	}

	if found == true {
		t.Error("Command was not deleted as expected")
	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}

func TestDeleteCommandUsingCommandName(t *testing.T) {
	cmd := Command{"abc", "cp", "comment a", time.Now(), time.Now()}
	cmd2 := Command{"def", "mv", "comment b", time.Now(), time.Now()}
	cmd3 := Command{"ghk", "sleep", "comment c", time.Now(), time.Now()}

	WriteCmdHistoryFile(testdataDir, cmd)
	WriteCmdHistoryFile(testdataDir, cmd2)
	WriteCmdHistoryFile(testdataDir, cmd3)

	foundIndex := DeleteCmd(testdataDir, "sleep", "commandString")

	if foundIndex == -1 {
		t.Error("Command was not deleted")
	}

	ret, err := ReadCmdHistoryFile(testdataDir)

	if err != nil {
		t.Error("Received an error")
	}

	found := false
	for _, cmd := range ret {
		if cmd.CmdString == "sle" {
			found = true

		}
	}

	if found == true {
		t.Error("Command was not deleted as expected")
	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}

func TestOverwriteCmdHistoryFile(t *testing.T) {

	cmd := NewCommand("cp", "comment a")
	cmd2 := NewCommand("mv", "comment b")
	cmd3 := NewCommand("sleep", "comment c")

	if WriteCmdHistoryFile(testdataDir, cmd) == false {
		t.Error("Unable to write " + cmd.CmdString)
	}

	time.Sleep(1 * time.Second)
	if WriteCmdHistoryFile(testdataDir, cmd2) == false {
		t.Error("Unable to write " + cmd2.CmdString)
	}

	if WriteCmdHistoryFile(testdataDir, cmd3) == false {
		t.Error("Unable to write " + cmd3.CmdString)
	}

	cmds, err := ReadCmdHistoryFile(testdataDir)

	if err != nil {
		t.Error("Unable to read history file")
	}

	foundIndex := DeleteCmd(testdataDir, "sleep", "commandString")

	if foundIndex == -1 {
		t.Error("Failed to delete command")
	}

	cmds, err = ReadCmdHistoryFile(testdataDir)

	if err != nil {
		t.Error("Unable to read history file")
	}
	//updatedData, _ := json.MarshalIndent(cmds, "", "\t")
	//fmt.Println(string(updatedData))

	for _, cmd := range cmds {
		if cmd.CmdString == "sleep" {
			t.Error("Command was not deleted as expected")
		}
	}

	t.Cleanup(func() {
		os.Remove(testHistoryFile)
	})
}
