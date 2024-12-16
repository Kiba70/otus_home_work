package main

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fsDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(dir))

	for _, entry := range fsDir {
		if strings.Contains(entry.Name(), "=") {
			continue // "=" не должно быть в имени файла
		}

		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			slog.Error("Не могу открыть файл", "error", err)
			continue
		}
		// Не используем defer т.к. f необходимо закрыть до выхода из функции

		ev := EnvValue{}

		fStat, _ := f.Stat()
		if fStat.Size() == 0 {
			ev.NeedRemove = true
			env[entry.Name()] = ev
		} else {
			scanner := bufio.NewScanner(f)
			scanner.Scan()

			ev.Value = string(bytes.ReplaceAll([]byte(scanner.Text()), []byte{0}, []byte{10}))
			ev.Value = strings.TrimRight(ev.Value, "\t ")

			env[entry.Name()] = ev
		}

		f.Close()
	}

	return env, nil
}
