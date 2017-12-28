package models

import "time"

type BoilerAlarmHistory struct {
	MyUidObject

	Boiler			*Boiler				`orm:"rel(fk);index"`
	Parameter		*RuntimeParameter	`orm:"rel(fk);index"`

	StartDate		time.Time			`orm:"auto_now_add;type(datetime);index"`
	EndDate			time.Time			`orm:"null;type(datetime);index"`

	ConfirmedDate	time.Time			`orm:"null;type(datetime);index"`
	ConfirmedBy		*User				`orm:"rel(fk);null;index"`
	VerifiedDate	time.Time			`orm:"null;type(datetime);index"`
	VerifiedBy		*User				`orm:"rel(fk);null;index"`

	TriggerRule		*RuntimeAlarmRule	`orm:"rel(fk);null;index"`
	AlarmLevel		int32				`orm:"index"`
	Priority		int32				`orm:"index;default(0)"`

	Description		string
}

func (alarm *BoilerAlarmHistory) TableIndex() [][]string {
	return [][]string {
		[]string {"Boiler", "Parameter", "TriggerRule"},
	}
}

func (alarm *BoilerAlarmHistory) TableUnique() [][]string {
	return [][]string {
		[]string {"Boiler", "Parameter", "TriggerRule", "StartDate"},
	}
}
