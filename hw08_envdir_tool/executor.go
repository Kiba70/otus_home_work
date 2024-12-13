package main

import (
	"log/slog"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Правим переменные окружения
	for key, e := range env {
		if e.NeedRemove {
			os.Unsetenv(key)
			continue
		}
		os.Setenv(key, e.Value)
	}

	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	// c.Env = os.Environ() - Не обязательно - используются текущие (уже установленные) переменные окружения

	err := c.Run()
	if err != nil {
		slog.Error("Ошибка выполнения вызываемой программы", "error", err.Error())
	}
	returnCode = c.ProcessState.ExitCode()

	return
}
