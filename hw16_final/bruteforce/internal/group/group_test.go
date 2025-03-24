package group_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	tm "bruteforce/internal"
	bk "bruteforce/internal/bucket"
	gr "bruteforce/internal/group"
	"github.com/stretchr/testify/assert"
)

var (
	msgErrSetElp = "сдвиг времени установился не правильно"
	msgErrDeadL  = "DeadLine is bad"
)

func TestNewGroup(t *testing.T) {
	rate := 1.0
	capacity := int64(2)
	capMap := 4096
	queye := bk.NewQueyeBuckets(capMap)
	group := gr.NewBucketGroup(rate, capacity, int64(capMap), queye, true, bk.NewBucket, bk.NewQueyeBuckets)

	assert.False(t, group.IsNilBuckets(), "Backet's Map is Nil")
	assert.False(t, group.IsNilQueue(), "Queue is Nil")
	assert.True(t, group.IsLeakageRateEqRate(rate), "Bad 'rate'")
	assert.Equal(t, group.GetLeakageRate(), rate, "Bad 'rate'")

	group.ExistScheduller()
}

func TestGroupSimpleOk(t *testing.T) {
	tm.GetNowMu.Lock()
	tm.GetNow = func() time.Time { return tm.Start.Add(tm.GetElapsed()) }
	tm.GetNowMu.Unlock()

	rate := 1.0
	capacity := int64(5)
	id := "127.0.0.1"
	capMap := 4096
	bucket := bk.NewQueyeBuckets(capMap)
	group := gr.NewBucketGroup(rate, capacity, int64(capMap), bucket, false, bk.NewBucket, bk.NewQueyeBuckets)

	simple := func() {
		tm.SetElapsed(0)
		capacity = int64(5)
		drops := int64(1)
		assert.Equal(t, drops, group.AddDrops(id, drops))                             // "add"
		tm.SetElapsed(time.Nanosecond)                                                // "time-set"
		assert.Equal(t, time.Nanosecond, tm.GetElapsed(), msgErrSetElp)               // "time-set"
		deadLine, ok := group.GetDeadLine(id)                                         // "till"
		assert.True(t, ok)                                                            // "till"
		assert.Equal(t, time.Second-time.Nanosecond, deadLine, msgErrDeadL)           // "till"
		tm.SetElapsed(time.Second - time.Nanosecond)                                  // "time-set"
		assert.Equal(t, time.Second-time.Nanosecond, tm.GetElapsed(), msgErrSetElp)   // "time-set"
		deadLine, ok = group.GetDeadLine(id)                                          // "till"
		assert.True(t, ok)                                                            // "till"
		assert.Equal(t, time.Nanosecond, deadLine, msgErrDeadL)                       // "till"
		tm.SetElapsed(time.Second)                                                    // "time-set"
		assert.Equal(t, time.Second, tm.GetElapsed(), msgErrSetElp)                   // "time-set"
		deadLine, ok = group.GetDeadLine(id)                                          // "till"
		assert.True(t, ok)                                                            // "till"
		assert.Equal(t, time.Duration(0), deadLine, msgErrDeadL)                      // "till"
		assert.Equal(t, drops, group.AddDrops(id, drops))                             // "add"
		tm.AddToElapsed(time.Second / 2)                                              // "time-add"
		assert.Equal(t, time.Second*3/2, tm.GetElapsed(), msgErrSetElp)               // "time-add"
		deadLine, ok = group.GetDeadLine(id)                                          // "till"
		assert.True(t, ok)                                                            // "till"
		assert.Equal(t, time.Second/2, deadLine, msgErrDeadL)                         // "till"
		assert.Equal(t, drops, group.AddDrops(id, drops))                             // "add"
		tm.AddToElapsed(time.Second/2 - time.Nanosecond)                              // "time-add"
		assert.Equal(t, time.Second*2-time.Nanosecond, tm.GetElapsed(), msgErrSetElp) // "time-add"
		tm.AddToElapsed(time.Second * 5)                                              // "time-add"
		assert.Equal(t, 7*time.Second-time.Nanosecond, tm.GetElapsed(), msgErrSetElp) // "time-add"
		drops = int64(6)                                                              // "add"
		realAddDrops := group.AddDrops(id, drops)                                     // "add"
		assert.NotEqual(t, drops, realAddDrops)                                       // "add"
		assert.Equal(t, int64(5), realAddDrops)                                       // "add"
		deadLine, ok = group.GetDeadLine(id)                                          // "till"
		assert.True(t, ok)                                                            // "till"
		assert.Equal(t, time.Second*5, deadLine, msgErrDeadL)                         // "till"
	}
	simple()

	group.RemoveDrops(id)
	dropSum, ok := group.GetDropsSum(id)
	assert.False(t, ok)
	assert.LessOrEqual(t, dropSum, int64(0), "ID still has a count after RemoveDrops()?!")

	// Тестируем RemooveEmptyDrops() - удаление пустых корзин
	simple()
	group.RemooveEmpty()
	tm.SetElapsed(time.Hour)
	group.RemooveEmpty()

	// Тестируем Reset().
	simple()
	group.Reset()
	dropSum, ok = group.GetDropsSum(id)
	assert.False(t, ok)
	assert.LessOrEqual(t, dropSum, int64(0), "ID still has a count after RemoveDrops()?!")
	deadLine, ok := group.GetDeadLine(id)
	assert.False(t, ok)
	assert.Zero(t, deadLine, "ID still has a count after RemoveDrops()?!")

	// Удаление 'левой' корзины.
	group.RemoveDrops("fake")
	dropSum, ok = group.GetDropsSum("fake")
	assert.False(t, ok)
	assert.LessOrEqual(t, dropSum, int64(0), "'fake'-корзина существовала?!")

	// group.ExistScheduller()
}

