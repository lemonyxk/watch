package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Lemo-yxk/lemo/console"

	"github.com/Lemo-yxk/go-watch/app"
	"github.com/Lemo-yxk/go-watch/vars"
)

func init() {

	console.Bold.Println("Welcome use go watch")
	console.Bold.Println("version:1.4")

	flag.StringVar(&vars.ListenPath, "path", ".", "path")
	flag.IntVar(&vars.Sig, "sig", 0x2, "sig")
	flag.Parse()

	info, err := os.Stat(vars.ListenPath)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		console.FgRed.Println(vars.ListenPath, "is not dir")
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
