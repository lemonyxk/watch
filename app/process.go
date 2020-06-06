// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package app

import (
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/Lemo-yxk/lemo/console"
	"github.com/fsnotify/fsnotify"

	"github.com/Lemo-yxk/go-watch/vars"
)

func (w *Watch) StopProcess() {

	console.Bold.Println("stop process", syscall.Signal(vars.Sig))

	for i := 0; i < len(w.commands); i++ {
		if !w.commands[i].status {
			continue
		}

		var err = killGroup(w.commands[i].cmd)
		if err != nil {
			console.FgRed.Println(err)
		}

		console.Bold.Println(w.commands[i].cmd.Process.Pid, "kill success")
	}

	w.commands = nil
}

func (w *Watch) startProcess(event fsnotify.Event) {

	var start = time.Now()

	console.Bold.Println("start process", event)

	for i := 0; i < len(w.config.command); i++ {
		var cmdArray = strings.Split(w.config.command[i], ": ")
		var p = ""
		var c = ""
		if len(cmdArray) == 0 {
			c = ""
		} else if len(cmdArray) == 1 {
			c = cmdArray[0]
		} else {
			p = cmdArray[0]
			if !path.IsAbs(cmdArray[0]) {
				p = path.Join(w.listenPath, cmdArray[0])
			}
			c = strings.Join(cmdArray[1:len(cmdArray)], ": ")
		}

		if p == "" {
			p = w.listenPath
		}

		if strings.HasSuffix(p, "*") {
			p = p[0 : len(p)-1]
		}

		if event.Name != "init" {
			if !strings.HasPrefix(event.Name, p) {
				continue
			}
		}

		if c == "" {
			continue
		}

		var cmd = newCmd(c)

		cmd.Dir = w.listenPath
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		var cmdInfo = &CmdInfo{cmd: cmd, status: true}
		w.commands = append(w.commands, cmdInfo)

		go func() {
			_, err = cmdInfo.cmd.Process.Wait()
			if err != nil {
				console.FgRed.Println(err)
			}
			cmdInfo.status = false
		}()

		console.Bold.Println(cmd.Process.Pid, "run success")
	}

	console.Bold.Println("time", float64(time.Now().Sub(start).Milliseconds())/1000)

}
