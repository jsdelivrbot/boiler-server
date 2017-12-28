package models

import "time"

type BoilerMaintainSchedule struct {
	InspectInnerDateNext	time.Time		`orm:"null;type(datetime);"`
	InspectOuterDateNext	time.Time		`orm:"null;type(datetime);"`
	InspectValveDateNext	time.Time		`orm:"null;type(datetime);"`
	InspectGaugeDateNext	time.Time		`orm:"null;type(datetime);"`
}
