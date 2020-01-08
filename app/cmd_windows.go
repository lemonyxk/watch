package app

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/gookit/color"
)

func (w *Watch) StopProcess() {

	for _, cmd := range w.commands {

		var err error
		err = killWindows(cmd)
		if err != nil {
			color.Red.Println(err)
		}

		_, err = cmd.Process.Wait()
		if err != nil {
			color.Red.Println(err)
		}

		color.Bold.Println(cmd.Process.Pid, "kill success")
	}

	w.commands = nil
}

func killWindows(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

func (w *Watch) startProcess() {

	for _, v := range w.config.start {

		var cmd *exec.Cmd

		cmd = exec.Command("cmd", "/C", v)

		cmd.Dir = w.listenPath
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
