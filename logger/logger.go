package logger

import (
	"errors"
	"sync/atomic"

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

type LogFunc func(*Record)

type Logger struct {
	f   LogFunc
	sub *pubsub.Subscriber
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
			if record, err := l.sub.NonBlockRecv(); err != nil {
				return
			} else {
				l.f(record.(*Record))
			}
		}
	}
}

func NewLogger(f LogFunc) (*Logger, error) {
	sub := logChannel.NewSubscriber()
	if sub == nil {
		return nil, errors.New("new logger error")
	}

	atomic.AddInt32(&loggerN, 1)
	l := &Logger{
		f:   f,
		sub: sub,
		Closer: closer.NewCloser(func() error {
			atomic.AddInt32(&loggerN, -1)
			return nil
		}),
	}
	go l.serv()
	return l, nil
}

func init() {
	logChannel = pubsub.NewChannel(32)
	logPub = logChannel.NewPublisher()
	loggerN = 0
}

func log(level int, depth int, format string, argv ...interface{}) {
	if loggerN > 0 {
		record := newRecord(level, depth, format, argv...)
		logPub.Send(record, 0)
	}
}

func Finestf(format string, argv ...interface{}) {
	log(FINEST, 3, format, argv...)
}

func Finest(argv ...interface{}) {
	log(FINEST, 3, "", argv...)
}

func Finef(format string, argv ...interface{}) {
	log(FINE, 3, format, argv...)
}

func Fine(argv ...interface{}) {
	log(FINE, 3, "", argv...)
}

func Debugf(format string, argv ...interface{}) {
	log(DEBUG, 3, format, argv...)
}

func Debug(argv ...interface{}) {
	log(DEBUG, 3, "", argv...)
}

func Tracef(format string, argv ...interface{}) {
	log(TRACE, 3, format, argv...)
}

func Trace(argv ...interface{}) {
	log(TRACE, 3, "", argv...)
}

func Infof(format string, argv ...interface{}) {
	log(INFO, 3, format, argv...)
}

func Info(argv ...interface{}) {
	log(INFO, 3, "", argv...)
}

func Warnf(format string, argv ...interface{}) {
	log(WARNING, 3, format, argv...)
}

func Warn(argv ...interface{}) {
	log(WARNING, 3, "", argv...)
}

func Errorf(format string, argv ...interface{}) {
	log(ERROR, 3, format, argv...)
}

func Error(argv ...interface{}) {
	log(ERROR, 3, "", argv...)
}

func Critcalf(format string, argv ...interface{}) {
	log(CRITICAL, 3, format, argv...)
}

func Critcal(argv ...interface{}) {
	log(CRITICAL, 3, "", argv...)
}
