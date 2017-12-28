package models

type BoilerMedium struct {
	MyIdNoAutoObject

	RuntimeParameters 	[]*RuntimeParameter		`orm:"reverse(many)"`
}