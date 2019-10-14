package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Event represents a log event
type Event struct {
	id      int
	message string
}

var (
	errorMessage = Event{0, "%s %s"}
	initMessage  = Event{1, "Initializing %s with %s"}
)

func init() {

	if os.Getenv("REPLICATOR_JSON") == "1" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		formatter := &logrus.TextFormatter{
			FullTimestamp: true,
		}
		logrus.SetFormatter(formatter)
	}

}

// Error is a function to pring all error related stuff
func Error(msg string, err error) {
	logrus.Errorf(errorMessage.message, msg, err)
}

// InitMessage logs initialization messages
func InitMessage(msg, with string) {
	logrus.Infof(initMessage.message, msg, with)
}

// Info is a generic function to log to stdout
func Info(args ...interface{}) {
	logrus.Info(args...)
}
