package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type CmdInfo struct {
	cmd    *exec.Cmd
	status bool
}

type Watch struct {
	watch      *fsnotify.Watcher
	listenPath string
	config     Config
	cache      map[string]string
	task       []string
	mux        sync.RWMutex
	commands   []*CmdInfo
	isRun      bool
}

type Config struct {
	ignore  Ignore
	command []string
}

type Ignore struct {
	paths  []string
	files  []string
	others []string
}

func (w *Watch) Run() {

	w.cache = make(map[string]string)

	w.CreateWatch()

	w.GetConfig()

	w.WatchPathExceptIgnore()

	w.Listen()

	w.Loop()

	w.RunTask()

	w.Block()

	defer func() { _ = w.watch.Close() }()
}

func (w *Watch) RunTask() {
	w.DelayTask()
	startChan <- fsnotify.Event{Name: "init"}
}

func (w *Watch) Block() {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sign := <-signalChan

	fmt.Println("waiting close...")

	w.StopProcess()

	fmt.Println("close success", sign)
}

func (w *Watch) Listen() {
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:

				// filter match file
				if w.MatchFile(ev.Name) {
					break
				}

				// filter match path
				if w.MatchPath(ev.Name) {
					break
				}

				// filter regex
				if w.MatchOthers(ev.Name) {
					break
				}

				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("create", ev.Name)
					// if is dir, add watch
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Add(ev.Name)
						fmt.Println("add watch", ev.Name)
					}
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("delete", ev.Name)
					// if delete file is dir, remove watch
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Remove(ev.Name)
						fmt.Println("delete watch", ev.Name)
					}
				}

				// rename event
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("rename", ev.Name)
					fmt.Println("delete watch", ev.Name)
					// remove old watch
					_ = w.watch.Remove(ev.Name)
				}

				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					if w.IsUpdate(ev.Name) {
						w.Task(ev)
					}
				}

				// write event
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if w.IsUpdate(ev.Name) {
						w.Task(ev)
					}
				}

			case err := <-w.watch.Errors:
				fmt.Println("error", err)
			}
		}
	}()
}
