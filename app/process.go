package app

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/lemoyxk/watch/vars"
)

func (w *Watch) StopProcess() {

	fmt.Println("stop process", syscall.Signal(vars.Sig))

	for i := 0; i < len(w.commands); i++ {
		if !w.commands[i].status {
			continue
		}

		var err = killGroup(w.commands[i].cmd)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(w.commands[i].cmd.Process.Pid, "kill success")
	}

	w.commands = nil
}

func (w *Watch) startProcess(event fsnotify.Event) {

	var start = time.Now()

	fmt.Println("start process", event)

	for i := 0; i < len(w.config.command); i++ {
		var cmdArray = strings.Split(w.config.command[i], ": ")
		var p = ""
		var c = ""
		var wait = false
		if len(cmdArray) == 0 {
			c = ""
		} else if len(cmdArray) == 1 {
			c = cmdArray[0]
		} else {
			var pb = strings.Split(cmdArray[0], " ")
			p = pb[0]
			if !path.IsAbs(p) {
				p = path.Join(w.listenPath, p)
			}
			c = strings.Join(cmdArray[1:len(cmdArray)], ": ")
			if len(pb) > 1 {
				if pb[1] == "wait" {
					wait = true
				}
			}
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

		if wait {
			_, err = cmdInfo.cmd.Process.Wait()
			if err != nil {
				fmt.Println(err)
			}
			cmdInfo.status = false
		} else {
			go func() {
				_, err = cmdInfo.cmd.Process.Wait()
				if err != nil {
					fmt.Println(err)
				}
				cmdInfo.status = false
			}()
		}

		fmt.Println(cmd.Process.Pid, "run success")
	}

	fmt.Println("time", float64(time.Now().Sub(start).Milliseconds())/1000)

}
