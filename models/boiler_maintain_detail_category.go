package models

type BoilerMaintainDetailCategory struct {
	MyIdObject

	Items		[]*BoilerMaintainDetailItem		`orm:"reverse(many)"`
}
