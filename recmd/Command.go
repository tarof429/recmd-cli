package recmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/theckman/yacspin"
)

// Command represents a command and optionally a comment to document what the command does
type Command struct {
	CmdHash   string        `json:"commandHash"`
	CmdString string        `json:"commandString"`
	Comment   string        `json:"comment"`
	Duration  time.Duration `json:"duration"`
}

// ScheduledCommand represents a command that is scheduled to run
type ScheduledCommand struct {
	Command
	Coutput    string    `json:"coutput"`
	ExitStatus int       `json:"exitStatus"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

const historyFile = ".cmd_history.json"

// ReadCmdHistoryFile reads historyFile and generates a list of Command structs
func ReadCmdHistoryFile(dir string) ([]Command, error) {

	var (
		historyData []byte    // Data representing our history file
		cmds        []Command // List of commands produced after unmarshalling historyData
		err         error     // Any errors we might encounter
	)

	// Read the history file into historyData
	historyData, err = ioutil.ReadFile(dir + "/" + historyFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while reading historyfile: %v\n", err)
		return cmds, err
	}

	// Unmarshall historyData into a list of commands
	err = json.Unmarshal(historyData, &cmds)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while unmarshalling: %v\n", err)
	}

	return cmds, err

}

// SelectCmd returns a command
func SelectCmd(dir string, value string) (Command, error) {

	cmds, error := ReadCmdHistoryFile(dir)

	if error != nil {
		return Command{}, error
	}

	for _, cmd := range cmds {

		if strings.Index(cmd.CmdHash, value) == 0 {
			return cmd, nil
		}
	}

	return Command{}, nil
}

// DeleteCmd deletes a command. It's best to pass in the commandHash
// because commands may look similar.
func DeleteCmd(dir string, value string) int {

	cmds, error := ReadCmdHistoryFile(dir)

	if error != nil {
		return -1
	}

	foundIndex := -1

	for index, cmd := range cmds {
		if strings.Index(cmd.CmdHash, value) == 0 {
			foundIndex = index
			break
		}
	}

	if foundIndex != -1 {
		//fmt.Println("Found command. Found index was " + strconv.Itoa(foundIndex))
		// We may want to do more investigation to know why this works...
		cmds = append(cmds[:foundIndex], cmds[foundIndex+1:]...)

		// Return whether we are able to overwrite the history file
		OverwriteCmdHistoryFile(dir, cmds)
	}

	return foundIndex
}

// OverwriteCmdHistoryFile overwrites the history file with []Command passed in as a parameter
func OverwriteCmdHistoryFile(dir string, cmds []Command) bool {

	mode := int(0644)

	updatedData, _ := json.MarshalIndent(cmds, "", "\t")

	error := ioutil.WriteFile(dir+"/"+historyFile, updatedData, os.FileMode(mode))

	return error == nil
}

// CreateCmdHistoryFile creates an empty history file
func CreateCmdHistoryFile(dir string) bool {

	// Check if the file does not exist. If not, then create it and add our first command to it.
	f, err := os.Open(dir + "/" + historyFile)

	// Immediately close the file since we plan to write to it
	defer f.Close()

	// Check if the file doesn't exist and if so, then write it.
	if err != nil {

		mode := int(0644)

		error := ioutil.WriteFile(dir+"/"+historyFile, []byte(nil), os.FileMode(mode))

		return error == nil
	}
	return true
}

// UpdateCommandDuration updates a command with the same hash in the history file
func UpdateCommandDuration(dir string, cmd Command, duration time.Duration) bool {

	// Check if the file does not exist. If not, then create it and add our first command to it.
	f, err := os.Open(dir + "/" + historyFile)

	// Immediately close the file since we plan to write to it
	f.Close()

	// Check if the file doesn't exist and if so, then write it.
	if err != nil {

		// The array of commands
		var cmds []Command

		// Set the duration
		cmd.Duration = duration

		cmds = append(cmds, cmd)

		mode := int(0644)

		updatedData, _ := json.MarshalIndent(cmds, "", "\t")

		error := ioutil.WriteFile(dir+"/"+historyFile, updatedData, os.FileMode(mode))

		return error == nil
	}

	// Update the command in the history file

	// The array of commands
	var cmds []Command

	// Read history file
	data, err := ioutil.ReadFile(dir + "/" + historyFile)

	// An error occured while reading historyFile.
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return false
	}

	if err := json.Unmarshal(data, &cmds); err != nil {
		fmt.Fprintf(os.Stderr, "JSON unmarshalling failed: %s\n", err)
		return false
	}

	//fmt.Println("Updating duration")

	var found bool
	var foundIndex int

	// Update the duration for the command
	for index, c := range cmds {
		if c.CmdHash == cmd.CmdHash {
			//fmt.Println("Found command")
			foundIndex = index
			found = true
			//c.Duration = cmd.Duration
			//fmt.Fprintf(os.Stderr, "Command hash already exists: %s\n", cmd.CmdString)
			break
			//return false
		}
	}

	if found == true {
		cmds[foundIndex].Duration = duration
		//fmt.Println(cmds[foundIndex])
	}

	// Convert the struct to JSON
	updatedData, updatedDataErr := json.MarshalIndent(cmds, "", "\t")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", updatedDataErr)
	}

	mode := int(0644)

	error := ioutil.WriteFile(dir+"/"+historyFile, updatedData, os.FileMode(mode))

	return error == nil
}

// WriteCmdHistoryFile writes a command to the history file
func WriteCmdHistoryFile(dir string, cmd Command) bool {

	// Check if the file does not exist. If not, then create it and add our first command to it.
	f, err := os.Open(dir + "/" + historyFile)

	// Immediately close the file since we plan to write to it
	f.Close()

	// Check if the file doesn't exist and if so, then write it.
	if err != nil {
		// The array of commands
		var cmds []Command

		cmds = append(cmds, cmd)

		mode := int(0644)

		updatedData, _ := json.MarshalIndent(cmds, "", "\t")

		error := ioutil.WriteFile(dir+"/"+historyFile, updatedData, os.FileMode(mode))

		return error == nil
	}

	// Update the command in the history file

	// The array of commands
	var cmds []Command

	// Read history file
	data, err := ioutil.ReadFile(dir + "/" + historyFile)

	// An error occured while reading historyFile.
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return false
	}

	if err := json.Unmarshal(data, &cmds); err != nil {
		fmt.Fprintf(os.Stderr, "JSON unmarshalling failed: %s\n", err)
		return false
	}

	// Check if the command hash alaready exists, and prevent the user from adding the same command
	for _, c := range cmds {
		if c.CmdHash == cmd.CmdHash {
			// c.Duration = cmd.Duration
			fmt.Fprintf(os.Stderr, "Command hash already exists: %s\n", cmd.CmdString)
			//break
			return false
		}
	}

	cmds = append(cmds, cmd)

	// Convert the struct to JSON
	updatedData, updatedDataErr := json.MarshalIndent(cmds, "", "\t")

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", updatedDataErr)
	}

	mode := int(0644)

	error := ioutil.WriteFile(dir+"/"+historyFile, updatedData, os.FileMode(mode))

	return error == nil

}

// NewCommand creates a new Command struct and populates the fields
func NewCommand(cmdString string, cmdComment string) Command {

	formattedHash := func() string {
		h := sha1.New()
		h.Write([]byte(cmdString))
		return fmt.Sprintf("%.15x", h.Sum(nil))
	}()

	cmd := Command{formattedHash,
		strings.Trim(cmdString, ""),
		strings.Trim(cmdComment, ""),
		-1}

	return cmd
}

// ScheduleCommand runs a Command based on a function passed in as the second parameter.
// This gives the ability to run Commands in multiple ways; for example, as a "mock" command
// (RunMockCommand) or a shell script command (RunShellScriptCommand).
func ScheduleCommand(cmd Command, f func(*ScheduledCommand, chan int)) ScheduledCommand {
	var sc ScheduledCommand

	sc.CmdHash = cmd.CmdHash
	sc.CmdString = cmd.CmdString
	sc.Comment = cmd.Comment
	sc.Duration = -1

	// Create a channel to hold exit status
	c := make(chan int)

	// Set the start time
	sc.StartTime = time.Now()

	// Run the command in a goroutine
	go f(&sc, c)

	// Receive the exit status of the command
	status := <-c

	now := time.Now()

	// Set end time after we receive from the channel
	sc.EndTime = now

	// Calculate the duration and store it
	sc.Duration = now.Sub(sc.StartTime)

	// The main reason why this code exists is to use the value received from the channel.
	if status != 0 {
		fmt.Fprintf(os.Stderr, "\nError: command failed.\n")
	}

	return sc
}

// RunMockCommand runs a mock command
func RunMockCommand(sc *ScheduledCommand, c chan int) {
	time.Sleep(1 * time.Second)
	sc.ExitStatus = 99
	sc.Coutput = "Mock stdout message"
	c <- sc.ExitStatus
}

// RunShellScriptCommand runs a command written to a temporary file
func RunShellScriptCommand(sc *ScheduledCommand, c chan int) {

	tempFile, err := ioutil.TempFile(os.TempDir(), "recmd-")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to create temp file: %d\n", err)
	}

	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("#!/bin/sh\n\n" + sc.CmdString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Errror: unable to write script to temp file: : %s\n", err)
	}

	cmd := exec.Command("sh", tempFile.Name())

	// We may want to make this configurable in the future.
	// For now, all commands will be run from the user's home directory
	cmd.Dir, err = os.UserHomeDir()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to obtain home directory: %s\n", err)
	}

	//out, err := cmd.Output()

	combinedOutput, combinedOutputErr := cmd.CombinedOutput()

	// fmt.Fprintf(os.Stdout, "\nError: %s error 2: %v\n", string(combinedOutput), err2)

	if combinedOutputErr != nil {
		sc.ExitStatus = -1
	}

	sc.Coutput = string(combinedOutput)

	c <- sc.ExitStatus
}

// RunShellScriptCommandWithSpinner runs a command with a spinner
func RunShellScriptCommandWithSpinner(sc *ScheduledCommand, c chan int) {

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[14],
		Suffix:          " Scheduling commmand ",
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, _ := yacspin.New(cfg)

	spinner.Start()

	RunShellScriptCommand(sc, c)

	spinner.Stop()
}
