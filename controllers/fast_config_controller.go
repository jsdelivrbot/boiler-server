package controllers

import (
	"github.com/AzureTech/goazure"
	"time"
	"github.com/AzureRelease/boiler-server/models"
	"encoding/json"
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
	"strconv"
	"github.com/pborman/uuid"
	"github.com/AzureTech/goazure/orm"
)

type FastConfigController struct {
	MainController
}

type TempToCur struct {
	TemplateUid string `json:"template"`
	Code  int  `json:"code""`
}

type termCombined struct {
	BoilerUid  string    `json:"boiler_uid"`
	Code       int64    `json:"code"`
}

type FastChannelConfig struct {
	Chan  Chans				`json:"Chan"`
	Param CommParam			`json:"Param"`
	Code  int  `json:"Code"`
}

type Chans struct {
   Analogue  []Analogue	`json:"Analogue"`
   Switch    []Switch	`json:"Switch"`
   Range    []Range		`json:"Range"`
}

type Analogue struct {
	Parameter Parameter		`json:"Parameter"`
	Func int		`json:"Func"`
	Byte int		`json:"Byte"`
	Modbus int		`json:"Modbus"`
	ChannelNumber int		`json:"ChannelNumber"`
	SequenceNumber int32		`json:"SequenceNumber"`
	Status int32			`json:"Status"`
	SwitchStatus int32		`json:"SwitchStatus"`
}

type Switch struct {
	Parameter Parameter		`json:"Parameter"`
	Func int			`json:"Func"`
	Modbus int           `json:"Modbus"`
	BitAddress int 		`json:"BitAddress"`
	ChannelNumber int		`json:"ChannelNumber"`
	SequenceNumber int32		`json:"SequenceNumber"`
	Status int32			`json:"Status"`
	SwitchStatus int32		`json:"SwitchStatus"`
}

type Range struct {
	Parameter Parameter  	`json:"Parameter"`
	Func int 		`json:"Func"`
	Byte int			`json:"Byte"`
	Modbus int			`json:"Modbus"`
	Ranges []Ranges		`json:"Ranges"`
	ChannelNumber  int  `json:"ChannelNumber"`
	SequenceNumber int32  `json:"SequenceNumber"`
	Status   int32        `json:"Status"`
	SwitchStatus  int32   `json:"SwitchStatus"`

}

type Parameter struct {
	Name string   `json:"Name"`
	Scale float32	`json:"Scale"`
	Unit string		`json:"Unit"`
}

type Ranges struct {
	Min int64		`json:"Min"`
	Max int64		`json:"Max"`
	Name string     `json:"Name"`
	Value int64     `json:"Value"`
}

