package app

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"syscall"
)

var startChan = make(chan fsnotify.Event)
var stopChan = make(chan syscall.Signal)

func (w *Watch) Task(event fsnotify.Event) {
	stopChan <- syscall.SIGINT
	startChan <- event
}

func (w *Watch) Loop() {
	go func() {
		for {
			select {
			case event := <-startChan:
				log.Println("start process", event)
				w.startProcess()
				//log.Println(w.HasStartSuccess())
			case sig := <-stopChan:
				log.Println("stop process", sig)
				w.StopProcess()
			}
		}
	}()
}
