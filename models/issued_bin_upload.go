package models

import "time"

type IssuedBinUpload struct {
	Name string                   `orm:"pk;size(50)"`
	CreateTime time.Time          `orm:"type(datetime)"`
	UpdateTime time.Time          `orm:"type(datetime)"`
	Organization  *Organization   `orm:"rel(fk);null"`
	BinPath string                `orm:"size(50)"`
	Status bool
}
