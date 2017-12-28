package models

import "time"

type UserLogin struct {
	MyUidObject

	User			*User			`orm:"rel(fk);null"`
	LoginPassword		string		`orm:"size(60)"`
	IsSuccess		bool

	Sessions		[]*UserSession	`orm:"reverse(many)"`

	LoginMethod		string
	LoginIp			string			`orm:"size(20)"`
	LoginDate		time.Time		`orm:"auto_now_add;type(datetime)"`
}
