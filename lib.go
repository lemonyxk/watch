package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (w *Watch) CreateListenPath(pathName string) {
	w.listenPath = pathName
}

func (w *Watch) createWatch() {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	w.watch = watch
}

func (w *Watch) GetConfig() {

	var watchPathConfig = path.Join(w.listenPath, "watch")

	file, err := os.OpenFile(watchPathConfig, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(watchPathConfig, "is not found")
		os.Exit(0)
	}

	defer file.Close()

	var reader = bufio.NewReader(file)

	var key = ""

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		var rule = strings.Trim(string(line), " ")

		if rule == "" {
			continue
		}

		if strings.HasPrefix(rule, "#") {
			continue
		}

		if strings.HasPrefix(rule, "[") && strings.HasSuffix(rule, "]") {
			key = rule[1 : len(rule)-1]
			continue
		}

		switch key {
		// get path
		case "ignore":
			var asbPath = path.Join(w.listenPath, rule)
			info, err := os.Stat(asbPath)
			if err != nil {
				continue
			}
			if info.IsDir() {
				w.config.ignore.paths = append(w.config.ignore.paths, asbPath)
			} else {
				w.config.ignore.files = append(w.config.ignore.files, asbPath)
			}
		case "start":
			w.config.start = append(w.config.start, rule)
		}

	}
}

func (w *Watch) MatchPath(pathName string, source []string) bool {
	for _, v := range source {
		if strings.HasPrefix(pathName, v) {
			return true
		}
	}
	return false
}

func (w *Watch) SetCache(pathName string, size int) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.cache[pathName] = size
}

func (w *Watch) GetCache(pathName string) int {
	w.mux.RLock()
	defer w.mux.RUnlock()
	if content, ok := w.cache[pathName]; ok {
		return content
	}
	return 0
}

func (w *Watch) GetSize(pathName string) (int, error) {
	info, err := os.Stat(pathName)
	if err != nil {
		return 0, err
	}
	return int(info.Size()), err
}

func (w *Watch) GetContent(pathName string) (string, error) {

	bytes, err := ioutil.ReadFile(pathName)
	if err != nil {
		return "", err
	}

	var content = string(bytes)
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "")

	return content, nil
}

func (w *Watch) IsUpdate(pathName string) bool {

	// Just wait file release lock from fsnotify
	time.Sleep(200 * time.Millisecond)

	size, err := w.GetSize(pathName)
	if err != nil {
		return false
	}

	var cache = w.GetCache(pathName)

	w.SetCache(pathName, size)

	return cache != size
}

func (w *Watch) WatchPathExceptIgnore() {
	filepath.Walk(w.listenPath, func(pathName string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			return err
		}

		if w.MatchPath(pathName, w.config.ignore.paths) {
			return err
		}

		err = w.watch.Add(pathName)
		if err != nil {
			return err
		}

		w.AddTask(pathName)

		log.Println("watch dir", pathName)

		return err
	})
}

func (w *Watch) AddTask(pathName string) {
	if pathName == w.listenPath {
		return
	}
	w.task = append(w.task, pathName)
}

func (w *Watch) DelayTask() {
	for _, dir := range w.task {
		filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {

			if i.IsDir() {
				return err
			}

			// log.Println("watch file", p)

			size, err := w.GetSize(p)

			if err != nil {
				return err
			}

			w.SetCache(p, size)

			return err
		})
	}

	log.Println(w.cache)
}
