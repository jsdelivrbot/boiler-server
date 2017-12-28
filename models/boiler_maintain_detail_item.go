package models

type BoilerMaintainDetailItem struct {
	MyIdObject

	Category	*BoilerMaintainDetailCategory		`orm:"rel(fk);index"`
	Details		[]*BoilerMaintainDetail			`orm:"reverse(many)"`
}
