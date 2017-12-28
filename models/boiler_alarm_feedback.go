package models

type BoilerAlarmFeedback struct {
	MyUidObject

	State		int32		`orm:"index"`
	Alarm		*BoilerAlarm	`orm:"rel(fk);index"`
	Content		string
}

const (
	ALARM_FEEDBACK_TYPE_DEFAULT = 0
	ALARM_FEEDBACK_TYPE_CONFIRM = 1
	ALARM_FEEDBACK_TYPE_REJECT = 2
	ALARM_FEEDBACK_TYPE_VERIFIED = 3
)