func TestGroupVariedOk(t *testing.T) {
	tm.GetNowMu.Lock()
	tm.GetNow = func() time.Time { return tm.Start.Add(tm.GetElapsed()) }
	tm.GetNowMu.Unlock()
	rate := 60.0
	capacity := int64(1000)
	id := "127.0.0.1"
	capMap := 4096
	bucket := bk.NewQueyeBuckets(capMap)
	group := gr.NewBucketGroup(rate, capacity, int64(capMap), bucket, false, bk.NewBucket, bk.NewQueyeBuckets)

	varied := func() {
		tm.SetElapsed(0)
		capacity = int64(1000)
		drops := int64(100)
		assert.Equal(t, drops, group.AddDrops(id, drops))               // "add"
		tm.SetElapsed(time.Nanosecond)                                  // "time-set"
		assert.Equal(t, time.Nanosecond, tm.GetElapsed(), msgErrSetElp) // "time-set"
		assert.Equal(t, int64(900), group.AddDrops(id, int64(1000)))    // "add"
		assert.Equal(t, int64(0), group.AddDrops(id, int64(1)))         // "add
		tm.SetElapsed(time.Second)                                      // "time-set"
		assert.Equal(t, time.Second, tm.GetElapsed(), msgErrSetElp)     // "time-set"
		dropSum, ok := group.GetDropsSum(id)
		assert.True(t, ok)
		assert.Equal(t, int64(940), dropSum)
	}
	varied()

	group.RemoveDrops(id)
	dropSum, ok := group.GetDropsSum(id)
	assert.False(t, ok)
	assert.LessOrEqual(t, dropSum, int64(0), "ID still has a count after RemoveDrops()?!")

	// Тестируем RemooveEmptyDrops() - удаление пустых корзин
	varied()
	group.RemooveEmpty()
	tm.SetElapsed(time.Hour)
	group.RemooveEmpty()

	// Тестируем Reset().
	varied()
	group.Reset()
	dropSum, ok = group.GetDropsSum(id)
	assert.False(t, ok)
	assert.LessOrEqual(t, dropSum, int64(0), "ID still has a count after RemoveDrops()?!")
	deadLine, ok := group.GetDeadLine(id)
	assert.False(t, ok)
	assert.Zero(t, deadLine, "ID still has a count after RemoveDrops()?!")

	// Удаление 'левой' корзины.
	group.RemoveDrops("fake")
	dropSum, ok = group.GetDropsSum("fake")
	assert.False(t, ok)
	assert.LessOrEqual(t, dropSum, int64(0), "'fake'-корзина существовала?!")

	// group.ExistScheduller()
}

func TestPeriodicPrune(t *testing.T) {
	tm.SetElapsed(0)
	id := "localhost"

	capMap := 4096
	bucket := bk.NewQueyeBuckets(capMap)
	group := gr.NewBucketGroup(1e7, 8, int64(capMap), bucket, false, bk.NewBucket, bk.NewQueyeBuckets)
	group.Scheduller(time.Microsecond)
	assert.Equal(t, int64(8), group.AddDrops(id, 100), "Didn't fill bucket?!")
	deadLine, _ := group.GetDeadLine(id)
	fmt.Printf("DeadLine(): %v\n", deadLine)

	// Wait for the periodic prune.
	time.Sleep(time.Millisecond)
	tm.SetElapsed(time.Millisecond)

	dropSum, _ := group.GetDropsSum(id)
	assert.Zerof(t, dropSum, "Key's bucket is not empty: %d?!", dropSum)

	group.ExistScheduller()
}

func TestMain(m *testing.M) {
	// Выдает значение стартового времени + сдвиг Elapsed
	tm.GetNowMu.Lock()
	tm.GetNow = func() time.Time { return tm.Start.Add(tm.GetElapsed()) }
	tm.GetNowMu.Unlock()

	os.Exit(m.Run())
}
