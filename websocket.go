package main

import (
	"io"
	"os/exec"
)

func execute(sh string, in io.Reader, out io.Writer, down chan error) {
	cmd := exec.Command(sh)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Start()
	if err != nil {
		down <- err
		return
	}
	down <- cmd.Wait()
}
