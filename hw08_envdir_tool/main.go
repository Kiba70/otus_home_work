package main

import (
	"log/slog"
	"os"
)

func main() {
	exitCode := 111
	defer func() {
		os.Exit(exitCode)
	}()

	// Учитывая, что формат командной строки  жёстко задан, flags не используем
	if len(os.Args) < 3 {
		slog.Error("Incorrect input\n" + os.Args[0] + " PATH COMMAND arg1 arg2...\n")
	}

	env, err := ReadDir(os.Args[1]) // os.Args[1] - path
	if err != nil {
		slog.Error("Ошибка чтения каталога с переменными окружения", "error", err.Error())
		return
	}

	exitCode = RunCmd(os.Args[2:], env) // os.Args[2:] - команда на запуск с параметрами
}
