package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"

	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"

	"fmt"
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
	var provinces []orm.Params

	qs := dba.BoilerOrm.QueryTable("location")
	qsProvince := qs.Filter("LocationId__lt", 100)
	if num, err := qsProvince.Filter("IsDeleted", false).OrderBy("LocationId").Values(&provinces, "LocationId", "SuperId", "Name", "NameEn", "LocationName"); (err != nil ||num == 0) {
		fmt.Println("Read Location-Provinces Error:", err, num);
	}

	for _, p := range provinces {
		if p["LocationId"] == 0 || p["Name"] == "全国" {
			continue
		}
		var cities []orm.Params
		qsCity := qs.Filter("SuperId", p["LocationId"])
		if num, err := qsCity.Filter("IsDeleted", false).OrderBy("LocationId").Values(&cities, "LocationId", "SuperId", "Name", "NameEn", "LocationName"); (err != nil || num == 0) {
			fmt.Println("Read Location-Cities Error:", err, num);
		}

		for _, c := range cities {
			var regions []orm.Params
			qsRegion := qs.Filter("SuperId", c["LocationId"])
			if num, err := qsRegion.Filter("IsDeleted", false).OrderBy("LocationId").Values(&regions, "LocationId", "SuperId", "Name", "NameEn", "LocationName"); (err != nil || num == 0) {
				fmt.Println("Read Location-Regions Error:", err, num);
			}
			c["regions"] = regions;
		}

		p["cities"] = cities;
	}

	ctl.Data["json"] = provinces
	ctl.ServeJSON()
}

func (ctl *LocationController) LocationListWeixin() {
	var provinces []orm.Params

	qs := dba.BoilerOrm.QueryTable("location")
	qsProvince := qs.Filter("SuperId", 0)
	if num, err := qsProvince.Filter("IsDeleted", false).OrderBy("LocationId").Values(&provinces, "LocationId", "SuperId", "Name", "NameEn"); (err != nil ||num == 0) {
		goazure.Error("Read Location-Provinces Error:", err, num);
	}

	ctl.Data["json"] = provinces
	ctl.ServeJSON()
}

func (ctl *LocationController) GetLocationName(local *models.Location) error {
	if local == nil {
		return errors.New("Can not be Nil Location!")
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
		return err;
	}

	for _, l := range locals {
		ctl.GetLocationName(l)
		DataCtl.UpdateData(l)
	}

	return nil
}

func (localCtl *LocationController) GetLocationWithId(localId int64) *models.Location {
	local := models.Location{ LocationId: localId }
	err := DataCtl.ReadData(&local, "LocationId")
	if err != nil {
		fmt.Println("Read Location Error: ", err)
	}
	return &local
}