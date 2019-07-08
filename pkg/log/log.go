package log

import (
	logrus "github.com/sirupsen/logrus"
	"os"
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

func ErrorLog(msg string, err error) {
	logrus.Errorf("%s: %s", msg, err)
}

func Info(msg string) {
	logrus.Info(msg)
}
