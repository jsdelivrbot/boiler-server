package caches

import (
	"github.com/AzureRelease/boiler-server/models"
)

type BoilerRuntimeCache struct {
	models.MyIdObject

	Runtime				*models.BoilerRuntime		`orm:"rel(fk);null;index"`

	Boiler				*models.Boiler				`orm:"rel(fk);index"`

	Parameter			*models.RuntimeParameter	`orm:"rel(fk);index"`
	ParameterName		string
	Unit				string
	//Fix				int32
	Value				float64

	Alarm				*models.BoilerAlarm			`orm:"rel(fk);null;index"`
	AlarmLevel			int							`orm:"null;index"`
	AlarmDescription	string						`orm:"null;index"`
}

func (cache *BoilerRuntimeCache) GetCache() interface{} {
	return cache
}

func (cache *BoilerRuntimeCache) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}
