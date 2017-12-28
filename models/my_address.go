package models

type MyAddress struct {
	Location		*Location		`orm:"rel(fk)"`
	Address			string
	ZipCode			string			`orm:"size(20)"`

	Longitude		float64					//经度
	Latitude		float64					//纬度
}

func (addr MyAddress) GetLocation() *Location {
	return addr.Location
}