package recmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Command represents a command and optionally a comment to document what the command does
type Command struct {
	CmdHash          string    `json:"commandHash"`
	CmdString        string    `json:"commandString"`
	Comment          string    `json:"comment"`
	Creationtime     time.Time `json:"creationTime"`
	Modificationtime time.Time `json:"lastExecutionTime"`
}

// ScheduledCommand represents a command that is scheduled to run
type ScheduledCommand struct {
	Command
	Stdout     string    `json:"stdout"`
	Stderr     string    `json:"stderr"`
	ExitStatus int       `json:"exitStatus"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}

// // Config represents global configuration
// type Config struct {
// 	Salt string `json:"salt"`
// }

const historyFile = ".cmd_history.json"

//const configFile = ".recmd.json"

// LoadConfigFile reads configFile and initializes recmd.
// func LoadConfigFile(prefix string) Config {
// 	var config Config

// 	data, err := ioutil.ReadFile(prefix + "/" + configFile)

// 	// An error occured while reading historyFile.
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		return config
// 	}

// 	if err := json.Unmarshal(data, &config); err != nil {
// 		log.Fatalf("JSON unmarshalling failed: %s\n", err)
// 	}

// 	return config
// }

// GenerateConfig generates a config file
// func GenerateConfig(prefix string) Config {

// 	// Generate a new hash
// 	h := sha1.New()

// 	// Set the string we want to hash
// 	rand.Seed(time.Now().UTC().UnixNano())
// 	rand := rand.Int()
// 	io.WriteString(h, strconv.Itoa(rand))

// 	// Create Config with our hash
// 	var config = Config{hex.EncodeToString(h.Sum(nil))}

// 	return config
// }

// WriteConfig writes Config to a fiel
// func WriteConfig(path string, config Config) {
// 	// Convert the struct to JSON
// 	data, err := json.MarshalIndent(config, "", "\t")

// 	if err != nil {
// 		fmt.Printf("%s\n", err)
// 	}

// 	mode := int(0644)

// 	ioutil.WriteFile(path+"/"+configFile, data, os.FileMode(mode))
// }

// ReadCmdHistoryFile reads historyFile and generates a list of Command structs
func ReadCmdHistoryFile(dir string) ([]Command, error) {

	var cmds []Command

	data, err := ioutil.ReadFile(dir + "/" + historyFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while reading historyfile: %v\n", err)
		return cmds, err
	}

	if err := json.Unmarshal(data, &cmds); err != nil {
		//fmt.Fprintf(os.Stderr, "JSON unmarshalling failed: %s\n", err)
		return cmds, err
	}

	return cmds, nil
}

// SelectCmd returns a command
func SelectCmd(dir string, field string, value string) ([]Command, error) {

	cmds, error := ReadCmdHistoryFile(dir)

	var ret []Command

	if error != nil {
		return ret, error
	}

	for _, cmd := range cmds {
		switch field {
		case "commandString":
			if strings.Index(cmd.CmdString, value) == 0 {
				ret = append(ret, cmd)
			}
		case "commandHash":
			if strings.Index(cmd.CmdHash, value) == 0 {
				ret = append(ret, cmd)
			}
		}
	}

	return ret, nil
}

// DeleteCmd deletes a command. It's best to pass in the commandHash
// because commands may look similar.
func DeleteCmd(dir string, value string, field string) int {

	cmds, error := ReadCmdHistoryFile(dir)

	if error != nil {
		return -1
	}

	foundIndex := -1

	for index, cmd := range cmds {
		//fmt.Println(cmd)

		switch field {
		case "commandString":

			if strings.Index(cmd.CmdString, value) == 0 {
				foundIndex = index
				break
			}
		case "commandHash":

			if strings.Index(cmd.CmdHash, value) == 0 {
				foundIndex = index
				break
			}
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

	// Add the command to the history file

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

	// Loop through cmds to check whether the command already exists.
	for _, c := range cmds {
		if c.CmdHash == cmd.CmdHash {
			fmt.Fprintf(os.Stderr, "Command hash already exists: %s\n", cmd.CmdString)
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
		time.Now(),
		time.Now()}

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
	sc.StartTime = time.Now()
	sc.Modificationtime = time.Now()

	// Create a channel to hold exit status
	c := make(chan int)

	// Run the command in a goroutine
	go f(&sc, c)

	// Receive the exit status of the command
	status := <-c

	// Set end time after we receive from the channel
	sc.EndTime = time.Now()

	if status != 0 {
		log.Printf("Command status for %s: %d\n", sc.CmdHash, status)
	}

	return sc
}

// RunMockCommand runs a mock command
func RunMockCommand(sc *ScheduledCommand, c chan int) {
	time.Sleep(1 * time.Second)
	sc.ExitStatus = 99
	sc.Stdout = "Mock stdout message"
	sc.Stderr = "Mock stderr message"
	c <- sc.ExitStatus
}

// RunShellScriptCommand runs a command written to a temporary file
func RunShellScriptCommand(sc *ScheduledCommand, c chan int) {

	tempFile, err := ioutil.TempFile(os.TempDir(), "recmd-")

	if err != nil {
		fmt.Fprintf(os.Stdout, "Command status: %d\n", err)
	}

	defer os.Remove(tempFile.Name())

	//fmt.Fprintf(os.Stdout, "Created "+tempFile.Name()+"\n")

	_, err = tempFile.WriteString("#!/bin/sh\n\n" + sc.CmdString)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Errror while writing file: : %s\n", err)
	}

	out, err := exec.Command("sh", tempFile.Name()).Output()

	if err == nil {
		sc.ExitStatus = 0
	} else {
		sc.ExitStatus = -1

		if err.Error() != "" {
			sc.Stderr = err.Error()
		}
	}

	if out == nil {
		sc.Stdout = ""
	} else {
		sc.Stdout = string(out)
	}

	c <- sc.ExitStatus
}