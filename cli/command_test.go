package cli

import (
	"fmt"
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

func TestGenerateDummySecret(t *testing.T) {

	secret := GenerateDummySecret()
	fmt.Println(secret)
}
