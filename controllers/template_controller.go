package controllers

import (
	"github.com/AzureTech/goazure"
	"encoding/json"
	"fmt"
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/pborman/uuid"
	"github.com/AzureRelease/boiler-server/models"
	"strconv"
	"github.com/AzureTech/goazure/orm"
	"github.com/AzureRelease/boiler-server/conf"
	"bytes"
)

type TemplateController struct {
	MainController
}

var TempCtrl *TemplateController = &TemplateController{}
type Chan struct {
	ParameterId		int			`json:"parameter_id"`
	ChannelType		int			`json:"channel_type"`
	ChannelNumber	int			`json:"channel_number"`
	FcodeId       int          `json:"fcodeName"`  //功能码
	BitAddress      int          `json:"bitAddress"`  //位地址
	TermByteId        int          `json:"termByte"`    //高低字节
	Modbus          int          `json:"modbus"`      //modbus
	Status			int32		`json:"status"`
	SequenceNumber	int32		`json:"sequence_number"`
	SwitchValue		int32		`json:"switch_status"`
	Ranges			[]bRange 	`json:"ranges"`
}

type TemplateUpdate struct {
	Chan   []Chan            `json:"chan"`
	Uid    string            `json:uid`
	Name          string      `json:"name"`
	Param     Param          `json:"param"`
}

type IssuedTemplateUpdate struct {
	TemplateUpdate    TemplateUpdate   `json:"TemplateUpdate"`
}

type Template struct {
	Uid             string          `json:"uid"`
	Name            string          `json:"name"`
	OrganizationUid string          `json:"organizationUid"`
	Channel         []ChannelIssued `json:"channel"`
	Param           Param           `json:"param"`
}

type TemplateConfig struct {
	Terminals    []models.Terminal
	TemplateUid string
}
type GroupConfig struct {
	Start string    `json:"start"`
	End   string	`json:"end"`
	Template string		`json:"template"`
}
type TemplateGroupConfig struct {
	GroupConfig  []GroupConfig  `json:"groupConfig"`
}
type TempReturnErr struct {
	TerminalCode string
	MsgErr string
}

func (ctl *TemplateController) FastTemplateToCur() {
	usr := ctl.GetCurrentUser()
	var wTempl TempToCur
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &wTempl); err != nil {
		goazure.Error("Unmarshal JSON Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		return
	}
	var tempConfig TemplateConfig
	var terminal models.Terminal
	qs := dba.BoilerOrm.QueryTable("terminal")
	if usr.IsOrganizationUser() {
		qs = qs.Filter("Organization__Uid", usr.Organization.Uid)
	}
	if err:=qs.Filter("IsDeleted", false).Filter("TerminalCode",wTempl.Code).One(&terminal);err!=nil{
		goazure.Error("Query TerminalCode Error:",err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("非法终端"))
		return
	}
	var terminals []models.Terminal
	terminals = append(terminals,terminal )
	tempConfig.Terminals = terminals
	tempConfig.TemplateUid = wTempl.TemplateUid
	fmt.Println("cccccc:",tempConfig.Terminals[0].TerminalCode)
	fmt.Println("tempalteUid:",tempConfig.TemplateUid)
	go ctl.IssuedTemplateToCurr(tempConfig)
}

func (ctl *TemplateController) TemplateGroupConfig() {
	usr:=ctl.GetCurrentUser()
	var terminals []models.Terminal
	var tempConfig TemplateConfig
	var tempGroupConfig TemplateGroupConfig
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &tempGroupConfig); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Comment Error", err)
		return
	} else {
		if conf.BatchFlag {
			ctl.Ctx.Output.SetStatus(200)
			ctl.Ctx.Output.Body([]byte("正在配置中..."))
		} else {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("服务器忙，还在处理上次批量下发！"))
			return
		}
	}
	go func () {
		conf.BatchFlag = false
		for _,c := range tempGroupConfig.GroupConfig {
			start,err:=strconv.Atoi(c.Start)
			if err!=nil {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("终端编号输入错误!"))
				goazure.Error("Format Int Error",err)
			}
			end,error := strconv.Atoi(c.End)
			if error!=nil {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("终端编号输入错误!"))
				goazure.Error("Format Int Error",err)
			}
			qs := dba.BoilerOrm.QueryTable("terminal")
			if usr.IsOrganizationUser() {
				qs = qs.Filter("Organization__Uid", usr.Organization.Uid)
			}
			if _,err := qs.Filter("IsDeleted", false).Filter("TerminalCode__gte",start).Filter("TerminalCode__lte",end).OrderBy("-UpdatedDate").All(&terminals);err!=nil{
				goazure.Error("Query terminal Error",err)
			}
			tempConfig.Terminals = terminals
			tempConfig.TemplateUid = c.Template
			fmt.Println("需要更新的终端:",tempConfig.Terminals)
			fmt.Println("更新的模板：",tempConfig.TemplateUid)
			ctl.IssuedTemplateToCurr(tempConfig)
		}
		conf.BatchFlag =true
	}()

}

