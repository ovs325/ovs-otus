package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type FromStage struct {
	numData int
	data    any
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	var num int
	var mu sync.Mutex
	// Слайс каналов, для автопостройки Pipeline
	chList := getChList(len(stages) + 1)
	// Сортировщик и Closer
	go closer(chList, out, done, &mu, &num)
	// Pipeline - Запуск цепочки Stages
	pipeline(stages, chList, done)
	// Горутина для нумерации входных данных
	go numeratorIn(chList, in, &mu, &num)
	return out
}

// Отдаёт слайс каналов.
func getChList(num int) []Bi {
	chList := make([]Bi, num)
	// Слайс каналов, для автопостройки Pipeline
	for i := 0; i < num; i++ {
		chList[i] = make(Bi)
	}
	return chList
}

// Ф-ция, скидывающая сортированные данные в выходной канал out
// и закрывающая все остальные каналы горутин для закрытия их самих.
func closer(chList []Bi, out Bi, done Out, mu *sync.Mutex, num *int) {
	dataList := make(map[int]any)
	for {
		select {
		case <-done:
			mu.Lock()
			sortDataSender(dataList, out)
			mu.Unlock()
			return
		case resNum, ok := <-chList[len(chList)-1]:
			if !ok {
				mu.Lock()
				sortDataSender(dataList, out)
				mu.Unlock()
				return
			}
			res := resNum.(FromStage)
			mu.Lock()
			dataList[res.numData] = res.data
			if len(dataList) == *num {
				// закрываем все каналы горутин (сами горутины)
				channelsCloser(chList)
			}
			mu.Unlock()
		}
	}
}

// Отправка итогового результата в out-канал.
func sortDataSender(dataList map[int]any, out Bi) {
	for i := 0; i < len(dataList); i++ {
		out <- dataList[i]
	}
	close(out)
}

// Гаситель каналов.
func channelsCloser(chList []Bi) {
	for _, ch := range chList {
		close(ch)
	}
}

// Запуск списка Stage-й.
func pipeline(stages []Stage, chList []chan any, done In) {
	for i, stage := range stages {
		go stageRun(i, stage, chList, done)
	}
}

// Запуск Stage.
func stageRun(numCh int, stage Stage, chList []chan any, done In) {
	for {
		toStage := make(Bi)
		select {
		case itemNum, ok := <-chList[numCh]:
			if ok {
				item := itemNum.(FromStage)
				go stageHandler(chList[numCh+1], stage(toStage), item.numData)
				toStage <- item.data
			}
		case <-done:
			// Закрываем внешнюю горутину Stage-а
			close(toStage)
			return
		}
		close(toStage)
	}
}

// Сброс нумерованных данных, полученных из Stage, в выходной канал.
func stageHandler(ch Bi, chOut Out, num int) {
	ch <- addNum(chOut, num)
}

// нумерует выходные данные Stages.
func addNum(inCh In, num int) any {
	for res := range inCh {
		return FromStage{
			numData: num,
			data:    res,
		}
	}
	return FromStage{}
}

// Нумератор порядка входных данных, поступающих на вход Pipeline.
func numeratorIn(chList []Bi, in In, mu *sync.Mutex, num *int) {
	for data := range in {
		mu.Lock()
		chList[0] <- FromStage{
			numData: *num,
			data:    data,
		}
		*num++
		mu.Unlock()
	}
}
