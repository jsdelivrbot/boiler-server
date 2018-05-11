package models

type BoilerTemplate struct {
	MyUidObject
	TemplateId		int32			`orm:"index"`
	Organization *Organization    `orm:"rel(fk);null"`
}

func (obj *BoilerTemplate) TableUnique() [][]string {
	return [][]string {
		//[]string{"Boiler", "CreatedDate"},
		[]string{"TemplateId"},
	}
}