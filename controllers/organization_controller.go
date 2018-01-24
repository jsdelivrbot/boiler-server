package controllers

import (
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"reflect"
	"fmt"
	"strconv"
	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"
	"encoding/json"
	"time"
	"errors"
)

type OrganizationController struct {
	MainController
}

var OrgCtrl *OrganizationController = &OrganizationController{}

const organizationDefautlsPath string = "models/properties/organization_defaults/"

func (ctl *OrganizationController) OrganizationList() {
	usr := ctl.GetCurrentUser()
	var orgs []models.Organization
	var tid int64
	var method string = ""
	if ctl.Input()["tid"] != nil && len(ctl.Input()["tid"]) > 0 {
		tid, _ = strconv.ParseInt(ctl.Input()["tid"][0], 10, 32)
	}
	if ctl.Input()["scope"] != nil && len(ctl.Input()["scope"]) > 0 &&
		ctl.Input()["scope"][0] == "register" {
		method = ctl.Input()["scope"][0]
	}

	if usr == nil || usr.IsCommonUser() {
		if method != "register" {
			return
		}
	}

	qs := dba.BoilerOrm.QueryTable("organization")
	qs = qs.RelatedSel("Address__Location").RelatedSel("Type")
	if 	method != "register" && usr.IsOrganizationUser() {
		orCond := orm.NewCondition().Or("Uid", usr.Organization.Uid).Or("SuperOrganization__Uid", usr.Organization.Uid).Or("CreatedBy__Uid", usr.Uid)
		cond := orm.NewCondition().AndCond(orCond)
		qs = qs.SetCond(cond)
	}
	if tid > 0 {
		qs = qs.Filter("Type__TypeId", int32(tid))
	}
	if num, err := qs.Filter("IsDeleted", false).OrderBy("Type__TypeId", "-IsSupervisor").All(&orgs); err != nil || num == 0 {
		goazure.Error("Read Orgs Error:", err, num)
	}

	ctl.Data["json"] = orgs
	ctl.ServeJSON()
}

func (ctl *OrganizationController) OrganizationTypeList() {
	var types []orm.Params
	qs := dba.BoilerOrm.QueryTable("organization_type")
	qs = qs.Filter("TypeId__gte", 1)
	if num, err := qs.OrderBy("TypeId").Values(&types, "TypeId", "Name"); err != nil {
		goazure.Error("Read Orgs Error:", err, num)
	}

	ctl.Data["json"] = types
	ctl.ServeJSON()
}

type Org struct {
	Uid						string		`json:"uid"`
	TypeId					int			`json:"type_id"`
	
	Name					string		`json:"name"`
	
	LocationId				int			`json:"location_id"`
	Address					string		`json:"address"`

	ShowBrand				bool		`json:"show_brand"`
	BrandName				string		`json:"brand_name"`

	IsSuper					bool		`json:"is_super"`
	Supervisor  			string		`json:"supervisor"`

	GenerateSampleBoilers	bool		`json:"generate_sample_boilers"`
	GenerateSampleData		bool		`json:"generate_sample_data"`
}

