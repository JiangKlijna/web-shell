// +build windows

package main

import (
	"github.com/iamacarpet/go-winpty"
)

type windowsPTY struct {
	cmd string
	pty *winpty.WinPTY
}

func (wp windowsPTY) Close() {
	wp.pty.Close()
}

func (wp windowsPTY) SetSize(w, h uint16) {
	wp.pty.SetSize(uint32(w), uint32(h))
}

func (wp windowsPTY) Read(p []byte) (n int, err error) {
	return wp.pty.StdOut.Read(p)
}

func (wp windowsPTY) Write(p []byte) (n int, err error) {
	return wp.pty.StdIn.Write(p)
}

func OpenPty(cmd string) (PTY, error) {
	pty, err := winpty.Open("winpty", cmd)
	if err != nil {
		return nil, err
	}
	return &windowsPTY{cmd, pty}, nil
}
