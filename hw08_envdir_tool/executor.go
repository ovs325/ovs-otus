package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd запускает команду + аргументы (cmd) с переменными окружения из env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	envList := make([]string, 0, len(env))
	for name, value := range env {
		os.Unsetenv(name)
		if !value.NeedRemove {
			envList = append(envList, fmt.Sprintf("%s=%s", name, value.Value))
		}
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = os.Environ()
	command.Env = append(command.Env, envList...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return e.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		return 1
	}
	return 0
}
