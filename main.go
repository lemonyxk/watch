package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const interval = time.Second

type Watch struct {
	watch      *fsnotify.Watcher
	listenPath string
	config     Config
	cache      map[string]string
	mux        sync.RWMutex
	isInterval bool
	cmd        *exec.Cmd
	isRun      bool
	pid        int
}

type Config struct {
	ignore Ignore
	start  []string
}

type Ignore struct {
	paths []string
	files []string
}

func init() {
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)

	var listenPath string
	flag.StringVar(&listenPath, "path", ".", "path")
	flag.Parse()

	log.Println(listenPath)
}

func main() {

	var pathName = "/Users/lemo/test/bairenniuniu/"

	var watch = &Watch{}

	watch.CreateListenPath(pathName)

	watch.Run()

}

func (w *Watch) Run() {

	go func() {
		w.onInterval()
		startChan <- fsnotify.Event{Name: "init"}
	}()

	w.cache = make(map[string]string)

	w.createWatch()

	w.getConfig()

	log.Println(w.config.ignore.files)

	w.watchPathExceptIgnore()

	w.listen()

	w.loop()

	w.block()
}

func (w *Watch) block() {
	// 创建信号
	signalChan := make(chan os.Signal, 1)
	// 通知
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞
	sign := <-signalChan

	for {
		log.Println("waiting close...")
		w.stopProcess()
		time.Sleep(time.Second)
		break
	}

	w.watch.Close()

	log.Println("close success", sign)
}

func (w *Watch) onInterval() {

	w.isInterval = true

	time.AfterFunc(interval, func() {
		w.isInterval = false
	})
}

func (w *Watch) listen() {
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:

				if w.isInterval {
					break
				}

				w.onInterval()

				// 排除 IGNORE 文件
				var match = false
				for _, f := range w.config.ignore.files {
					if strings.HasPrefix(ev.Name, f) {
						log.Println("ignore files", ev.Name)
						match = true
					}
				}
				if match {
					break
				}

				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("create", ev.Name)
					// 这里获取新创建文件的信息，如果是目录，则加入监控中
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						w.watch.Add(ev.Name)
						log.Println("add watch", ev.Name)
					}
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("delete", ev.Name)
					// 如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						w.watch.Remove(ev.Name)
						log.Println("delete watch", ev.Name)
					}
				}

				// 重命名文件 删除监听
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("rename", ev.Name)
					log.Println("delete watch", ev.Name)
					// 获取不到旧文件的资料 直接移除
					w.watch.Remove(ev.Name)
				}

				// 修改权限
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					if w.isUpdate(ev.Name) {
						w.task(ev)
					}
				}

				// 写入文件
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if w.isUpdate(ev.Name) {
						w.task(ev)
					}
				}

			case err := <-w.watch.Errors:
				log.Println("error", err)
			}
		}
	}()
}
