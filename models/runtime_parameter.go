package models

import (
	"fmt"
	"github.com/AzureRelease/boiler-server/common"
)

type RuntimeParameter struct {
	MyIdNoAutoObject

	NameShort			string							`orm:"size(60)"`
	IsDefault           bool		                     `orm:"index"`
	ParamId       		int32							`orm:"index"`
	Category      		*RuntimeParameterCategory		`orm:"rel(fk);index"`
	Medium        		*RuntimeParameterMedium			`orm:"rel(fk);index"`
	Organization        *Organization					`orm:"rel(fk);index"`
	BoilerMediums 		[]*BoilerMedium					`orm:"rel(m2m)"`
	Length        		int32							`orm:"default(2)"`
	Scale         		float32							`orm:"default(1.0)"`
	Fix					int32							`orm:"default(2)"`
	Unit          		string
}

func (param *RuntimeParameter)AddBoilerMedium(medId int64) error {
	var err error
	var num int64

	m2m := common.BoilerOrm.QueryM2M(param, "BoilerMediums")

	var boilerTypes []BoilerMedium
	num, err = common.BoilerOrm.QueryTable("boiler_medium").OrderBy("Id").All(&boilerTypes)

	bType := boilerTypes[medId]
	fmt.Println("BoilerMedium Id:", medId)
	fmt.Println("Parameters: ", param)
	fmt.Println("Type: ", bType)
	if !m2m.Exist(&bType) {
		num, err = m2m.Add(&bType)
		fmt.Println("Added M2M Num:", num)
	}

	return err
}