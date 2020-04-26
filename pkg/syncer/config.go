package syncer

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

var (
	DefaultConfig = Config{
		SyncLogger:         zap.NewNop().Sugar(),
		EnableChecksums:    false,
		AggregationTimeMin: time.Second,
		AggregationTimeMax: time.Second * 10,
	}
)

type SyncLogger interface {
	Debugf(fmt string, args ...interface{})
}

type Config struct {
	SyncLogger         SyncLogger
	EnableChecksums    bool
	AggregationTimeMin time.Duration
	AggregationTimeMax time.Duration
}

func NewConfig(opts ...Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

func (cfg Config) Validate() error {
	if cfg.AggregationTimeMax < cfg.AggregationTimeMin {
		return fmt.Errorf("cfg.AggregationTimeMax (%v) < cfg.AggregationTimeMin (%v)",
			cfg.AggregationTimeMax, cfg.AggregationTimeMin)
	}
	return nil
}
