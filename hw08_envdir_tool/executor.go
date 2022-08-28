package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	var setEnvData []string
	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return 1
			}
		} else {
			setEnvData = append(setEnvData, strings.Join([]string{key, val.Value}, "="))
		}
	}
	command.Env = append(os.Environ(), setEnvData...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	var waitStatus syscall.WaitStatus
	if err := command.Run(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		var errorExit *exec.ExitError
		if errors.As(err, &errorExit) {
			waitStatus = errorExit.Sys().(syscall.WaitStatus)
			returnCode = waitStatus.ExitStatus()
		} else {
			returnCode = 112
		}
	} else {
		// Command was successful
		waitStatus = command.ProcessState.Sys().(syscall.WaitStatus)
		returnCode = waitStatus.ExitStatus()
	}
	return returnCode
}
