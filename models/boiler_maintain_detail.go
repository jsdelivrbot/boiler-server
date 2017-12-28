package models

type BoilerMaintainDetail struct {
	MyUidObject

	Template	*BoilerMaintainDetailTemplate		`orm:"rel(fk);index"`
}
