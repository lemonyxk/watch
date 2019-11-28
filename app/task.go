package app

import (
	"syscall"

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
				w.startProcess()
			case sig := <-stopChan:
				color.Bold.Println("stop process", sig)
				w.StopProcess()
			}
		}
	}()
}
