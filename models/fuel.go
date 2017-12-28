package models

type Fuel struct {
	MyUidObject

	FuelId		int32			`orm:"index"`
	Type		*FuelType		`orm:"rel(fk)"`
	QNetVArMin	int64
	QNetVArMax	int64
}

func (obj *Fuel) TableUnique() [][]string {
	return [][]string {
		[]string{"FuelId"},
	}
}
