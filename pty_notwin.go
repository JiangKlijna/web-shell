// +build !windows

package main

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
)

type notwinPTY struct {
	cmd string
	f   *os.File
}

func (np notwinPTY) Read(p []byte) (n int, err error) {
	return np.f.Read(p)
}

func (np notwinPTY) Write(p []byte) (n int, err error) {
	return np.f.Write(p)
}

func (np notwinPTY) Close() {
	np.f.Close()
}

func (np notwinPTY) SetSize(w, h uint16) {
	pty.Setsize(np.f, &pty.Winsize{w, h, 8, 8})
}

func OpenPty(cmd string) (PTY, error) {
	f, err := pty.Start(exec.Command(cmd))
	if err != nil {
		return nil, err
	}
	return &notwinPTY{cmd, f}, nil
}
