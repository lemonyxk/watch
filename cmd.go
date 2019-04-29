package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func (w *Watch) StopProcess() {

	if w.pid == 0 {
		return
	}

	err := syscall.Kill(-w.pid, syscall.SIGINT)
	if err != nil {
		log.Println(err)
	}

	w.pid = 0

	log.Println(w.cmd.Process.Pid, "kill success")
}

func (w *Watch) GetStartCommand() string {
	return fmt.Sprintf("cd %s && %s", w.listenPath, strings.Join(w.config.start, " && "))
}

func (w *Watch) startProcess() {

	if w.pid != 0 {
		return
	}

	w.cmd = exec.Command("bash", "-c", w.GetStartCommand())

	w.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	w.cmd.Stderr = os.Stderr
	w.cmd.Stdin = os.Stdin
	w.cmd.Stdout = os.Stdout

	err := w.cmd.Start()
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	w.pid = w.cmd.Process.Pid

	log.Println(w.cmd.Process.Pid, "run success")
}

func (w *Watch) HasStartSuccess() (string, error) {

	process, err := os.FindProcess(w.cmd.Process.Pid)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%d is running", process.Pid), nil

	//cmd := exec.Command("bash", "-c", "ps axu | grep -v grep | grep "+w.GetStartCommand())
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	return "", err
	//}
	//
	//return string(out), nil
}
