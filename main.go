package recmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Command represents a command and optionally a comment to document what the command does
type Command struct {
	CmdHash   string `json:"commandHash"`
	CmdString string `json:"commandString"`
	Comment   string `json:"comment"`
}

// // Config represents global configuration
// type Config struct {
// 	Salt string `json:"salt"`
// }

const historyFile = "cmd_history.json"

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
func ReadCmdHistoryFile() []Command {

	var cmds []Command

	data, err := ioutil.ReadFile(historyFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred whiel reading historyfile: %v\n", err)
		return cmds
	}

	if err := json.Unmarshal(data, &cmds); err != nil {
		fmt.Fprintf(os.Stderr, "JSON unmarshalling failed: %s\n", err)
	}

	return cmds
}

// WriteCmdHistoryFile writes a command to the history file
func WriteCmdHistoryFile(cmd Command) bool {

	// Check if the file does not exist. If not, then create it and add our first command to it.
	f, err := os.Open(historyFile)

	// Immediately close the file since we plan to write to it
	f.Close()

	// Check if the file doesn't exist and if so, then write it.
	if err != nil {
		// The array of commands
		var cmds []Command

		cmds = append(cmds, cmd)

		mode := int(0644)

		updatedData, _ := json.MarshalIndent(cmds, "", "\t")

		error := ioutil.WriteFile(historyFile, updatedData, os.FileMode(mode))

		return error == nil
	}

	// Add the command to the history file

	// The array of commands
	var cmds []Command

	// Read history file
	data, err := ioutil.ReadFile(historyFile)

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

	error := ioutil.WriteFile(historyFile, updatedData, os.FileMode(mode))

	return error == nil

}

// NewCommand creates a new Command struct and populates the fields
func NewCommand(cmdString string, cmdComment string) Command {

	formattedHash := func() string {
		h := sha1.New()
		h.Write([]byte(cmdString))
		return fmt.Sprintf("%.15x", h.Sum(nil))
	}()

	cmd := Command{formattedHash, cmdString, cmdComment}

	return cmd
}
