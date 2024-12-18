package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("test with ignore errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	// Какая из двух версий тестов ниже более правильна?
	// Та, где переменная runConcurentlyTasks объявлена как atomic или где к ней обращаемся по ссылке?
	t.Run("tasks without errors by require.Eventually", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var (
			runTasksCount       int32
			runConcurentlyTasks atomic.Int32
			sumTime             time.Duration
			wg                  sync.WaitGroup
		)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				runConcurentlyTasks.Add(1)
				defer runConcurentlyTasks.Add(-1)
				time.Sleep(taskSleep) // В любом случае необходимо какое-то время на выполнение задачи
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
		}()

		require.Eventually(t, func() bool { return runConcurentlyTasks.Load() > 1 },
			sumTime, time.Millisecond, "not parallel")
		wg.Wait()

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("tasks without errors by require.Eventually version &runConcurentlyTasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var (
			runTasksCount       int32
			runConcurentlyTasks int32
			sumTime             time.Duration
			wg                  sync.WaitGroup
		)

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				atomic.AddInt32(&runConcurentlyTasks, 1)
				defer atomic.AddInt32(&runConcurentlyTasks, -1)
				time.Sleep(taskSleep) // В любом случае необходимо какое-то время на выполнение задачи
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
		}()

		require.Eventually(t, func() bool { return atomic.LoadInt32(&runConcurentlyTasks) > 1 },
			sumTime, time.Millisecond, "not parallel")
		wg.Wait()

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}
