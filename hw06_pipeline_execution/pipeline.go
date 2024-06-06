package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func Pipeline(in In, stages ...Stage) Out {
	pipeOut := in
	for _, stage := range stages {
		pipeOut = stage(pipeOut)
	}
	return pipeOut
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var wg sync.WaitGroup
	out := make(Bi)

	for {
		data, ok := <-in
		if !ok {
			break
		}
		wg.Add(1)
		toPipe := make(Bi)

		go func() {
			select {
			case <-done:
				wg.Done()
				return
			case res := <-Pipeline(toPipe, stages...):
				out <- res
			}
			wg.Done()
		}()
		toPipe <- data
		close(toPipe)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
