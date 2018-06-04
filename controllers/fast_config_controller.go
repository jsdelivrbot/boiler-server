package controllers

import (
	"github.com/AzureTech/goazure"
	"time"
	"github.com/AzureRelease/boiler-server/models"
	"encoding/json"
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
)

type FastConfigController struct {
	MainController
}

func (ctl *FastConfigController) FastBoilerAdd() {
	usr := ctl.GetCurrentUser()

	/*if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		return nil, errors.New(e)
	}*/

	var info 	BoilerInfo
	var boiler 	models.Boiler

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &info); err != nil {
		goazure.Error("Unmarshal BoilerInfo JSON Error", err)
	}

	if len(info.Uid) > 0 {
		boiler.Uid = info.Uid
		if err := DataCtl.ReadData(&boiler); err != nil {
			goazure.Warn("Read BoilerInfo Failed!", err)
		}
	} else {
		boiler.CreatedBy = usr
	}

	if len(info.Name) > 0 { boiler.Name = info.Name }
	if len(info.RegisterCode) > 0 { boiler.RegisterCode = info.RegisterCode }
	if len(info.DeviceCode) > 0 { boiler.DeviceCode = info.DeviceCode }
	if len(info.ModelCode) > 0 { boiler.ModelCode = info.ModelCode }
	if info.EvaporatingCapacity > 0 { boiler.EvaporatingCapacity = info.EvaporatingCapacity }

	if len(info.FactoryNumber) > 0 { boiler.FactoryNumber = info.FactoryNumber }
	if len(info.CertificateNumber) > 0 { boiler.CertificateNumber = info.CertificateNumber }

	var usage models.BoilerUsage
	var med models.BoilerMedium
	var fuel models.Fuel
	var form models.BoilerTypeForm
	var template models.BoilerTemplate
	usage.Id = 1
	med.Id = info.MediumId
	fuel.Uid = info.FuelId
	form.Id = info.FormId
	if err := DataCtl.ReadData(&usage); err == nil { boiler.Usage = &usage }
	if err := DataCtl.ReadData(&med); err == nil { boiler.Medium = &med }
	if err := DataCtl.ReadData(&fuel); err == nil { boiler.Fuel = &fuel }
	if err := DataCtl.ReadData(&form); err == nil { boiler.Form = &form }
	if err :=dba.BoilerOrm.QueryTable("boiler_template").Filter("TemplateId",info.TemplateId).One(&template);err == nil {boiler.Template = &template}
	var enterprise, factory, maintainer, supervisor models.Organization
	enterprise.Uid = info.EnterpriseId
	factory.Uid = info.FactoryId
	maintainer.Uid = info.MaintainerId
	supervisor.Uid = info.SupervisorId
	if err := DataCtl.ReadData(&enterprise); err == nil { boiler.Enterprise = &enterprise }
	if err := DataCtl.ReadData(&factory); err == nil { boiler.Factory = &factory }
	if err := DataCtl.ReadData(&maintainer); err == nil { boiler.Maintainer = &maintainer }
	if err := DataCtl.ReadData(&supervisor); err == nil { boiler.Supervisor = &supervisor }

	if boiler.InspectInnerDateNext.IsZero() { boiler.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30) }
	if boiler.InspectOuterDateNext.IsZero() { boiler.InspectOuterDateNext = time.Now().Add(time.Hour * 24 * 30) }
	if boiler.InspectValveDateNext.IsZero() { boiler.InspectValveDateNext = time.Now().Add(time.Hour * 24 * 30) }
	if boiler.InspectGaugeDateNext.IsZero() { boiler.InspectGaugeDateNext = time.Now().Add(time.Hour * 24 * 30) }

	boiler.UpdatedBy = usr
	if err := DataCtl.AddData(&boiler, true); err != nil {
		goazure.Error("Add Boiler Error!", err)
	}

	if	num, err := dba.BoilerOrm.QueryTable("boiler_organization_linked").
		Filter("Boiler__Uid", boiler.Uid).Delete(); err != nil {
		goazure.Warn("Deleted Old Links Error:", err, num)
	}

	for _, li := range info.Links {
		var linked 	models.BoilerOrganizationLinked
		var og 		models.Organization
		var ogType	models.OrganizationType

		og.Uid = li.Uid
		ogType.TypeId = li.Type

		linked.Boiler = &boiler
		if 	err := DataCtl.ReadData(&og, "Uid"); err == nil {
			linked.Organization = &og
		} else {
			goazure.Error("Org Uid Is Not Valid!", err)
			continue
		}

		if 	err := DataCtl.ReadData(&ogType, "TypeId"); err == nil {
			linked.OrganizationType = &ogType
		} else {
			goazure.Error("OrgType Is Not Valid!", err)
			continue
		}

		if num, err := dba.BoilerOrm.InsertOrUpdate(&linked); err != nil {
			goazure.Error("M2M OrganizationLinked Add Error:", err, num)
		}
	}
	go CalcCtl.InitBoilerCalculateParameter([]*models.Boiler{&boiler})
	ctl.Data["json"] =boiler
	ctl.ServeJSON()
	goazure.Info("Updated Boiler:", boiler, info)
	fmt.Println("hahahhahahahahhahhaha")
}

func (ctl *FastConfigController) FastTermCombined() {

}

func (ctl *FastConfigController) FastTermChannelConfig() {

}
