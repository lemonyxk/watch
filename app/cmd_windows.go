// +build windows

package app

import (
	"os"
	"os/exec"
	"strconv"
)

func newCmd(command string) *exec.Cmd {
	var cmd = exec.Command("cmd", "/C", command)
	return cmd
}

func killGroup(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}
