package caches

//1201
type BoilerRuntimeCacheHeat struct {
	BoilerRuntimeCache
}

func (cache *BoilerRuntimeCacheHeat) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "CreatedDate"},
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}