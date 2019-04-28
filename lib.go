package main

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

func (w *Watch) getConfig() {

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

func (w *Watch) matchPath(pathName string, source []string) bool {
	for _, v := range source {
		if strings.HasPrefix(pathName, v) {
			return true
		}
	}
	return false
}

func (w *Watch) setCache(pathName string, content string) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.cache[pathName] = content
}

func (w *Watch) getCache(pathName string) string {
	w.mux.RLock()
	defer w.mux.RUnlock()
	if content, ok := w.cache[pathName]; ok {
		return content
	}
	return ""
}

func (w *Watch) getContent(pathName string) (string, error) {

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

func (w *Watch) isUpdate(pathName string) bool {

	// Just wait file release lock from fsnotify
	time.Sleep(100 * time.Millisecond)

	content, err := w.getContent(pathName)
	if err != nil {
		return false
	}

	var cache = w.getCache(pathName)

	w.setCache(pathName, content)

	return cache != content
}

func (w *Watch) watchPathExceptIgnore() {
	filepath.Walk(w.listenPath, func(pathName string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			return err
		}

		if w.matchPath(pathName, w.config.ignore.paths) {
			return err
		}

		err = w.watch.Add(pathName)
		if err != nil {
			return err
		}

		log.Println("watch", pathName)

		return err
	})
}
