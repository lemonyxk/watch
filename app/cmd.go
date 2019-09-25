package app

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gookit/color"
)

func (w *Watch) StopProcess() {

	for _, cmd := range w.commands {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		if err != nil {
			color.Red.Println(err)
		}

		_, _ = cmd.Process.Wait()

		color.Bold.Println(cmd.Process.Pid, "kill success")
	}

	w.commands = nil
}

func (w *Watch) startProcess() {

	for _, v := range w.config.start {

		var args = strings.Split(v, " ")
		var cmd = exec.Command(args[0], args[1:]...)
		cmd.Dir = w.listenPath

		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		w.commands = append(w.commands, cmd)

		color.Bold.Println(cmd.Process.Pid, "run success")
	}

}
