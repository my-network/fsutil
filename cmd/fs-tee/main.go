package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/my-network/fsutil/pkg/file"
	"github.com/my-network/fsutil/pkg/file/cached"
	"github.com/my-network/fsutil/pkg/file/localfs"
	"github.com/my-network/fsutil/pkg/syncer"
)

func syntaxExit() {
	flag.Usage()
	os.Exit(int(syscall.EINVAL))
}

func assertNoError(err error) {
	if err == nil {
		return
	}

	log.Panic(err)
}

func walkErrorHandler(err error) error {
	switch err := err.(type) {
	case file.ErrWalkNotDir:
		return nil
	case file.ErrWalkOpen:
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func watchErrorHandler(err error) error {
	switch err := err.(type) {
	case file.ErrWatchMark:
		if os.IsNotExist(err) {
			return nil
		}
	default:
		return walkErrorHandler(err)
	}
	return err
}

func main() {
	profile := flag.String("profile", "", "enable a profile: \"huge-latency-on-dst\""+
		" (effectively: -checksum -cache-data-dst=1000000 -cache-metadata-dst=1000000 -keep-open-dst=1000)")
	skipInitialSync := flag.Bool("skip-initial-sync", false, "do not start re-syncing everything on start")
	aggregationTimeMin := flag.String("aggregation-time-min", "1s",
		`minimal time to wait for more events on a file`)
	aggregationTimeMax := flag.String("aggregation-time-max", "30s",
		`miximal time to wait for more events on a file`)
	checksum := flag.Bool("checksum", false,
		`enable checking the checksum of files and do not sync if file has not changed`)
	cacheDataDst := flag.Uint("cache-data-dst", 0,
		`cache checksum of the destination data to avoid extra copyings for the specified amount of files.`+
			`The destination data should not be changed bypass the fs-tee instance!`)
	cacheMetadataDst := flag.Uint("cache-metadata-dst", 0,
		`cache file metadata of the destination to avoid extra scannings and copyings for the specified amount of files/directories. `+
			`The destination data should not be changed bypass the fs-tee instance!`)
	keepOpenDst := flag.Uint("keep-open-dst", 0,
		`keep files of the destination opened to avoid extra syscalls (open()/close()) for the specified amount of files.`+
			`The destination data should not be changed bypass the fs-tee instance!`)
	flag.Parse()

	if flag.NArg() != 2 {
		syntaxExit()
	}

	var syncerOpts []syncer.Option

	switch *profile {
	case "":
	case "huge-latency-on-dst":
		*checksum = true
		*cacheDataDst = 1000000
		*cacheMetadataDst = 1000000
		*keepOpenDst = 1000
	}

	if *aggregationTimeMin != "" {
		duration, err := time.ParseDuration(*aggregationTimeMin)
		assertNoError(err)
		syncerOpts = append(syncerOpts, syncer.OptionAggregationTimeMin{Value: duration})
	}

	if *aggregationTimeMax != "" {
		duration, err := time.ParseDuration(*aggregationTimeMax)
		assertNoError(err)
		syncerOpts = append(syncerOpts, syncer.OptionAggregationTimeMax{Value: duration})
	}

	if *checksum {
		syncerOpts = append(syncerOpts, syncer.OptionChecksum{Enable: true})
	}

	var dstStorageOpts []cached.Option

	if *cacheDataDst > 0 {
		dstStorageOpts = append(dstStorageOpts, cached.OptionCacheDataDst{AmountOfFiles: *cacheDataDst})
	}

	if *cacheMetadataDst > 0 {
		dstStorageOpts = append(dstStorageOpts, cached.OptionCacheMetadataDst{AmountOfFiles: *cacheMetadataDst})
	}

	if *keepOpenDst > 0 {
		dstStorageOpts = append(dstStorageOpts, cached.OptionKeepOpenDst{AmountOfFiles: *keepOpenDst})
	}

	pathSrc := flag.Arg(0)
	pathDst := flag.Arg(1)

	srcStorage := localfs.NewStorage(pathSrc)

	dstStorageBackend := localfs.NewStorage(pathDst)
	dstStorage := cached.NewStorage(dstStorageBackend, dstStorageOpts...)

	syncerCfg := syncer.NewConfig(syncerOpts...)
	assertNoError(syncerCfg.Validate())

	ctx := context.Background()

	syncerInstance, err := syncer.NewSyncer(ctx, srcStorage, dstStorage, syncerCfg)
	assertNoError(err)

	eventEmitter, err := srcStorage.Watch(nil, nil, nil, nil, watchErrorHandler)
	assertNoError(err)

	go func() {
		for {
			select {
			case ev := <-eventEmitter.C():
				fmt.Println("EVENT", ev.Path().LocalPath(), ev.Timestamp(), ev.TypeMask())
				fileInfo, err := srcStorage.Stat(ctx, nil, ev.Path(), true)
				assertNoError(err)
				if fileInfo.IsDir() {
					err := eventEmitter.Watch(nil, ev.Path(), nil, nil, watchErrorHandler)
					if !os.IsNotExist(err) {
						assertNoError(err)
					}

					err = syncerInstance.QueueRecursive(ctx, ev.Path(), nil, walkErrorHandler)
					assertNoError(err)
				} else {
					err = syncerInstance.Queue(ev.Path())
					assertNoError(err)
				}
			}
		}
	}()

	if !*skipInitialSync {
		err := syncerInstance.QueueRecursive(ctx, nil, nil, walkErrorHandler)
		assertNoError(err)
	}

	syncerInstance.Wait()
}
