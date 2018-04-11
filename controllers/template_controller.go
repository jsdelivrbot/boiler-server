package controllers

import (
	"github.com/AzureTech/goazure"
	"encoding/json"
	"fmt"
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/pborman/uuid"
)

type TemplateController struct {
	goazure.Controller
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
	sql:="insert into issued_template(uid,name,organization_id) values(?,?,?)"
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
			//channelSql:="insert into issued"
		}
	}

}

