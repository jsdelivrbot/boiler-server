package caches

import (
	"time"
	"github.com/AzureRelease/boiler-server/models"
)

type BoilerRuntimeCacheStatusRunning struct {
	models.MyIdObject

	Runtime			*models.BoilerRuntime		`orm:"rel(fk);null;index"`

	Boiler			*models.Boiler			`orm:"rel(fk);index"`

	Parameter		*models.RuntimeParameter	`orm:"rel(fk);index"`
	ParameterName		string
	Unit			string
	//Fix			int32
	Value			float64

	Alarm			*models.BoilerAlarm		`orm:"rel(fk);null;index"`
	AlarmLevel		int				`orm:"null;index"`
	AlarmDescription	string				`orm:"null;index"`

	Date			time.Time		`orm:"type(datetime);index"`
	Duration		int64
}

func (cache *BoilerRuntimeCacheStatusRunning) TableUnique() [][]string {
	return [][]string {
		//[]string{"Boiler", "CreatedDate"},
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}