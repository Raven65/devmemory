//go:build windows

package server

import "os/exec"

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = nil // use default console behavior
	return cmd.Start()
}
