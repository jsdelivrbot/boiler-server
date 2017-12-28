package caches

import (
	"github.com/AzureTech/goazure"

	"github.com/AzureRelease/boiler-server/models"
	"encoding/json"
)

type BoilerRuntimeHistory struct {
	models.MyIdObject

	Boiler			*models.Boiler		`orm:"rel(fk);index"`

	JsonData		string				`orm:"type(text)"`
	Histories		[]*History			`orm:"-"`
}

type History struct {
	ParameterId		int64				`json:"pid"`
	Value			float64				`json:"val"`
	Alarm			int					`json:"alm"`
}

func (history *BoilerRuntimeHistory) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "CreatedDate"},
	}
}

func (history *BoilerRuntimeHistory) Marshal() {
	if js, err := json.Marshal(&history.Histories); err != nil {
		goazure.Error("History Json Marshal Error:", err, history.Histories)
	} else {
		history.JsonData = string(js)
	}
}

func (history *BoilerRuntimeHistory) Unmarshal() {
	his := []*History{}
	if err := json.Unmarshal([]byte(history.JsonData), &his); err != nil {
		goazure.Error("History Json Unmarshal Error:", history.JsonData)
	} else {
		history.Histories = his
	}
}