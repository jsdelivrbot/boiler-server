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

type termCombined struct {
	BoilerUid  string    `json:"boiler_uid"`
	Code       int64    `json:"code"`
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
	ctl.Data["json"] =boiler.Uid
	ctl.ServeJSON()
	goazure.Info("Updated Boiler:", boiler, info)
}

func (ctl *FastConfigController) FastTermCombined() {
	usr := ctl.GetCurrentUser()
	var wCombined termCombined
	var terminal models.Terminal
	if err:=json.Unmarshal(ctl.Ctx.Input.RequestBody,&wCombined);err!=nil{
		goazure.Error("Unmarshal JSON Error",err)
	}
	fmt.Println("combined:",wCombined)
	qs := dba.BoilerOrm.QueryTable("terminal").RelatedSel("Organization")
	if !usr.IsAdmin() {
			qs = qs.Filter("Organization",usr.Organization.Uid)
	}
	if err:=qs.Filter("TerminalCode",wCombined.Code).Filter("IsDeleted",false).One(&terminal);err!=nil{
		goazure.Error("UNFind Terminal:",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("非法终端!"))
		return
	}
	fmt.Println("terminal:",terminal)
	if dba.BoilerOrm.QueryTable("boiler_terminal_combined").Filter("Terminal__Uid",terminal.Uid).Exist(){
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端已被占用!"))
		return
	}

	sql:="insert into boiler_terminal_combined(boiler_id,terminal_id,terminal_code,terminal_set_id) values(?,?,?,?) on duplicate key update terminal_id=?,terminal_code = ?"
	if _,err:=dba.BoilerOrm.Raw(sql,wCombined.BoilerUid,terminal.Uid,terminal.TerminalCode,1,terminal.Uid,terminal.TerminalCode).Exec();err!=nil{
		goazure.Error("Insert boiler_terminal_combined Error",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端绑定失败"))
		return
	}
}

func (ctl *FastConfigController) FastTermUnbind() {
	usr := ctl.GetCurrentUser()
	var termBind termCombined
	var terminal models.Terminal
	if err:= json.Unmarshal(ctl.Ctx.Input.RequestBody,&termBind);err!=nil{
		goazure.Error("Unmarshal JSON Error",err)
	}
	qs := dba.BoilerOrm.QueryTable("terminal").RelatedSel("Organization")
	if !usr.IsAdmin() {
		qs = qs.Filter("Organization",usr.Organization.Uid)
	}
	if err:=qs.Filter("TerminalCode",termBind.Code).Filter("IsDeleted",false).One(&terminal);err!=nil{
		goazure.Error("UNFind Terminal:",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("非法终端!"))
		return
	}
	if _,err:=dba.BoilerOrm.QueryTable("boiler_terminal_combined").
	Filter("Boiler__Uid",termBind.BoilerUid).Filter("Terminal__Uid",terminal.Uid).Delete();err!=nil{
		goazure.Error("Unbind Boiler/Terminal Error",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("解绑失败"))
	}
}

func (ctl *FastConfigController) FastTermChannelConfig() {

}
