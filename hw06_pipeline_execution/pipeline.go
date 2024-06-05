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

	for data := range Pipeline(in, stages...) {
		select {
		case <-done:
			close(out)
			return out
		default:
			wg.Add(1)
			go func(sendData any) {
				select {
				case out <- sendData:
					wg.Done()
				case <-done:
				}
			}(data)
		}
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