//组装模板的模拟量一报文
func (ctl *TemplateController) IssuedTemplateAnalogOne(Uid string) ([]byte) {
	var temp = 1
	Byte := make([]byte, 0)
	var templateAnalog []models.IssuedChannelConfigTemplate
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").
		RelatedSel("Func").RelatedSel("Byte").
		Filter("Template__Uid", Uid).Filter("ChannelType", models.CHANNEL_TYPE_TEMPERATURE).OrderBy("ChannelNumber").All(&templateAnalog); err != nil {
		goazure.Error("insert issued_channel_config_template", err)
	}
	fmt.Println("templateAnalog：", len(templateAnalog))
	L := len(templateAnalog)
	if L == 0 {
		for c := 0; c < 12; c++ {
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteTwo(0)...)
		}
	} else {
		for i := 0; i < L; i++ {
			value := templateAnalog[i]
			if temp == int(value.ChannelNumber) {
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteOne(int32(value.Byte.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
			} else {
				for c := 0; c < (int(value.ChannelNumber) - temp); c++ {
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteTwo(0)...)
				}
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteOne(int32(value.Byte.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				temp = int(value.ChannelNumber)
			}
			if i == L-1 {
				if temp != 12 {
					for c := 0; c < (12 - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
				}
			}
			temp++
		}
	}
	return Byte
}

//组模拟量通道二
func (ctl *TemplateController) IssuedTemplateAnalogTwo(Uid string) ([]byte) {
	var temp = 1
	Byte := make([]byte, 0)
	var templateAnalog []models.IssuedChannelConfigTemplate
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").
		RelatedSel("Func").RelatedSel("Byte").
		Filter("Template__Uid", Uid).Filter("ChannelType", models.CHANNEL_TYPE_ANALOG).OrderBy("ChannelNumber").All(&templateAnalog); err != nil {
		goazure.Error("insert issued_channel_config_template", err)
	}
	fmt.Println("templateAnalog：", len(templateAnalog))
	L := len(templateAnalog)
	if L == 0 {
		for c := 0; c < 12; c++ {
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteTwo(0)...)
		}
	} else {
		for i := 0; i < L; i++ {
			value := templateAnalog[i]
			if temp == int(value.ChannelNumber) {
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteOne(int32(value.Byte.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
			} else {
				for c := 0; c < (int(value.ChannelNumber) - temp); c++ {
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteTwo(0)...)
				}
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteOne(int32(value.Byte.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				temp = int(value.ChannelNumber)
			}
			if i == L-1 {
				if temp != 12 {
					for c := 0; c < (12 - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
				}
			}
			temp++
		}
	}
	return Byte
}

//组开关量
func (ctl *TemplateController) IssuedTemplateSwitch(Uid string) ([]byte) {
	var temp = 1
	Byte := make([]byte, 0)
	var templateAnalog []models.IssuedChannelConfigTemplate
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").
		RelatedSel("Func").
		Filter("Template__Uid", Uid).Filter("ChannelType", models.CHANNEL_TYPE_SWITCH).OrderBy("ChannelNumber").All(&templateAnalog); err != nil {
		goazure.Error("insert issued_channel_config_template", err)
	}
	fmt.Println("templateAnalog：", len(templateAnalog))
	L := len(templateAnalog)
	if L == 0 {
		for c := 0; c < 48; c++ {
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteTwo(0)...)
		}
	} else {
		for i := 0; i < L; i++ {
			value := templateAnalog[i]
			if temp == int(value.ChannelNumber) {
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				Byte = append(Byte, IntToByteOne(int32(value.BitAddress))...)
			} else {
				for c := 0; c < (int(value.ChannelNumber) - temp); c++ {
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteTwo(0)...)
					Byte = append(Byte, IntToByteOne(0)...)
				}
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				Byte = append(Byte, IntToByteOne(int32(value.BitAddress))...)
				temp = int(value.ChannelNumber)
			}
			if i == L-1 {
				if temp != 48 {
					for c := 0; c < (48 - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
					}
				}
			}
			temp++
		}
	}
	return Byte
}

//组装状态量
func (ctl *TemplateController) IssuedTemplateRange(Uid string) ([]byte) {
	var temp = 1
	Byte := make([]byte, 0)
	var templateAnalog []models.IssuedChannelConfigTemplate
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").
		RelatedSel("Func").RelatedSel("Byte").
		Filter("Template__Uid", Uid).Filter("ChannelType", models.CHANNEL_TYPE_RANGE).OrderBy("ChannelNumber").All(&templateAnalog); err != nil {
		goazure.Error("insert issued_channel_config_template", err)
	}
	fmt.Println("templateAnalog：", len(templateAnalog))
	L := len(templateAnalog)
	if L == 0 {
		for c := 0; c < 12; c++ {
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteTwo(0)...)
		}
	} else {
		for i := 0; i < L; i++ {
			value := templateAnalog[i]
			if temp == int(value.ChannelNumber) {
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteOne(int32(value.Byte.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
			} else {
				for c := 0; c < (int(value.ChannelNumber) - temp); c++ {
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteOne(0)...)
					Byte = append(Byte, IntToByteTwo(0)...)
				}
				Byte = append(Byte, IntToByteOne(int32(value.Func.Value))...)
				Byte = append(Byte, IntToByteOne(int32(value.Byte.Value))...)
				Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				temp = int(value.ChannelNumber)
			}
			if i == L-1 {
				if temp != 12 {
					for c := 0; c < (12 - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
				}
			}
			temp++
		}
	}
	return Byte
}

//组通信参数
func (ctl *TemplateController) IssuedTemplateCommunication(Uid string) ([]byte) {
	Byte := make([]byte, 0)
	var commuTemp models.IssuedCommunicationTemplate
	if err := dba.BoilerOrm.QueryTable("issued_communication_template").
		RelatedSel("BaudRate").RelatedSel("DataBit").RelatedSel("StopBit").RelatedSel("CheckBit").
		RelatedSel("CorrespondType").RelatedSel("SubAddress").RelatedSel("HeartBeat").Filter("Template__Uid", Uid).One(&commuTemp); err != nil {
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
		goazure.Error("Query issued_communication_template Error", err)
	} else {
		Byte = append(Byte, IntToByteOne(int32(commuTemp.BaudRate.Value))...)
		Byte = append(Byte, IntToByteOne(int32(commuTemp.DataBit.Value))...)
		Byte = append(Byte, IntToByteOne(int32(commuTemp.StopBit.Value))...)
		Byte = append(Byte, IntToByteOne(int32(commuTemp.CheckBit.Value))...)
		Byte = append(Byte, IntToByteOne(int32(commuTemp.CorrespondType.Value))...)
		Byte = append(Byte, IntToByteOne(int32(commuTemp.SubAddress.Value))...)
		Byte = append(Byte, IntToByteOne(int32(commuTemp.HeartBeat.Value))...)
	}
	return Byte
}

//组装模板报文
func (ctl *TemplateController) IssuedTemplateMessage(Uid string) ([]byte) {
	Byte := make([]byte, 0)
	Byte = append(Byte, ctl.IssuedTemplateAnalogOne(Uid)...)
	fmt.Println("模拟量1长度:", len(Byte))
	Byte = append(Byte, ctl.IssuedTemplateAnalogTwo(Uid)...)
	fmt.Println("模拟量2长度:", len(Byte))
	Byte = append(Byte, ctl.IssuedTemplateSwitch(Uid)...)
	fmt.Println("开关位长度:", len(Byte))
	Byte = append(Byte, ctl.IssuedTemplateRange(Uid)...)
	fmt.Println("状态量长度:", len(Byte))
	Byte = append(Byte, ctl.IssuedTemplateCommunication(Uid)...)
	fmt.Println("数据长度:", len(Byte))
	return Byte
}

func (ctl *TemplateController) TemplateChannelConfigDelete(cnf *models.RuntimeParameterChannelConfig, uid string) error {
	if cnf.ChannelType == models.CHANNEL_TYPE_RANGE {
		if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config_range").
			Filter("ChannelConfig__Uid", cnf.Uid).
			Delete(); err != nil {
			goazure.Error("Delete Ranges Error:", err, num)
		}
	}
	if !(cnf.ChannelType == models.CHANNEL_TYPE_SWITCH && (cnf.ChannelNumber == 1 || cnf.ChannelNumber == 2)){
		if num, err := dba.BoilerOrm.QueryTable("issued_analogue_switch").Filter("Channel", cnf.Uid).Delete(); err != nil {
			goazure.Error("Delete issued analogue Error:", err, num)
			return err
		}
	}
	if num, err := dba.BoilerOrm.Delete(cnf); err != nil {
		goazure.Error("Delete Channel Config Error:", err, num)
		return err
	}

	return nil
}

func (ctl *TemplateController) IssuedTemplateToCurr(temp TemplateConfig) {
	usr:=ctl.GetCurrentUser()
	ip:= IssuedCtl.IssuedGetIp(ctl.Ctx.Request.RemoteAddr)
	Byte:=ctl.IssuedTemplateMessage(temp.TemplateUid)
	var configTemp []models.IssuedChannelConfigTemplate
	var confSwitchTemps []models.IssuedChannelConfigTemplate
	var commTemp models.IssuedCommunicationTemplate
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").
		Filter("Template", temp.TemplateUid).Filter("ChannelType__in",models.CHANNEL_TYPE_TEMPERATURE,models.CHANNEL_TYPE_ANALOG,models.CHANNEL_TYPE_RANGE).
		OrderBy("ChannelType").All(&configTemp); err != nil {
		goazure.Error("Query issued channel config template Error", err)
		return
	}
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").
		Filter("Template", temp.TemplateUid).Filter("ChannelType",models.CHANNEL_TYPE_SWITCH).
		OrderBy("ChannelNumber").All(&confSwitchTemps); err != nil {
		goazure.Error("Query issued channel config template Error", err)
		return
	}
	if err:=dba.BoilerOrm.QueryTable("issued_communication_template").Filter("Template__Uid",temp.TemplateUid).One(&commTemp);err!=nil{
		goazure.Error("Query issued_communication_template Error",err)
		return
	}
	for _, t := range temp.Terminals {
		var aCnf []models.RuntimeParameterChannelConfig
		if _, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
			Filter("Terminal__Uid", t.Uid).Filter("IsDefault", false).
			All(&aCnf); err != nil {
			goazure.Error("Get ChannelConfig To Delete Error:", err)
		}
		for _, a := range aCnf {
			ctl.TemplateChannelConfigDelete(&a, t.Uid)
		}
		if num, err := dba.BoilerOrm.QueryTable("issued_switch_default").Filter("Terminal__Uid", t.Uid).Delete(); err != nil {
			goazure.Error("Delete issued switch Error:", err, num)
		}
		//**删除结束**
		var cnf models.RuntimeParameterChannelConfig
		for _, c := range configTemp {
			cnf.Terminal = &t
			cnf.ChannelType = c.ChannelType
			cnf.ChannelNumber = c.ChannelNumber
			cnf.IsDefault = false
			//TODO: Default Signed & Threshold
			cnf.Signed = true
			cnf.NegativeThreshold = 32768
			cnf.Parameter = c.Parameter
			cnf.Status = c.Status
			if cnf.Status == models.CHANNEL_STATUS_SHOW {
				cnf.SequenceNumber = c.SequenceNumber
			} else {
				cnf.SequenceNumber = -1
			}
			cnf.Name = cnf.Parameter.Name
			cnf.Length = cnf.Parameter.Length
			cnf.Uid=uuid.New()
			if cnf.ChannelType == models.CHANNEL_TYPE_SWITCH {
				if c.SwitchStatus == 0 {
					c.SwitchStatus = 1
				}
				cnf.SwitchStatus = c.SwitchStatus
			}
				if _,err:=dba.BoilerOrm.Insert(&cnf);err!=nil{
					goazure.Error("insert runtime parameter channel config Error",err)
				}
				analoguesql := "insert into issued_analogue_switch(channel_id,create_time,function_id,byte_id,modbus,bit_address) values(?,now(),?,?,?,?)"
				if _, err := dba.BoilerOrm.Raw(analoguesql, cnf.Uid, c.Func.Id, c.Byte.Id, c.Modbus,c.BitAddress).Exec(); err != nil {
					goazure.Error("Insert issued_analogue Error", err)
				}

			if cnf.ChannelType == models.CHANNEL_TYPE_RANGE {
				var configRangeTemp []models.IssuedChannelConfigRangeTemplate
				if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_range_template").Filter("ChannelConfig", c.Uid).OrderBy("Value").All(&configRangeTemp); err != nil {
					goazure.Error("Query issued channel config range template Error", err)
				}
				var aRanges []*models.RuntimeParameterChannelConfigRange
				for i, r := range configRangeTemp {
					var rg models.RuntimeParameterChannelConfigRange
					if r.Max < r.Min {
						e := fmt.Sprintln("状态值范围设定有误！")
						goazure.Error(e)
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte(e))
					}
					for _, ar := range aRanges {
						if ar.Max >= r.Min {
							e := fmt.Sprintln("状态值范围设定有误，请不要设定重复的范围区间！")
							goazure.Error(e)
							ctl.Ctx.Output.SetStatus(400)
							ctl.Ctx.Output.Body([]byte(e))
						}
					}
					rg.ChannelConfig = &cnf
					rg.Min = r.Min
					rg.Max = r.Max
					rg.Name = r.Name
					rg.Value = int64(i)

					if err := DataCtl.AddData(&rg, true); err != nil {
						goazure.Error("Added ChannelConfig Range Error", rg, err)
					}
				}
			}
		}
		for _, c := range confSwitchTemps {
			cnf.Terminal = &t
			cnf.ChannelType = c.ChannelType
			cnf.ChannelNumber = c.ChannelNumber
			cnf.IsDefault = false
			//TODO: Default Signed & Threshold
			cnf.Signed = true
			cnf.NegativeThreshold = 32768
			cnf.Parameter = c.Parameter
			cnf.Status = c.Status
			if cnf.Status == models.CHANNEL_STATUS_SHOW {
				cnf.SequenceNumber = c.SequenceNumber
			} else {
				cnf.SequenceNumber = -1
			}
			cnf.Name = cnf.Parameter.Name
			cnf.Length = cnf.Parameter.Length
			cnf.Uid=uuid.New()
			if cnf.ChannelType == models.CHANNEL_TYPE_SWITCH {
				if c.SwitchStatus == 0 {
					c.SwitchStatus = 1
				}
				cnf.SwitchStatus = c.SwitchStatus
			}
			if cnf.ChannelType==models.CHANNEL_TYPE_SWITCH && (cnf.ChannelNumber == 1 || cnf.ChannelNumber ==2) {
				switchburnsql := "insert into issued_switch_default(uid,terminal_id,create_time,channel_type,channel_number,function_id,modbus,bit_address) values(uuid(),?,now(),?,?,?,?,?)"
				if _, err := dba.BoilerOrm.Raw(switchburnsql, t.Uid, c.ChannelType,c.ChannelNumber,c.Func.Id, c.Modbus, c.BitAddress).Exec(); err != nil {
					goazure.Error("Insert issued_switch_default Error", err)
				}
			} else {
				if _,err:=dba.BoilerOrm.Insert(&cnf);err!=nil{
					goazure.Error("insert runtime parameter channel config Error",err)
				}
				analoguesql := "insert into issued_analogue_switch(channel_id,create_time,function_id,byte_id,modbus,bit_address) values(?,now(),?,?,?,?)"
				if _, err := dba.BoilerOrm.Raw(analoguesql, cnf.Uid, c.Func.Id, c.Byte.Id, c.Modbus,c.BitAddress).Exec(); err != nil {
					goazure.Error("Insert issued_analogue Error", err)
				}
			}
		}
		//插入通信参数
		commSql:="replace into issued_communication(terminal_id,baud_rate_id,data_bit_id,stop_bit_id,check_bit_id,correspond_type_id,sub_address_id,heart_beat_id) values(?,?,?,?,?,?,?,?)"
		if _, er := dba.BoilerOrm.Raw(commSql, t.Uid, commTemp.BaudRate.Id, commTemp.DataBit.Id, commTemp.StopBit.Id, commTemp.CheckBit.Id, commTemp.CorrespondType.Id, commTemp.SubAddress.Id, commTemp.HeartBeat.Id).Exec(); er != nil {
			goazure.Error("Insert issued_communication Error", er)
		}
		fmt.Println("插入通信参数结束")
		//插入报文
		var param []orm.Params
		var ver int32
		fmt.Println("code:",t.TerminalCode)
		code:=strconv.FormatInt(t.TerminalCode,10)
		verSql:="select ver from issued_message where sn=?"
		if num,err:=dba.BoilerOrm.Raw(verSql,t.TerminalCode).Values(&param);err!=nil || num==0 {
			goazure.Error("Query issued_message Error",err)
			ver = 1
		} else {
			a:=fmt.Sprintf("%s",param[0]["ver"])
			v,err:=strconv.Atoi(a)
			ver=int32(v) + 1
			if ver >=32769 {
				ver = 1
			}
			if err!=nil {
				goazure.Error("ParseInt Error")
			}
		}
		message:=SocketCtrl.c0(Byte,code,ver)
		messageSql:="insert into issued_message(sn,ver,create_time,update_time,curr_message) values(?,?,now(),now(),?) on duplicate key update ver=?,update_time=now(),curr_message=?"
		if _,err:=dba.BoilerOrm.Raw(messageSql,code,ver,string(message),ver,string(message)).Exec();err!=nil{
			goazure.Error("Insert issued_message Error",err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("组包失败"))
		}
		//插入配置的模板信息
		tempStatusSql:="insert into issued_term_temp_status(sn,create_time,update_time,template_id) values(?,now(),now(),?) on duplicate key update update_time=now(),template_id=?"
		if _,err:=dba.BoilerOrm.Raw(tempStatusSql,code,temp.TemplateUid,temp.TemplateUid).Exec();err!=nil{
			goazure.Error("Insert issued_term_temp_status Error",err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("存入失败"))
		}
		/*reqBuf :=IssuedCtl.ReqMessage(code)
		if reqBuf == "" {
			return
		}
		if len(reqBuf) != 362 {
			return
		}
		SocketCtrl.SocketBatchConfigSend(reqBuf)*/
		go func(){
				reqBuf :=IssuedCtl.ReqMessage(code)
				if reqBuf == "" {
					goazure.Error("报文为空")
					return
				}
				if len(reqBuf) != 362 {
					goazure.Error("报文长度错误")
					return
				}
			if temp:=IssuedCtl.IssuedVerController(code);!temp{
				goazure.Error(code+"终端版本号低，不支持改功能!")
				return
			}
			if conf.ContentLogsFlag {
				IssuedCtl.IssuedContentLogs(usr.Username,ip,code,conf.TermConfig,temp.TemplateUid+":"+fmt.Sprintf("%x",reqBuf))
			} else {
				IssuedCtl.IssuedContentLogs(usr.Username,ip,code,conf.TermConfig,temp.TemplateUid)
			}
				buf:=SocketCtrl.SocketConfigSend(reqBuf)
			if buf == nil {
				goazure.Error("发送报文失败")
				return
			} else if bytes.Equal(conf.TermNoRegist, buf) {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("终端还未连接平台"))
				return
			} else if bytes.Equal(conf.TermTimeout, buf) {
				goazure.Error("终端返回信息超时")
				return
			} else if len(buf) > 4 {
				switch buf[15] {
				case 16:
					newVer := ByteToIntTwo(buf[13:15])
					termCode := fmt.Sprintf("%s", buf[7:13])
					sql := "insert into issued_version(sn,ver,create_time,update_time) values(?,?,now(),now()) on duplicate key update ver=?,update_time=now()"
					if _, err := dba.BoilerOrm.Raw(sql, termCode, newVer, newVer).Exec(); err != nil {
						goazure.Error("insert issued_version Error", err)
					}
					fmt.Println(termCode+"终端插入终端版本成功")
				case 1:
					termCode := fmt.Sprintf("%s", buf[7:13])
					goazure.Error(termCode+"终端CRC校验错误")
				case 2:
					termCode := fmt.Sprintf("%s", buf[7:13])
					goazure.Error(termCode+"终端SN不一致")
				case 3:
					termCode := fmt.Sprintf("%s", buf[7:13])
					goazure.Error(termCode+"终端配置错误")
				case 4:
					termCode := fmt.Sprintf("%s", buf[7:13])
					goazure.Error(termCode+"终端PLC未连接")
				default:
					termCode := fmt.Sprintf("%s", buf[7:13])
					goazure.Error(termCode+"终端配置错误")
				}
			} else {
				termCode := fmt.Sprintf("%s", buf[7:13])
				goazure.Error(termCode+"终端返回报文信息错误")
			}
		}()
	}
}

func (ctl *TemplateController) IssuedTemplate() {
	var template Template
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &template); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("Name:", template.Name)
	fmt.Println("OUid:", template.OrganizationUid)
	template.Uid = uuid.New()
	sql := "insert into issued_template(uid,name,create_time,update_time,organization_id) values(?,?,now(),now(),?)"
	if _, err := dba.BoilerOrm.Raw(sql, template.Uid, template.Name, template.OrganizationUid).Exec(); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("数据库异常!"))
		goazure.Error("insert issued_template Error", err)
		return
	}
	commuSql := "insert into issued_communication_template(uid,template_id,baud_rate_id,data_bit_id,stop_bit_id,check_bit_id,correspond_type_id,sub_address_id,heart_beat_id) values(uuid(),?,?,?,?,?,?,?,?)"
	if _, err := dba.BoilerOrm.Raw(commuSql, template.Uid, template.Param.BaudRateId, template.Param.DataBitId, template.Param.StopBitId, template.Param.CheckDigitId, template.Param.CommunInterfaceId, template.Param.SubAddressId, template.Param.HeartBitId).Exec(); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("数据库异常!"))
		goazure.Error("insert issued_template Error", err)
		return
	}
	for i, c := range template.Channel {
		goazure.Warn("[", i, "]ChannelIssuedUpdate: ", c)
		if c.ParameterId <= 0 {
			continue
		} else {
			channelUid := uuid.New()
			channelSql := "insert into issued_channel_config_template(uid,create_time,parameter_id,template_id,channel_type,channel_number,status,sequence_number,switch_status,func_id,byte_id,bit_address,modbus) " +
				"values(?,now(),?,?,?,?,?,?,?,?,?,?,?)"
			if _, err := dba.BoilerOrm.Raw(channelSql, channelUid, c.ParameterId, template.Uid, c.ChannelType, c.ChannelNumber, c.Status, c.SequenceNumber, c.SwitchValue, c.FcodeId, c.TermByteId, c.BitAddress, c.Modbus).Exec(); err != nil {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("数据库异常!"))
				goazure.Error("insert issued_channel_config_template Error", err)
				return
			}

			if c.ChannelType == models.CHANNEL_TYPE_RANGE {
				for i, r := range c.Ranges {
					if r.Max < r.Min {
						e := fmt.Sprintln("状态值范围设定有误！")
						goazure.Error(e)
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte(e))
						return
					}
					rangeUid := uuid.New()
					rangeSql := "insert into issued_channel_config_range_template(uid,name,create_time,channel_config_id,min,max,value) " +
						"values(?,?,now(),?,?,?,?)"
					if _, err := dba.BoilerOrm.Raw(rangeSql, rangeUid, r.Name, channelUid, r.Min, r.Max, i).Exec(); err != nil {
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte("数据库异常!"))
						goazure.Error("insert issued_channel_config_range_template Error", err)
						return
					}
				}
			}

		}
	}
}

