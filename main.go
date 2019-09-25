package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Lemo-yxk/go-watch/app"
)

var ListenPath = "."

func init() {

	fmt.Println("Welcome use go watch")
	fmt.Println("version:1.0")

	log.SetFlags(log.Ltime | log.Ldate)

	flag.StringVar(&ListenPath, "path", ".", "path")
	flag.Parse()

	info, err := os.Stat(ListenPath)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	if !info.IsDir() {
		log.Println(ListenPath, "is not a dir")
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
