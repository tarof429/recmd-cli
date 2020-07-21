package recmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
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

// func TestReadCmdHistoryFile(t *testing.T) {

// 	cmd := Command{"abc", "ls", "list files", -1}
// 	cmd2 := Command{"def", "df", "Show disk usage", -1}

// 	cmds := []Command{cmd, cmd2}

// 	// Convert the struct to JSON
// 	data, err := json.MarshalIndent(cmds, "", "\t")

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	mode := int(0644)

// 	ioutil.WriteFile(testHistoryFile, data, os.FileMode(mode))

// 	readCmds, err := ReadCmdHistoryFile(testdataDir)

// 	if err != nil {
// 		t.Error("an error occured while reading command history")
// 	}

// 	// Check to make sure that the commands we wrote are the ones we created by comparing hashes
// 	for i := 0; i < len(cmds); i++ {
// 		if cmds[i].CmdHash != readCmds[i].CmdHash {
// 			t.Error("The command histories are not equal")
// 		}

// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

// func TestCreateCmdHistoryFiel(t *testing.T) {
// 	if CreateCmdHistoryFile(testdataDir) == false {
// 		t.Fail()
// 	}

// }

// func TestWriteCmdHistoryFile(t *testing.T) {

// 	cmd := NewCommand("find / -type f -perm 0777 -print -exec chmod 644 {} \\;", "Find all 777 permission files and use chmod command to set permissions to 644")
// 	cmd2 := NewCommand("df --print-type --total --human-readable /home /dev/sda6", "Display Total Information of Partitions in Human Readable Terms")

// 	cmds := []Command{cmd, cmd2}

// 	// Convert the struct to JSON
// 	data, err := json.MarshalIndent(cmds, "", "\t")

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	mode := int(0644)

// 	// Write two commands to cmd_history.json
// 	ioutil.WriteFile(testHistoryFile, data, os.FileMode(mode))

// 	// Define a new command
// 	cmd3 := NewCommand("rsync -avz rpmpkgs/ root@192.168.0.101:/home/", "Copy a Directory from Local Server to a Remote Server")

// 	WriteCmdHistoryFile(testdataDir, cmd3)

// 	readCmds, err := ReadCmdHistoryFile(testdataDir)

// 	found1, found2, found3 := false, false, false

// 	for _, c := range readCmds {
// 		if c.CmdString == "find / -type f -perm 0777 -print -exec chmod 644 {} \\;" {
// 			found1 = true
// 		} else if c.CmdString == "df --print-type --total --human-readable /home /dev/sda6" {
// 			found2 = true
// 		} else if c.CmdString == "rsync -avz rpmpkgs/ root@192.168.0.101:/home/" {
// 			found3 = true
// 		}
// 	}

// 	if found1 == false {
// 		t.Error("The ls command was missing from history")
// 	}

// 	if found2 == false {
// 		t.Error("The df command was missing from history")
// 	}

// 	if found3 == false {
// 		t.Error("The top command was missing from history")
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

// func TestWriteSameCommands(t *testing.T) {

// 	cmd := NewCommand("ls", "list files")
// 	cmd2 := NewCommand("ls", "Show disk usage")

// 	if WriteCmdHistoryFile(testdataDir, cmd) == false {
// 		t.Error("Unable to write " + cmd.CmdString)
// 	}

// 	if WriteCmdHistoryFile(testdataDir, cmd2) == true {
// 		t.Error("Wrongfully wrote the same command to history " + cmd.CmdString)
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})

// }

// func TestNewCommand(t *testing.T) {

// 	cmd := NewCommand("df /usr", "Find disk usage")

// 	if cmd.CmdString != "df /usr" {
// 		t.Error(cmd.CmdString + " was not expected")
// 	}

// 	cmd2 := NewCommand("free", "Find memory usage")

// 	if cmd2.CmdString != "free" {
// 		t.Error(cmd.CmdString + " was not expected")
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

// func TestMultipleNewCommand(t *testing.T) {