//模板列表
func (ctl *TemplateController) TemplateList() {
	var template []models.IssuedTemplate
	usr := ctl.GetCurrentUser()
	qs := dba.BoilerOrm.QueryTable("issued_template")
	qs = qs.RelatedSel("Organization")
	if usr.IsOrganizationUser() {
		qs = qs.Filter("Organization__Uid", usr.Organization.Uid)
	}
	if num, err := qs.Filter("IsDeleted", false).OrderBy("-UpdateTime").All(&template); err != nil {
		fmt.Printf("Returned Rows Num: %d, %s", num, err)
	}
	ctl.Data["json"] = template
	ctl.ServeJSON()
}

//获取模拟量通道一
func (ctl *TemplateController) TemplateAnalogOne() {
	var temp models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：", temp)
	if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").RelatedSel("Template").Filter("Template", temp.Uid).Filter("ChannelType", models.CHANNEL_TYPE_TEMPERATURE).All(&cTemps); err != nil {
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix := make([]models.IssuedChannelConfigTemplate, 12)
	fmt.Println("channel:", cTemps)
	for _, c := range cTemps {
		matrix[c.ChannelNumber-1] = c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取模拟量通道二
func (ctl *TemplateController) TemplateAnalogTwo() {
	var temp models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：", temp)
	if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").RelatedSel("Template").Filter("Template", temp.Uid).Filter("ChannelType", models.CHANNEL_TYPE_ANALOG).All(&cTemps); err != nil {
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix := make([]models.IssuedChannelConfigTemplate, 12)
	fmt.Println("channel:", cTemps)
	for _, c := range cTemps {
		matrix[c.ChannelNumber-1] = c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取开关量一
func (ctl *TemplateController) TemplateSwitchOne() {
	var temp models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：", temp)
	if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Template").Filter("Template", temp.Uid).Filter("ChannelType", models.CHANNEL_TYPE_SWITCH).Filter("ChannelNumber__gte", 1).Filter("ChannelNumber__lte", 16).All(&cTemps); err != nil {
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix := make([]models.IssuedChannelConfigTemplate, 16)
	fmt.Println("channel:", cTemps)
	for _, c := range cTemps {
		matrix[c.ChannelNumber-1] = c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取开关量二
func (ctl *TemplateController) TemplateSwitchTwo() {
	var temp models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：", temp)
	if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Template").Filter("Template", temp.Uid).Filter("ChannelType", models.CHANNEL_TYPE_SWITCH).Filter("ChannelNumber__gte", 17).Filter("ChannelNumber__lte", 32).All(&cTemps); err != nil {
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix := make([]models.IssuedChannelConfigTemplate, 16)
	fmt.Println("channel:", cTemps)
	for _, c := range cTemps {
		matrix[(c.ChannelNumber-1)%16] = c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取开关量三
func (ctl *TemplateController) TemplateSwitchThree() {
	var temp models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：", temp)
	if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Template").Filter("Template", temp.Uid).Filter("ChannelType", models.CHANNEL_TYPE_SWITCH).Filter("ChannelNumber__gte", 33).Filter("ChannelNumber__lte", 48).All(&cTemps); err != nil {
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix := make([]models.IssuedChannelConfigTemplate, 16)
	fmt.Println("channel:", cTemps)
	for _, c := range cTemps {
		matrix[(c.ChannelNumber-1)%16] = c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取状态量
func (ctl *TemplateController) TemplateRange() {
	var temp models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：", temp)
	if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").RelatedSel("Template").Filter("Template", temp.Uid).Filter("ChannelType", models.CHANNEL_TYPE_RANGE).All(&cTemps); err != nil {
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix := make([]models.IssuedChannelConfigTemplate, 12)
	fmt.Println("channel:", cTemps)
	for _, c := range cTemps {
		var ranges []*models.IssuedChannelConfigRangeTemplate
		if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_range_template").Filter("ChannelConfig__Uid", c.Uid).OrderBy("Value").All(&ranges); err != nil {
			goazure.Error("Get Template ChannelConfig Range Error:", err, num)
		} else {
			c.Ranges = ranges
		}
		matrix[c.ChannelNumber-1] = c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取通信参数
func (ctl *TemplateController) TemplateCommunication() {
	var temp models.IssuedTemplate
	var cTemps models.IssuedCommunicationTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	if err := dba.BoilerOrm.QueryTable("issued_communication_template").RelatedSel("BaudRate").RelatedSel("CheckBit").RelatedSel("CorrespondType").RelatedSel("DataBit").RelatedSel("HeartBeat").RelatedSel("StopBit").RelatedSel("SubAddress").RelatedSel("Template").Filter("Template__Uid", temp.Uid).One(&cTemps); err != nil {
		goazure.Error("Query issued_communication_template Error:", err)
	}
	ctl.Data["json"] = cTemps
	ctl.ServeJSON()
}

//修改模板信息
func (ctl *TemplateController) TemplateUpdate() {
	var template IssuedTemplateUpdate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &template); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("template:",template)
	var confRTemp []models.IssuedChannelConfigTemplate
	if _,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").Filter("Template__Uid",template.TemplateUpdate.Uid).Filter("ChannelType",models.CHANNEL_TYPE_RANGE).All(&confRTemp);err!=nil{
		goazure.Error("Query issued_channel_config_template Error",err)
	}
	for _,c := range confRTemp {
		if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_range_template").Filter("ChannelConfig__Uid",c.Uid).Delete();err!=nil{
			goazure.Error("Delete Ranges Error:", err, num)
		}
	}
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").Filter("Template__Uid",template.TemplateUpdate.Uid).Delete();err!=nil{
		goazure.Error("Delete Ranges Error:", err, num)
	}
	if _,err := dba.BoilerOrm.QueryTable("issued_communication_template").Filter("Template__Uid",template.TemplateUpdate.Uid).Delete();err!=nil{
		goazure.Error("Query IssuedCommunicationTemplate Error", err)
	}
	//插入数据
	for i, c := range template.TemplateUpdate.Chan {
		goazure.Warn("[", i, "]ChannelIssuedUpdate: ", c)
		if c.ParameterId <= 0 {
			continue
		} else {
			channelUid := uuid.New()
			channelSql := "insert into issued_channel_config_template(uid,create_time,parameter_id,template_id,channel_type,channel_number,status,sequence_number,switch_status,func_id,byte_id,bit_address,modbus) " +
				"values(?,now(),?,?,?,?,?,?,?,?,?,?,?)"
			if _, err := dba.BoilerOrm.Raw(channelSql, channelUid, c.ParameterId, template.TemplateUpdate.Uid, c.ChannelType, c.ChannelNumber, c.Status, c.SequenceNumber, c.SwitchValue, c.FcodeId, c.TermByteId, c.BitAddress, c.Modbus).Exec(); err != nil {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("数据库异常!"))
				goazure.Error("insert issued_channel_config_template Error", err)
				return
			}

			if c.ChannelType == models.CHANNEL_TYPE_RANGE {
				for i, r := range c.Ranges {
					if r.Max < r.Min {
						e := fmt.Sprintln("状态值范围设定有误！")
						goazure.Error(e)
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte(e))
						return
					}
					rangeUid := uuid.New()
					rangeSql := "insert into issued_channel_config_range_template(uid,name,create_time,channel_config_id,min,max,value) " +
						"values(?,?,now(),?,?,?,?)"
					if _, err := dba.BoilerOrm.Raw(rangeSql, rangeUid, r.Name, channelUid, r.Min, r.Max, i).Exec(); err != nil {
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte("数据库异常!"))
						goazure.Error("insert issued_channel_config_range_template Error", err)
						return
					}
				}
			}

		}
	}
	commuSql := "insert into issued_communication_template(uid,template_id,baud_rate_id,data_bit_id,stop_bit_id,check_bit_id,correspond_type_id,sub_address_id,heart_beat_id) values(uuid(),?,?,?,?,?,?,?,?)"
	if _, err := dba.BoilerOrm.Raw(commuSql, template.TemplateUpdate.Uid, template.TemplateUpdate.Param.BaudRateId, template.TemplateUpdate.Param.DataBitId, template.TemplateUpdate.Param.StopBitId, template.TemplateUpdate.Param.CheckDigitId, template.TemplateUpdate.Param.CommunInterfaceId, template.TemplateUpdate.Param.SubAddressId, template.TemplateUpdate.Param.HeartBitId).Exec(); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("数据库异常!"))
		goazure.Error("insert issued_template Error", err)
		return
	}
	tempSql:="update issued_template set name=?,update_time=now() where uid=? and is_deleted=false"
	if _,err :=dba.BoilerOrm.Raw(tempSql,template.TemplateUpdate.Name,template.TemplateUpdate.Uid).Exec();err!=nil{
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("数据库异常!"))
		goazure.Error("insert issued_template Error", err)
		return
	}
	//删除终端模板配置的状态
	if _,err:=dba.BoilerOrm.QueryTable("issued_term_temp_status").Filter("Template__Uid",template.TemplateUpdate.Uid).Delete();err!=nil{
		goazure.Error("Delete IssuedCommunicationTemplate Error", err)
	}
}

//删除模板
func (ctl *TemplateController) TemplateDelete() {
	var template models.IssuedTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &template); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	template.IsDeleted = true
	if err := dba.BoilerOrm.Read(&template); err != nil {
		goazure.Error("Query IssuedTemplate Error", err)
	}
	template.IsDeleted = true
	num, err := dba.BoilerOrm.Update(&template)
	if err != nil {
		fmt.Printf("update Rows Num: %d, %s", num, err)
	}
	var channelTemplate []models.IssuedChannelConfigTemplate
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").Filter("Template__Uid", template.Uid).Filter("ChannelType",models.CHANNEL_TYPE_RANGE).All(&channelTemplate); err != nil {
		goazure.Error("Query IssuedChannelConfigTemplate Error", err)
	}
	for _, ct := range channelTemplate {
			if num, err := dba.BoilerOrm.QueryTable("issued_channel_config_range_template").Filter("ChannelConfig__Uid", ct.Uid).Delete(); err != nil {
				goazure.Error("Delete Old Ranges Error:", err, num)
			}
	}
	if _, err := dba.BoilerOrm.QueryTable("issued_channel_config_template").Filter("Template__Uid", template.Uid).Delete(); err != nil {
		goazure.Error("Delete IssuedChannelConfigTemplate Error", err)
	}
	if _,err := dba.BoilerOrm.QueryTable("issued_communication_template").Filter("Template__Uid",template.Uid).Delete();err!=nil{
		goazure.Error("Delete IssuedCommunicationTemplate Error", err)
	}
	//删除终端模板配置的状态
	if _,err:=dba.BoilerOrm.QueryTable("issued_term_temp_status").Filter("Template__Uid",template.Uid).Delete();err!=nil{
		goazure.Error("Delete IssuedCommunicationTemplate Error", err)
	}
}


