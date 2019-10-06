package log

import (
	logrus "github.com/sirupsen/logrus"
	"os"
)

// Event struct to represent log message
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

// Error is a function print error values
func Error(msg string, err error) {
	logrus.Errorf(errorMessage.message, msg, err)
}

// Fatal is a function that kills the program and logs ins case of an error
func Fatal(err error) {
	logrus.Fatal(err)
}

// InitMessage is a function to print initialization messages
func InitMessage(msg, with string) {
	logrus.Infof(initMessage.message, msg, with)
}

// Info only shows generic info
func Info(args ...interface{}) {
	logrus.Info(args...)
}
