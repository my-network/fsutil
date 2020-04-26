package syncer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/rand/mathrand"
)

func TestTaskHeap(t *testing.T) {
	rand.Seed(0)

	h := taskHeap{}

	var removeTask, fixTask *task
	for i := 0; i < 900; i++ {
		task := &task{
			Config:       DefaultConfig,
			FirstEventTS: time.Unix(rand.Int63(), rand.Int63()),
			LastEventTS:  time.Unix(rand.Int63(), rand.Int63()),
		}
		h.Push(task)
		if i == 300 {
			removeTask = task
		}
		if i == 600 {
			fixTask = task
		}
	}

	h.Remove(removeTask)

	fixTask.FirstEventTS = time.Unix(rand.Int63(), rand.Int63())
	fixTask.LastEventTS = time.Unix(rand.Int63(), rand.Int63())
	h.Fix(fixTask)

	for i := 0; i < 101; i++ {
		task := &task{
			Config:       DefaultConfig,
			FirstEventTS: time.Unix(rand.Int63(), rand.Int63()),
			LastEventTS:  time.Unix(rand.Int63(), rand.Int63()),
		}
		h.Push(task)
	}

	var prevTask *task
	for i := 0; i < 1000; i++ {
		nextTask := h.Pop()
		if prevTask == nil {
			prevTask = nextTask
			continue
		}
		require.True(t, !nextTask.ExpirationDeadline().Before(prevTask.ExpirationDeadline()))
		prevTask = nextTask
	}
}

func BenchmarkTaskHeap_Fix(b *testing.B) {
	prng := mathrand.New()
	randTime := make([]time.Time, 0, b.N)
	for i := 0; i < 1000; i++ {
		randTime = append(randTime, time.Unix(int64(prng.Uint64AddRotateMultiply()>>1), 0))
	}

	for _, itemCount := range []int{10, 1000, 1000000} {
		h := taskHeap{}
		for i := 0; i < itemCount; i++ {
			task := &task{
				Config:       DefaultConfig,
				FirstEventTS: time.Unix(rand.Int63(), rand.Int63()),
				LastEventTS:  time.Unix(rand.Int63(), rand.Int63()),
			}
			h.Push(task)
		}
		b.Run(fmt.Sprintf("itemCount%d", itemCount), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				randIdx := mathrand.ReduceUint32(prng.Uint32AddRotateMultiply(), uint32(itemCount))
				task := (*h.int())[randIdx]
				task.FirstEventTS = randTime[i%500]
				task.LastEventTS = randTime[500+i%500]
				h.Fix(task)
			}
		})
	}
}
