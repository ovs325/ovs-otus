package buckets

import (
	"sync"
	"sync/atomic"
	"time"
)

// Секунда в актуальном формате.
const Sec64 = float64(time.Second)

var (
	GetNowMu sync.RWMutex
	GetNow   = time.Now
	// Стартовое время.
	Start   = time.Date(1990, 1, 2, 0, 0, 0, 0, time.UTC).Round(0)
	Elapsed int64
)

// Получить  сдвиг времени.
func GetElapsed() time.Duration {
	return time.Duration(atomic.LoadInt64(&Elapsed))
}

// Установить сдвиг времени.
func SetElapsed(v time.Duration) {
	atomic.StoreInt64(&Elapsed, int64(v))
}

// Увеличить сдвиг времени.
func AddToElapsed(v time.Duration) {
	atomic.AddInt64(&Elapsed, int64(v))
}
