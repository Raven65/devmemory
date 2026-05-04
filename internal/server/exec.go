//go:build !windows

package server

import "os/exec"

func execCommand(name string, args ...string) error {
	return exec.Command(name, args...).Start()
}
