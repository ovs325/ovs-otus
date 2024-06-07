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

func StagesChaine(in In, stages ...Stage) Out {
	pipeOut := in
	for _, stage := range stages {
		pipeOut = stage(pipeOut)
	}
	return pipeOut
}

type FromCh struct {
	num  int
	data any
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var wg sync.WaitGroup
	out := make(Bi)

	fromCh := make(Bi)

	go func() {
		dataList := map[int]any{}
		for {
			select {
			case <-done:
				close(out)
				return
			case res, ok := <-fromCh:
				if !ok {
					for i := 0; i < len(dataList); i++ {
						out <- dataList[i]
					}
					close(out)
					return
				}
				item := res.(FromCh)
				dataList[item.num] = item.data
			}
		}
	}()

	i := 0
	for {
		data, ok := <-in
		if !ok {
			break
		}
		wg.Add(1)
		toPipe := make(Bi)

		go func(n int) {
			isSend := true
			var res any
			select {
			case <-done:
				isSend = false
			case res = <-StagesChaine(toPipe, stages...):
			}
			if isSend {
				fromCh <- FromCh{num: n, data: res}
			}
			wg.Done()
		}(i)
		select {
		case <-done:
		case toPipe <- data:
		}
		close(toPipe)
		i++
	}

	go func() {
		wg.Wait()
		close(fromCh)
	}()
	return out
}
