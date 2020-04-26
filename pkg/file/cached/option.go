package cached

type Option interface {
	apply(*Config)
}

type OptionCacheDataDst struct {
	AmountOfFiles uint
}

func (opt OptionCacheDataDst) apply(cfg *Config) {
	cfg.CacheDataDst = opt.AmountOfFiles
}

type OptionCacheMetadataDst struct {
	AmountOfFiles uint
}

func (opt OptionCacheMetadataDst) apply(cfg *Config) {
	cfg.CacheMetadataDst = opt.AmountOfFiles
}

type OptionKeepOpenDst struct {
	AmountOfFiles uint
}

func (opt OptionKeepOpenDst) apply(cfg *Config) {
	cfg.KeepOpenDst = opt.AmountOfFiles
}