// 	cmd := NewCommand("df /usr fdfdfsaasf fsfadf", "Find disk usage")
// 	cmd2 := NewCommand("df /usr fdfdfsaasf fsfadf", "Find disk usage")

// 	// Test whether the hashes are the same. They should be because the command line is the same.
// 	if cmd.CmdHash != cmd2.CmdHash {

// 		// fmt.Println(cmd.CmdHash)
// 		// fmt.Println(cmd2.CmdHash)

// 		t.Error("The hashes for the two commands were not the same")
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

// func TestRunMockCommand(t *testing.T) {

// 	cmd := NewCommand("ls", "List files")

// 	sc := ScheduleCommand(cmd, RunMockCommand)

// 	if sc.ExitStatus != 99 {
// 		t.Error("The exit status of the command was not 0")
// 	}
// 	// data, _ := json.MarshalIndent(sc, "", "\t")
// 	// fmt.Println(string(data))

// 	time.Sleep(1 * time.Second)

// 	sc = ScheduleCommand(cmd, RunMockCommand)

// 	if sc.ExitStatus != 99 {
// 		t.Error("The exit status of the command was not 0")
// 	}

// 	fmt.Printf("%0.0f seconds\n", sc.Duration.Seconds())
// 	// data, _ = json.MarshalIndent(sc, "", "\t")
// 	// fmt.Println(string(data))

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

func TestRunShellScriptCommand(t *testing.T) {

	cmd := NewCommand("sleep 2", "Sleep 2 seconds")

	ret := WriteCmdHistoryFile(testdataDir, cmd)

	if ret == false {
		fmt.Println("Unable to write history file")
	}

	cmd2 := NewCommand("ls /", "List files under root")

	ret = WriteCmdHistoryFile(testdataDir, cmd2)

	if ret == false {
		fmt.Println("Unable to write history file")
	}

	cmds, _ := ReadCmdHistoryFile(testdataDir)

	var found bool
	var foundHash string

	found = false

	for _, cmd = range cmds {
		if cmd.CmdString == "sleep 2" {
			found = true
			foundHash = cmd.CmdHash
			break
		}
	}

	if found == false {
		fmt.Println("Unable to find command")
	}

	selectedCmd, _ := SelectCmd(testdataDir, foundHash)

	sc := ScheduleCommand(selectedCmd, RunShellScriptCommand)

	if sc.ExitStatus != 0 {
		t.Error("The exit status of the command was not 0")
	}

	ret = UpdateCommandDuration(testdataDir, cmd, sc.Duration)

	if ret != true {
		t.Error("Error while updating command")
	}
	cmds, _ = ReadCmdHistoryFile(testdataDir)

	data, _ := json.MarshalIndent(cmds, "", "\t")
	fmt.Println(string(data))

	// t.Cleanup(func() {
	// 	os.Remove(testHistoryFile)
	// })

}

// func TestRunShellScriptCommandInvalid(t *testing.T) {

// 	cmd := NewCommand("lslsls", "List files")

// 	sc := ScheduleCommand(cmd, RunShellScriptCommand)

// 	if sc.ExitStatus == 0 {
// 		t.Error("The exit status of the command was 0")
// 	}
// 	// data, _ := json.MarshalIndent(sc, "", "\t")
// 	// fmt.Println(string(data))

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})

// }

// func TestRunShellScriptCommandMultiple(t *testing.T) {

// 	cmd := NewCommand("cd /; ls; cd /home; ls", "List files")

// 	sc := ScheduleCommand(cmd, RunShellScriptCommand)

// 	if sc.ExitStatus != 0 {
// 		t.Error("The exit status of the command was not 0")
// 	}

// 	// data, _ := json.MarshalIndent(sc, "", "\t")
// 	// fmt.Println(string(data))

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})

// }

// func TestRunShellScriptCommandLongRunning(t *testing.T) {

// 	cmd := NewCommand("sleep 1", "Take a brief nap")

// 	sc := ScheduleCommand(cmd, RunShellScriptCommand)

// 	if sc.ExitStatus != 0 {
// 		t.Error("The exit status of the command was not 0")
// 	}

