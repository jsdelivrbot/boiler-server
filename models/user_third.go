package models

import "time"

type UserThird struct {
	MyUidObject

	User		*User		`orm:"rel(fk);null;index"`

	Application	*Application	`orm:"rel(fk);index"`

	Platform	string		`orm:"size(60);index"`
	App		string		`orm:"size(60);index"`
	Identity	string		`orm:"size(60);index"`

	Language	string
	Sex		int
	Province	string
	City		string
	Country		string
	HeadImageUrl	string
	//Privilege	[]string

	OpenId		string		`orm:"index"`
	UnionId		string		`orm:"index"`
	AccessToken	string		`orm:"index"`
	RefreshToken	string		`orm:"index"`
	SessionKey	string
	ExpiresIn	time.Time	`orm:"null;index"`
}