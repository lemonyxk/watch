package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gookit/color"

	"github.com/Lemo-yxk/go-watch/app"
)

var ListenPath = "."

func init() {

	color.Bold.Println("Welcome use go watch")
	color.Bold.Println("version:1.1")

	flag.StringVar(&ListenPath, "path", ".", "path")
	flag.Parse()

	info, err := os.Stat(ListenPath)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		color.Red.Println(ListenPath, "is not dir")
		os.Exit(0)
	}

	l, _ := filepath.Abs(ListenPath)

	ListenPath = l
}

func main() {

	var watch = &app.Watch{}

	watch.CreateListenPath(ListenPath)

	watch.Run()

}
