package models

type RuntimeAlarmRule struct {
	MyUidObject

	Parameter			*RuntimeParameter		`orm:"rel(fk);index"`
	Organization		*Organization			`orm:"rel(fk);index"`

	BoilerForm			*BoilerTypeForm			`orm:"rel(fk);null;index"`
	BoilerMedium		*BoilerMedium			`orm:"rel(fk);null;index"`
	BoilerFuelType		*FuelType				`orm:"rel(fk);null;index"`
	BoilerCapacityMin	int32					`orm:"index"`
	BoilerCapacityMax	int32					`orm:"index"`

	Normal				float32					// 基准值
	Warning				float32					// 告警值
	Danger				float32					// 危险值

	Priority			int32					`orm:"index;default(0)"`
	Scope				int32					`orm:"index"` 				// 告警作用域：系统/企业

	NeedSend			bool					// 是否推送报警
	Delay				int64					// 延迟报警时间，单位分

	Description			string
}

const (
	RUNTIME_ALARM_LEVEL_UNDEFINED = -1
	RUNTIME_ALARM_LEVEL_NORMAL = 0
	RUNTIME_ALARM_LEVEL_WARNING = 1
	RUNTIME_ALARM_LEVEL_DANGER = 2
)

const (
	RUNTIME_ALARM_PRIORITY_LOW = 0
	RUNTIME_ALARM_PRIORITY_NORMAL = 1
	RUNTIME_ALARM_PRIORITY_HIGH = 2
)

const (
	RUNTIME_ALARM_SCOPE_GLOBAL = 0
	RUNTIME_ALARM_SCOPE_SYSTEM = 1
	RUNTIME_ALARM_SCOPE_ENTERPRISE = 2
)
