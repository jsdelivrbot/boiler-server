package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
	"github.com/AzureRelease/boiler-server/util"
	"github.com/AzureTech/goazure/orm"
	"encoding/json"
)

type IssuedController struct {
	MainController
}
type Code struct {
	Uid string `json:"uid"`
}
//重启
func (ctl *IssuedController) TerminalRestart() {
	var code Code
	var terminal models.Terminal
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &code); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	//fmt.Println(code.Uid)
	if err := dba.BoilerOrm.QueryTable("terminal").RelatedSel("organization").Filter("Uid", terminal.Uid).One(&terminal); err != nil {
		e := fmt.Sprintln("Read Terminal Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	//SocketCtrl.SocketTerminalRestart(fmt.Sprintf("%d", terminal.TerminalCode))
}
type AppBinInfo struct {
	Uid string `json:"uid"`
	Path string `json:"path"`
}
//升级配置
func (ctl *IssuedController)UpgradeConfiguration() {
	fmt.Println("aaaaaaaaaaaa")
	var appBinInfo AppBinInfo
	var terminal models.Terminal
	if err:= json.Unmarshal(ctl.Ctx.Input.RequestBody,&appBinInfo);err!=nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	fmt.Println("appBinInfo:",appBinInfo)
	terminal.Uid = appBinInfo.Uid
	if err := dba.BoilerOrm.QueryTable("terminal").RelatedSel("organization").Filter("Uid", terminal.Uid).One(&terminal); err != nil {
		e := fmt.Sprintln("Read Terminal Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	_,er:=dba.BoilerOrm.Raw("REPLACE INTO appBinInfo(sn,path) values(?,?)",terminal.TerminalCode,appBinInfo.Path).Exec()
	if er == nil {
		goazure.Info("Insert appBinInfo success")
	} else {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("文件添加失败"))
		goazure.Error("Insert appBinInfo fail:",er)
		return
	}
	var maps []orm.Params
	num, err := orm.NewOrm().Raw("SELECT * FROM appBinStatus where sn =?", terminal.TerminalCode).Values(&maps)
	if err == nil && num > 0 {
		fmt.Println("map:", maps[0]["burnedStatus"], maps[0]["burningStatus"])
		if maps[0]["burnedStatus"] == "1" && maps[0]["burningStatus"] == "2" {
			_, err := orm.NewOrm().Raw("UPDATE appBinStatus set burnedStatus =0,updateTime = now() where sn = ?", terminal.TerminalCode).Exec()
			if err != nil {
				goazure.Error("Update appBinStatus Error:", err)
			}
		} else if maps[0]["burnedStatus"] == "0" && maps[0]["burningStatus"] == "2" {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("终端未在线,等待终端连接升级"))
		} else if maps[0]["burnedStatus"] == "0" && maps[0]["burningStatus"] == "1" {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("正在升级中。。。"))
		} else if maps[0]["burnedStatus"] == "0" && maps[0]["burningStatus"] == "0" {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("等待升级"))
		}
	} else {
		fmt.Println("organized:", terminal.Organization.Name)
		_, err := orm.NewOrm().Raw("INSERT INTO appBinStatus(sn,company,createTime,updateTime,burnedStatus,burningStatus,burningPackages,burnNum)"+
			" values(?,?,now(),now(),0,0,0,0)", terminal.TerminalCode, terminal.Organization.Name).Exec()
		if err != nil {
			goazure.Error("Insert appBinSatus Error:", err)
		} else {

		}
	}

}

//列出bin文件
func (ctl *IssuedController)BinFileList() {
	usr := ctl.GetCurrentUser()
	var binUploads []*models.IssuedBinUpload
	qs := dba.BoilerOrm.QueryTable("issued_bin_upload")
	if usr.IsOrganizationUser() {
		qs = qs.Filter("Organization__Uid", usr.Organization.Uid)
	}
	if num,err := qs.OrderBy("UpdateTime").All(&binUploads); err !=nil || num ==0 {
		goazure.Error("Read binFile List Error:",err,num)
	} else {
		goazure.Info("Returned binFile RowNum:",num)
	}
	fmt.Println("binUploads:",binUploads)
	ctl.Data["json"] = binUploads
	ctl.ServeJSON()
}


//文件上传
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
		basePath := "E:\\log\\"
		filePath := basePath + fileName
		if err := ctl.SaveToFile("file", filePath); err != nil {
			e := fmt.Sprintln("Save File Error:", err, fileName)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
		if _, err := dba.MyORM.Raw("insert into issued_bin_upload(name,create_time,update_time,organization_id,bin_path,status ) values(?,now(),now(),?,?,0) on duplicate key update update_time=now(),organization_id=organization_id", fileName,orgUid,binPath+fileName).Exec(); err != nil {
			goazure.Error("Insert issuedBinUpload error")
			ctl.Ctx.Output.Body([]byte("插入数据失败"))
		} else {
		  ctl.Ctx.Output.Body([]byte("上传成功"))
		}
		if b:=util.FtpClient(fileName); b{
			goazure.Info(" ftp client success")
		} else {
			goazure.Error("Insert issuedBinUpload fail")
			ctl.Ctx.Output.Body([]byte("传输文件失败"))
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
