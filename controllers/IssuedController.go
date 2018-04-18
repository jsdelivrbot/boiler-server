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
type BoilerStatus struct {
	Uid string `json:"uid"`
	Value bool  `json:"value"`
}

type TermErrCode struct {
	StartDate    time.Time  `json:"startDate"`
	EndDate      time.Time  `json:"endDate"`
	Sn           string     `json:"sn"`
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
type BoilerIssuedMini struct {
	BoilerId string `json:"boiler_id"`
	Code  string    `json:"terminal_code"`
	Value int  `json:"value"`
}

type BoilerIssued struct {
	BoilerId string `json:"boiler_id"`
	TermId string    `json:"terminal_id"`
	Value   int   `json:"value"`
}
type VersionIssued struct {
	Sn         string
	Ver         int32
	UpdateTime  time.Time
}

//第一个模拟通道
func (ctl *IssuedController) IssuedAnalogOne(Uid string)([]byte) {
	var temp int32 =1
	var end int32 =12
	Byte:=make([]byte,0)
	var anaOne []models.IssuedAnalogueSwitch
	if num,err:=dba.BoilerOrm.QueryTable("issued_analogue_switch").
	RelatedSel("Function").RelatedSel("Byte").RelatedSel("Channel").
	Filter("Channel__Terminal__Uid",Uid).Filter("Channel__ChannelType",models.CHANNEL_TYPE_TEMPERATURE).
	OrderBy("Channel__ChannelNumber").All(&anaOne);err!=nil{
		goazure.Error("Query issued_analogue_switch Error",err)
	} else {
		//组模拟通道1
		L:=int(num)
		if L ==0 {
			for c := temp - 1; c < end; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=anaOne[i]
				if temp == value.Channel.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				} else {
					for c := 0; c < int(value.Channel.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					temp=value.Channel.ChannelNumber
				}
				if i==L-1 {
					if temp != end {
						for c := 0; c < int(end - temp); c++ {
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
	fmt.Println("模拟量一的Byte:",Byte)
	return Byte
}
//第二个模拟通道
func (ctl *IssuedController) IssuedAnalogTwo(Uid string)([]byte) {
	var temp int32 =1
	var end int32 =12
	Byte:=make([]byte,0)
	var anaTwo []models.IssuedAnalogueSwitch
	if num,err:=dba.BoilerOrm.QueryTable("issued_analogue_switch").
		RelatedSel("Function").RelatedSel("Byte").RelatedSel("Channel").
		Filter("Channel__Terminal__Uid",Uid).Filter("Channel__ChannelType",models.CHANNEL_TYPE_ANALOG).
		OrderBy("Channel__ChannelNumber").All(&anaTwo);err!=nil{
		goazure.Error("Query issued_analogue_switch Error",err)
	} else {
		//组模拟通道2
		L:=int(num)
		if L ==0 {
			for c := temp -1; c < end; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=anaTwo[i]
				if temp == value.Channel.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				} else {
					for c := 0; c < int(value.Channel.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					temp=value.Channel.ChannelNumber
				}
				if i==L-1 {
					if temp !=end {
						for c := 0; c < int(end - temp); c++ {
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
	return Byte
}
//组开关量默认位
func (ctl *IssuedController) IssuedSwitchBurn(Uid string)([]byte) {
	Byte:=make([]byte,0)
	var swi   []models.IssuedSwitchDefault
	if num,err:=dba.BoilerOrm.QueryTable("issued_switch_default").RelatedSel("Function").Filter("Terminal__Uid",Uid).OrderBy("ChannelNumber").All(&swi);err!=nil{
		goazure.Error("Query issued_switch_default",err)
	} else {
		if num == 0 {
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteTwo(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
			Byte = append(Byte, IntToByteTwo(0)...)
			Byte = append(Byte, IntToByteOne(0)...)
		} else {
			for _,c := range swi {
				Byte = append(Byte, IntToByteOne(int32(c.Function.Id))...)
				Byte = append(Byte, IntToByteTwo(int32(c.Modbus))...)
				Byte = append(Byte, IntToByteOne(int32(c.BitAddress))...)
			}
		}
	}
	return Byte
}
//组剩余开关量
func (ctl *IssuedController) IssuedSwitch(Uid string)([]byte) {
	Byte:=make([]byte,0)
	var temp int32 =3
	var end int32 =48
	var switchs []models.IssuedAnalogueSwitch
	if num,err:=dba.BoilerOrm.QueryTable("issued_analogue_switch").
	RelatedSel("Function").RelatedSel("Channel").
	Filter("Channel__ChannelType",models.CHANNEL_TYPE_SWITCH).
	Filter("Channel__Terminal__Uid",Uid).
	OrderBy("Channel__ChannelNumber").All(&switchs);err!=nil{
		goazure.Error("Query issued_analogue_switch Error",err)
	} else {
		L:= int(num)
		if L ==0 {
			for c := temp -1; c < end; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=switchs[i]
				if temp == value.Channel.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					Byte = append(Byte, IntToByteOne(int32(value.BitAddress))...)
				} else {
					for c := 0; c < int(value.Channel.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					Byte = append(Byte, IntToByteOne(int32(value.BitAddress))...)
					temp=value.Channel.ChannelNumber
				}
				if i==L-1 {
					if temp !=end {
						for c := 0; c < int(end - temp); c++ {
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
	fmt.Println("开关位的Byte:",Byte)
	return Byte
}
//组状态量
func (ctl *IssuedController) IssuedRange(Uid string)([]byte) {
	Byte:=make([]byte,0)
	var temp int32 =1
	var end int32 =12
	var anaThree []models.IssuedAnalogueSwitch
	if num,err:=dba.BoilerOrm.QueryTable("issued_analogue_switch").
		RelatedSel("Function").RelatedSel("Byte").RelatedSel("Channel").
		Filter("Channel__Terminal__Uid",Uid).Filter("Channel__ChannelType",models.CHANNEL_TYPE_RANGE).
		OrderBy("Channel__ChannelNumber").All(&anaThree);err!=nil{
		goazure.Error("Query issued_analogue_switch Error",err)
	}else {
		L:=int(num)
		if L ==0 {
			for c := temp - 1; c < end; c++ {
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteOne(0)...)
				Byte = append(Byte, IntToByteTwo(0)...)
			}
		} else {
			for i:=0;i<L;i++{
				value:=anaThree[i]
				if temp == value.Channel.ChannelNumber {
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
				} else {
					for c := 0; c < int(value.Channel.ChannelNumber - temp); c++ {
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteOne(0)...)
						Byte = append(Byte, IntToByteTwo(0)...)
					}
					Byte = append(Byte, IntToByteOne(int32(value.Function.Id))...)
					Byte = append(Byte, IntToByteOne(int32(value.Byte.Id))...)
					Byte = append(Byte, IntToByteTwo(int32(value.Modbus))...)
					temp=value.Channel.ChannelNumber
				}
				if i==L-1 {
					if temp !=end {
						for c := 0; c < int(end - temp); c++ {
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
	return Byte
}
//组装通信参数
func (ctl *IssuedController) IssuedCommunication(Uid string)([]byte) {
	Byte:=make([]byte,0)
	var communication models.IssuedCommunication
	if err:=dba.BoilerOrm.QueryTable("issued_communication").RelatedSel("BaudRate").
	RelatedSel("DataBit").RelatedSel("StopBit").RelatedSel("CheckBit").
	RelatedSel("CorrespondType").RelatedSel("SubAddress").RelatedSel("HeartBeat").Filter("Terminal__Uid",Uid).One(&communication);err!=nil{
		goazure.Error("Query issued_communication",err)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
		Byte = append(Byte,IntToByteOne(0)...)
	} else {
		Byte = append(Byte,IntToByteOne(int32(communication.BaudRate.Id))...)
		Byte = append(Byte,IntToByteOne(int32(communication.DataBit.Id))...)
		Byte = append(Byte,IntToByteOne(int32(communication.StopBit.Id))...)
		Byte = append(Byte,IntToByteOne(int32(communication.CheckBit.Id))...)
		Byte = append(Byte,IntToByteOne(int32(communication.CorrespondType.Id))...)
		Byte = append(Byte,IntToByteOne(int32(communication.SubAddress.Id))...)
		Byte = append(Byte,IntToByteOne(int32(communication.HeartBeat.Id))...)
	}
	return Byte
}
//查询终端版本号
func (ctl *IssuedController) IssuedVersion(Code string)(int32)  {
	var versionIssued VersionIssued
	var version int32
	versionSql:="select sn,ver from issued_version where sn=?"
	if err:=dba.BoilerOrm.Raw(versionSql,Code).QueryRow(&versionIssued);err!=nil{
		goazure.Error("Query issued_version Error",err)
		version =1
	} else {
		fmt.Println("versionIssued:",versionIssued)
		version = versionIssued.Ver + 1
		if version >=32769 {
			version = 1
		}
	}
	return version
}
//组装现有配置的报文
func (ctl *IssuedController) IssuedMessage(Uid string)([]byte) {
	Byte:=make([]byte,0)
	Byte = append(Byte,ctl.IssuedAnalogOne(Uid)...)
	fmt.Println("模拟量1长度:",len(Byte))
	Byte = append(Byte,ctl.IssuedAnalogTwo(Uid)...)
	fmt.Println("模拟量2长度:",len(Byte))
	Byte = append(Byte,ctl.IssuedSwitchBurn(Uid)...)
	fmt.Println("开关量点火位长度:",len(Byte))
	Byte = append(Byte,ctl.IssuedSwitch(Uid)...)
	fmt.Println("开关位长度:",len(Byte))
	Byte = append(Byte,ctl.IssuedRange(Uid)...)
	fmt.Println("状态量长度:",len(Byte))
	Byte = append(Byte,ctl.IssuedCommunication(Uid)...)
	fmt.Println("数据长度:",len(Byte))
	return Byte
}

//根据Code获取报文
func (ctl *IssuedController)ReqMessage(Code string)(string) {
	var info Info
	info.Sn=Code
	fmt.Println("infoSn",info.Sn)
	sql:="select curr_message from issued_message where sn=?"
	if err:=dba.BoilerOrm.Raw(sql,Code).QueryRow(&info);err!=nil{
		goazure.Error("Query issued_message Error",err)
		return ""
	}
	fmt.Println("下发的报文:",info.CurrMessage)
	return info.CurrMessage
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
	reqBuf:=ctl.ReqMessage(confIssued.Code)
	if reqBuf=="" {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("还未保存配置"))
		return
	}
	fmt.Println("要下发的buf:",reqBuf)
	buf:=SocketCtrl.SocketConfigSend(reqBuf)
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
			fmt.Println("返回的配置状态:",buf[15])
			newVer:=ByteToIntTwo(buf[13:15])
			fmt.Println("14到15字节：",buf[13:15])
			fmt.Println("返回的版本号:",newVer)
			termCode:=fmt.Sprintf("%s",buf[7:13])
			fmt.Println("返回的终端版本号:",termCode)
			if newVer==1 {
				sql:="insert into issued_version(sn,ver,create_time,update_time) values(?,?,now(),now())"
				if _,err:=dba.BoilerOrm.Raw(sql,termCode,newVer).Exec();err!=nil{
					goazure.Error("insert issued_version Error",err)
					fmt.Println("终端版本号插入失败")
				}
				fmt.Println("终端版本号插入成功")
			} else {
				sql:="update issued_version set ver=?,update_time=now() where sn=?"
				if _,err:=dba.BoilerOrm.Raw(sql,newVer,termCode).Exec();err!=nil{
					goazure.Error("update issued_version Error",err)
					fmt.Println("终端版本号修改失败")
				}
				fmt.Println("终端版本号修改成功")
			}
		} else {
			fmt.Println("返回的配置状态:",buf[15])
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("终端配置错误"))
		}
	} else {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("返回报文信息错误"))
	}
}

func (ctl *IssuedController) TerminalErrorList() {
	var termErr TermErrCode
	var errAlarm []models.IssuedPlcAlarm
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &termErr); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	fmt.Println("startTime:",termErr.StartDate)
	fmt.Println("endTime:",termErr.EndDate)
	if _,err:=dba.BoilerOrm.QueryTable("issued_plc_alarm").RelatedSel("Err").Filter("Sn",termErr.Sn).
	Filter("CreateTime__gte",termErr.StartDate).Filter("CreateTime__lte",termErr.EndDate).OrderBy("-CreateTime").All(&errAlarm);err!=nil{
		goazure.Error("Query Issued Plc Alarm Error",err)
	}
	ctl.Data["json"] = errAlarm
	ctl.ServeJSON()
}

//获取锅炉重启状态
func (ctl *IssuedController)IssuedBoilerStatus() {
	var bStatus BoilerStatus
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &bStatus); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	var param []orm.Params
	var status bool
	if num,err:=dba.BoilerOrm.QueryTable("issued_boiler_status").Filter("Boiler__Uid",bStatus.Uid).Values(&param,"Status");err!=nil || num ==0 {
		goazure.Warning("Read Boiler Burning Status Error!", err, num)
		status=false
	} else {
		status = param[0]["Status"].(bool)
	}
	ctl.Data["json"] = status
	ctl.ServeJSON()
}
//修改锅炉重启
func (ctl *IssuedController)IssuedBoilerUpdate() {
	var bStatus BoilerStatus
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &bStatus); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	sql:="insert into issued_boiler_status(boiler_id,create_time,update_time,status) values(?,now(),now(),?) on duplicate key update update_time=now(),status=?"
	if _,err:=dba.BoilerOrm.Raw(sql,bStatus.Uid,bStatus.Value,bStatus.Value).Exec();err!=nil{
		goazure.Error("Insert issued_boiler_status Error",err)
	}
}

//小程序锅炉重启
func (ctl *IssuedController) IssuedBoilerMini() {
	var bMini BoilerIssuedMini
	var combined models.BoilerTerminalCombined
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &bMini); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	if err:=dba.BoilerOrm.QueryTable("boiler_terminal_combined").RelatedSel("Terminal").Filter("Boiler__Uid",bMini.BoilerId).
	Filter("TerminalCode",bMini.Code).One(&combined);err!=nil{
		goazure.Error("Query Boiler Terminal Combined Error",err)
	}
	if combined.Terminal.IsOnline {
		buf:=SocketBoilerSend(bMini.Code,combined.TerminalSetId,bMini.Value)
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
			if buf[13]!=16 {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("终端配置错误"))
				return
			}
		} else {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("返回报文信息错误"))
		}
	} else {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端未在线!"))
		return
	}
}

//锅炉重启
func (ctl *IssuedController) IssuedBoiler() {
	var boilerIssued BoilerIssued
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &boilerIssued); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}
	bindSql:="select terminal_code,terminal_set_id from boiler_terminal_combined where boiler_id=? and terminal_id=?"
	var Code string
	var TermSetId int32
	if err:=dba.BoilerOrm.Raw(bindSql,boilerIssued.BoilerId,boilerIssued.TermId).QueryRow(&Code,&TermSetId);err!=nil{
		goazure.Error("Query terminal_code Error",err)
	}
	fmt.Println("Code:",Code)
	fmt.Println("TermId:",boilerIssued.TermId)
	fmt.Println("TermSetId:",TermSetId)
	fmt.Println("value:",boilerIssued.Value)
	//查询终端是否在线
	var Online bool
	OnlineSql:="select is_online from terminal where terminal_code = ?"
	if err:=dba.BoilerOrm.Raw(OnlineSql,Code).QueryRow(&Online);err!=nil {
		goazure.Error("Query terminal Error",err)
	}
	fmt.Println("Online:",Online)
	if Online {
		buf:=SocketBoilerSend(Code,TermSetId,boilerIssued.Value)
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
			if buf[13]!=16 {
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte("终端配置错误"))
				return
			}
		} else {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("返回报文信息错误"))
		}
	} else {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("终端未在线!"))
		return
	}
}

//终端重启
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
		if buf[13]!=16 {
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
