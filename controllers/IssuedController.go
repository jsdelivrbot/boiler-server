package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
)

type IssuedController struct {
	goazure.Controller
}
func (ctl *IssuedController)BinUpload() {
	binPath := "C:\\ftpServer\\"
	var orgUid string
	if len(ctl.Input()["orgUid"][0]) >0 {
		orgUid = ctl.Input()["orgUid"][0]
	} else {
		goazure.Error("orgUid is null")
	}
	if file, header, er := ctl.GetFile("file"); er != nil {
		fmt.Println("解析错误")
	} else {
		if file == nil {
			fmt.Println("文件为空")
		}
		fileName := header.Filename
		basePath := "/home/apps/bin/"
		filePath := basePath + fileName
		if err := ctl.SaveToFile("file", filePath); err != nil {
			e := fmt.Sprintln("Save File Error:", err, fileName)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
		if res, err := dba.MyORM.Raw("insert into issued_bin_upload(name,create_time,update_time,organization_id,bin_path,status ) values(?,now(),now(),?,?,0) on duplicate key update update_time=now(),organization_id=organization_id", fileName,orgUid,binPath+fileName).Exec(); err != nil {
			fmt.Println(err)
			fmt.Println(res)
		} else {
		  ctl.Ctx.Output.Body([]byte("上传成功"))
		}
		goazure.Info("Save Done:", header.Filename, fileName)
	}
}

func (ctl *IssuedController)FunctionCodeList(){
	var functionCodes []models.IssuedFunctionCode
	if _,err:=dba.MyORM.QueryTable("issued_function_code").OrderBy("Id").All(&functionCodes);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=functionCodes
	ctl.ServeJSON()
}
func (ctl *IssuedController)ByteCodeList(){
	var bytes []models.IssuedByte
	if _,err:=dba.MyORM.QueryTable("issued_byte").OrderBy("Id").All(&bytes);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=bytes
	ctl.ServeJSON()
}
func (ctl *IssuedController)BaudRateList() {
	var baudRates []models.IssuedBaudRate
	if _,err:=dba.MyORM.QueryTable("issued_baud_rate").OrderBy("Id").All(&baudRates);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=baudRates
	ctl.ServeJSON()
}
func (ctl *IssuedController)CorrespondTypeList() {
	var correspondTypes []models.IssuedCorrespondType
	if _,err:=dba.MyORM.QueryTable("issued_correspond_type").OrderBy("Id").All(&correspondTypes);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=correspondTypes
	ctl.ServeJSON()
}
func (ctl *IssuedController)DateBitList() {
	var dateBits []models.IssuedDateBit
	if _,err:=dba.MyORM.QueryTable("issued_date_bit").OrderBy("Id").All(&dateBits);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=dateBits
	ctl.ServeJSON()
}
func (ctl *IssuedController)HeartbeatPacketList() {
	var Heartbeats []models.IssuedHeartbeatPacket
	if _,err:=dba.MyORM.QueryTable("issued_heartbeat_packet").OrderBy("Id").All(&Heartbeats);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=Heartbeats
	ctl.ServeJSON()
}
func (ctl *IssuedController)ParityBitList() {
	var parityBits []models.IssuedParityBit
	if _,err:=dba.MyORM.QueryTable("issued_parity_bit").OrderBy("Id").All(&parityBits);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=parityBits
	ctl.ServeJSON()
}
func (ctl *IssuedController)SlaveAddressList() {
	var slaveAddresses []models.IssuedSlaveAddress
	if _,err:=dba.MyORM.QueryTable("issued_slave_address").OrderBy("Id").All(&slaveAddresses);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=slaveAddresses
	ctl.ServeJSON()
}
func (ctl *IssuedController)StopBitList() {
	var stopBits []models.IssuedStopBit
	if _,err:=dba.MyORM.QueryTable("issued_stop_bit").OrderBy("Id").All(&stopBits);err!=nil{
		goazure.Error("Fetch ByteCodes List Error: ",err)
	}
	ctl.Data["json"]=stopBits
	ctl.ServeJSON()
}
