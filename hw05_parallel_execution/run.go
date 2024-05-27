package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run запускает задачи в n подпрограммах и останавливает свою работу при получении m ошибок от задач.
func Run(tasks []Task, n, m int) error {
	if m < 0 {
		return ErrErrorsLimitExceeded
	}
	lenTasks := len(tasks)
	if lenTasks == 0 {
		return nil
	}
	var wg sync.WaitGroup
	responseChan := make(chan error)
	defer close(responseChan)
	taskChan := make(chan Task)
	defer close(taskChan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task := <-taskChan:
					select {
					case <-ctx.Done():
						return
					case responseChan <- task():
					}
				}
			}
		}(ctx)
	}

	fail := 0
	if n > 0 { // Запускаем задачи и считываем тезультат
		success := 0
		taskChan <- tasks[0]
		taskRun := 1
		for fail < m && fail+success < taskRun && taskRun < len(tasks) {
			select {
			case res := <-responseChan:
				if res == nil {
					success++
					continue
				}
				fail++
			case taskChan <- tasks[taskRun]:
				taskRun++
			}
		}
	}
	cancel()
	wg.Wait()
	if fail >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
