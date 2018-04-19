package models

import "time"

type TerminalIssued struct {
	Terminal *Terminal
	IsOnline  bool
	TermVer    int32
	TermUpdateTime time.Time
	PlatVer    int32
	PlatUpdateTime time.Time
	IssuedTermTempStatus  IssuedTermTempStatus
}

type Terminal struct {
	MyUidObject

	TerminalCode		int64			`orm:"index"`
	Organization		*Organization	`orm:"rel(fk);null;index"`
	Boilers				[]*Boiler		`orm:"reverse(many);null;index;rel_through(BoilerGo/models.BoilerTerminalCombined)"`

	LocalIp				string			`orm:"size(60);null"`
	RemoteIp			string			`orm:"size(60);null"`
	RemotePort			int
	SimNumber			string			`orm:"size(20)"`

	InstalledBy			string			`orm:"size(60);null"`
	InstalledDate		time.Time		`orm:"type(datetime);null"`

	UploadFlag			bool			`orm:"default(true)"`
	UploadPeriod		int64			`orm:"default(45)"`

	IsOnline			bool			`orm:"index"`

	TerminalSetId		int32			`orm:"-"`
}

func (term *Terminal) TableUnique() [][]string {
	return [][]string {
		[]string{"TerminalCode"},
	}
}
