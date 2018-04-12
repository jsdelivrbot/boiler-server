package controllers

import (
	"github.com/AzureTech/goazure"
	"encoding/json"
	"fmt"
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/pborman/uuid"
	"github.com/AzureRelease/boiler-server/models"
)

type TemplateController struct {
	MainController
}
type Template struct {
	Uid     string           `json:"uid"`
	Name    string           `json:"name"`
	OrganizationUid string    `json:"organizationUid"`
	Channel []ChannelIssued  `json:"channel"`
	Param   Param            `json:"param"`
}
func (ctl *TemplateController) IssuedTemplate() {
	var template Template
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &template); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("Name:",template.Name)
	fmt.Println("OUid:",template.OrganizationUid)
	template.Uid=uuid.New()
	sql:="insert into issued_template(uid,name,create_time,update_time,organization_id) values(?,?,now(),now(),?)"
	if _,err:=dba.BoilerOrm.Raw(sql,template.Uid,template.Name,template.OrganizationUid).Exec();err!=nil{
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("数据库异常!"))
		goazure.Error("insert issued_template Error", err)
		return
	}
	commuSql:="insert into issued_communication_template(uid,template_id,baud_rate_id,data_bit_id,stop_bit_id,check_bit_id,correspond_type_id,sub_address_id,heart_beat_id) values(uuid(),?,?,?,?,?,?,?,?)"
	if _,err:=dba.BoilerOrm.Raw(commuSql,template.Uid,template.Param.BaudRateId,template.Param.DataBitId,template.Param.StopBitId,template.Param.CheckDigitId,template.Param.CommunInterfaceId,template.Param.SubAddressId,template.Param.HeartBitId).Exec();err!=nil{
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("数据库异常!"))
		goazure.Error("insert issued_template Error", err)
		return
	}
	for i, c := range template.Channel {
		goazure.Warn("[", i, "]ChannelIssuedUpdate: ", c)
		if c.ParameterId <=0 {
			continue
		} else {
			channelUid:=uuid.New()
			channelSql:="insert into issued_channel_config_template(uid,create_time,parameter_id,template_id,channel_type,channel_number,status,sequence_number,switch_status,func_id,byte_id,bit_address,modbus) "+
				"values(?,now(),?,?,?,?,?,?,?,?,?,?,?)"
			if _,err:=dba.BoilerOrm.Raw(channelSql,channelUid,c.ParameterId,template.Uid,c.ChannelType,c.ChannelNumber,c.Status,c.SequenceNumber,c.SwitchValue,c.FcodeId,c.TermByteId,c.BitAddress,c.Modbus).Exec();err!=nil{
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("数据库异常!"))
				goazure.Error("insert issued_channel_config_template Error", err)
				return
			}

			if c.ChannelType==models.CHANNEL_TYPE_RANGE {
				for i, r := range c.Ranges {
					if r.Max < r.Min {
						e := fmt.Sprintln("状态值范围设定有误！")
						goazure.Error(e)
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte(e))
						return
					}
					rangeUid:=uuid.New()
					rangeSql:="insert into issued_channel_config_range_template(uid,name,create_time,channel_config_id,min,max,value) " +
						"values(?,?,now(),?,?,?,?)"
						if _,err:=dba.BoilerOrm.Raw(rangeSql,rangeUid,r.Name,channelUid,r.Min,r.Max,i).Exec();err!=nil{
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
	qs:=dba.BoilerOrm.QueryTable("issued_template")
	qs=qs.RelatedSel("Organization")
	if usr.IsOrganizationUser(){
		qs.Filter("Organization__Uid",usr.Organization.Uid)
	}
	if num, err := qs.Filter("IsDeleted", false).OrderBy("-UpdateTime").All(&template);err!=nil{
		fmt.Printf("Returned Rows Num: %d, %s", num, err)
	}
	ctl.Data["json"] = template
	ctl.ServeJSON()
}

//获取模拟量通道一
func (ctl *TemplateController) TemplateAnalogOne() {
	var temp  models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：",temp)
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").RelatedSel("Template").Filter("Template",temp.Uid).Filter("ChannelType",models.CHANNEL_TYPE_TEMPERATURE).All(&cTemps);err!=nil{
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix:=make([]models.IssuedChannelConfigTemplate,12)
	fmt.Println("channel:",cTemps)
	for _,c := range cTemps {
		matrix[c.ChannelNumber-1]=c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取模拟量通道二
func (ctl *TemplateController) TemplateAnalogTwo() {
	var temp  models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：",temp)
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").RelatedSel("Template").Filter("Template",temp.Uid).Filter("ChannelType",models.CHANNEL_TYPE_ANALOG).All(&cTemps);err!=nil{
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix:=make([]models.IssuedChannelConfigTemplate,12)
	fmt.Println("channel:",cTemps)
	for _,c := range cTemps {
		matrix[c.ChannelNumber-1]=c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取开关量一
func (ctl *TemplateController) TemplateSwitchOne() {
	var temp  models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：",temp)
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Template").Filter("Template",temp.Uid).Filter("ChannelType",models.CHANNEL_TYPE_SWITCH).Filter("ChannelNumber__gte",1).Filter("ChannelNumber__lte",16).All(&cTemps);err!=nil{
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix:=make([]models.IssuedChannelConfigTemplate,16)
	fmt.Println("channel:",cTemps)
	for _,c := range cTemps {
		matrix[c.ChannelNumber-1]=c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取开关量二
func (ctl *TemplateController) TemplateSwitchTwo() {
	var temp  models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：",temp)
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Template").Filter("Template",temp.Uid).Filter("ChannelType",models.CHANNEL_TYPE_SWITCH).Filter("ChannelNumber__gte",17).Filter("ChannelNumber__lte",32).All(&cTemps);err!=nil{
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix:=make([]models.IssuedChannelConfigTemplate,16)
	fmt.Println("channel:",cTemps)
	for _,c := range cTemps {
		matrix[(c.ChannelNumber-1)%16]=c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

//获取开关量三
func (ctl *TemplateController) TemplateSwitchThree() {
	var temp  models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：",temp)
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Template").Filter("Template",temp.Uid).Filter("ChannelType",models.CHANNEL_TYPE_SWITCH).Filter("ChannelNumber__gte",33).Filter("ChannelNumber__lte",48).All(&cTemps);err!=nil{
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix:=make([]models.IssuedChannelConfigTemplate,16)
	fmt.Println("channel:",cTemps)
	for _,c := range cTemps {
		matrix[(c.ChannelNumber-1)%16]=c
	}
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}


//获取状态量
func (ctl *TemplateController) TemplateRange() {
	var temp  models.IssuedTemplate
	var cTemps []models.IssuedChannelConfigTemplate
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &temp); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	fmt.Println("cTemp：",temp)
	if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").RelatedSel("Parameter").RelatedSel("Func").RelatedSel("Byte").RelatedSel("Template").Filter("Template",temp.Uid).Filter("ChannelType",models.CHANNEL_TYPE_RANGE).All(&cTemps);err!=nil{
		goazure.Error("Query issued_channel_config_template Error:", err, num)
	}
	matrix:=make([]models.IssuedChannelConfigTemplate,12)
	fmt.Println("channel:",cTemps)
	for _,c := range cTemps {
		var ranges []*models.IssuedChannelConfigRangeTemplate
		if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_range_template").Filter("ChannelConfig__Uid",c.Uid).OrderBy("Value").All(&ranges);err!=nil{
			goazure.Error("Get Template ChannelConfig Range Error:", err, num)
		} else {
			c.Ranges=ranges
		}
		matrix[c.ChannelNumber-1]=c
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
	if err:=dba.BoilerOrm.QueryTable("issued_communication_template").RelatedSel("BaudRate").RelatedSel("CheckBit").RelatedSel("CorrespondType").RelatedSel("DataBit").RelatedSel("HeartBeat").RelatedSel("StopBit").RelatedSel("SubAddress").RelatedSel("Template").Filter("Template__Uid",temp.Uid).One(&cTemps);err!=nil{
		goazure.Error("Query issued_communication_template Error:", err)
	}
	ctl.Data["json"] = cTemps
	ctl.ServeJSON()
}


//修改模板信息
func (ctl *TemplateController) TemplateUpdate() {

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
	template.IsDeleted=true
	if err:=dba.BoilerOrm.Read(&template);err!=nil {
		goazure.Error("Query IssuedTemplate Error", err)
	}
	template.IsDeleted=true
	 num,err:=dba.BoilerOrm.Update(&template)
	if err!=nil{
		fmt.Printf("update Rows Num: %d, %s", num, err)
	}
	var channelTemplate []models.IssuedChannelConfigTemplate
	if _,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").Filter("Template__Uid",template.Uid).All(&channelTemplate);err!=nil{
		goazure.Error("Query IssuedChannelConfigTemplate Error", err)
	}
	for _,ct:= range channelTemplate {
		if ct.ChannelType==models.CHANNEL_TYPE_RANGE{
			if num,err:=dba.BoilerOrm.QueryTable("issued_channel_config_range_template").Filter("ChannelConfig__Uid",ct.Uid).Delete();err!=nil{
				goazure.Error("Delete Old Ranges Error:", err, num)
			}
		}
	}
	if _,err:=dba.BoilerOrm.QueryTable("issued_channel_config_template").Filter("Template__Uid",template.Uid).Delete();err!=nil{
		goazure.Error("Query IssuedChannelConfigTemplate Error", err)
	}
}

