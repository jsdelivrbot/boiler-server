package models

type BoilerTemplate struct {
	MyUidObject

	TemplateId		int32			`orm:"index"`
}

func (obj *BoilerTemplate) TableUnique() [][]string {
	return [][]string {
		//[]string{"Boiler", "CreatedDate"},
		[]string{"TemplateId"},
	}
}