package cli

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/theckman/yacspin"
)

type CommandStatus string

const (
	Idle      CommandStatus = "Idle"
	Running   CommandStatus = "Running"
	Completed CommandStatus = "Completed"

	recmdDmn string = "recmd-dmn" // name of the command we want to run
	recmdPid string = "recmd-dmn.pid"
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
	// Directory containing configuration and command history
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

	wd, _ := os.Getwd()

	//parentDir := filepath.Dir(wd)

	confDirPath := filepath.Join(wd, "conf")

	recmdSecretFilePath = filepath.Join(confDirPath, recmdSecretFile)

	if _, err := os.Stat(recmdSecretFilePath); err != nil {
		ioutil.WriteFile(recmdSecretFilePath, GenerateDummySecret(), os.FileMode(0644))
	}

}

func GenerateDummySecret() []byte {
	ret := make([]byte, 40)

	return ret
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
		return cmds, err
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmds)

	return cmds, err

}

// Queue gets the list of commands that are in the queue
func Queue() ([]Command, error) {
	var (
		historyData []byte    // Data representing our history file
		cmds        []Command // List of commands in the queue
		err         error     // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())

	url := "http://localhost:8999/secret/" + encodedSecret + "/queue"

	resp, err := http.Get(url)

	if err != nil {
		return cmds, err
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
		return cmd, err
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
		return cmds, err
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
		return cmds, err
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(historyData, &cmds)

	return cmds, err
}

// StatusCmd gets the status of the daemon.
func StatusCmd() (bool, error) {

	var (
		statusData []byte // Data representing our history file
		status     bool   // List of commands produced after unmarshalling historyData
		err        error  // Any errors we might encounter
	)

	encodedSecret := getBase64(GetSecret())

	url := "http://localhost:8999/secret/" + encodedSecret + "/status"

	resp, err := http.Get(url)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	statusData, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(statusData, &status)

	return status, err
}

// RunCmd run s a command
func RunCmd(value string, background bool) (ScheduledCommand, error) {

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

		return cmd, nil
	}

	spinner.Start()

	resp, err := http.Get(url)

	//log.Println(resp)

	if err != nil {
		return cmd, err
	}

	defer resp.Body.Close()

	historyData, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(historyData, &cmd)

	spinner.Stop()

	return cmd, nil
}

// AddCmd adds a command.
func AddCmd(command string, description string, workingDirectory string) (string, error) {

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
		return status, err
	}

	defer resp.Body.Close()

	data, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(data, &status)

	return status, nil
}

// StartCmd starts the daemon. In order for this to work, the daemon
// must be colocated with the CLI. For some reason, if recmd-dmn is started twice,
// it cannot serve requests. So there is logic in this code to 1) check if a PID
// file exists in the current directory 2) If it exists, kill the process, delete the
// PID file 3) Start recmd-dmn
func StartCmd() error {

	log.Println("Starting command")

	dir, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	// binDir := filepath.Join(dir, "bin")
	//log.Printf("command.go: Starting recmd in %v\n", dir)

	pidFilePath := filepath.Join(dir, recmdPid)

	cmd := exec.Command("bin/" + recmdDmn)
	cmd.Dir = dir
	err = cmd.Start()
	//fmt.Println(cmd.Process.Pid)
	pid := cmd.Process.Pid
	mode := int(0644)
	data := []byte(strconv.Itoa(pid))
	ioutil.WriteFile(pidFilePath, data, os.FileMode(mode))

	// Check if the process is available
	i := 0
	for i < 3 {
		time.Sleep(time.Second)
		_, err = List()

		if err == nil {
			break
		}
		i = i + 1
	}

	return err
}

// StopCmd stops the daemon by finding the process indicated by the PID file
func StopCmd() error {

	log.Println("Stopping command")

	dir, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	pidFilePath := filepath.Join(dir, recmdPid)

	if _, err := os.Stat(pidFilePath); err != nil {
		return nil
	}

	currentPid, err := ioutil.ReadFile(pidFilePath)

	if err != nil {
		log.Println("Found a PID file but unable to read contents")
	}

	currentPidAsInt, _ := strconv.Atoi(string(currentPid))

	//log.Printf("Found PID: %v\n", currentPidAsInt)

	p, err := os.FindProcess(currentPidAsInt)

	if err == nil {
		//log.Println("Stopping process")
		err = p.Signal(os.Interrupt)
		p.Wait()
		// _, err := p.Wait()
		time.Sleep(time.Second)
		// if err != nil {
		// 	log.Println("Stopped")
		// }
	}

	return nil
}

// StopThenStartCmd stops the dameon before starting it. The intention
// is to provide a safe way to start the dameon.
func StopThenStartCmd() error {

	func() error {
		err := StopCmd()
		if err != nil {
			return err
		}
		return nil
	}()

	func() error {

		err := StartCmd()
		if err != nil {
			return err
		}
		return nil
	}()
	return nil
}
