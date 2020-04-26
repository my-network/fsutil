package syncer

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/golang/go/src/encoding/base64"
	"github.com/my-network/fsutil/pkg/file"
	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/rand/mathrand"
)

type testTask struct {
	Path      file.Path
	TouchTime []time.Time
}

func readableKey(path file.Path) string {
	return base64.StdEncoding.EncodeToString([]byte(path.Key()))
}

func TestTaskStorageScheduler(t *testing.T) {
	rand.Seed(0)
	stor := &taskStorage{}
	stor.initFields(DefaultConfig)
	stor.ExpiredChan = make(chan *task, 1000)
	stor.taskAddOrRefreshChan = make(chan *task, 1000)

	tasks := make([]*testTask, 100)
	for idx := range tasks {
		task := &testTask{}
		for i := 0; i < rand.Intn(4); i++ {
			b := make([]byte, rand.Intn(10))
			_, err := rand.Read(b)
			require.NoError(t, err)
			task.Path = append(task.Path, string(b))
		}
		tasks[idx] = task
	}

	now := time.Now()

	touchedPaths := map[string]struct{}{}
	for i := 0; i < 100; i++ {
		randIdx := rand.Intn(len(tasks))
		task := tasks[randIdx]
		touchTime := time.Unix(now.Unix()-1, rand.Int63()%1000000000) // already expired time to get results instantly
		task.TouchTime = append(task.TouchTime, touchTime)
		if debug {
			fmt.Printf("%v %s\n", touchTime, readableKey(task.Path))
		}
		stor.AddOrRefresh(task.Path, touchTime)
		touchedPaths[task.Path.Key()] = struct{}{}
	}

	stor.initTaskScheduler()

	expiredTasks := make([]*task, 0, len(touchedPaths))
	for len(expiredTasks) < cap(expiredTasks) {
		expiredTask := <-stor.ExpiredChan
		if debug {
			fmt.Printf("%3d %3d %v %s\n", len(expiredTasks), cap(expiredTasks),
				expiredTask.LastEventTS, readableKey(expiredTask.Path))
		}
		expiredTasks = append(expiredTasks, expiredTask)
	}

	expiredPathMap := map[string]struct{}{}
	var prevExpiredTask *task
	for _, expiredTask := range expiredTasks {
		_, ok := expiredPathMap[expiredTask.Path.Key()]
		require.False(t, ok, readableKey(expiredTask.Path))
		expiredPathMap[expiredTask.Path.Key()] = struct{}{}
		if prevExpiredTask == nil {
			prevExpiredTask = expiredTask
			continue
		}
		require.False(t, prevExpiredTask.ExpirationDeadline().After(expiredTask.ExpirationDeadline()),
			fmt.Sprintf("%v %v > %v %v",
				prevExpiredTask.ExpirationDeadline(), readableKey(prevExpiredTask.Path),
				expiredTask.ExpirationDeadline(), readableKey(expiredTask.Path)))
	}

	err := stor.Close()
	require.NoError(t, err)
	expiredTask, ok := <-stor.ExpiredChan
	require.False(t, ok, fmt.Sprintf("got unexpected task: %v (also %d tasks left; total: %d, should be: %d)",
		expiredTask, len(stor.ExpiredChan), 1+len(stor.ExpiredChan)+len(expiredTasks), len(touchedPaths)))
}

func BenchmarkTaskStorageScheduler(b *testing.B) {
	prng := mathrand.New()

	stor, err := newTaskStorage(DefaultConfig)
	if err != nil {
		panic(err)
	}
	now := time.Now()

	paths := make([]file.Path, 256)
	for idx := range paths {
		path := file.Path{}
		for i := 0; i < rand.Intn(4); i++ {
			b := make([]byte, rand.Intn(10))
			_, _ = rand.Read(b)
			path = append(path, string(b))
		}
		paths[idx] = path
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, ok := <-stor.ExpiredChan
			if !ok {
				return
			}
		}
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		touchTime := time.Unix(now.Unix()-1, int64(prng.Uint32AddRotateMultiply()%1000000000))
		randIdx := prng.Uint32AddRotateMultiply() & 0xff
		stor.AddOrRefresh(paths[randIdx], touchTime)
	}
	_ = stor.Close()
	wg.Wait()
}
