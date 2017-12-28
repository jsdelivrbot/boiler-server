package models

type UserSession struct {
	MyUidObject

	User		*User		`orm:"rel(fk);index"`
	Login		*UserLogin	`orm:"rel(fk);index"`
	Application	*Application	`orm:"rel(fk);null;index"`

	IsActived	bool		`orm:"index"`
	Token		string		`orm:"size(60);index"`
}
