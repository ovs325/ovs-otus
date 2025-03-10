package group

import (
	tp "bruteforce/internal/types"
	"container/heap"

	"runtime"
	"sync"
	"time"
)

type bucketMap map[string]BucketItem

type Queue interface {
	heap.Interface
	// Получить первый элемент очереди (с наименьшим приорететом)
	GetFirstDrops() any
}

type BucketItem interface {
	// Объем недокапанного запаса в каплях.
	DropsSum() int64
	// Промежуток времени до опустошения корзины.
	DeadLine() time.Duration
	// Добавление содержимиго в корзину в каплях с учетом вместимости.
	AddDrops(addNum int64) (int64, bool)
	GetID() string
	GetPriority() int
	// Нужно ли обнулять корзину?
	IsReset() (ok bool)
	// Является ли Нулевым указателем типа корзины
	IsNilPointer(i any) bool
	GetAllBucketParams() *tp.AllBucketParams
}

// Группа для корзин с одинаковыми характеристиками
// (Общей ёмкости группы не существует)
// Капли - алиас для понятия 'Условные еденицы'
type BucketGroup struct {
	lock        sync.Mutex // Для потокобезопастности
	bucketsMap  bucketMap  // Карта корзин группы
	lenQueue    int64      // Длинна очереди (новой карты)
	queue       Queue      // Приорететная Очередь
	leakageRate float64    // Скорость утечки в каплях в секунду (для каждой корзины в группе)
	capacity    int64      // Ёмкость в каплях (для каждой корзины в группе)
	exitCh      chan bool  // Канал для сброса планировщика (чистильщика)

	newItem  func(id string, leakageRate float64, capacity int64) BucketItem // Фабрика новых экземпляров (Bucket)
	newQueue func(lenQueue int) Queue                                        // Фабрика новых экземпляров (Queye)
}

func NewBucketGroup(
	leakageRate float64,
	capacity, lenQueue int64,
	queue Queue,
	isRunSheduller bool,
	nwFunc func(id string, leakageRate float64, capacity int64) BucketItem,
	nwQueue func(lenQueue int) Queue,
) *BucketGroup {
	group := &BucketGroup{
		lenQueue:    lenQueue,
		bucketsMap:  make(bucketMap),
		queue:       queue,
		leakageRate: leakageRate,
		capacity:    capacity,
		exitCh:      make(chan bool),

		newItem:  nwFunc,
		newQueue: nwQueue,
	}
	if isRunSheduller {
		group.Scheduller(time.Second)
	}

	return group
}

func DelBucketGroup(group *BucketGroup) {
	group.ExistScheduller()
	group.bucketsMap = nil
	group = nil

	runtime.GC()
}

// Остановка планировщика удаления
func (g *BucketGroup) ExistScheduller() {
	// c.Reset()
	close(g.exitCh)
}

// Очистить очередь и пересоздать карту корзин
func (c *BucketGroup) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.bucketsMap = make(bucketMap)
	c.queue = c.newQueue(int(c.lenQueue))
}

func (g *BucketGroup) GetCapacity() int64 {
	return g.capacity
}

func (g *BucketGroup) GetLeakageRate() float64 {
	return g.leakageRate
}

func (g *BucketGroup) GetLenQueue() int64 {
	return g.lenQueue
}

// Объем недокапанного запаса в каплях по id корзины.
func (g *BucketGroup) GetDropsSum(id string) (int64, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	bucket, ok := g.bucketsMap[id]
	if bucket == nil || !ok {
		return 0, false
	}

	return bucket.DropsSum(), true
}

// Промежуток времени до опустошения по id корзины.
func (g *BucketGroup) GetDeadLine(id string) (time.Duration, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	bucket, ok := g.bucketsMap[id]
	if bucket == nil || !ok {
		return 0, false
	}
	return bucket.DeadLine(), true
}

// Пополнить запас капель в по id корзине
func (g *BucketGroup) AddDrops(id string, drops int64) int64 {
	g.lock.Lock()
	defer g.lock.Unlock()

	bucket, ok := g.bucketsMap[id]
	if bucket == nil || !ok {
		// Это место (внутри фабрики), где реальный экземпляр Типа присваивается интерфейсу Bucket-а! И всё!!
		bucket = g.newItem(id, g.leakageRate, g.capacity)
		g.bucketsMap[id] = bucket
		g.queue.Push(bucket)
	}
	if realAddDrops, ok := bucket.AddDrops(drops); ok {
		// Упорядочивание элементов в очереди
		heap.Fix((g.queue).(heap.Interface), bucket.GetPriority())
		return realAddDrops
	}
	return 0
}

// Разрешен ли запрос?
func (g *BucketGroup) IsAllowed(id string, drops int64) (isAllowed bool) {
	if realAddDrops := g.AddDrops(id, drops); drops == realAddDrops {
		return true
	}
	return false
}

// Удаление корзины по по id
func (g *BucketGroup) RemoveDrops(id string) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if bucket, ok := g.bucketsMap[id]; bucket != nil || ok {
		priority := bucket.GetPriority()
		delete(g.bucketsMap, id)
		heap.Remove((g.queue).(heap.Interface), priority)
	}
}

// Удаление пустых корзин
func (g *BucketGroup) RemooveEmpty() {
	g.lock.Lock()
	defer g.lock.Unlock()
	for b := g.queue.GetFirstDrops().(BucketItem); !b.IsNilPointer(b); b = g.queue.GetFirstDrops().(BucketItem) {
		if !b.IsReset() {
			break
		}
		priority := b.GetPriority()
		delete(g.bucketsMap, b.GetID())
		heap.Remove((g.queue).(heap.Interface), priority)
	}
}

// Планировщик удаления пустых корзин
func (g *BucketGroup) Scheduller(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-g.exitCh:
				ticker.Stop()
				return
			case <-ticker.C:
				g.RemooveEmpty()
			}
		}
	}()
}

// Функции для тестирования

func (g *BucketGroup) IsNilQueue() bool {
	return g.queue == nil
}

func (g *BucketGroup) IsNilBuckets() bool {
	return g.bucketsMap == nil
}

func (g *BucketGroup) IsLeakageRateEqRate(rate float64) bool {
	return g.leakageRate == rate
}

func (g *BucketGroup) GetAllGroupData() *tp.AllGroupData {
	group := tp.AllGroupData{}
	group.LenMap = len(g.bucketsMap)
	if bac := g.queue.GetFirstDrops(); !bac.(BucketItem).IsNilPointer(bac) {
		if params := bac.(BucketItem).GetAllBucketParams(); params != nil {
			group.FirstDrops = params.Id
		}
	}
	group.BacketsParams = make(map[string]tp.AllBucketParams, group.LenMap)
	for name, backet := range g.bucketsMap {
		if params := backet.GetAllBucketParams(); params != nil {
			params.IsReset = backet.IsReset()
			group.BacketsParams[name] = *params
		}
	}
	return &group
}
