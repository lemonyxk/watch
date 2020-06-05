package app

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

const Interval = 300 * time.Millisecond

var startChan = make(chan fsnotify.Event)

func (w *Watch) Task(event fsnotify.Event) {
	startChan <- event
}

func (w *Watch) Loop() {
	var intervalTimer *time.Timer
	var running = false
	go func() {
		for {
			select {
			case event := <-startChan:

				if running {
					continue
				}

				if intervalTimer != nil {
					intervalTimer.Stop()
				}

				intervalTimer = time.AfterFunc(Interval, func() {
					running = true
					w.StopProcess()
					w.startProcess(event)
					running = false
				})

			}
		}
	}()
}
