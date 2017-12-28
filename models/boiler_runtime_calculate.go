package models

type BoilerRuntimeCalculate struct {
	MyUidObject

	Boiler			*Boiler			`orm:"rel(fk);index"`
	Parameter		*RuntimeParameter	`orm:"rel(fk);index"`
	Value			int64
}
