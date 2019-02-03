package logger

import (
	"fmt"
	"testing"
	"time"
)

type LP struct{}

func (lp LP) L(r *Record) error { fmt.Printf("%s", r); return nil }
func (lp LP) Close() error      { return nil }

func TestLogger(t *testing.T) {
	l, err := NewLogger(LP{})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	Debugln("hello")
	Warnln("hello")
	Errorln("hello")

	time.Sleep(1 * time.Second)
}
