package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/logger"
)

func TestLogger(t *testing.T) {
	logger.NewLogger(func(r *logger.Record) {
		fmt.Printf("%s", r)
	})

	logger.Debug("hello")
	logger.Warn("hello")
	logger.Error("hello")

	time.Sleep(1 * time.Second)
}
