package buckets

import (
	"fmt"
	"testing"
	"time"

	sr "bruteforce/internal"
	"github.com/stretchr/testify/assert"
)

var (
	msgErrSetElp = "сдвиг времени установился не правильно"
	msgErrDeadL  = "DeadLine is bad"
	msgErrDropS  = "DropsSum is bad"
)

func TestNewLeakyBucket(t *testing.T) {
	sr.SetElapsed(0)
	sr.GetNowMu.Lock()
	sr.GetNow = func() time.Time { return sr.Start.Add(sr.GetElapsed()) }
	sr.GetNowMu.Unlock()
	rate := 1.0
	capacity := int64(5)
	bucket := NewBucket("test", rate, capacity).(*Bucket)

	assert.Equal(t, bucket.expiry, sr.GetNow(), "Bad 'expiry'")
	assert.Equal(t, bucket.leakageRate, rate, "Bad 'rate'")
	assert.Equal(t, bucket.GetLeakageRate(), rate, "Bad 'rate'")
	assert.Equal(t, bucket.capacity, capacity, "Bad 'capacity'")
	assert.Equal(t, bucket.GetCapacity(), capacity, "Bad 'capacity'")
	sr.GetNow = time.Now
}

func TestBucketOk(t *testing.T) {
	sr.SetElapsed(0)
	sr.GetNowMu.Lock()
	sr.GetNow = func() time.Time { return sr.Start.Add(sr.GetElapsed()) }
	sr.GetNowMu.Unlock()
	rate := 1.0
	capacity := int64(5)
	bucket := NewBucket("test", rate, capacity).(*Bucket)

	drops := int64(1)
	realAddDrops, ok := bucket.AddDrops(drops)
	assert.True(t, ok)
	assert.Equal(t, drops, realAddDrops)

	fmt.Printf("elapsed: %v - сдвиг времени\n", sr.GetElapsed())
	fmt.Printf("DeadLine: %v - Промежуток времени до опустошения корзины\n", bucket.DeadLine())
	fmt.Printf("DropsSum: %v - Объем недокапанного запаса в каплях\n", bucket.DropsSum())

	dropSum := bucket.DropsSum()
	assert.Equal(t, drops, dropSum, "Объем недокапанного запаса в каплях не совпадает с ожидаемым")
	assert.GreaterOrEqual(t, capacity, dropSum, "Объем недокапанного запаса в каплях больше емкости корзины")
	assert.Equal(t, capacity-drops, bucket.GetFreeCapacity(), "Объем недокапанного запаса в каплях больше емкости корзины")
	// Изменения текущего времени
	sr.SetElapsed(time.Nanosecond)                                  // Установить сдвиг времени
	assert.Equal(t, time.Nanosecond, sr.GetElapsed(), msgErrSetElp) //
	assert.Equal(
		t, time.Second-time.Nanosecond, bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	sr.SetElapsed(time.Second - time.Nanosecond)                                // Установить сдвиг времени
	assert.Equal(t, time.Second-time.Nanosecond, sr.GetElapsed(), msgErrSetElp) //
	assert.Equal(
		t, time.Nanosecond, bucket.DeadLine(), msgErrDeadL,
	) // Проверить Промежуток времени до опустошения корзины
	sr.SetElapsed(time.Second)                                  // Установить сдвиг времени
	assert.Equal(t, time.Second, sr.GetElapsed(), msgErrSetElp) //
	assert.Equal(
		t, time.Duration(0), bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	// Увеличение сдвига времени
	drops = int64(1)                                                // колличество добавляемых капель
	realAddDrops, ok = bucket.AddDrops(drops)                       // колличество реально добавленных капель
	assert.True(t, ok)                                              //
	assert.Equal(t, drops, realAddDrops)                            //
	sr.AddToElapsed(time.Second / 2)                                // Увеличить сдвиг времени
	assert.Equal(t, time.Second*3/2, sr.GetElapsed(), msgErrSetElp) //
	assert.Equal(
		t, time.Second/2, bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	realAddDrops, ok = bucket.AddDrops(drops)        // колличество реально добавленных капель
	assert.True(t, ok)                               //
	assert.Equal(t, drops, realAddDrops)             //
	sr.AddToElapsed(time.Second/2 - time.Nanosecond) // Увеличить сдвиг времени
	assert.Equal(t, time.Second*2-time.Nanosecond, sr.GetElapsed(), msgErrSetElp)
	assert.Equal(
		t, time.Second+time.Nanosecond, bucket.DeadLine(), msgErrDeadL,
	) // Проверить Промежуток времени до опустошения корзины
	// Изменение ёмкости корзины
	bucket.ChangeCapacity(int64(6)) //
	_ = bucket.GetCapacity()        //
	assert.Equal(
		t, time.Second+time.Nanosecond, bucket.DeadLine(), msgErrDeadL,
	) // Проверить Промежуток времени до опустошения корзины
	bucket.ChangeCapacity(int64(4)) //
	_ = bucket.GetCapacity()        //
	assert.Equal(
		t, time.Second+time.Nanosecond, bucket.DeadLine(), msgErrDeadL,
	) // Проверить Промежуток времени до опустошения корзины
	bucket.ChangeCapacity(int64(1)) //
	_ = bucket.GetCapacity()        //
	assert.Equal(
		t, time.Second, bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	bucket.ChangeCapacity(int64(4)) //
	capacity = bucket.GetCapacity() //
	assert.Equal(
		t, time.Second, bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	assert.Equal(t, int64(1), bucket.DropsSum(), msgErrDropS)                     //
	assert.Equal(t, time.Second*2-time.Nanosecond, sr.GetElapsed(), msgErrSetElp) //
	// Превышение свободного объема
	sr.AddToElapsed(time.Second * 5)                                              // Увеличить сдвиг времени
	assert.Equal(t, 7*time.Second-time.Nanosecond, sr.GetElapsed(), msgErrSetElp) //
	assert.Equal(
		t, int64(0), bucket.DropsSum(), msgErrDropS) //
	assert.Equal(
		t, time.Duration(0), bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	drops = int64(5)                          //
	realAddDrops, ok = bucket.AddDrops(drops) //
	assert.True(t, ok)                        //
	assert.NotEqual(t, drops, realAddDrops)   //
	assert.Equal(t, int64(4), realAddDrops)   //
	assert.Equal(
		t, time.Second*4, bucket.DeadLine(), msgErrDeadL) // Проверить Промежуток времени до опустошения корзины
	msg := "Колличество реально добавленных капель не равно свободному объему, но меньше количества добавлявшихся капель"
	assert.GreaterOrEqual(t, realAddDrops, capacity-bucket.DropsSum(), msg)
	sr.GetNow = time.Now
}

func TestOverFlow(t *testing.T) {
	sr.SetElapsed(0)
	sr.GetNowMu.Lock()
	sr.GetNow = func() time.Time { return sr.Start.Add(sr.GetElapsed()) }
	sr.GetNowMu.Unlock()

	rate := 60.0
	capacity := int64(1000)
	bucket := NewBucket("test", rate, capacity).(*Bucket)

	drops := int64(100)
	realAddDrops, ok := bucket.AddDrops(drops)
	assert.True(t, ok)
	assert.Equal(t, drops, realAddDrops)
	// Изменения текущего времени
	sr.SetElapsed(time.Nanosecond)
	assert.Equal(t, time.Nanosecond, sr.GetElapsed(), msgErrSetElp)
	drops = int64(1000)
	realAddDrops, ok = bucket.AddDrops(drops)
	assert.True(t, ok)
	assert.Equal(t, int64(900), realAddDrops)
	drops = int64(1)
	realAddDrops, ok = bucket.AddDrops(drops)
	assert.False(t, ok)
	assert.Equal(t, int64(0), realAddDrops)
	sr.SetElapsed(time.Second)
	assert.Equal(t, time.Second, sr.GetElapsed(), msgErrSetElp)
	assert.GreaterOrEqual(t, int64(940), capacity-bucket.DropsSum(), "")
	sr.GetNow = time.Now
}