// 	// data, _ := json.MarshalIndent(sc, "", "\t")
// 	// fmt.Println(string(data))

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})

// }

// // TODO: Need to change test to read the hash and use it to select it
// func TestScheduleCommand(t *testing.T) {

// 	cmd := NewCommand("uname", "Show my name")

// 	result := WriteCmdHistoryFile(testdataDir, cmd)

// 	if result != true {
// 		t.Error("Unable to write history file")
// 	}

// 	cmds, _ := ReadCmdHistoryFile(testdataDir)

// 	var foundHash string
// 	for _, c := range cmds {
// 		if c.CmdString == "uname" {
// 			foundHash = c.CmdHash
// 			break
// 		}
// 	}

// 	if foundHash == "" {
// 		t.Error("Command not found")
// 	}

// 	ret, err := SelectCmd(testdataDir, foundHash)

// 	if err != nil {
// 		t.Error("Unable to read history file")
// 	}

// 	sc := ScheduleCommand(ret, RunShellScriptCommand)

// 	if sc.ExitStatus != 0 {
// 		t.Error("The exit status of the command was not 0")
// 	}

// 	data, _ := json.MarshalIndent(sc, "", "\t")
// 	fmt.Println(string(data))

// 	if sc.ExitStatus != 0 {
// 		t.Error("Exit status was not 0")
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

// func TestDeleteCommandUsingCommandHash(t *testing.T) {
// 	cmd := Command{"abc", "cp", "comment a", -1}
// 	cmd2 := Command{"def", "mv", "comment b", -1}
// 	cmd3 := Command{"ghk", "sleep", "comment c", -1}

// 	WriteCmdHistoryFile(testdataDir, cmd)
// 	WriteCmdHistoryFile(testdataDir, cmd2)
// 	WriteCmdHistoryFile(testdataDir, cmd3)

// 	foundIndex := DeleteCmd(testdataDir, "def")

// 	if foundIndex == -1 {
// 		t.Error("Command was not deleted")
// 	}

// 	ret, err := ReadCmdHistoryFile(testdataDir)

// 	if err != nil {
// 		t.Error("Received an error")
// 	}

// 	found := false
// 	for _, cmd := range ret {
// 		if cmd.CmdHash == "def" {
// 			found = true

// 		}
// 	}

// 	if found == true {
// 		t.Error("Command was not deleted as expected")
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }

// func TestOverwriteCmdHistoryFile(t *testing.T) {

// 	cmd := NewCommand("cp", "comment a")
// 	cmd2 := NewCommand("mv", "comment b")
// 	cmd3 := NewCommand("sleep", "comment c")

// 	if WriteCmdHistoryFile(testdataDir, cmd) == false {
// 		t.Error("Unable to write " + cmd.CmdString)
// 	}

// 	if WriteCmdHistoryFile(testdataDir, cmd2) == false {
// 		t.Error("Unable to write " + cmd2.CmdString)
// 	}

// 	if WriteCmdHistoryFile(testdataDir, cmd3) == false {
// 		t.Error("Unable to write " + cmd3.CmdString)
// 	}

// 	cmds, err := ReadCmdHistoryFile(testdataDir)

// 	if err != nil {
// 		t.Error("Unable to read history file")
// 	}
// 	//updatedData, _ := json.MarshalIndent(cmds, "", "\t")
// 	//fmt.Println(string(updatedData))

// 	found := false
// 	var foundHash string

// 	for _, cmd := range cmds {
// 		if cmd.CmdString == "sleep" {
// 			foundHash = cmd.CmdHash
// 			found = true
// 			break
// 		}
// 	}

// 	if found == false {
// 		t.Error("Command was not deleted as expected")
// 	}

// 	cmd, _ = SelectCmd(testdataDir, foundHash)

// 	if cmd.CmdString != "sleep" {
// 		t.Error("Unable to fnd the command")
// 	}

// 	t.Cleanup(func() {
// 		os.Remove(testHistoryFile)
// 	})
// }
