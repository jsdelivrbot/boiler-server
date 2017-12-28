package caches

import "time"

//1003
type BoilerRuntimeCacheFlowDaily struct {
	BoilerRuntimeCache

	Date		time.Time		`orm:"type(datetime);index"`
	Hours		int32
}

func (cache *BoilerRuntimeCacheFlowDaily) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Date"},
		[]string{"Boiler", "Parameter", "Date"},
	}
}