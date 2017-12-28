package caches

import "time"

//1201
type BoilerRuntimeCacheHeatDaily struct {
	BoilerRuntimeCache

	Date		time.Time		`orm:"type(datetime);index"`
}

func (cache *BoilerRuntimeCacheHeatDaily) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Date"},
		[]string{"Boiler", "Parameter", "Date"},
	}
}
