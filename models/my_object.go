package models

import "time"

type MyObject struct {
	Name			string		`orm:"size(60);index"`
	NameEn			string		`orm:"size(60)"`
	//NameShort		string		`orm:"size(60)"`
	Remark			string

	CreatedDate		time.Time	`orm:"type(datetime);auto_now_add;index"`
	CreatedBy		*User		`orm:"rel(fk);null"`
	UpdatedDate		time.Time	`orm:"type(datetime);auto_now;index"`
	UpdatedBy		*User		`orm:"rel(fk);null"`
	IsDemo			bool		`orm:"index"`
	IsDeleted		bool		`orm:"index"`
}