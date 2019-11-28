package app

import (
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"

	"github.com/Lemo-yxk/go-watch/vars"
)

var startChan = make(chan fsnotify.Event)
var stopChan = make(chan syscall.Signal)

func (w *Watch) Task(event fsnotify.Event) {
	stopChan <- syscall.Signal(vars.Sig)
	startChan <- event
}

func (w *Watch) Loop() {
	go func() {
		for {
			select {
			case event := <-startChan:
				color.Bold.Println("start process", event)
				var start = time.Now()
				w.startProcess()
				color.Bold.Println("time", float64(time.Now().Sub(start).Milliseconds())/1000)
			case sig := <-stopChan:
				color.Bold.Println("stop process", sig)
				w.StopProcess()
			}
		}
	}()
}
