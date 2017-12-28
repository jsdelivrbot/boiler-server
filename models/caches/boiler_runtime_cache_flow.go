package caches

//1003
type BoilerRuntimeCacheFlow struct {
	BoilerRuntimeCache
}

func (cache *BoilerRuntimeCacheFlow) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "CreatedDate"},
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}