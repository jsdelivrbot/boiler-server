package caches

//1002
type BoilerRuntimeCacheSteamPressure struct {
	BoilerRuntimeCache
}

func (cache *BoilerRuntimeCacheSteamPressure) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "CreatedDate"},
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}