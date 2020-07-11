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
	// 创建信号
	signalChan := make(chan os.Signal, 1)
	// 通知
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞
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

				// 排除 IGNORE 文件
				// 用于新添加的文件
				if w.MatchFile(ev.Name) {
					break
				}

				// 排除目录
				if w.MatchPath(ev.Name) {
					break
				}

				// 排除正则
				if w.MatchOthers(ev.Name) {
					break
				}

				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("create", ev.Name)
					// 这里获取新创建文件的信息，如果是目录，则加入监控中
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Add(ev.Name)
						fmt.Println("add watch", ev.Name)
					}
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("delete", ev.Name)
					// 如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Remove(ev.Name)
						fmt.Println("delete watch", ev.Name)
					}
				}

				// 重命名文件 删除监听
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("rename", ev.Name)
					fmt.Println("delete watch", ev.Name)
					// 获取不到旧文件的资料 直接移除
					_ = w.watch.Remove(ev.Name)
				}

				// 修改权限
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					if w.IsUpdate(ev.Name) {
						w.Task(ev)
					}
				}

				// 写入文件
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
