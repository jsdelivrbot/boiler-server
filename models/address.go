package models

type Address struct {
	MyUidObject

	MyAddress

	Boilers			[]*Boiler		`orm:"reverse(many)"`
	Organizations		[]*Organization		`orm:"reverse(many)"`
}
