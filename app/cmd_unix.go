// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package app

import (
	"os/exec"
	"syscall"

	"github.com/Lemo-yxk/go-watch/vars"
)

func killGroup(cmd *exec.Cmd) error {
	return syscall.Kill(-cmd.Process.Pid, syscall.Signal(vars.Sig))
}

func newCmd(command string) *exec.Cmd {
	var cmd = exec.Command("bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
