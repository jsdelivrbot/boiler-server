package models

import "time"

type Terminal struct {
	MyUidObject

	TerminalCode		int64
	Organization		*Organization	`orm:"rel(fk);null;index"`
	Boilers				[]*Boiler		`orm:"reverse(many)"`

	LocalIp				string			`orm:"size(60);null"`
	RemoteIp			string			`orm:"size(60);null"`
	RemotePort			int
	SimNumber			string			`orm:"size(20)"`

	InstalledBy			string			`orm:"size(60);null"`
	InstalledDate		time.Time		`orm:"type(datetime);null"`

	UploadFlag			bool			`orm:"default(true)"`
	UploadPeriod		int64			`orm:"default(45)"`

	IsOnline			bool
}

func (term *Terminal) TableUnique() [][]string {
	return [][]string {
		[]string{"TerminalCode"},
	}
}
