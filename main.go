package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lemoyxk/watch/app"
	"github.com/lemoyxk/watch/vars"
)

func init() {

	fmt.Println("Welcome use go watch")
	fmt.Println("version:1.4")

	flag.StringVar(&vars.ListenPath, "path", ".", "path")
	flag.IntVar(&vars.Sig, "sig", 0x2, "sig")
	flag.Parse()

	info, err := os.Stat(vars.ListenPath)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		fmt.Println(vars.ListenPath, "is not dir")
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
