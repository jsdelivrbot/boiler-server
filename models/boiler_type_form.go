package models

type BoilerTypeForm struct {
	MyIdNoAutoObject

	Type		*BoilerType		`orm:"rel(fk);index"`
}