package hw06pipelineexecution

import "sync"

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

type FromStage struct {
	numData int
	data    any
}

func withNum(inCh In, num int) any {
	for res := range inCh {
		return FromStage{
			numData: num,
			data:    res,
		}
	}
	return FromStage{}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	chList := make([]Bi, len(stages)+1)
	// Слайс каналов, для автопостройки Pipeline
	for i := 0; i < len(stages); i++ {
		chList[i] = make(Bi)
	}
	chList[len(stages)] = make(Bi)

	out := make(Bi)
	var num int
	var mu sync.Mutex
	dataList := make(map[int]any)
	// Горутина для сортировки исходных данных
	go func() {
		for {
			select {
			case <-done:
				mu.Lock()
				for i := 0; i < len(dataList); i++ {
					out <- dataList[i]
				}
				mu.Unlock()
				close(out)
				return
			case resNum, ok := <-chList[len(stages)]:
				if !ok {
					mu.Lock()
					for i := 0; i < len(dataList); i++ {
						out <- dataList[i]
					}
					mu.Unlock()
					close(out)
					return
				}
				res := resNum.(FromStage)
				mu.Lock()
				dataList[res.numData] = res.data
				if len(dataList) == num {
					for _, ch := range chList {
						close(ch)
					}
				}
				mu.Unlock()
			}
		}
	}()
	// Горутины для работы со Stage-ами
	StRanner(stages, chList, done)

	// Горутина для нумерации входных данных
	go func() {
		for data := range in {
			mu.Lock()
			chList[0] <- FromStage{
				numData: num,
				data:    data,
			}
			num++
			mu.Unlock()
		}
	}()
	return out
}

// Горутины для работы со Stage-ами.
func StRanner(stages []Stage, chList []chan any, done In) {
	for i, stage := range stages {
		go func(numCh int, stage Stage) {
			for {
				toStage := make(Bi)
				select {
				case itemNum, ok := <-chList[numCh]:
					if ok {
						item := itemNum.(FromStage)
						go func() {
							chList[numCh+1] <- withNum(stage(toStage), item.numData)
						}()
						toStage <- item.data
					}
				case <-done:
					close(toStage)
					return
				}
				close(toStage)
			}
		}(i, stage)
	}
}
