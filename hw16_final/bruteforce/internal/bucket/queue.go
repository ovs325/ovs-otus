package buckets

import (
	gr "bruteforce/internal/group"
)

// https://golang.org/pkg/container/heap/
//
// Приорететная очередь корзин
type queyeBuckets []*Bucket

func NewQueyeBuckets(volume int) gr.Queue {
	queye := make(queyeBuckets, 0, volume)
	return &queye
}

func (p queyeBuckets) Len() int { return len(p) }

func (p queyeBuckets) Less(i, j int) bool { return p[i].expiry.Before(p[j].expiry) }

func (p queyeBuckets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].priority = i
	p[j].priority = j
}

func (p *queyeBuckets) Push(x any) {
	b := x.(*Bucket)
	b.priority = p.Len()
	*p = append(*p, b)
}

func (p *queyeBuckets) Pop() any {
	queueOld := *p
	numDel := queueOld.Len() - 1
	delItem := queueOld[numDel]
	queueOld[numDel] = nil // Это для GS
	delItem.priority = -1  // Сохраняет данные
	*p = queueOld[0:numDel]
	return delItem
}

func (p *queyeBuckets) Reset(lenQueue int) {
	if lenQueue == 0 {
		lenQueue = cap(*p)
	}
	newQueue := make(queyeBuckets, 0, lenQueue)
	*p = newQueue
}

// Корзины, которые пусты или появятся раньше всего, находятся на вершине кучи.
// (Окажутся в начале Приорететной очереди с индексом 0 и будут выданы функцией GetFirstDrops()).
// Это позволяет быстро сокращать количество пустых корзин, что очень хорошо масштабируется.
// Приоритет настраивается каждый раз, когда в очередь добавляется какое-либо количество через Push()
func (p queyeBuckets) GetFirstDrops() any {
	if p.Len() > 0 {
		bucket := p[0]
		return bucket
	}
	var empty *Bucket = nil
	return empty
}
