package watcher

import (
	"time"
)

func Start(pollingRate int64) {
	for {
		time.Sleep(time.Duration(pollingRate) * time.Millisecond)
	}
}

