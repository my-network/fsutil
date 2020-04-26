package syncer

import (
	"fmt"
	"sync"
	"time"

	"github.com/my-network/fsutil/pkg/file"
)

const (
	debug = false
)

type taskStorage struct {
	config Config

	taskMap              map[string]*task
	taskWaitHeap         taskHeap
	ExpiredChan          chan *task
	taskAddOrRefreshChan chan *task
	waitingTask          *task
	wg                   sync.WaitGroup
}

func newTaskStorage(cfg Config) (*taskStorage, error) {
	storage := &taskStorage{}
	if err := storage.init(cfg); err != nil {
		return nil, err
	}
	return storage, nil
}

func (storage *taskStorage) init(cfg Config) error {
	storage.initFields(cfg)
	storage.initTaskScheduler()
	return nil
}

func (storage *taskStorage) initFields(cfg Config) {
	storage.config = cfg
	storage.taskMap = map[string]*task{}
}

func (storage *taskStorage) initTaskScheduler() {
	if storage.taskAddOrRefreshChan == nil { // storage.taskAddOrRefreshChan could also be set by unit-tests
		storage.taskAddOrRefreshChan = make(chan *task, 100)
	}
	if storage.ExpiredChan == nil { // storage.ExpiredChan could also be set by unit-tests
		storage.ExpiredChan = make(chan *task, 100)
	}

	storage.wg.Add(1)
	go func() {
		defer storage.wg.Done()
		defer func() {
			close(storage.ExpiredChan)
		}()

		storage.taskSchedulerLoop()
	}()
}

func (storage *taskStorage) taskSchedulerLoop() {
	for {
		var waitChan <-chan time.Time

		if storage.waitingTask != nil {
			deadline := storage.waitingTask.ExpirationDeadline()
			waitChan = time.After(time.Until(deadline))
		} else {
			waitChan = nil // a nil channel will never fire
		}

		select {
		case task, ok := <-storage.taskAddOrRefreshChan:
			if !ok {
				return
			}
			storage.addOrRefresh(task)
			if debug {
				if len(storage.taskMap)-1 != storage.taskWaitHeap.Len() {
					panic(fmt.Sprintf("%d %d", len(storage.taskMap), storage.taskWaitHeap.Len()))
				}
			}

		case <-waitChan:
			storage.waitingTask.IsExpired = true
			storage.ExpiredChan <- storage.waitingTask
			delete(storage.taskMap, storage.waitingTask.Path.Key())
			storage.waitingTask = nil

			if storage.taskWaitHeap.Len() > 0 {
				storage.waitingTask = storage.taskWaitHeap.Pop()
			}
		}
	}
}

func (storage *taskStorage) AddOrRefresh(path file.Path, touchTime time.Time) {
	task := &task{
		Config:       storage.config,
		FirstEventTS: touchTime,
		LastEventTS:  touchTime,
	}
	task.Path = make(file.Path, len(path))
	copy(task.Path, path)
	storage.taskAddOrRefreshChan <- task
}

func (storage *taskStorage) addOrRefresh(task *task) {
	oldTask := storage.taskMap[task.Path.Key()]
	if oldTask != nil {
		oldTask.Merge(task)
		task = oldTask

		if task == storage.waitingTask {
			// storage.waitingTask is not actual anymore (as the waited task):
			storage.taskWaitHeap.Push(storage.waitingTask)
			storage.waitingTask = nil
		} else {
			storage.taskWaitHeap.Fix(task)
		}
	} else {
		storage.taskMap[task.Path.Key()] = task
		storage.taskWaitHeap.Push(task)
	}

	earliestTask := storage.taskWaitHeap.Pop()
	if storage.waitingTask != nil &&
		!earliestTask.ExpirationDeadline().Before(storage.waitingTask.ExpirationDeadline()) {

		storage.taskWaitHeap.Push(earliestTask)
		return
	}

	// we found a task which expires earlier, let's wait for it
	// instead of currently waited task:

	if storage.waitingTask != nil {
		storage.taskWaitHeap.Push(storage.waitingTask)
	}
	storage.waitingTask = earliestTask
}

func (storage *taskStorage) Close() error {
	close(storage.taskAddOrRefreshChan)
	storage.wg.Wait()
	return nil
}
