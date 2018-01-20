package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"
	"fmt"
	"github.com/AzureRelease/boiler-server/dba"

	"github.com/AzureRelease/boiler-server/models"
	"errors"
)

type LocationController struct {
	goazure.Controller
}

var LocalCtl *LocationController

type ALocal struct {
	models.LocationInterface
	LocationName		string
}

func (al ALocal) GetLocation() *models.Location {
	return al.LocationInterface.GetLocation()
}

func (ctl *LocationController) LocationList() {
	var locations []orm.Params

	var provinces, cities, regions []orm.Params

	if num, err := dba.BoilerOrm.QueryTable("location").
		Filter("IsDeleted", false).OrderBy("LocationId").Limit(-1).Values(&locations, "LocationId", "SuperId", "Name", "NameEn", "LocationName"); err != nil {
			goazure.Error("Get Locations Error:", err, num)
	}

	for _, loc := range locations {
		//goazure.Info(i, "Location:", loc)
		if loc["LocationId"].(int64) < 100 {
			if loc["LocationId"].(int64) > 0 {
				loc["cities"] = []orm.Params{}
			}
			provinces = append(provinces, loc)
		} else if loc["LocationId"].(int64) < 10000 {
			for _, p := range provinces {
				if loc["SuperId"] == p["LocationId"] {
					p["cities"] = append(p["cities"].([]orm.Params), loc)
				}
			}
			loc["regions"] = []orm.Params{}
			cities = append(cities, loc)
		} else {
			for _, c := range cities {
				if loc["SuperId"] == c["LocationId"] {
					c["regions"] = append(c["regions"].([]orm.Params), loc)
				}
			}
			regions = append(regions, loc)
		}
	}

	//goazure.Warn("Locations:", provinces)
	//panic(0)

	ctl.Data["json"] = provinces
	ctl.ServeJSON()
}

func (ctl *LocationController) LocationListWeixin() {
	var provinces []orm.Params

	qs := dba.BoilerOrm.QueryTable("location")
	qsProvince := qs.Filter("SuperId", 0)
	if num, err := qsProvince.Filter("IsDeleted", false).OrderBy("LocationId").Values(&provinces, "LocationId", "SuperId", "Name", "NameEn"); err != nil ||num == 0 {
		goazure.Error("Read Location-Provinces Error:", err, num)
	}

	ctl.Data["json"] = provinces
	ctl.ServeJSON()
}

func (ctl *LocationController) GetLocationName(local *models.Location) error {
	if local == nil {
		return errors.New("can not be nil location")
	}
	l := local
	var localName string
	if l.LocationId == 0 {
		localName = "全国"
	} else {
		for l.LocationId != 0 {
			localName = l.Name + " " + localName
			l = LocalCtl.GetLocationWithId(l.SuperId)
		}
	}
	local.LocationName = localName
	//fmt.Println("Location", local)

	return nil
}

func (ctl *LocationController) UpdateLocalName() (error) {
	var locals []*models.Location
	ql := dba.BoilerOrm.QueryTable("location")
	if num, err := ql.Filter("IsDeleted", false).Limit(-1).All(&locals); err != nil || num == 0 {
		goazure.Error("Get Locations Error:", err, num)
		return err
	}

	for _, l := range locals {
		ctl.GetLocationName(l)
		DataCtl.UpdateData(l)
	}

	return nil
}

func (ctl *LocationController) GetLocationWithId(localId int64) *models.Location {
	local := models.Location{ LocationId: localId }
	err := DataCtl.ReadData(&local, "LocationId")
	if err != nil {
		fmt.Println("Read Location Error: ", err)
	}
	return &local
}