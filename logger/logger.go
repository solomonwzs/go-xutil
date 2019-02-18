package logger

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/solomonwzs/goxutil/closer"
	"github.com/solomonwzs/goxutil/pubsub"
)

const (
	FINEST = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

type LogProcessor interface {
	L(r *Record) error
	Close() error
}

type MsgLogProcessor interface {
	L(r *Record) error
	M(msg interface{}) error
	Close() error
}

type Logger struct {
	lp  LogProcessor
	sub *pubsub.Subscriber
	closer.Closer
}

type MsgLogger struct {
	mlp MsgLogProcessor
	sub *pubsub.Subscriber
	ch  chan interface{}
	closer.Closer
}

var (
	logChannel *pubsub.Channel
	logPub     *pubsub.Publisher
	loggerN    int32
)

func (l *Logger) serv() {
	for {
		select {
		case <-l.Done():
			return
		case <-l.sub.Ready():
			if record, err := l.sub.NonBlockRecv(); err == nil {
				l.lp.L(record.(*Record))
			}
		}
	}
}

func NewLogger(lp LogProcessor) (*Logger, error) {
	sub := logChannel.NewSubscriber()
	if sub == nil {
		return nil, errors.New("new logger error")
	}

	atomic.AddInt32(&loggerN, 1)
	l := &Logger{
		lp:  lp,
		sub: sub,
		Closer: closer.NewCloser(func() error {
			atomic.AddInt32(&loggerN, -1)
			return lp.Close()
		}),
	}
	go l.serv()
	return l, nil
}

func (ml *MsgLogger) serv() {
	for {
		select {
		case <-ml.Done():
			return
		case msg := <-ml.ch:
			ml.mlp.M(msg)
		case <-ml.sub.Ready():
			if record, err := ml.sub.NonBlockRecv(); err == nil {
				ml.mlp.L(record.(*Record))
			}
		}
	}
}

func NewMsgLogger(mlp MsgLogProcessor) (*MsgLogger, error) {
	sub := logChannel.NewSubscriber()
	if sub == nil {
		return nil, errors.New("new logger error")
	}

	atomic.AddInt32(&loggerN, 1)
	ml := &MsgLogger{
		mlp: mlp,
		sub: sub,
		Closer: closer.NewCloser(func() error {
			atomic.AddInt32(&loggerN, -1)
			return mlp.Close()
		}),
	}
	go ml.serv()
	return ml, nil
}

func init() {
	logChannel = pubsub.NewChannel(32)
	logPub = logChannel.NewPublisher()
	loggerN = 0
}

func log(level int, depth int, msg string) {
	if loggerN > 0 {
		record := newRecord(level, depth, msg)
		logPub.Send(record, 0)
	}
}

func Finestf(format string, argv ...interface{}) {
	log(FINEST, 3, fmt.Sprintf(format, argv...))
}

func Finestln(argv ...interface{}) {
	log(FINEST, 3, fmt.Sprintln(argv...))
}

func Finest(msg string) {
	log(FINEST, 3, msg)
}

func Finef(format string, argv ...interface{}) {
	log(FINE, 3, fmt.Sprintf(format, argv...))
}

func Fineln(argv ...interface{}) {
	log(FINE, 3, fmt.Sprintln(argv...))
}

func Fine(msg string) {
	log(FINE, 3, msg)
}

func Debugf(format string, argv ...interface{}) {
	log(DEBUG, 3, fmt.Sprintf(format, argv...))
}

func Debugln(argv ...interface{}) {
	log(DEBUG, 3, fmt.Sprintln(argv...))
}

func Debug(msg string) {
	log(DEBUG, 3, msg)
}

func Tracef(format string, argv ...interface{}) {
	log(TRACE, 3, fmt.Sprintf(format, argv...))
}

func Traceln(argv ...interface{}) {
	log(TRACE, 3, fmt.Sprintln(argv...))
}

func Trace(msg string) {
	log(TRACE, 3, msg)
}

func Infof(format string, argv ...interface{}) {
	log(INFO, 3, fmt.Sprintf(format, argv...))
}

func Infoln(argv ...interface{}) {
	log(INFO, 3, fmt.Sprintln(argv...))
}

func Info(msg string) {
	log(INFO, 3, msg)
}

func Warnf(format string, argv ...interface{}) {
	log(WARNING, 3, fmt.Sprintf(format, argv...))
}

func Warnln(argv ...interface{}) {
	log(WARNING, 3, fmt.Sprintln(argv...))
}

func Warn(msg string) {
	log(WARNING, 3, msg)
}

func Errorf(format string, argv ...interface{}) {
	log(ERROR, 3, fmt.Sprintf(format, argv...))
}

func Errorln(argv ...interface{}) {
	log(ERROR, 3, fmt.Sprintln(argv...))
}

func Error(msg string) {
	log(ERROR, 3, msg)
}

func Critcalf(format string, argv ...interface{}) {
	log(CRITICAL, 3, fmt.Sprintf(format, argv...))
}

func Critcalln(argv ...interface{}) {
	log(CRITICAL, 3, fmt.Sprintln(argv...))
}

func Critcal(msg string) {
	log(CRITICAL, 3, msg)
}

func DPrintf(format string, argv ...interface{}) {
	msg := fmt.Sprintf(format, argv...)
	dPrint(msg)
}

func DPrintln(argv ...interface{}) {
	msg := fmt.Sprintln(argv...)
	dPrint(msg)
}

func dPrint(msg string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	fmt.Printf("%s=%d= %s [%s:%d] \033[0m%s",
		_COLOR_GREEN, os.Getpid(), time.Now().Format("15:04:05"),
		ShortFilename(file), line, msg)
}
