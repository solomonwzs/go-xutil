package csignal

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func infinite(x int) int {
	return infinite(x) + 1
}

func TestCsignal(t *testing.T) {
	ch := make(chan os.Signal)
	Notify(ch, syscall.SIGINT, syscall.SIGSEGV)

	go infinite(0)

	sig := <-ch
	fmt.Println(">>", sig)
}
