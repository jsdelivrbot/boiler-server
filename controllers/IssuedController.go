package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
	"github.com/AzureRelease/boiler-server/util"
	"github.com/AzureTech/goazure/orm"
	"encoding/json"
	"github.com/AzureRelease/boiler-server/conf"
	"encoding/binary"
	"bytes"
	"time"
)

type IssuedController struct {
	MainController
}
var IssuedCtl *IssuedController = &IssuedController{}
type Code struct {
	Uid string `json:"uid"`
}
type ConfIssued struct {
	Uid string  `json:"uid"`
	Code string  `json:"code"`
}

func IntToByteOne(Int int32)([]byte){
	b_buf := bytes.NewBuffer([]byte{})
	err := binary.Write(b_buf, binary.BigEndian, Int)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	r_buf := []byte{b_buf.Bytes()[3]}
	return r_buf
}
func IntToByteTwo(Int int32)([]byte) {
	b_buf := bytes.NewBuffer([]byte{})
	err := binary.Write(b_buf, binary.BigEndian, Int)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	r_buf := []byte{b_buf.Bytes()[2],b_buf.Bytes()[3]}
	return r_buf
}

type AnalogueIssued struct {
	ChannelType int
	ChannelNumber int
	Func    int
	Byte    int
	Modbus  int
}

type SwitchIssued struct {
	ChannelType int
	ChannelNumber int
	Func int
	Modbus int
	BitAddress int
}

type CommunicationIssued struct {
	BaudRate   int
	DataBit    int
	StopBit    int
	CheckBit   int
	CorrespondType      int
	SubAddress      int
	HeartBeat    int
}

type VersionIssued struct {
	Sn         string
	Ver         int32
	UpdateTime  time.Time
}

