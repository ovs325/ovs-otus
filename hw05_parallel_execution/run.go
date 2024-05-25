package hw05parallelexecution

import (
	"context"
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var taskLauncher = func(ctx context.Context, tsk Task, tCh chan error) {
	select {
	case <-ctx.Done():
		tCh <- ErrErrorsLimitExceeded
		return
	case tCh <- tsk():
		return
	}
}

// Run запускает задачи в n подпрограммах и останавливает свою работу при получении m ошибок от задач.
func Run(tasks []Task, n, m int) error {
	if m < 0 || n <= 0 {
		return ErrErrorsLimitExceeded
	}
	lenTasks := len(tasks)
	if lenTasks == 0 {
		return nil
	}
	tChan := make(chan error)
	exceedErr := 0
	commonErr := 0
	success := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if lenTasks > n {
		for _, task := range tasks[:n] {
			go taskLauncher(ctx, task, tChan)
		}
		tasks = tasks[n:]
	} else {
		for _, task := range tasks {
			go taskLauncher(ctx, task, tChan)
		}
		tasks = []Task{}
	}

	for commonErr+success+exceedErr+len(tasks) < lenTasks {
		select {
		case res := <-tChan:
			if res != nil {
				if errors.Is(res, ErrErrorsLimitExceeded) {
					exceedErr++
					continue
				}
				commonErr++
				if commonErr >= m {
					cancel()
					continue
				}
				if len(tasks) > 0 {
					go taskLauncher(ctx, tasks[0], tChan)
					tasks = tasks[1:]
				}
				continue
			}
			success++
			if len(tasks) > 0 {
				go taskLauncher(ctx, tasks[0], tChan)
				tasks = tasks[1:]
			}
		default:
		}
	}
	if commonErr >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