type CommParam struct {
	BaudRate int	`json:"BaudRate"`
	DataBit int		`json:"DataBit"`
	StopBit int		`json:"StopBit"`
	CheckDigit int  `json:"CheckDigit"`
	CommunInterface int  `json:"CommunInterface"`
	SlaveAddress int   `json:"SlaveAddress"`
	HeartBeat int   `json:"HeartBeat"`
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
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		return
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
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		return
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
	if dba.BoilerOrm.QueryTable("boiler_terminal_combined").Filter("Terminal__Uid",terminal.Uid).Exist(){
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端已被占用!"))
		return
	}
	boiler := models.Boiler{}
	boiler.Uid = wCombined.BoilerUid
	errB := DataCtl.ReadData(&boiler)
	if errB != nil {
		e := fmt.Sprintln("Read Boiler Error:", errB)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	if boiler.Terminal == nil {
		boiler.Terminal = &terminal

		code := strconv.FormatInt(terminal.TerminalCode, 10)
		if len(code) < 6 {
			for l := len(code); l < 6; l++ {
				code = "0" + code
			}
		}
		boiler.TerminalCode = code
		boiler.TerminalSetId = 1

		if err := DataCtl.UpdateData(&boiler); err != nil {
			e := fmt.Sprintln("Boiler Bind Error:", err, boiler, terminal)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
	}
	var sql string =""
	if dba.BoilerOrm.QueryTable("boiler_terminal_combined").Filter("Boiler__Uid",wCombined.BoilerUid).Filter("TerminalSetId",1).Exist() {
		sql="update boiler_terminal_combined set terminal_id= ? ,terminal_code = ? where boiler_id=? and terminal_set_id = 1"
		if _,err:=dba.BoilerOrm.Raw(sql,terminal.Uid,terminal.TerminalCode,wCombined.BoilerUid).Exec();err!=nil{
			goazure.Error("Update boiler_terminal_combined Error",err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("终端绑定失败"))
			return
		}
	} else {
		sql ="insert into boiler_terminal_combined(boiler_id,terminal_id,terminal_code,terminal_set_id) values(?,?,?,?)"
		if _,err:=dba.BoilerOrm.Raw(sql,wCombined.BoilerUid,terminal.Uid,terminal.TerminalCode,1).Exec();err!=nil{
			goazure.Error("Insert boiler_terminal_combined Error",err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("终端绑定失败"))
			return
		}
	}
	go BlrCtl.RefreshGlobalBoilerList()
}

func (ctl *FastConfigController) FastTermUnbind() {
	usr := ctl.GetCurrentUser()
	var termBind termCombined
	var terminal models.Terminal
	if err:= json.Unmarshal(ctl.Ctx.Input.RequestBody,&termBind);err!=nil{
		goazure.Error("Unmarshal JSON Error",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		return
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
	boiler := models.Boiler{}
	boiler.Uid = termBind.BoilerUid
	errB := DataCtl.ReadData(&boiler)
	if errB != nil {
		e := fmt.Sprintln("Read Unbind Data Error:", errB)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	if boiler.Terminal.Uid == terminal.Uid {
		boiler.Terminal = nil
		boiler.TerminalCode = ""
		boiler.TerminalSetId = 0

		if err := DataCtl.UpdateData(&boiler); err != nil {
			e := fmt.Sprintln("Boiler Unbind Error:", err, boiler, terminal)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
		}
	}
	if _,err:=dba.BoilerOrm.QueryTable("boiler_terminal_combined").
	Filter("Boiler__Uid",termBind.BoilerUid).Filter("Terminal__Uid",terminal.Uid).Delete();err!=nil{
		goazure.Error("Unbind Boiler/Terminal Error",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("解绑失败"))
		return
	}
	go BlrCtl.RefreshGlobalBoilerList()
}

func (ctl *FastConfigController) FastTermChannelConfig() {
	var config FastChannelConfig

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &config); err != nil {
		goazure.Error("Unmarshal JSON Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		return
	}
	fmt.Println("config:",config)
	var terminal models.Terminal
	if err := dba.BoilerOrm.QueryTable("Terminal").Filter("TerminalCode", config.Code).One(&terminal); err != nil {
		goazure.Error("Query Terminal Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("找不到终端"))
		return
	}
	var aCnf []models.RuntimeParameterChannelConfig
	if _, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		Filter("Terminal__Uid", terminal.Uid).Filter("IsDefault", false).
		All(&aCnf); err != nil {
		goazure.Error("Get ChannelConfig To Delete Error:", err)
	}
	for _, a := range aCnf {
		TempCtrl.TemplateChannelConfigDelete(&a, terminal.Uid)
	}
	if num, err := dba.BoilerOrm.QueryTable("issued_switch_default").Filter("Terminal__Uid", terminal.Uid).Delete(); err != nil {
		goazure.Error("Delete issued switch Error:", err, num)
	}

	for _, analogue := range config.Chan.Analogue {
		var param models.RuntimeParameter
		var cnf models.RuntimeParameterChannelConfig
		var max int32
		sql := "select max(param_id) as param_id from runtime_parameter where category_id=? and is_deleted =false"
		if err := dba.BoilerOrm.Raw(sql, models.RUNTIME_PARAMETER_CATEGORY_ANALOG).QueryRow(&max); err != nil {
			goazure.Error("Query runtime_parameter Error", err)
		}
		if max >= 100 {
			param.ParamId = max + 1
		} else {
			param.ParamId = 100
		}
		a := fmt.Sprintf("%d%d", models.RUNTIME_PARAMETER_CATEGORY_ANALOG, param.ParamId)
		param.Id, _ = strconv.ParseInt(a, 10, 64)
		param.Length = 2
		param.Fix = 2
		param.Category = runtimeParameterCategory(models.RUNTIME_PARAMETER_CATEGORY_ANALOG)
		param.Name = analogue.Parameter.Name
		param.Unit = analogue.Parameter.Unit
		param.Scale = analogue.Parameter.Scale
		param.Medium = runtimeParameterMedium(0)
		param.AddBoilerMedium(0)
		param.Organization = terminal.Organization
		param.CreatedBy = ctl.GetCurrentUser()
		param.UpdatedBy = ctl.GetCurrentUser()
		param.IsDeleted = false
		if _,err:=dba.BoilerOrm.Insert(&param);err!=nil {
			e := fmt.Sprintln("Add/Update Parameter Error", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
		cnf.Uid = uuid.New()
		cnf.Terminal = &terminal
		if analogue.ChannelNumber <= 12 && analogue.ChannelNumber >=1 {
			cnf.ChannelType = int32(models.CHANNEL_TYPE_TEMPERATURE)
			cnf.ChannelNumber = int32(analogue.ChannelNumber)
		} else if analogue.ChannelNumber <=24 {
			cnf.ChannelType = int32(models.CHANNEL_TYPE_ANALOG)
			cnf.ChannelNumber = int32(analogue.ChannelNumber - 12)
		} else {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("模拟量通道数输入m错误"))
			return
		}

		cnf.IsDefault = false

		//TODO: Default Signed & Threshold
		cnf.Signed = true
		cnf.NegativeThreshold = 32768
		cnf.Parameter = &param
		cnf.Status = analogue.Status
		if cnf.Status == models.CHANNEL_STATUS_SHOW {
			cnf.SequenceNumber = analogue.SequenceNumber
		} else {
			cnf.SequenceNumber = -1
		}

		cnf.Name = cnf.Parameter.Name
		cnf.Length = cnf.Parameter.Length
		if _, err := dba.BoilerOrm.Insert(&cnf); err != nil {
			goazure.Error("Insert Analogue Error:", err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("添加模拟量失败"))
			return
		}
		analoguesql := "insert into issued_analogue_switch(channel_id,create_time,function_id,byte_id,modbus,bit_address) values(?,now(),?,?,?,?)"
		if _, err := dba.BoilerOrm.Raw(analoguesql, cnf.Uid, analogue.Func, analogue.Byte, analogue.Modbus, 0).Exec(); err != nil {
			goazure.Error("Insert issued_analogue Error", err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("添加模拟量功能码等失败"))
			return
		}
	}
	for _, swi := range config.Chan.Switch {
		var max int32
		var param models.RuntimeParameter
		var cnf models.RuntimeParameterChannelConfig
		sql := "select max(param_id) as param_id from runtime_parameter where category_id=? and is_deleted =false"
		if err := dba.BoilerOrm.Raw(sql, models.RUNTIME_PARAMETER_CATEGORY_SWITCH).QueryRow(&max); err != nil {
			goazure.Error("Query runtime_parameter Error", err)
		}
		if max >= 100 {
			param.ParamId = max + 1
		} else {
			param.ParamId = 100
		}
		a := fmt.Sprintf("%d%d", models.RUNTIME_PARAMETER_CATEGORY_SWITCH, param.ParamId)
		param.Id, _ = strconv.ParseInt(a, 10, 64)
		param.Length = 1
		param.Fix = 0
		param.Category = runtimeParameterCategory(models.RUNTIME_PARAMETER_CATEGORY_SWITCH)
		param.Name = swi.Parameter.Name
		//param.Unit = ""
		param.Scale = 1
		param.Medium = runtimeParameterMedium(0)
		param.AddBoilerMedium(0)
		param.Organization = terminal.Organization
		param.CreatedBy = ctl.GetCurrentUser()
		param.UpdatedBy = ctl.GetCurrentUser()
		param.IsDeleted = false
		if _,err:=dba.BoilerOrm.Insert(&param);err!=nil {
			e := fmt.Sprintln("Add/Update Parameter Error", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
		cnf.Uid = uuid.New()
		cnf.Terminal = &terminal
		cnf.ChannelType = int32(models.CHANNEL_TYPE_SWITCH)
		cnf.ChannelNumber = int32(swi.ChannelNumber)
		cnf.IsDefault = false

		//TODO: Default Signed & Threshold
		cnf.Signed = true
		cnf.NegativeThreshold = 32768
		cnf.Parameter = &param
		cnf.Status = swi.Status
		if cnf.Status == models.CHANNEL_STATUS_SHOW {
			cnf.SequenceNumber = swi.SequenceNumber
		} else {
			cnf.SequenceNumber = -1
		}
		if swi.SwitchStatus == 0 {
			swi.SwitchStatus = 1
		}
		cnf.SwitchStatus = swi.SwitchStatus

		cnf.Name = cnf.Parameter.Name
		cnf.Length = cnf.Parameter.Length
		if cnf.ChannelNumber == 1 || cnf.ChannelNumber == 2 {
			switchburnsql := "insert into issued_switch_default(uid,terminal_id,create_time,channel_type,channel_number,function_id,modbus,bit_address) values(uuid(),?,now(),?,?,?,?,?)"
			if _, err := dba.BoilerOrm.Raw(switchburnsql, terminal.Uid, models.CHANNEL_TYPE_SWITCH, swi.ChannelNumber, swi.Func, swi.Modbus, swi.BitAddress).Exec(); err != nil {
				goazure.Error("Insert issued_switch_default Error", err)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("添加开关点火plc失败"))
				return
			}
		} else {
			if _, err := dba.BoilerOrm.Insert(&cnf); err != nil {
				goazure.Error("Insert Switch Error:", err)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("添加开关量失败"))
				return
			}
			analoguesql := "insert into issued_analogue_switch(channel_id,create_time,function_id,byte_id,modbus,bit_address) values(?,now(),?,?,?,?)"
			if _, err := dba.BoilerOrm.Raw(analoguesql, cnf.Uid, swi.Func, 0, swi.Modbus, swi.BitAddress).Exec(); err != nil {
				goazure.Error("Insert issued_analogue Error", err)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("添加模拟量功能码等失败"))
				return
			}
		}
	}
	for _, ran := range config.Chan.Range {
		var max int32
		var param models.RuntimeParameter
		var cnf models.RuntimeParameterChannelConfig
		sql := "select max(param_id) as param_id from runtime_parameter where category_id=? and is_deleted =false"
		if err := dba.BoilerOrm.Raw(sql, models.RUNTIME_PARAMETER_CATEGORY_RANGE).QueryRow(&max); err != nil {
			goazure.Error("Query runtime_parameter Error", err)
		}
		if max >= 100 {
			param.ParamId = max + 1
		} else {
			param.ParamId = 100
		}
		a := fmt.Sprintf("%d%d", models.RUNTIME_PARAMETER_CATEGORY_RANGE, param.ParamId)
		param.Id, _ = strconv.ParseInt(a, 10, 64)
		param.Length = 1
		param.Fix = 0
		param.Category = runtimeParameterCategory(models.RUNTIME_PARAMETER_CATEGORY_RANGE)
		param.Name = ran.Parameter.Name
		//param.Unit = ""
		param.Scale = 1
		param.Medium = runtimeParameterMedium(0)
		param.AddBoilerMedium(0)
		param.Organization = terminal.Organization
		param.CreatedBy = ctl.GetCurrentUser()
		param.UpdatedBy = ctl.GetCurrentUser()
		param.IsDeleted = false
		if _,err:=dba.BoilerOrm.Insert(&param);err!=nil {
			e := fmt.Sprintln("Add/Update Parameter Error", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
		cnf.Uid = uuid.New()
		cnf.Terminal = &terminal
		cnf.ChannelType = int32(models.CHANNEL_TYPE_RANGE)
		cnf.ChannelNumber = int32(ran.ChannelNumber)
		cnf.IsDefault = false

		//TODO: Default Signed & Threshold
		cnf.Signed = true
		cnf.NegativeThreshold = 32768
		cnf.Parameter = &param
		cnf.Status = ran.Status
		if cnf.Status == models.CHANNEL_STATUS_SHOW {
			cnf.SequenceNumber = ran.SequenceNumber
		} else {
			cnf.SequenceNumber = -1
		}

		cnf.Name = cnf.Parameter.Name
		cnf.Length = cnf.Parameter.Length
		if _, err := dba.BoilerOrm.Insert(&cnf); err != nil {
			goazure.Error("Insert Range Error:", err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("添加状态量失败"))
			return
		}
		for _, r := range ran.Ranges {
			var rg models.RuntimeParameterChannelConfigRange
			if r.Max < r.Min {
				e := fmt.Sprintln("状态值范围设定有误！")
				goazure.Error(e)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte(e))
				return
			}
			rg.Uid = uuid.New()
			rg.ChannelConfig = &cnf
			rg.Min = r.Min
			rg.Max = r.Max
			rg.Name = r.Name
			rg.Value = r.Value
			if _,err:=dba.BoilerOrm.Insert(&rg);err!=nil{
				goazure.Error("insert channel_config_range error:",err)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("添加状态量状态失败"))
				return
			}
		}
			analoguesql := "insert into issued_analogue_switch(channel_id,create_time,function_id,byte_id,modbus,bit_address) values(?,now(),?,?,?,?)"
			if _, err := dba.BoilerOrm.Raw(analoguesql, cnf.Uid, ran.Func, ran.Byte, ran.Modbus, 0).Exec(); err != nil {
				goazure.Error("Insert issued_analogue Error", err)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("添加状态量功能码等失败"))
				return
			}
	}
	go ParamCtrl.RefreshParameters()
	//插入表通信参数
	sql := "replace into issued_communication(terminal_id,baud_rate_id,data_bit_id,stop_bit_id,check_bit_id,correspond_type_id,sub_address_id,heart_beat_id) values(?,?,?,?,?,?,?,?)"
	if _, err := dba.BoilerOrm.Raw(sql, terminal.Uid, config.Param.BaudRate, config.Param.DataBit, config.Param.StopBit, config.Param.CheckDigit,config.Param.CommunInterface, config.Param.SlaveAddress,config.Param.HeartBeat).Exec(); err != nil {
		goazure.Error("Insert issued_communication Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("添加通信参数失败"))
		return
	}
	var param []orm.Params
	var ver int32
	verSql := "select ver from issued_message where sn=?"
	if num, err := dba.BoilerOrm.Raw(verSql, terminal.TerminalCode).Values(&param); err != nil || num == 0 {
		goazure.Error("Query issued_message Error", err)
		ver = 1
	} else {
		a := fmt.Sprintf("%s", param[0]["ver"])
		v, err := strconv.Atoi(a)
		ver = int32(v) + 1
		if ver >= 32769 {
			ver = 1
		}
		if err != nil {
			goazure.Error("ParseInt Error")
			return
		}
	}
	Byte := IssuedCtl.IssuedMessage(terminal.Uid)
	Code := strconv.FormatInt(terminal.TerminalCode, 10)
	message := SocketCtrl.c0(Byte, Code, ver)
	messageSql := "insert into issued_message(sn,ver,create_time,update_time,curr_message) values(?,?,now(),now(),?) on duplicate key update ver=?,update_time=now(),curr_message=?"
	if _, err := dba.BoilerOrm.Raw(messageSql, terminal.TerminalCode, ver, string(message), ver, string(message)).Exec(); err != nil {
		goazure.Error("Insert issued_message Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("组包失败"))
		return
	}
}
