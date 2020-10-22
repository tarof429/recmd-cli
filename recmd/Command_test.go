package recmd

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

const testdataDir = "testdata"
const testHistoryFile = testdataDir + "/.cmd_history.json"

func TestMain(m *testing.M) {
	fmt.Println("Running tests...")

	status := m.Run()

	os.Exit(status)
}

func TestReadCmdHistoryFile(t *testing.T) {

	secret := "MApoj3G3bJnFJ9Is6ahMkFo99pfqLlA"
	command := "ls /foo/bar"

	//secret/{secret}/add/command/{command}/description/{description}"

	baseURL, err := url.Parse("http://locahost:8999")
	baseURL.Path += "/secret/" + secret
	baseURL.Path += "/add"
	baseURL.Path += "/command/" + command
	///secret/{secret}/delete/cmdHash/{cmdHash}")

	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
	}

	fmt.Println(baseURL)
	// cmd := Command{"abc", "ls", "list files", -1}
	// cmd2 := Command{"def", "df", "Show disk usage", -1}

	// cmds := []Command{cmd, cmd2}

	// // Convert the struct to JSON
	// data, err := json.MarshalIndent(cmds, "", "\t")

	// if err != nil {
	// 	t.Error(err)
	// }

	// mode := int(0644)

	// ioutil.WriteFile(testHistoryFile, data, os.FileMode(mode))

	// readCmds, err := ReadCmdHistoryFile(testdataDir)

	// if err != nil {
	// 	t.Error("an error occured while reading command history")
	// }

	// // Check to make sure that the commands we wrote are the ones we created by comparing hashes
	// for i := 0; i < len(cmds); i++ {
	// 	if cmds[i].CmdHash != readCmds[i].CmdHash {
	// 		t.Error("The command histories are not equal")
	// 	}

	// }

	// t.Cleanup(func() {
	// 	os.Remove(testHistoryFile)
	// })

}
