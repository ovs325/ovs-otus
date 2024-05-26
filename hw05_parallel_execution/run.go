package hw05parallelexecution

import (
	"context"
	"errors"
	"runtime"
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

	var err error
	if n > 0 {
		err = taskLauncher(tasks, taskChan, responseChan, m)
	}
	cancel()
	wg.Wait()
	return err
}

func taskLauncher(tasks []Task, taskChan chan<- Task, responseChan chan error, m int) error {
	success := 0
	taskRun := 1
	fail := 0
	taskChan <- tasks[0]
	tasks = tasks[1:]

	for fail+success < taskRun {
		select {
		case res := <-responseChan:
			switch res {
			case nil:
				success++
			default:
				fail++
			}
		default:
			if len(tasks) > 0 && fail < m {
				select {
				case taskChan <- tasks[0]:
					tasks = tasks[1:]
					taskRun++
				default:
					runtime.Gosched()
				}
			} else {
				runtime.Gosched()
			}
		}
	}
	if fail >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
