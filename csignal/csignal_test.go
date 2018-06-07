package csignal

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestCsignal(t *testing.T) {
	ch := make(chan os.Signal)
	Notify(ch, syscall.SIGINT)

	sig := <-ch
	fmt.Println(">>", sig)
}
