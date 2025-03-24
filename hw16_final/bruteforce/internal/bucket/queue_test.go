package buckets

import (
	"testing"
	"time"

	sr "bruteforce/internal"
	"github.com/stretchr/testify/assert"
)

var (
	maxLoop  = 5
	lenQueue = 4096
)

func TestPush(t *testing.T) {
	queue := make(queyeBuckets, 0, lenQueue)
	for i := 0; i < maxLoop; i++ {
		bucket := NewBucket("test", 1.0, 5).(*Bucket) // Заполняем очередь
		queue.Push(bucket)
		assert.Equal(t, queue[len(queue)-1], bucket, "Push должен добавлять корзины в конец очереди")
	}
}

func TestPop(t *testing.T) {
	queue := make(queyeBuckets, 0, lenQueue)
	for i := 1; i <= maxLoop; i++ {
		queue.Push(NewBucket("test", 1.0, 5).(*Bucket)) // Заполняем очередь
	}
	for i := 1; i <= maxLoop; i++ {
		bucket := queue[len(queue)-1]
		assert.Equal(t, bucket, queue.Pop(), "Pop должен удалять из конца очереди")
	}
}

func TestLen(t *testing.T) {
	queue := make(queyeBuckets, 0, lenQueue)
	assert.Zero(t, queue.Len(), "Queue не пуста!")
	for i := 1; i <= maxLoop; i++ {
		queue.Push(NewBucket("test", 1.0, 5).(*Bucket)) // Заполняем очередь
		assert.Equal(t, i, queue.Len(), "Ожидалось length %d, получено %d", i, queue.Len())
	}
	for i := 4; i >= 0; i-- {
		queue.Pop()
		assert.Equal(t, i, queue.Len(), "Ожидалось length %d, получено %d", i, queue.Len())
	}
}

func TestLess(t *testing.T) {
	queue := make(queyeBuckets, 0, lenQueue)

	for i := 0; i < maxLoop; i++ {
		b := NewBucket("test", 1.0, 5).(*Bucket)
		b.expiry = sr.GetNow().Add(time.Duration(i))
		queue.Push(b)
	}
	for i, j := 0, 4; i < maxLoop; i, j = i+1, j-1 {
		assert.False(t, i < j && !queue.Less(i, j), "Больше перепутано с меньше!")
	}
}

func TestSwap(t *testing.T) {
	queue := make(queyeBuckets, 0, lenQueue)
	for i := 0; i < maxLoop; i++ {
		queue.Push(NewBucket("test", 1.0, 5).(*Bucket))
	}
	i, j := 2, 4
	bucketI, bucketJ := queue[i], queue[j]
	queue.Swap(i, j)
	assert.False(t, bucketI != queue[j] || bucketJ != queue[i], "Элементы не были поменяны местами")
}

func TestReset(t *testing.T) {
	queue := make(queyeBuckets, 0, lenQueue)
	for i := 1; i <= maxLoop; i++ {
		queue.Push(NewBucket("test", 1.0, 5).(*Bucket)) // Заполняем очередь
	}
	assert.Equal(t, maxLoop, queue.Len())
	queue.Reset(lenQueue)
	assert.Zero(t, queue.Len(), "Queue не пуста!")
}
