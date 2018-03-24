package logs

import "github.com/AzureRelease/boiler-server/models"

type BoilerRuntimeLog struct {
	models.MyIdObject

	Runtime				*models.BoilerRuntime	`orm:"rel(fk);null;index"`

	TableName			string					`orm:"index"`
	Query				string
	Status				int						`orm:"index"`
	Duration			float64
	DurationTotal		float64
}

func (bLog *BoilerRuntimeLog) GetLog() interface{} {
	return bLog
}

const (
	BOILER_RUNTIME_LOG_STATUS_INIT	= 0
	BOILER_RUNTIME_LOG_STATUS_READY	= 1
	BOILER_RUNTIME_LOG_STATUS_DONE	= 2
)