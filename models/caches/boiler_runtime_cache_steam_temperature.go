package caches

//1001
type BoilerRuntimeCacheSteamTemperature struct {
	BoilerRuntimeCache
}

func (cache *BoilerRuntimeCacheSteamTemperature) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "CreatedDate"},
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}