package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gookit/color"

	"github.com/Lemo-yxk/go-watch/app"
	"github.com/Lemo-yxk/go-watch/vars"
)

func init() {

	color.Bold.Println("Welcome use go watch")
	color.Bold.Println("version:1.3")

	flag.StringVar(&vars.ListenPath, "path", ".", "path")
	flag.IntVar(&vars.Sig, "sig", 0x2, "sig")
	flag.Parse()

	info, err := os.Stat(vars.ListenPath)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		color.Red.Println(vars.ListenPath, "is not dir")
		os.Exit(0)
	}

	l, _ := filepath.Abs(vars.ListenPath)

	vars.ListenPath = l
}

func main() {

	var watch = &app.Watch{}

	watch.CreateListenPath(vars.ListenPath)

	watch.Run()

}
