package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	_COLOR_RED        = "\033[31m"
	_COLOR_GREEN      = "\033[32m"
	_COLOR_YELLOW     = "\033[33m"
	_COLOR_BLUE       = "\033[34m"
	_COLOR_PURPLE     = "\033[35m"
	_COLOR_LIGHT_BLUE = "\033[36m"
	_COLOR_GRAY       = "\033[37m"
	_COLOR_BLACK      = "\033[30m"
)

var (
	_SHORT_LEVEL_NAME map[int]string = map[int]string{
		FINEST:   "N",
		FINE:     "F",
		DEBUG:    "D",
		TRACE:    "T",
		INFO:     "I",
		WARNING:  "W",
		ERROR:    "E",
		CRITICAL: "C",
	}

	_LEVEL_COLOR map[int]string = map[int]string{
		FINEST:   _COLOR_BLACK,
		FINE:     _COLOR_BLUE,
		DEBUG:    _COLOR_GREEN,
		TRACE:    _COLOR_LIGHT_BLUE,
		INFO:     _COLOR_GRAY,
		WARNING:  _COLOR_YELLOW,
		ERROR:    _COLOR_RED,
		CRITICAL: _COLOR_PURPLE,
	}
)

type Record struct {
	File    string
	Line    int
	Message string
	Created time.Time
	Level   int
}

func newRecord(level int, calldepth int, msg string) *Record {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	return &Record{
		File:    file,
		Line:    line,
		Message: msg,
		Created: time.Now(),
		Level:   level,
	}
}

func (r *Record) String() string {
	return fmt.Sprintf("%s%d %s [%s %s:%d] \033[0m%s",
		_LEVEL_COLOR[r.Level], os.Getpid(),
		r.Created.Format("2006-01-02 15:04:05"), _SHORT_LEVEL_NAME[r.Level],
		ShortFilename(r.File), r.Line, r.Message)
}

func ShortFilename(filename string) string {
	for i := len(filename) - 1; i > 0; i-- {
		if filename[i] == '/' {
			return filename[i+1:]
		}
	}
	return filename
}
