package models

import "time"

type UserLogin struct {
	MyUidObject

	User			*User			`orm:"rel(fk);null"`
	IsLogin			bool			`orm:"index"`
	IsSuccess		bool			`orm:"index"`

	LoginPassword	string			`orm:"size(60)"`
	LoginMethod		string
	LoginIp			string			`orm:"size(20)"`
	LoginDate		time.Time		`orm:"auto_now_add;type(datetime)"`

	Sessions		[]*UserSession	`orm:"reverse(many)"`
}
