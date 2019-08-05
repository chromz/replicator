package log

import (
	logrus "github.com/sirupsen/logrus"
	"os"
)

type LogEvent struct {
	id int
	message string
}

var (
	errorMessage = LogEvent{0, "%s %s"}
	initMessage = LogEvent{1, "Initializing %s with %s"}
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

func Error(msg string, err error) {
	logrus.Errorf(errorMessage.message, msg, err)
}

func InitMessage(msg, with string) {
	logrus.Infof(initMessage.message, msg, with)
}

func Info(msg string) {
	logrus.Info(msg)
}

