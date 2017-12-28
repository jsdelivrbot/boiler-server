package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"
)

type DBController struct {
	goazure.Controller
}

var DBCtl *DBController

func (dbCtl *DBController) GetStringFromMap(m orm.Params, defaults string, col string, cols ...string) string {
	if m[col] != nil {
		return m[col].(string)
	} else {
		for _, c := range cols {
			if m[c] != nil {
				return m[c].(string)
			}
		}
		return defaults
	}
}