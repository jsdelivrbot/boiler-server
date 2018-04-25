package caches

import "github.com/AzureRelease/boiler-server/models"

type BoilerRuntimeHistoryArchived struct {
	models.MyIdObject

	Boiler			*models.Boiler				`orm:"rel(fk);index"`
	Parameter		*models.RuntimeParameter	`orm:"rel(fk);index"`

	Value			int64
}

func (his *BoilerRuntimeHistoryArchived) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Parameter", "CreatedDate"},
	}
}