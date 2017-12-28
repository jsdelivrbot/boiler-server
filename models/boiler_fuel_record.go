package models

import "time"

type BoilerFuelRecord struct {
	MyUidObject

	Boiler		*Boiler		`orm:"rel(fk);index"`
	Fuel		*Fuel		`orm:"rel(fk);index"`

	StartDate	time.Time	`orm:"index"`
	EndDate		time.Time	`orm:"index"`

	TotalFlow	float64
	FuelAmount	float64
	Rate		float64		`orm:"index"`
}
