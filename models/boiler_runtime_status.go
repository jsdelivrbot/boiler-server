package models

type BoilerRuntimeStatus struct {
	MyUidObject

	Boiler			*Boiler			`orm:"rel(fk);index"`
	Parameter		*RuntimeParameter	`orm:"rel(fk);index"`
	Alarm			*BoilerAlarm		`orm:"rel(fk);null;index"`


}
