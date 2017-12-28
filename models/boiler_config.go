package models

type BoilerConfig struct {
	MyUidObject

	Boiler			*Boiler		`orm:"rel(fk);index"`

	IsGenerateData		bool		`orm:"index"`
}

func (conf *BoilerConfig) TableUnique() [][]string {
	return [][]string {
		[]string {"Boiler"},
	}
}