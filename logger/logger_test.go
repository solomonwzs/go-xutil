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

type FLP struct {
	F *RotateFile
}

func (lp FLP) L(r *Record) error { fmt.Fprintf(lp.F, "%s", r); return nil }
func (lp FLP) Close() error      { return lp.F.Close() }

func TestLogger1(t *testing.T) {
	rf, err := NewRotateFile(RotateFileOptions{
		Name:    "xx",
		Path:    "/tmp",
		MaxSize: 130,
	})
	if err != nil {
		t.Fatal(err)
	}

	l, err := NewLogger(FLP{
		F: rf,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	Debugln("hello")
	Warnln("hello")
	Errorln("hello")

	time.Sleep(1 * time.Second)
}
