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
	num, err := qs.Filter("IsDeleted", false).OrderBy("-UpdateTime").All(&template)
	fmt.Printf("Returned Rows Num: %d, %s", num, err)
	ctl.Data["json"] = template
	ctl.ServeJSON()
}

//获取模板信息
func (ctl *TemplateController) TemplateInfo() {

}

//修改模板信息
func (ctl *TemplateController) TemplateUpdate() {

}
//删除模板
func (ctl *TemplateController) TemplateDelete() {

}

