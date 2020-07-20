package recmd

import (
	"fmt"
	"testing"
	"time"
)

func TestCurrentTime(t *testing.T) {
	currentTime := time.Now()
	fmt.Println("######################################")
	fmt.Println(currentTime.Format("2006-01-02 15:04:05"))
	fmt.Println(currentTime.Format("Mon 01 02 2006 15:04:05"))
}
