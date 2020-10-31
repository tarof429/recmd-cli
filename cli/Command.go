package cli

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/theckman/yacspin"
)

type CommandStatus string

const (
	Idle      CommandStatus = "Idle"
	Running   CommandStatus = "Running"
	Completed CommandStatus = "Completed"
)

// Command represents a command and optionally a description to document what the command does
type Command struct {
	CmdHash     string        `json:"commandHash"`
	CmdString   string        `json:"commandString"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Status      CommandStatus `json:"status"`
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

const (
	// D(rectory containing configuration and command history
	recmdDir = ".recmd"

	// The secret file
	recmdSecretFile = "recmd_secret"

	// The command history file
	recmdHistoryFile = "recmd_history.json"

	// Length of scret string
	secretLength = 40
)

// Global variables
var (
	recmdSecretFilePath string
	secretData          string
)

// InitTool initializes the tool
func InitTool() {

	// Create ~/.recmd if it doesn't exist
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Error, unable to obtain home directory path %v\n", err)
	}

	recmdDirPath := filepath.Join(homeDir, recmdDir)

	fileInfo, statErr := os.Stat(recmdDirPath)

	if os.IsNotExist((statErr)) {
		if err != nil {
			log.Fatalf("Error, please start recmd-dmn first: %v\n", err)
		}
	} else if !fileInfo.IsDir() {
		log.Fatalf("Error, ~/.recmd is not a directory")
	}

	recmdSecretFilePath = filepath.Join(recmdDirPath, recmdSecretFile)

}

// GetSecret gets the secret from the file system
func GetSecret() string {
	secretData, err := ioutil.ReadFile(recmdSecretFilePath)

	if err != nil {
		log.Fatalf("Error, unable to read secret from file %v\n", err)
	}

	if len(secretData) != secretLength {
		log.Fatalf("Error, invalid secret length %v\n", err)
	}

	return string(secretData)
}

func getBase64(line string) string {

	lineData := []byte(line)
	return base64.StdEncoding.EncodeToString(lineData)
}

// List lists the commands
func List() ([]Command, error) {
	var (
		historyData []byte    // Data representing our history file
		cmds        []Command // List of commands produced after unmarshalling historyData
		err         error     // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())

	url := "http://localhost:8999/secret/" + encodedSecret + "/list"

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmds)

	return cmds, err

}

// Status gets the list of commands that are in the queue
func Status() ([]Command, error) {
	var (
		historyData []byte    // Data representing our history file
		cmds        []Command // List of commands in the queue
		err         error     // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())

	url := "http://localhost:8999/secret/" + encodedSecret + "/status"

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmds)

	return cmds, err

}

// SelectCmd returns a command
func SelectCmd(value string) (Command, error) {

	var (
		historyData []byte  // Data representing our history file
		cmd         Command // List of commands produced after unmarshalling historyData
		err         error   // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())
	encodedCommandHash := getBase64(value)

	url := "http://localhost:8999/secret/" + encodedSecret + "/select/cmdHash/" + encodedCommandHash

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmd)

	return cmd, err
}

// SearchCmd returns a command by name
func SearchCmd(value string) ([]Command, error) {

	var (
		historyData []byte    // Data representing our history file
		cmds        []Command // List of commands produced after unmarshalling historyData
		err         error     // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())
	encodedDescription := getBase64(value)

	url := "http://localhost:8999/secret/" + encodedSecret + "/search/description/" + encodedDescription

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmds)

	return cmds, err
}

// DeleteCmd deletes a command. It's best to pass in the commandHash
// because commands may look similar.
func DeleteCmd(value string) ([]Command, error) {

	var (
		historyData []byte    // Data representing our history file
		cmds        []Command // List of commands produced after unmarshalling historyData
		err         error     // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())
	encodedCommandHash := getBase64(value)

	url := "http://localhost:8999/secret/" + encodedSecret + "/delete/cmdHash/" + encodedCommandHash

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(historyData, &cmds)

	return cmds, err
}

// RunCmd run s a command
func RunCmd(value string, background bool) ScheduledCommand {

	var (
		historyData []byte           // Data representing our history file
		cmd         ScheduledCommand // List of commands produced after unmarshalling historyData
		err         error            // Any errors we might encounter
	)

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[14],
		Suffix:          " Scheduling commmand ",
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	encodedSecret := getBase64(GetSecret())
	encodedCommandHash := getBase64(value)

	url := "http://localhost:8999/secret/" + encodedSecret + "/run/cmdHash/" + encodedCommandHash

	spinner, _ := yacspin.New(cfg)

	if background == true {

		spinner.Start()

		go http.Get(url)

		time.Sleep(1 * time.Second)

		spinner.Stop()

		return cmd
	}

	spinner.Start()

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmd)

	spinner.Stop()

	return cmd
}

// AddCmd adds a command.
func AddCmd(command string, description string, workingDirectory string) string {

	var (
		data   []byte // Data representing status
		status string // Status
		err    error  // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())
	encodedCommand := getBase64(command)
	encodedDescription := getBase64(description)
	encodedWorkingDirectory := getBase64(workingDirectory)

	url := "http://localhost:8999/secret/" + encodedSecret + "/add/command/" + encodedCommand + "/description/" + encodedDescription + "/workingDirectory/" + encodedWorkingDirectory

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	data, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(data, &status)

	return status
}
