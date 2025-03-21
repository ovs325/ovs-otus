package buckets

import (
	"math"
	"time"

	sr "bruteforce/internal"
	gr "bruteforce/internal/group"
	tp "bruteforce/internal/types"
)

// Объект 'Корзина'.
type Bucket struct {
	id          string    // id корзины для поиска по карте.
	capacity    int64     // Емкость в каплях
	leakageRate float64   // Скорость утечки (капли.сек)
	expiry      time.Time // время опустошения корзины.
	priority    int       // Приоритет корзины в очереди кучи
}

// Создает новую Корзину с требуемой скоростью утечки и ёмкостью.
func NewBucket(name string, r float64, cap int64) gr.BucketItem {
	return &Bucket{
		id:          name,
		leakageRate: r,
		capacity:    cap,
		expiry:      sr.GetNow(),
	}
}

// Период капания в наносекундах
func (b *Bucket) PeriodNs() float64 {
	return sr.Sec64 / b.leakageRate
}

// Объем недокапанного запаса в каплях.
func (b *Bucket) DropsSum() int64 {
	if b.IsReset() {
		return 0
	}
	return int64(math.Ceil(float64(b.DeadLine()) / b.PeriodNs()))
}

// Скорость утечки в каплях в секунду.
func (b *Bucket) GetLeakageRate() float64 {
	return b.leakageRate
}

// Ёмкость в каплях.
func (b *Bucket) GetCapacity() int64 {
	return b.capacity
}

// Свободная емкость корзины в каплях.
func (b *Bucket) GetFreeCapacity() int64 {
	return b.capacity - b.DropsSum()
}

// Изменение ёмкости корзины.
func (b *Bucket) ChangeCapacity(cap int64) {
	if cap < b.capacity && b.DropsSum() > cap {
		b.expiry = sr.GetNow().Add(time.Duration(b.PeriodNs() * float64(cap)))
	}
	b.capacity = cap
}

// Промежуток времени до опустошения корзины.
func (b *Bucket) DeadLine() time.Duration {
	res := b.expiry.Sub(sr.GetNow())
	if res < 0 {
		return 0
	}
	return res
}

// Корзина должна быть обнулена?
func (b *Bucket) IsReset() bool {
	sr.GetNowMu.Lock()
	nw := sr.GetNow()
	sr.GetNowMu.Unlock()
	return !nw.Before(b.expiry)
}

func (b *Bucket) IsNilPointer(i any) bool {
	return i.(*Bucket) == nil
}

// Добавление содержимиго в корзину в каплях с учетом вместимости.
// Возвращает колличество добавленных капель (сколько влезло до полного)
// и произошло ли добавление вообще (колличесво капель должно быть > 0)
func (b *Bucket) AddDrops(addNum int64) (int64, bool) {
	freeCap := b.GetFreeCapacity()
	switch {
	case freeCap <= 0:
		return 0, false
	case addNum > freeCap:
		addNum = freeCap
	}
	if b.IsReset() {
		b.expiry = sr.GetNow()
	}
	period := b.PeriodNs()
	b.expiry = b.expiry.Add(time.Duration((float64(addNum) * period)))
	return addNum, true
}

func (b *Bucket) GetID() string {
	return b.id
}

func (b *Bucket) GetPriority() int {
	return b.priority
}

func (b *Bucket) GetAllBucketParams() *tp.AllBucketParams {
	params := tp.AllBucketParams{}
	params.Id = b.id
	params.Capacity = b.capacity
	params.FreeCapacity = b.GetFreeCapacity()
	params.DropsSum = b.DropsSum()
	params.LeakageRate = b.leakageRate
	params.Expiry = b.expiry
	params.DeadLine = b.DeadLine()
	params.QPriority = b.priority
	params.IsReset = b.IsReset()
	params.IsAllowed = true
	return &params
}
