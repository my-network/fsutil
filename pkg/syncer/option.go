package syncer

import (
	"time"
)

type Option interface {
	apply(*Config)
}

type OptionChecksum struct {
	Enable bool
}

func (opt OptionChecksum) apply(cfg *Config) {
	cfg.EnableChecksums = opt.Enable
}

type OptionSyncLogger struct {
	SyncLogger SyncLogger
}

func (opt OptionSyncLogger) apply(cfg *Config) {
	cfg.SyncLogger = opt.SyncLogger
}

type OptionAggregationTimeMin struct {
	Value time.Duration
}

func (opt OptionAggregationTimeMin) apply(cfg *Config) {
	cfg.AggregationTimeMin = opt.Value
}

type OptionAggregationTimeMax struct {
	Value time.Duration
}

func (opt OptionAggregationTimeMax) apply(cfg *Config) {
	cfg.AggregationTimeMax = opt.Value
}
