package syncer

import (
	"time"

	"github.com/my-network/fsutil/pkg/file"
)

type task struct {
	Config Config
	Path   file.Path

	FirstEventTS time.Time
	LastEventTS  time.Time
	IsExpired    bool

	HeapIdx *int
}

func (t *task) Merge(addTask *task) {
	if t.IsExpired {
		return
	}
	if addTask.IsExpired {
		return
	}
	addTask.IsExpired = true

	if addTask.FirstEventTS.Before(t.FirstEventTS) {
		t.FirstEventTS = addTask.FirstEventTS
	}
	if addTask.LastEventTS.After(t.LastEventTS) {
		t.LastEventTS = addTask.LastEventTS
	}
}

func (t *task) ExpirationDeadline() time.Time {
	if t.IsExpired {
		return time.Time{}
	}

	deadline := t.LastEventTS.Add(t.Config.AggregationTimeMin)
	maxDeadline := t.FirstEventTS.Add(t.Config.AggregationTimeMax)
	if maxDeadline.Before(deadline) {
		deadline = maxDeadline
	}

	return deadline
}
