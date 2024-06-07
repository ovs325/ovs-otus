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

		go func(n int) {
			results := make([]any, len(stages)+1)
			results[0] = data
			for m, stage := range stages {
				toStage := make(Bi)
				go func(mm int) {
					toStage <- results[mm]
				}(m)
				val, ok := <-stage(toStage)
				if ok {
					results[m+1] = val
				}
				close(toStage)
			}
			fromCh <- FromCh{num: n, data: results[len(stages)]}

			wg.Done()
		}(i)
		i++
	}

	go func() {
		wg.Wait()
		close(fromCh)
	}()
	return out
}
