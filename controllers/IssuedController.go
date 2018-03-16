package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
)

type IssuedController struct {
	goazure.Controller
}
func (ctl *IssuedController)FunctionCodeList(){
	var functionCodes []models.TermFunctionCode
	if _,err:=dba.MyORM.QueryTable("term_function_code").OrderBy("Id").All(&functionCodes);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=functionCodes
	ctl.ServeJSON()
}
func (ctl *IssuedController)ByteCodeList(){
	var bytes []models.TermByte
	if _,err:=dba.MyORM.QueryTable("term_byte").OrderBy("Id").All(&bytes);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=bytes
	ctl.ServeJSON()
}
func (ctl *IssuedController)BaudRateList() {
	var baudRates []models.BaudRate
	if _,err:=dba.MyORM.QueryTable("baud_rate").OrderBy("Id").All(&baudRates);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=baudRates
	ctl.ServeJSON()
}
func (ctl *IssuedController)CorrespondTypeList() {
	var correspondTypes []models.CorrespondType
	if _,err:=dba.MyORM.QueryTable("correspond_type").OrderBy("Id").All(&correspondTypes);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=correspondTypes
	ctl.ServeJSON()
}
func (ctl *IssuedController)DateBitList() {
	var dateBits []models.DateBit
	if _,err:=dba.MyORM.QueryTable("date_bit").OrderBy("Id").All(&dateBits);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=dateBits
	ctl.ServeJSON()
}
func (ctl *IssuedController)HeartbeatPacketList() {
	var Heartbeats []models.HeartbeatPacket
	if _,err:=dba.MyORM.QueryTable("heartbeat_packet").OrderBy("Id").All(&Heartbeats);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=Heartbeats
	ctl.ServeJSON()
}
func (ctl *IssuedController)ParityBitList() {
	var parityBits []models.ParityBit
	if _,err:=dba.MyORM.QueryTable("parity_bit").OrderBy("Id").All(&parityBits);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=parityBits
	ctl.ServeJSON()
}
func (ctl *IssuedController)SlaveAddressList() {
	var slaveAddresses []models.SlaveAddress
	if _,err:=dba.MyORM.QueryTable("slave_address").OrderBy("Id").All(&slaveAddresses);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=slaveAddresses
	ctl.ServeJSON()
}
func (ctl *IssuedController)StopBitList() {
	var stopBits []models.StopBit
	if _,err:=dba.MyORM.QueryTable("stop_bit").OrderBy("Id").All(&stopBits);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=stopBits
	ctl.ServeJSON()
}
