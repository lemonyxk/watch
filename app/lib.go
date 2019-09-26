package app

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
)

func (w *Watch) CreateListenPath(pathName string) {
	w.listenPath = pathName
}

func (w *Watch) CreateWatch() {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	w.watch = watch
}

func (w *Watch) GetConfig() {

	var watchPathConfig = path.Join(w.listenPath, ".watch")

	file, err := os.OpenFile(watchPathConfig, os.O_RDONLY, 0666)
	if err != nil {

		var yes = "Y"

		for {

			color.Bold.Println(watchPathConfig, "is not found, create .watch file now : [Y/N]")

			if _, err := fmt.Scanf("%s", &yes); err != nil {
				break
			}

			yes = strings.ToUpper(yes)

			if yes != "N" && yes != "Y" {
				os.Exit(0)
			}

			break

		}

		if yes == "N" {
			os.Exit(0)
		}

		f, err := os.Create(watchPathConfig)
		if err != nil {
			panic(err)
		}

		defer func() { _ = f.Close() }()

		tf := strings.NewReader(strings.Join(Template, "\r\n"))

		_, err = io.Copy(f, tf)
		if err != nil {
			panic(err)
		}

		file, _ = os.OpenFile(watchPathConfig, os.O_RDONLY, 0666)

	}

	defer func() { _ = file.Close() }()

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

			if strings.HasSuffix(asbPath, "*") {
				w.config.ignore.others = append(w.config.ignore.others, asbPath)
				continue
			}

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

func (w *Watch) MatchOthers(pathName string) bool {
	for _, v := range w.config.ignore.others {
		if strings.HasPrefix(pathName, v[:len(v)-1]) {
			return true
		}
	}
	return false
}

func (w *Watch) MatchPath(pathName string) bool {
	for _, v := range w.config.ignore.paths {
		if strings.HasPrefix(pathName, v) {
			return true
		}
	}
	return false
}

func (w *Watch) MatchFile(pathName string) bool {
	for _, v := range w.config.ignore.files {
		if pathName == v {
			return true
		}
	}
	return false
}

func (w *Watch) GetMd5(pathName string) string {

	f, err := os.Open(pathName)
	if err != nil {
		panic(err)
	}

	defer func() { _ = f.Close() }()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

func (w *Watch) SetMd5ToCache(pathName string, value string) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.cache[pathName] = w.GetMd5(pathName)
}

func (w *Watch) GetMd5FromCache(pathName string) string {
	w.mux.RLock()
	defer w.mux.RUnlock()
	if content, ok := w.cache[pathName]; ok {
		return content
	}
	return ""
}

func (w *Watch) GetSize(pathName string) (int, error) {
	info, err := os.Stat(pathName)
	if err != nil {
		return 0, err
	}
	return int(info.Size()), err
}

func (w *Watch) SetModTimeToCache(pathName string, value string) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.cache[pathName] = w.GetModTime(pathName)
}

func (w *Watch) GetModTimeFromCache(pathName string) string {
	w.mux.RLock()
	defer w.mux.RUnlock()
	if content, ok := w.cache[pathName]; ok {
		return content
	}
	return ""
}

func (w *Watch) GetModTime(pathName string) string {
	info, err := os.Stat(pathName)
	if err != nil {
		panic(err)
	}

	var modTime = info.ModTime()

	return modTime.String()
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
	time.Sleep(300 * time.Millisecond)

	// size, err := w.GetSize(pathName)
	// if err != nil {
	// 	return false
	// }

	var cache = w.GetMd5FromCache(pathName)

	var value = w.GetMd5(pathName)

	w.SetMd5ToCache(pathName, value)

	return cache != value
}

func (w *Watch) WatchPathExceptIgnore() {
	_ = filepath.Walk(w.listenPath, func(pathName string, info os.FileInfo, err error) error {

		if w.MatchOthers(pathName) {
			return err
		}

		if !info.IsDir() {
			return err
		}

		if w.MatchPath(pathName) {
			return err
		}

		err = w.watch.Add(pathName)
		if err != nil {
			return err
		}

		w.AddTask(pathName)

		color.Bold.Println("watch dir", pathName)

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

	var s = 0

	for _, dir := range w.task {
		_ = filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {

			if i.IsDir() {
				return err
			}

			size, err := w.GetSize(p)

			if err != nil {
				return err
			}

			s += size

			var value = w.GetMd5(p)

			w.SetMd5ToCache(p, value)

			return err
		})
	}

	color.Bold.Println(fmt.Sprintf("go-watch cache size is %d KB", s/1024))
}
