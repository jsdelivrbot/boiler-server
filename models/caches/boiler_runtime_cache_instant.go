package caches

type BoilerRuntimeCacheInstant struct {
	BoilerRuntimeCache

	IsValid			bool
}

func (cache *BoilerRuntimeCacheInstant) TableIndex() [][]string {
	return [][]string {
		[]string{"Boiler", "Parameter", "UpdatedDate"},
	}
}

func (cache *BoilerRuntimeCacheInstant) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Parameter"},
	}
}