package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"syscall"
)

var startChan = make(chan fsnotify.Event)
var stopChan = make(chan syscall.Signal)

func (w *Watch) task(event fsnotify.Event) {
	stopChan <- syscall.SIGINT
	startChan <- event
}

func (w *Watch) loop() {
	go func() {
		for {
			select {
			case event := <-startChan:
				log.Println("start process", event)
				w.startProcess()
				//log.Println(w.hasStartSuccess())
			case sig := <-stopChan:
				log.Println("stop process", sig)
				w.stopProcess()
			}
		}
	}()
}
