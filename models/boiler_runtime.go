package models

type BoilerRuntime struct {
	MyIdObject
	//Uid 			string				`orm:"size(36);index"`

	Boiler			*Boiler				`orm:"rel(fk);index"`
	Parameter		*RuntimeParameter	`orm:"rel(fk);index"`
	Alarm			*BoilerAlarm		`orm:"rel(fk);null;index"`
	Value			int64

	Status			int					`orm:"index;default(1)"`
}

const (
	RUNTIME_STATUS_DEFAULT		= 0
	RUNTIME_STATUS_NEW			= 1
	RUNTIME_STATUS_NEEDRELOAD	= 2
)

func (runtime *BoilerRuntime) TableUnique() [][]string {
	return [][]string {
		[]string {"Boiler", "Parameter", "CreatedDate"},
	}
}