//下发配置
func (ctl *IssuedController) IssuedConfig() {
	var confIssued ConfIssued
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &confIssued); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	Byte:=make([]byte,0)
	var temp =1
	var anaOne []AnalogueIssued
	var anaTwo []AnalogueIssued
	var anaThree []AnalogueIssued
	var swi   SwitchIssued
	var switchs []SwitchIssued
	var communication CommunicationIssued
	sql:="select r.channel_type,r.channel_number, ib.value as byte,ifc.value as func,ia.modbus " +
	"from runtime_parameter_channel_config r,issued_analogue ia,issued_byte ib, issued_function_code ifc where r.terminal_id=? and r.uid=ia.channel_id and r.channel_type=1 and ia.function_id=ifc.id and ia.byte_id=ib.id" +
	" ORDER BY r.channel_number; "
	if _,err:=dba.BoilerOrm.Raw(sql,confIssued.Uid).QueryRows(&anaOne);err!=nil {
		goazure.Error("Query issued_analogue Error",err)
	} else {
		//组模拟通道1
		L:=len(anaOne)
		if L ==0 {
			for c := 0; c < 12; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=anaOne[i]
				if temp == value.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				} else {
					for c := 0; c < (value.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					temp=value.ChannelNumber
				}
				if i==L-1 {
					if temp !=12 {
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
	}
	anasql:="select r.channel_type,r.channel_number, ib.value as byte,ifc.value as func,ia.modbus " +
		"from runtime_parameter_channel_config r,issued_analogue ia,issued_byte ib, issued_function_code ifc where r.terminal_id=? and r.uid=ia.channel_id and r.channel_type=2 and ia.function_id=ifc.id and ia.byte_id=ib.id" +
		" ORDER BY r.channel_number; "
	if _,err:=dba.BoilerOrm.Raw(anasql,confIssued.Uid).QueryRows(&anaTwo);err!=nil {
		goazure.Error("Query issued_analogue Error",err)
	} else {
		//组模拟通道2
		temp=1
		L:=len(anaTwo)
		if L ==0 {
			for c := 0; c < 12; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=anaTwo[i]
				if temp == value.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				} else {
					for c := 0; c < (value.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					temp=value.ChannelNumber
				}
				if i==L-1 {
					if temp !=12 {
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
	}
	//组开关量点火位
	switchBurnSql:="select ifc.value as func,isb.modbus,isb.bit_address from issued_switch_burn isb,issued_function_code ifc where terminal_id=? and isb.function_id=ifc.id"
	if err:=dba.BoilerOrm.Raw(switchBurnSql,confIssued.Uid).QueryRow(&swi);err!=nil {
		Byte = append(Byte, IntToByteOne(0)...)
		Byte = append(Byte, IntToByteTwo(0)...)
		Byte = append(Byte, IntToByteOne(0)...)
	}else {
		Byte = append(Byte, IntToByteOne(int32(swi.Func))...)
		Byte = append(Byte, IntToByteTwo(int32(swi.Modbus))...)
		Byte = append(Byte, IntToByteOne(int32(swi.BitAddress))...)
	}
	//组剩余开关量
	switchsql:="select r.channel_type,r.channel_number,ifc.value as func,iswitch.modbus,iswitch.bit_address "+
		"from runtime_parameter_channel_config r,issued_switch iswitch, issued_function_code ifc where r.terminal_id=? and r.uid=iswitch.channel_id and iswitch.function_id=ifc.id and r.channel_type=3 "+
			" ORDER BY r.channel_number;"
	if _,err:=dba.BoilerOrm.Raw(switchsql,confIssued.Uid).QueryRows(&switchs);err!=nil {
		goazure.Error("Query issued_switch Error",err)
	} else {
		temp=2
		L:=len(switchs)
		if L ==0 {
			for c := 0; c < 47; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=switchs[i]
				if temp == value.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					Byte = append(Byte, IntToByteOne(int32(value.BitAddress))...)
				} else {
					for c := 0; c < (value.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					Byte = append(Byte, IntToByteOne(int32(value.BitAddress))...)
					temp=value.ChannelNumber
				}
				if i==L-1 {
					if temp !=48 {
						for c := 0; c < (48 - temp); c++ {
							Byte = append(Byte, IntToByteOne(0)...)
							Byte= append(Byte, IntToByteTwo(0)...)
							Byte = append(Byte, IntToByteOne(0)...)
						}
					}
				}
				temp++
			}
		}
	}
	//组状态量
	statussql:="select r.channel_type,r.channel_number, ib.value as byte,ifc.value as func,ia.modbus " +
		"from runtime_parameter_channel_config r,issued_analogue ia,issued_byte ib, issued_function_code ifc where r.terminal_id=? and r.uid=ia.channel_id and r.channel_type=5 and ia.function_id=ifc.id and ia.byte_id=ib.id" +
		" ORDER BY r.channel_number; "
	if _,err:=dba.BoilerOrm.Raw(statussql,confIssued.Uid).QueryRows(&anaThree);err!=nil {
		goazure.Error("Query issued_analogue Error",err)
	} else {
		//组状态量
		temp=1
		L:=len(anaThree)
		if L ==0 {
			for c := 0; c < 12; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=anaThree[i]
				if temp == value.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				} else {
					for c := 0; c < (value.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Func))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					temp=value.ChannelNumber
				}
				if i==L-1 {
					if temp !=12 {
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
	}
	//组装通信参数
	communicationsql:="select ibr.value as baud_rate , idb.value as data_bit, isb.value as stop_bit, ipb.value as check_bit, ict.value as correspond_type,isa.value as sub_address,ihp.value as heart_beat "+
		"from issued_communication ic,issued_baud_rate as ibr,issued_data_bit as idb,issued_stop_bit as isb,issued_parity_bit as ipb,issued_correspond_type as ict,issued_slave_address as isa,issued_heartbeat_packet as ihp " +
			"where ic.baud_rate_id=ibr.id and ic.data_bit_id=idb.id and ic.stop_bit_id=isb.id and ic.check_bit_id=ipb.id and ic.correspond_type_id=ict.id and ic.sub_address_id=isa.id and ic.heart_beat_id=ihp.id and ic.terminal_id=?"
	if err:=dba.BoilerOrm.Raw(communicationsql,confIssued.Uid).QueryRow(&communication);err!=nil{
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
	} else {
		Byte = append(Byte,IntToByteOne(int32(communication.BaudRate))...)
		Byte = append(Byte,IntToByteOne(int32(communication.DataBit))...)
		Byte = append(Byte,IntToByteOne(int32(communication.StopBit))...)
		Byte = append(Byte,IntToByteOne(int32(communication.CheckBit))...)
		Byte = append(Byte,IntToByteOne(int32(communication.CorrespondType))...)
		Byte = append(Byte,IntToByteOne(int32(communication.SubAddress))...)
		Byte = append(Byte,IntToByteOne(int32(communication.HeartBeat))...)
	}
	//查询终端版本号
	var versionIssued VersionIssued
	var version int32
	versionSql:="select sn,ver from issued_version where sn=?"
	if err:=dba.BoilerOrm.Raw(versionSql,confIssued.Code).QueryRow(&versionIssued);err!=nil{
		goazure.Error("Query issued_version Error",err)
		version =1
	} else {
		fmt.Println("versionIssued:",versionIssued)
		version = versionIssued.Ver + 1
		if version >=32769 {
			version = 1
		}
	}
	fmt.Println("终端版本号为：",version)
	buf:=SocketCtrl.SocketConfigSend(Byte,version,confIssued.Code)
	if buf==nil{
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("发送报文失败"))
		return
	} else if bytes.Equal(conf.TermNoRegist,buf) {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端还未连接平台"))
		return
	} else if bytes.Equal(conf.TermTimeout,buf) {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端返回信息超时"))
		return
	}else if len(buf) >4 {
		if buf[15]==16 {
			newVer:=ByteToInt(buf[13:15])
			termCode:=fmt.Sprintf("%s",buf[7:13])
			if newVer==1 {
				sql:="insert into issued_version(sn,ver,create_time,update_time) values(?,?,now(),now())"
				if _,err:=dba.BoilerOrm.Raw(sql,termCode,newVer).Exec();err!=nil{
					goazure.Error("insert issued_version Error",err)
				}
			} else {
				sql:="update issued_version set ver=?,update_time=now() where sn=?"
				if _,err:=dba.BoilerOrm.Raw(sql,newVer,termCode).Exec();err!=nil{
					goazure.Error("update issued_version Error",err)
				}
			}
		} else {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("终端配置错误"))
		}
	} else {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("返回报文信息错误"))
	}
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
	terminal.Uid =code.Uid
	if err := dba.BoilerOrm.QueryTable("terminal").RelatedSel("organization").Filter("Uid", terminal.Uid).One(&terminal); err != nil {
		e := fmt.Sprintln("Read Terminal Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	buf:=SocketCtrl.SocketTerminalRestart(fmt.Sprintf("%d", terminal.TerminalCode))
	if buf==nil{
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("发送报文失败"))
		return
	} else if bytes.Equal(conf.TermNoRegist,buf) {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端还未连接平台"))
		return
	} else if bytes.Equal(conf.TermTimeout,buf) {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端返回信息超时"))
		return
	} else if len(buf)>4 {
		if buf[15]!=16 {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端配置错误"))
		return
		}
	} else {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("返回报文信息错误"))
	}
}
type AppBinInfo struct {
	Uid string `json:"uid"`
	Path string `json:"path"`
}
//升级配置
func (ctl *IssuedController)UpgradeConfiguration() {
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
			ctl.Ctx.Output.Body([]byte("等待连接升级"))
		} else if maps[0]["burnedStatus"] == "0" && maps[0]["burningStatus"] == "1" {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("正在升级中。。。"))
		} else if maps[0]["burnedStatus"] == "0" && maps[0]["burningStatus"] == "0" {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("等待连接升级"))
		}
	} else {
		fmt.Println("organized:", terminal.Organization.Name)
		_, err := orm.NewOrm().Raw("INSERT INTO appBinStatus(sn,company,createTime,updateTime,burnedStatus,burningStatus,burningPackages,burnNum)"+
			" values(?,?,now(),now(),0,0,0,0)", terminal.TerminalCode, terminal.Organization.Name).Exec()
		if err != nil {
			goazure.Error("Insert appBinSatus Error:", err)
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
		basePath := conf.BinPath
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
func (ctl *IssuedController)DataBitList() {
	var dateBits []models.IssuedDataBit
	if _,err:=dba.MyORM.QueryTable("issued_data_bit").OrderBy("Id").All(&dateBits);err!=nil{
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