func (ctl *OrganizationController) OrganizationUpdate() {
	usr := ctl.GetCurrentUser()
	var o Org
	var organization models.Organization
	var addr models.Address

	var isCreated bool = false

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &o); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	goazure.Info("update org:", o)

	if len(o.Uid) > 0 {
		if err := dba.BoilerOrm.QueryTable("organization").RelatedSel("Address__Location").RelatedSel("Contact").Filter("Uid", o.Uid).One(&organization); err != nil {
			e := fmt.Sprintf("Read Organization for Update Error: %v", err)
			goazure.Warn(e)
			isCreated = true
		}
	} else {
		//goazure.Error("Read Organization for Update:", organization)
		isCreated = true
	}

	oType := models.OrganizationType{}
	oType.TypeId = int32(o.TypeId)

	if err := DataCtl.ReadData(&oType, "TypeId"); err != nil {
		e := fmt.Sprintf("Read OrganizationType Error: %v", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if 	!usr.IsAdmin() {
		if 	usr.Organization.Uid != organization.Uid &&
			(organization.CreatedBy == nil || usr.Uid != organization.CreatedBy.Uid) &&
			(usr.Organization == nil || organization.SuperOrganization == nil || usr.Organization.Uid != organization.SuperOrganization.Uid) {
			ctl.Ctx.Output.SetStatus(403)
			ctl.Ctx.Output.Body([]byte("Permission Denied!"))
			goazure.Error("Permission Denied!")
			return
		}
	}

	if err := dba.BoilerOrm.QueryTable("address").RelatedSel("Location").Filter("Address", o.Address).Filter("Location__LocationId", o.LocationId).One(&addr); err != nil {
		e := fmt.Sprintf("Read Organization Address Error: %v", err)
		goazure.Warn(e)

		local := models.Location{}
		local.LocationId = int64(o.LocationId)

		if err := DataCtl.ReadData(&local, "LocationId"); err != nil {
			e := fmt.Sprintf("Read Organization Location Error: %v", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}

		addr.Address = o.Address
		addr.Location = &local

		if err := DataCtl.AddData(&addr, true, "Address", "Location"); err != nil {
			e := fmt.Sprintf("Add New Address Error: %v", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
	}

	organization.Name = o.Name
	organization.Type = &oType
	organization.Address = &addr
	
	organization.UpdatedBy = usr

	if isCreated {
		organization.CreatedBy = usr
		organization.SuperOrganization = usr.Organization
	}

	if usr.IsAdmin() {
		organization.ShowBrand = o.ShowBrand
		organization.BrandName = o.BrandName

		organization.IsSupervisor = o.IsSuper
	}

	if err := DataCtl.AddData(&organization, true); err != nil {
		e := fmt.Sprintf("Add Organization Error: %v", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if isCreated && o.GenerateSampleBoilers {
		if boilers, err := ctl.GenerateSampleBoilers(&organization); err != nil {
			goazure.Error("Generate Sample Boilers Error:", err)
		} else if o.GenerateSampleData {
			for _, b := range boilers {
				var conf models.BoilerConfig
				qs := dba.BoilerOrm.QueryTable("boiler_config")
				if err := qs.Filter("Boiler__Uid", b.Uid).One(&conf); err != nil {
					goazure.Error("Read BoilerConfig Error", err)
				}

				conf.Boiler = b
				conf.Name = b.Name
				conf.IsGenerateData = true

				if err := DataCtl.AddData(&conf, true); err != nil {
					e := fmt.Sprintln("Update BoilerConfig JSON Error", err)
					goazure.Error(e)
				}
			}
		}

		BlrCtl.RefreshGlobalBoilerList()
	}

	goazure.Info("Updated Organization:", organization)
}

func (ctl *OrganizationController) GenerateSampleBoilers(org *models.Organization) ([]*models.Boiler, error) {
	goazure.Info("GenerateSampleBoilers")
	if org.Type.TypeId != 1 && org.Type.TypeId != 2 {
		return []*models.Boiler{}, errors.New("organization is not enterprise or factory")
	}

	var boilers []*models.Boiler
	for i := 0; i < 4; i++ {
		boilers = append(boilers, &models.Boiler{})
	}

	boilerCoal := boilers[0]
	boilerGas := boilers[1]
	boilerBiomass := boilers[2]
	boilerWater := boilers[3]

	usage := models.BoilerUsage{}
	usage.Id = 1

	form201 := models.BoilerTypeForm{}
	form202 := models.BoilerTypeForm{}
	form205 := models.BoilerTypeForm{}
	form201.Id = 201
	form202.Id = 202
	form205.Id = 205

	mediumSteam := models.BoilerMedium{}
	mediumWater := models.BoilerMedium{}
	mediumSteam.Id = 1
	mediumWater.Id = 2

	fuelCoal02 := models.Fuel{ FuelId: 102 }
	fuelCoal05 := models.Fuel{ FuelId: 105 }
	fuelGas := models.Fuel{ FuelId: 301 }
	fuelBiomass := models.Fuel{ FuelId: 401 }

	DataCtl.ReadData(&form201, "Id")
	DataCtl.ReadData(&form202, "Id")
	DataCtl.ReadData(&form205, "Id")
	DataCtl.ReadData(&mediumSteam, "Id")
	DataCtl.ReadData(&mediumWater, "Id")
	DataCtl.ReadData(&fuelCoal02, "FuelId")
	DataCtl.ReadData(&fuelCoal05, "FuelId")
	DataCtl.ReadData(&fuelGas, "FuelId")
	DataCtl.ReadData(&fuelBiomass, "FuelId")

	boilerCoal.Form = &form201
	boilerCoal.Medium = &mediumSteam
	boilerCoal.Fuel = &fuelCoal02
	boilerCoal.EvaporatingCapacity = 4

	boilerGas.Form = &form201
	boilerGas.Medium = &mediumSteam
	boilerGas.Fuel = &fuelGas
	boilerGas.EvaporatingCapacity = 12

	boilerBiomass.Form = &form202
	boilerBiomass.Medium = &mediumSteam
	boilerBiomass.Fuel = &fuelBiomass
	boilerBiomass.EvaporatingCapacity = 6

	boilerWater.Form = &form205
	boilerWater.Medium = &mediumWater
	boilerWater.Fuel = &fuelGas
	boilerWater.EvaporatingCapacity = 8

	for i, b := range boilers {
		b.Name = org.Name + "锅炉#" + strconv.FormatInt(int64(i), 10)
		b.Usage = &usage
		switch org.Type.TypeId {
		case 1:
			b.Factory = org
		case 2:
			b.Enterprise = org
		}

		b.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30)
		b.InspectOuterDateNext = time.Now().Add(time.Hour * 24 * 30)
		b.InspectValveDateNext = time.Now().Add(time.Hour * 24 * 30)
		b.InspectGaugeDateNext = time.Now().Add(time.Hour * 24 * 30)

		if err := DataCtl.AddData(b, false); err != nil {
			goazure.Error("Generate Sample Boilers Error:", err, b)
		}
	}

	go CalcCtl.InitBoilerCalculateParameter(boilers)

	return boilers, nil
}

func (ctl *OrganizationController) OrganizationDelete() {
	usr := ctl.GetCurrentUser()
	var o Org
	var organization models.Organization
	var users []*models.User

	if !usr.IsAdmin() {
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte("Permission Denied!"))
		goazure.Error("Permission Denied!")
		return
	}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &o); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	if err := dba.BoilerOrm.QueryTable("organization").RelatedSel("Address__Location").RelatedSel("Contact").Filter("Uid", o.Uid).One(&organization); err != nil {
		e := fmt.Sprintf("Read Organization for Delete Error: %v", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if num, err := dba.BoilerOrm.QueryTable("user").Filter("Organization__Uid", o.Uid).All(&users); err != nil {
		e := fmt.Sprintln("Read Organization Users for Delete Error:", err, num)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if err := DataCtl.DeleteData(&organization); err != nil {
		e := fmt.Sprintln("Delete Organization Error!", organization, err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	for _, u := range users {
		if err := DataCtl.DeleteData(u); err != nil {
			e := fmt.Sprintln("Delete Organization User Error!", u, err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
	}
}


func (ctl *OrganizationController) InitOrganizationDefaults() {
	var orgType models.OrganizationType

	DataCtl.GenerateDefaultData(reflect.TypeOf(orgType), organizationDefautlsPath, "type", nil)
}

func (ctl *OrganizationController) organizationWithUid(uid string) *models.Organization {
	o := models.Organization{ MyUidObject: models.MyUidObject{ Uid: uid} }
	if err := DataCtl.ReadData(&o); err != nil {
		goazure.Error("Read Organization Error", err)
		return nil
	}

	return &o
}