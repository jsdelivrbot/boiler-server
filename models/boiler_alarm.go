package models

import (
	//"time"
)
import "time"

type BoilerAlarm struct {
	MyUidObject

	Boiler			*Boiler				`orm:"rel(fk);index"`
	Parameter		*RuntimeParameter	`orm:"rel(fk);index"`
	Runtime			[]*BoilerRuntime	`orm:"reverse(many)"`

	StartDate		time.Time			`orm:"auto_now_add;type(datetime);index"`
	EndDate			time.Time			`orm:"null;type(datetime);index"`

	ConfirmedDate	time.Time			`orm:"null;type(datetime);index"`
	ConfirmedBy		*User				`orm:"rel(fk);null;index"`
	VerifiedDate	time.Time			`orm:"null;type(datetime);index"`
	VerifiedBy		*User				`orm:"rel(fk);null;index"`

	//Feedback		[]*BoilerAlarmFeedback	`orm:"reverse(many)"`

	TriggerRule		*RuntimeAlarmRule	`orm:"rel(fk);null;index"`
	AlarmLevel		int32				`orm:"index"`
	State			int32				`orm:"index"`
	Priority		int32				`orm:"index;default(0)"`
	NeedSend		bool				`orm:"index"`

	Description		string
}

const (
	BOILER_ALARM_STATE_DEFAULT = 0
	BOILER_ALARM_STATE_NEW = 1
	BOILER_ALARM_STATE_PENDING = 2
	BOILER_ALARM_STATE_CONFIRMED = 3
	BOILER_ALARM_STATE_REJECTED = 4
	BOILER_ALARM_STATE_VERIFIED = 5
	BOILER_ALARM_STATE_DONE = 10
)

func (alarm *BoilerAlarm) TableIndex() [][]string {
	return [][]string {
		[]string {"Boiler", "Parameter", "TriggerRule"},
	}
}

func (alarm *BoilerAlarm) TableUnique() [][]string {
	return [][]string {
		[]string {"Boiler", "Parameter", "TriggerRule", "StartDate"},
	}
}
