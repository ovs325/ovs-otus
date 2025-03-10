package buckets

import (
	"sync/atomic"
	"time"

	gr "bruteforce/internal/group"
)

// Секунда в актуальном формате
const Sec64 = float64(time.Second)

var GetNow = time.Now

// Стартовое время.
var Start = time.Date(1990, 1, 2, 0, 0, 0, 0, time.UTC).Round(0)
var Elapsed int64

// Получить  сдвиг времени.
func GetElapsed() time.Duration {
	return time.Duration(atomic.LoadInt64(&Elapsed))
}

// Установить сдвиг времени
func SetElapsed(v time.Duration) {
	atomic.StoreInt64(&Elapsed, int64(v))
}

// Увеличить сдвиг времени
func AddToElapsed(v time.Duration) {
	atomic.AddInt64(&Elapsed, int64(v))
}

// Сбросить  сдвиг времени
func reset(c *gr.BucketGroup) {
	c.Reset()
	SetElapsed(0)
}
