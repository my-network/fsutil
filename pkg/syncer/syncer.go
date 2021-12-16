package syncer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/my-network/fsutil/pkg/file"
	"github.com/my-network/fsutil/pkg/file/event"
)

type Syncer struct {
	config      Config
	ctx         context.Context
	src         file.Storage
	dst         file.Storage
	wg          sync.WaitGroup
	taskStorage *taskStorage
}

func NewSyncer(ctx context.Context, src, dst file.Storage, cfg *Config) (*Syncer, error) {
	if cfg == nil {
		cfg = &DefaultConfig
	}
	syncer := &Syncer{
		config: *cfg,
		ctx:    ctx,
		src:    src,
		dst:    dst,
	}
	err := syncer.init()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the syncer: %w", err)
	}
	return syncer, nil
}

func (syncer *Syncer) init() error {
	if err := syncer.config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	if syncer.config.AggregationTimeMin <= 0 {
		syncer.config.AggregationTimeMin = time.Nanosecond
	}
	if syncer.config.AggregationTimeMax < syncer.config.AggregationTimeMin {
		syncer.config.AggregationTimeMax = syncer.config.AggregationTimeMin
	}

	err := syncer.initTaskStorage()
	if err != nil {
		return fmt.Errorf("unable to initialize a task storage: %w", err)
	}

	err = syncer.initEventProcessor()
	if err != nil {
		return fmt.Errorf("unable to initialize an event processor: %w", err)
	}

	err = syncer.initCopier()
	if err != nil {
		return fmt.Errorf("unable to initialize a copier: %w", err)
	}

	return nil
}

func (syncer *Syncer) initTaskStorage() error {
	var err error
	syncer.taskStorage, err = newTaskStorage(syncer.config)
	if err != nil {
		return fmt.Errorf("unable to create a task storage: %w", err)
	}

	syncer.wg.Add(1)
	go func() {
		defer syncer.wg.Done()
		<-syncer.ctx.Done()
		_ = syncer.taskStorage.Close()
	}()
	return nil
}

func (syncer *Syncer) initCopier() error {
	syncer.wg.Add(1)
	go func() {
		defer syncer.wg.Done()
		syncer.copierLoop()
	}()

	return nil
}

func (syncer *Syncer) copierLoop() {
	for {
		select {
		case fileEvent := <-syncer.taskStorage.ExpiredChan:
			fmt.Println("SYNC", fileEvent)
		case <-syncer.ctx.Done():
			return
		}
	}
}

func (syncer *Syncer) Wait() {
	syncer.wg.Wait()
}

func (syncer *Syncer) Queue(path file.Path) error {
	now := time.Now()
	err := syncer.warmupForSync(path)
	if err != nil {
		return err
	}
	syncer.taskStorage.AddOrRefresh(path, now)
	return nil
}

func (syncer *Syncer) QueueRecursive(
	ctx context.Context,
	path file.Path,
	shouldWalkFn file.ShouldWalkFunc,
	errHandlerFn file.ErrorHandlerFunc,
) error {
	return file.Walk(
		ctx,
		syncer.src,
		nil,
		path,
		func(dir file.Directory, obj os.FileInfo) error {
			return syncer.Queue(dir.Path().Append(obj.Name()))
		},
		shouldWalkFn,
		errHandlerFn,
	)
}

type cachedStorage interface {
	OpenInBackground(
		path file.Path,
		mask file.OpenFlag,
		defaultPerm os.FileMode,
	)
}

func (syncer *Syncer) warmupForSync(path file.Path) error {
	dst, ok := syncer.dst.(cachedStorage)
	if !ok {
		return nil
	}
	obj, err := syncer.src.Open(syncer.ctx, nil, path, file.FlagRead, 0000)
	if err != nil {
		if file.IsNotExist(err) {
			syncer.config.SyncLogger.Debugf("file '%s' disappeared, skipping",
				path.LocalPath())
			return nil
		}
		return fmt.Errorf("unable to open src file '%s': %w",
			path.LocalPath(), err)
	}
	if obj.LastStat().Mode().IsDir() {
		syncer.config.SyncLogger.Debugf("'%s' is directory, skipping",
			path.LocalPath())
		return nil
	}
	dst.OpenInBackground(path, file.FlagReadWrite|file.FlagCreate|file.FlagAppend, obj.LastStat().Mode().Perm())
	return nil
}

func (syncer *Syncer) eventProcessorLoop(inChan chan event.Event) {
	for {
		select {
		case fileEvent := <-inChan:
			err := syncer.Queue(fileEvent.Path)
			if err != nil {
				panic(err)
			}
		case <-syncer.ctx.Done():
			return
		}
	}
}

func (syncer *Syncer) initEventProcessor() error {
	syncer.
		syncer.wg.Add(1)
	go func() {
		defer syncer.wg.Done()
		syncer.eventProcessorLoop()
	}()
}
