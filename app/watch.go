package app

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"

	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const Interval = 500 * time.Millisecond

type Watch struct {
	watch      *fsnotify.Watcher
	listenPath string
	config     Config
	cache      map[string]string
	task       []string
	mux        sync.RWMutex
	isInterval bool
	commands   []*exec.Cmd
	isRun      bool
}

type Config struct {
	ignore Ignore
	start  []string
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

	w.RunTask()

	w.WatchPathExceptIgnore()

	w.Listen()

	w.Loop()

	w.Block()

	defer func() { _ = w.watch.Close() }()
}

func (w *Watch) RunTask() {
	time.AfterFunc(Interval, func() {
		w.DelayTask()
		w.OnInterval()
		startChan <- fsnotify.Event{Name: "init"}
	})
}

func (w *Watch) Block() {
	// 创建信号
	signalChan := make(chan os.Signal, 1)
	// 通知
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	// 阻塞
	sign := <-signalChan

	for {
		color.Bold.Println("waiting close...")
		w.StopProcess()
		time.Sleep(100 * time.Millisecond)
		break
	}

	color.Bold.Println("close success", sign)
}

func (w *Watch) OnInterval() {

	w.isInterval = true

	time.AfterFunc(Interval, func() {
		w.isInterval = false
	})
}

func (w *Watch) Listen() {
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:

				if w.isInterval {
					break
				}

				w.OnInterval()

				// 排除 IGNORE 文件
				if w.MatchFile(ev.Name) {
					break
				}

				if w.MatchOthers(ev.Name) {
					break
				}

				if ev.Op&fsnotify.Create == fsnotify.Create {
					color.Bold.Println("create", ev.Name)
					// 这里获取新创建文件的信息，如果是目录，则加入监控中
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Add(ev.Name)
						color.Bold.Println("add watch", ev.Name)
					}
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					color.Bold.Println("delete", ev.Name)
					// 如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Remove(ev.Name)
						color.Bold.Println("delete watch", ev.Name)
					}
				}

				// 重命名文件 删除监听
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					color.Bold.Println("rename", ev.Name)
					color.Bold.Println("delete watch", ev.Name)
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
				color.Red.Println("error", err)
			}
		}
	}()
}
