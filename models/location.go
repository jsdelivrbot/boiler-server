package models

type Location struct {
	MyUidObject

	LocationId		int64			`orm:"index"`
	SuperId			int64			`orm:"index"`

	LocationName		string			//`orm:"-"`
}

type LocationInterface interface {
	GetLocation()	*Location
}