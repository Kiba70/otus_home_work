package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrErrorsIncorrectParams = errors.New("incorrect parameters")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		countErrors atomic.Int32
		wg          sync.WaitGroup
	)

	// Граничные значения
	if n <= 0 || len(tasks) == 0 {
		return ErrErrorsIncorrectParams
	}

	ch := make(chan Task) // Канал на передачу задач. Специально не буферизирован.

	// По условиям задачи сразу запускаем указанное количество горутин (или по количеству задач)
	for i := 0; i < min(n, len(tasks)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch { // Забираем очередную задачу из очереди
				if err := task(); err != nil {
					countErrors.Add(1) // Ошибка - увеличиваем счётчик ошибок
				}
			}
		}()
	}

	// Запускаем задачи
	for _, task := range tasks {
		if m > 0 && countErrors.Load() >= int32(m) { // Достугнут счётчик ошибок
			break
		}
		ch <- task
	}

	close(ch) // Закрываем канал - команда на завершение горутинам
	wg.Wait() // Ждём всех

	if m > 0 && countErrors.Load() >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
