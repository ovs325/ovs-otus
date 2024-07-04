package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m < 0 {
		return ErrErrorsLimitExceeded
	}
	lenTasks := len(tasks)
	if lenTasks == 0 || n <= 0 {
		return nil
	}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan Task)
	responseChan := make(chan error)
	defer close(responseChan)
	defer cancel()
	for ; n > 0; n-- {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			for task := range taskChan {
				select {
				case <-ctx.Done():
					return
				case responseChan <- task():
				}
			}
		}(ctx)
	}
	// Запускаем задачи и считываем результат
	for taskRun := 0; m > 0 && taskRun < lenTasks; {
		select {
		case res := <-responseChan:
			if res != nil {
				m--
			}
		case taskChan <- tasks[taskRun]:
			taskRun++
		}
	}
	cancel()
	close(taskChan)
	wg.Wait()
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
