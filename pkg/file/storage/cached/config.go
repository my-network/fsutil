package cached

type Config struct {
	CacheDataDst     uint
	CacheMetadataDst uint
	KeepOpenDst      uint
}

func NewConfig(opts ...Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}
