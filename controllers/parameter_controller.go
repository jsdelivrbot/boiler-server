package controllers

import (
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"

	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"

	"reflect"
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"encoding/csv"
	"strconv"
	"log"
	"time"
	"encoding/json"
	"sort"

	"github.com/pborman/uuid"
	"github.com/AzureRelease/boiler-server/models/logs"
	"math/rand"
)

type ParameterController struct {
	MainController

	ReloadLimit 	int64
	Parameters 		[]*models.RuntimeParameter
}

var ParamCtrl *ParameterController = &ParameterController{}

func init() {
	ParamCtrl.MainController = *MainCtrl

	go ParamCtrl.RefreshParameters()
}

//新增数据类型
type Channel struct {
	Channel []ChannelIssued `json:"channel"`
	Param   Param           `json:"param"`
}

//新增param
type Param struct {
	TerminalCode      string `json:"terminal_code"`   //终端编码
	BaudRateId        int    `json:"baudRate"`        //波特率
	DataBitId         int    `json:"dataBit"`         //数据位
	StopBitId         int    `json:"stopBit"`         //停止位
	CheckDigitId      int    `json:"checkDigit"`      //校检位
	CommunInterfaceId int    `json:"communInterface"` //通信接口地址
	SubAddressId      int    `json:"slaveAddress"`    //从机地址
	HeartBitId        int    `json:"heartbeat"`       //心跳包频率
}

//新增的通道下发
type ChannelIssued struct {
	TerminalCode string `json:"terminal_code"`
	ParameterId  int    `json:"parameter_id"`

	ChannelType    int32   `json:"channel_type"`
	ChannelNumber  int32   `json:"channel_number"`
	FcodeId        int   `json:"fcodeName"`  //功能码
	BitAddress     int   `json:"bitAddress"` //位地址
	TermByteId     int   `json:"termByte"`   //高低字节
	Modbus         int   `json:"modbus"`     //modbus
	Status         int32 `json:"status"`
	SequenceNumber int32 `json:"sequence_number"`

	SwitchValue int32 `json:"switch_status"`

	Scale float32 `json:"scale"`

	Ranges []bRange `json:"ranges"`

	IsDeleted bool `json:"is_deleted"`
}

type bRuntimeParameter struct {
	Id      int64   `json:"id"`
	CategoryId int64 `json:"category_id"`
	OrganizationId string `json:"organization_id"`
	Fix    int32  `json:"fix"`
	Length  int32   `json:"length"`
	Name   string  `json:"name"`
	ParamId  int32   `json:"param_id"`
	Scale   float32  `json:scale`
	Unit   string  `json:"unit"`
	Remark  string  `json:"remark"`
}

//传递前台通信参数
func (ctl *ParameterController) IssuedCommunication() {
	var issuedCommunication models.IssuedCommunication
	var param Param
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &param); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	if err := dba.BoilerOrm.QueryTable("issued_communication").RelatedSel("BaudRate").RelatedSel("DataBit").RelatedSel("StopBit").RelatedSel("CheckBit").RelatedSel("CorrespondType").
		RelatedSel("SubAddress").RelatedSel("HeartBeat").Filter("Terminal__TerminalCode", param.TerminalCode).One(&issuedCommunication); err != nil {
		goazure.Error("Query issued_communication Error", err)
	}
	ctl.Data["json"] = issuedCommunication
	ctl.ServeJSON()
}

func (ctl *ParameterController) ChannelIssuedUpdate() {
	var chanIssu Channel
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &chanIssu); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	var ter models.Terminal
	code, _ := strconv.ParseInt(chanIssu.Param.TerminalCode, 10, 32)
	if err := dba.BoilerOrm.QueryTable("Terminal").Filter("TerminalCode", code).Filter("IsDeleted", false).One(&ter); err != nil {
		goazure.Error("Get Channel's Terminal Error:", err)
	}
	var aCnf []models.RuntimeParameterChannelConfig
	if _, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		Filter("Terminal__Uid", ter.Uid).Filter("IsDefault", false).
		All(&aCnf); err != nil {
		goazure.Error("Get ChannelConfig To Delete Error:", err)
	}
	for _, a := range aCnf {
		TempCtrl.TemplateChannelConfigDelete(&a, ter.Uid)
	}
	if num, err := dba.BoilerOrm.QueryTable("issued_switch_default").Filter("Terminal__Uid", ter.Uid).Delete(); err != nil {
		goazure.Error("Delete issued switch Error:", err, num)
	}
	var cnf models.RuntimeParameterChannelConfig
	for _, c := range chanIssu.Channel {
		for _, p := range ParamCtrl.Parameters {
			if p.Id == int64(c.ParameterId) {
				cnf.Parameter = p
			}
		}
			cnf.Terminal = &ter
			cnf.ChannelType = c.ChannelType
			cnf.ChannelNumber = c.ChannelNumber
			cnf.IsDefault = false
			//TODO: Default Signed & Threshold
			cnf.Signed = true
			cnf.NegativeThreshold = 32768
			cnf.Status = c.Status
			if cnf.Status == models.CHANNEL_STATUS_SHOW {
				cnf.SequenceNumber = c.SequenceNumber
			} else {
				cnf.SequenceNumber = -1
			}
			cnf.Name = cnf.Parameter.Name
			cnf.Length = cnf.Parameter.Length
			cnf.Uid = uuid.New()
			if cnf.ChannelType == models.CHANNEL_TYPE_SWITCH {
				if c.SwitchValue == 0 {
					c.SwitchValue = 1
				}
				cnf.SwitchStatus = c.SwitchValue
			}
			if cnf.ChannelType == models.CHANNEL_TYPE_SWITCH && (cnf.ChannelNumber == 1 || cnf.ChannelNumber == 2){
				switchburnsql := "insert into issued_switch_default(uid,terminal_id,create_time,channel_type,channel_number,function_id,modbus,bit_address) values(uuid(),?,now(),?,?,?,?,?)"
				if _, err := dba.BoilerOrm.Raw(switchburnsql, ter.Uid,c.ChannelType,c.ChannelNumber,c.FcodeId, c.Modbus, c.BitAddress).Exec(); err != nil {
					goazure.Error("Insert issued_switch_default Error", err)
				}
			} else {
				if _, err := dba.BoilerOrm.Insert(&cnf); err != nil {
					goazure.Error("insert runtime parameter channel config Error", err)
				}
				analoguesql := "insert into issued_analogue_switch(channel_id,create_time,function_id,byte_id,modbus,bit_address) values(?,now(),?,?,?,?)"
				if _, err := dba.BoilerOrm.Raw(analoguesql, cnf.Uid, c.FcodeId, c.TermByteId, c.Modbus,c.BitAddress).Exec(); err != nil {
					goazure.Error("Insert issued_analogue Error", err)
				}
			}

			if cnf.ChannelType == models.CHANNEL_TYPE_RANGE {
				var aRanges []*models.RuntimeParameterChannelConfigRange
				sort.Sort(ByRMin(c.Ranges))
				for i, r := range c.Ranges {
					var rg models.RuntimeParameterChannelConfigRange
					if r.Max < r.Min {
						e := fmt.Sprintln("状态值范围设定有误！")
						goazure.Error(e)
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte(e))
						return
					}
					for _, ar := range aRanges {
						if ar.Max >= r.Min {
							e := fmt.Sprintln("状态值范围设定有误，请不要设定重复的范围区间！")
							goazure.Error(e)
							ctl.Ctx.Output.SetStatus(400)
							ctl.Ctx.Output.Body([]byte(e))
							return
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
	//插入表通信参数
	sql := "replace into issued_communication(terminal_id,baud_rate_id,data_bit_id,stop_bit_id,check_bit_id,correspond_type_id,sub_address_id,heart_beat_id) values(?,?,?,?,?,?,?,?)"
	if _, er := dba.BoilerOrm.Raw(sql, ter.Uid, chanIssu.Param.BaudRateId, chanIssu.Param.DataBitId, chanIssu.Param.StopBitId, chanIssu.Param.CheckDigitId, chanIssu.Param.CommunInterfaceId, chanIssu.Param.SubAddressId, chanIssu.Param.HeartBitId).Exec(); er != nil {
		goazure.Error("Insert issued_communication Error", er)
	}
	fmt.Println("插入通信参数结束")
	var param []orm.Params
	var ver int32
	verSql := "select ver from issued_message where sn=?"
	if num, err := dba.BoilerOrm.Raw(verSql, ter.TerminalCode).Values(&param); err != nil || num == 0 {
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
	Byte := IssuedCtl.IssuedMessage(ter.Uid)
	Code := strconv.FormatInt(ter.TerminalCode, 10)
	message := SocketCtrl.c0(Byte, Code, ver)
	messageSql := "insert into issued_message(sn,ver,create_time,update_time,curr_message) values(?,?,now(),now(),?) on duplicate key update ver=?,update_time=now(),curr_message=?"
	if _, err := dba.BoilerOrm.Raw(messageSql, ter.TerminalCode, ver, string(message), ver, string(message)).Exec(); err != nil {
		goazure.Error("Insert issued_message Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("组包失败"))
	}
	//删除终端模板配置的状态
	if _,err:=dba.BoilerOrm.QueryTable("issued_term_temp_status").Filter("Sn",ter.TerminalCode).Delete();err!=nil{
		goazure.Error("Delete IssuedCommunicationTemplate Error", err)
	}

}

func (ctl *ParameterController) RefreshParameters() {
	ParamCtrl.WaitGroup.Add(1)
	var params []*models.RuntimeParameter
	qs := dba.BoilerOrm.QueryTable("runtime_parameter")
	if num, err := qs.RelatedSel("Category").
		Filter("IsDeleted", false).OrderBy("Id").
		All(&params); err != nil || num == 0 {
		goazure.Error("Get RuntimeParameterList Error:", num, err)
	}

	for i, v := range params {
		if num, err := dba.BoilerOrm.LoadRelated(v, "BoilerMediums"); err != nil && num == 0 {
			goazure.Error("[", i, "]", v, num, err)
		}

		for _, b := range v.BoilerMediums {
			b.Name = strings.TrimSuffix(b.Name, "锅炉")
		}
	}

	ParamCtrl.Parameters = params

	ParamCtrl.WaitGroup.Done()
}

func (ctl *ParameterController) RuntimeParameterList() {
	goazure.Warning("Ready to Get RuntimeParameterList")
	ParamCtrl.WaitGroup.Wait()
	ctl.Data["json"] = ParamCtrl.Parameters
	ctl.ServeJSON()
}

func (ctl *ParameterController) RuntimeParameterIssuedList() {
	goazure.Warning("Ready to Get RuntimeParameterList")
	//ParamCtrl.bWaitGroup.Wait()
	usr:= ctl.GetCurrentUser()
	var issuedParams []*models.IssuedParameterOrganization
	qs := dba.BoilerOrm.QueryTable("issued_parameter_organization")
	qs = qs.RelatedSel("Parameter").RelatedSel("Organization").RelatedSel("Parameter__Category")
	if usr.IsOrganizationUser() {
		qs = qs.Filter("Organization__Uid",usr.Organization.Uid)
		//qs = qs.Filter("Scope", models.RUNTIME_ALARM_SCOPE_ENTERPRISE)
	}
	if num,err:=qs.Filter("IsDeleted",false).All(&issuedParams);err!=nil || num == 0 {
		goazure.Error("Get RuntimeParameterList Error:", num, err)
	}
	fmt.Println("qs:",fmt.Sprintf("%s",qs))
	for i, v := range issuedParams {
		if num, err := dba.BoilerOrm.LoadRelated(v.Parameter, "BoilerMediums"); err != nil && num == 0 {
			goazure.Error("[", i, "]", v, num, err)
		}

		for _, b := range v.Parameter.BoilerMediums {
			b.Name = strings.TrimSuffix(b.Name, "锅炉")
		}
	}
	ctl.Data["json"] = issuedParams
	ctl.ServeJSON()
}

func (ctl *ParameterController) ChannelDataReload(t time.Time) {
	var lgIn	logs.BoilerRuntimeLog
	nonce := rand.Intn(254)
	lgIn.Name = "ChannelDataReload()"
	lgIn.CreatedDate = t
	lgIn.Status = logs.BOILER_RUNTIME_LOG_STATUS_INIT
	LogCtl.AddReloadLog(&lgIn)

	data := ctl.DataListNeedReload(nonce)

	var lgRd	logs.BoilerRuntimeLog
	lgRd.Name = "DataListNeedReload()"
	lgRd.TableName = "boiler_m163"
	lgRd.Query = "SELECT"
	lgRd.CreatedDate = time.Now()
	lgRd.Duration = float64(lgRd.CreatedDate.Sub(t)) / float64(time.Second)
	lgRd.DurationTotal = lgRd.Duration
	lgRd.Status = logs.BOILER_RUNTIME_LOG_STATUS_DONE
	LogCtl.AddReloadLog(&lgRd)

	runtimeReload := func(rt orm.Params, bCnfs []*models.RuntimeParameterChannelConfig, b *models.Boiler, lastLog *logs.BoilerRuntimeLog) {
		for _, cnf := range bCnfs {
			var rtm models.BoilerRuntime
			var value int64

			chName := models.ChannelName(cnf.ChannelType, cnf.ChannelNumber)
			//goazure.Info("Channel_Name", chName, rt[chName], reflect.ValueOf(rt[chName]).Kind())
			switch reflect.ValueOf(rt[chName]).Kind() {
			case reflect.String:
				value, _ = strconv.ParseInt(rt[chName].(string), 10, 64)
				if value == 0 {
					v, _ := strconv.ParseFloat(rt[chName].(string), 10)
					value = int64(v)
				}
			case reflect.Int32:
				value = int64(rt[chName].(int32))
			case reflect.Int64:
				value = rt[chName].(int64)
			case reflect.Float32:
				value = int64(rt[chName].(float32))
			case reflect.Float64:
				value = int64(rt[chName].(float64))
			}

			if value > 65535 {
				continue
			}

			switch cnf.Parameter.Category.Id {
			case models.RUNTIME_PARAMETER_CATEGORY_ANALOG:
				if cnf.Signed == true && value > int64(cnf.NegativeThreshold) {
					value = int64(cnf.NegativeThreshold) - value
				}

				if cnf.Scale > 0 {
					value = int64(float32(value) * cnf.Scale)
				}
			case models.RUNTIME_PARAMETER_CATEGORY_SWITCH:
				idx := (cnf.ChannelNumber-1)%16 + 1
				var d int64 = 1
				for i := int32(1); i < idx; i++ {
					d *= 2
				}
				v := value / d
				v = v % 2

				value = v
			case models.RUNTIME_PARAMETER_CATEGORY_CALCULATE:
				fallthrough
			case models.RUNTIME_PARAMETER_CATEGORY_RANGE:
				var v int64 = -1
				var rmk string = ""
				for _, rg := range cnf.Ranges {
					if value >= rg.Min && value <= rg.Max {
						v = rg.Value
						rmk = rg.Name

						goazure.Warn("Channel Range Matched!", rg, value, v)
						break
					}
				}

				if v >= 0 {
					cnf.Remark = rmk
				}
			}

			rtm.Parameter = cnf.Parameter
			rtm.Boiler = b
			rtm.Value = value
			rtm.Remark = cnf.Remark

			rtm.Status = models.RUNTIME_STATUS_NEW

			if tm, err := time.ParseInLocation("2006-01-02 15:04:05", rt["TS"].(string), time.Local); err != nil {
				goazure.Error("Parse Time Error:", err)
			} else {
				//goazure.Info("Parse Time:", t, "||", tm)
				fmt.Println("m163 ts:",t)
				fmt.Println("BoilerRuntime表的CreateDate:", tm)
				rtm.CreatedDate = tm
			}

			if  err := DataCtl.AddData(&rtm, true); err != nil {
				goazure.Error("Reload Runtime Error:", err)
				//isSuccess = false
			} else {
				var lgd logs.BoilerRuntimeLog
				lgd.Name = "runtimeReload()"
				lgd.Runtime = &rtm
				lgd.TableName = "boiler_m163 -> boiler_runtime"
				lgd.Query = "INSERT"
				lgd.CreatedDate = time.Now()
				lgd.Duration = float64(lgd.CreatedDate.Sub(lastLog.CreatedDate)) / float64(time.Second)
				lgd.DurationTotal = lastLog.DurationTotal + lgd.Duration
				lgd.Status = logs.BOILER_RUNTIME_LOG_STATUS_DONE
				LogCtl.AddReloadLog(&lgd)

				go RtmCtl.RuntimeDataReload(&rtm, lgd.DurationTotal)
			}
		}
	}

	for _, d := range data {
		var tm time.Time
		code := d["Boiler_term_id"].(string)
		set := d["Boiler_boiler_id"].(string)
		t := d["TS"].(string)
		//ver := d["Boiler_data_fmt_ver"].(string)
		//sn := d["Boiler_sn"].(string)
		tm, _ = time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)

		rawDisUid :=
			"UPDATE	`boiler_m163` " +
			"SET	`need_reload` = 0 " +
			"WHERE	`uid` = " + "'" + d["uid"].(string) + "';"

		//var boiler	models.Boiler
		var combined	models.BoilerTerminalCombined

		co, _ := strconv.ParseInt(code, 10, 32)
		st, _ := strconv.ParseInt(set, 10, 32)

		if  err := dba.BoilerOrm.QueryTable("boiler_terminal_combined").
			RelatedSel("Boiler").RelatedSel("Terminal").
			Filter("TerminalCode", co).Filter("TerminalSetId", st).OrderBy("TerminalSetId").
			One(&combined);
			err != nil {
			goazure.Error("Get BoilerInfo Error:", err, co, st)

			go dba.BoilerOrm.Raw(rawDisUid).Exec()

			continue
		}

		startTime := time.Now()

		var lgr logs.BoilerRuntimeLog
		lgr.Name = "runtimeReload()"
		lgr.TableName = "boiler_m163 -> boiler_runtime"
		lgr.Query = d["uid"].(string)
		lgr.CreatedDate = startTime
		lgr.Duration = float64(startTime.Sub(tm)) / float64(time.Second)
		lgr.DurationTotal = lgr.Duration
		lgr.Status = logs.BOILER_RUNTIME_LOG_STATUS_READY
		LogCtl.AddReloadLog(&lgr)

		if 	combined.Boiler != nil &&
			len(combined.Boiler.Uid) > 0 {

			bCnfs := ctl.ChannelConfigList(code)

			go runtimeReload(d, bCnfs, combined.Boiler, &lgr)
		}

		go dba.BoilerOrm.Raw(rawDisUid).Exec()
	}

	/*
	rawDis :=
		"UPDATE	`boiler_m163` " +
		"SET	`need_reload` = 0 " +
		"WHERE	`need_reload` = " + strconv.FormatInt(int64(nonce), 10) + ";"
	go dba.BoilerOrm.Raw(rawDis).Exec()
	*/
}

func (ctl *ParameterController) DataListNeedReload(nonce int) []orm.Params {
	var data 	[]orm.Params

	/*
	rawReady :=
		"UPDATE	`boiler_m163` " +
		"SET	`need_reload` = " + strconv.FormatInt(int64(nonce), 10) + " " +
		"WHERE	`need_reload` = 1 " +
		"LIMIT	" + strconv.FormatInt(limit, 10) + ";"
	*/

	raw :=
		"SELECT	`rtm`.* " +
		//"FROM	`boiler_m163` AS `rtm`, `boiler_terminal_combined` AS `boiler` " +
		//"WHERE	`boiler`.`terminal_code` = CAST(`rtm`.`Boiler_term_id` AS SIGNED) " +
		//"AND 	`boiler`.`terminal_set_id` = CAST(`rtm`.`Boiler_boiler_id` AS SIGNED) " +
		//"WHERE	`boiler`.`terminal_code` = `rtm`.`Boiler_term_id` " +
		//"  AND 	`boiler`.`terminal_set_id` = `rtm`.`Boiler_boiler_id` " +
		//"  AND	`rtm`.`need_reload` = TRUE "
		"FROM	`boiler_m163` AS `rtm`"  +
		//"WHERE	`need_reload` = " + strconv.FormatInt(int64(nonce), 10) + ";"
		"WHERE	`need_reload` = 1;"
		//"WHERE	`TS` > '2018-03-27 00:00:00';"
		//"ORDER BY `TS` DESC; "

	var lg logs.BoilerRuntimeLog
	lg.Name = "DataListNeedReload()"
	lg.TableName = "boiler_163m"
	lg.Query = raw
	lg.CreatedDate = time.Now()
	lg.Status = logs.BOILER_RUNTIME_LOG_STATUS_READY
	LogCtl.AddReloadLog(&lg)

	/*
	if res, err := dba.BoilerOrm.Raw(rawReady).Exec(); err != nil {
		goazure.Error("Ready DataListNeedReload Error", err, res)
		return data
	}
	*/

	if  num, err := dba.BoilerOrm.Raw(raw).Values(&data); err != nil {
		goazure.Error("Get DataListNeedReload Error", err, num)
	}

	return data
}

func (ctl *ParameterController) ChannelConfigList(code interface{}) []*models.RuntimeParameterChannelConfig {
	//goazure.Warning("ChannelConfigList:", code)
	var dfConfs 	[]*models.RuntimeParameterChannelConfig
	var bConfs 		[]*models.RuntimeParameterChannelConfig
	var co			int64

	switch reflect.ValueOf(code).Kind() {
	case reflect.String:
		co, _ = strconv.ParseInt(code.(string), 10, 64)
	case reflect.Int:
		co = int64(code.(int))
	case reflect.Int8:
		co = int64(code.(int8))
	case reflect.Int16:
		co = int64(code.(int16))
	case reflect.Int32:
		co = int64(code.(int32))
	case reflect.Int64:
		co = code.(int64)
	}
	//查询出默认的通道
	if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		RelatedSel("Parameter").
		Filter("IsDefault", true).Filter("IsDeleted", false).All(&dfConfs); err != nil || num == 0 {
		goazure.Error("Default Channel Config is Missing!", err, num)
	} else {
		goazure.Info("Get Boiler Channel Config:", num, "\n", bConfs)
	}
	//查询出配置的通道
	if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		RelatedSel("Parameter").
		RelatedSel("Terminal").
		Filter("Terminal__TerminalCode", co).Filter("IsDeleted", false).All(&bConfs); err != nil || num == 0 {
		goazure.Error("Get Boiler Channel Config Error:", err, num)
	}
	//状态量查询

	for _, c := range bConfs {
		if c.ChannelType == models.CHANNEL_TYPE_RANGE {
			var bRanges []*models.RuntimeParameterChannelConfigRange
			if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config_range").
				Filter("ChannelConfig__Uid", c.Uid).Filter("IsDeleted", false).
				OrderBy("Value").
				All(&bRanges); err != nil {
				goazure.Error("Get Boiler ChannelConfig Range Error:", err, num)
			} else {
				//goazure.Info("Get Boiler ChannelConfig Range:", num, "\n", bRanges)
				c.Ranges = bRanges
			}
		}
	}

	hasIdMatched := func(list []int64, id int64) bool {
		for _, item := range list {
			if item == id {
				return true
			}
		}

		return false
	}

	var matchIds []int64
	//循环默认通道
	for _, dc := range dfConfs {
		var isMatched bool = false
		//循环配置的通道
		for _, c := range bConfs {
			//goazure.Warning("ChannelType:", c.ChannelType, dc.ChannelType)
			//goazure.Warning("ChannelNumber:", c.ChannelNumber, dc.ChannelNumber)
			//goazure.Warning("Parameter.Id:", c.Parameter.Id, dc.Parameter.Id)
			//判断是否配置了默认通道或者默认通道的位置是否被配置了，是的话，将id加入matchIds里
			if (c.ChannelType == dc.ChannelType &&
				c.ChannelNumber == dc.ChannelNumber) ||
				c.Parameter.Id == dc.Parameter.Id {

				isMatched = true

				matchIds = append(matchIds, c.Parameter.Id)
				break
			}
		}
		//通道类型没有配置默认通道，将默认通道加进去。已经配置了或者默认通道的端口被其他的配置了
		if !isMatched && !hasIdMatched(matchIds, dc.Parameter.Id) {
			bConfs = append(bConfs, dc)

			matchIds = append(matchIds, dc.Parameter.Id)
		}

		isMatched = false
	}

	/*
	for _, c := range bConfs {
		var isMatched bool = false
		for i, dc := range dfConfs {
			goazure.Warning("ChannelType:", c.ChannelType, dc.ChannelType)
			goazure.Warning("ChannelNumber:", c.ChannelNumber, dc.ChannelNumber)
			goazure.Warning("Parameter.Id:", c.Parameter.Id, dc.Parameter.Id)
			if (c.ChannelType 	== dc.ChannelType &&
				c.ChannelNumber == dc.ChannelNumber) ||
				c.Parameter.Id == dc.Parameter.Id {

				dfConfs[i] = c

				isMatched = true
				matchIds = append(matchIds, c.Parameter.Id)
				break
			}
		}

		if isMatched == false {
			dfConfs = append(dfConfs, c)
		}

		isMatched = false
	}
	*/

	//for i, c := range dfConfs {
	//	if hasIdMatched(matchIds, c.Parameter.Id) && c.IsDefault {
	//		dfConfs = append
	//	}
	//}

	var confs []*models.RuntimeParameterChannelConfig
	var showConfs, hideConfs, normalConfs, defaultConfs []*models.RuntimeParameterChannelConfig
	//循环配置的通道
	for _, c := range bConfs {
		//配置了默认通道
		if c.IsDefault {
			defaultConfs = append(defaultConfs, c)
		} else {
			switch c.Status {
			//选择位置展示
			case models.CHANNEL_STATUS_SHOW:
				showConfs = append(showConfs, c)
				//选择位置隐藏
			case models.CHANNEL_STATUS_HIDE:
				hideConfs = append(hideConfs, c)
				//选择位置默认
			case models.CHANNEL_STATUS_DEFAULT:
				fallthrough
			default:
				normalConfs = append(normalConfs, c)
			}
		}
	}

	sort.Sort(CnfBySeq(showConfs))
	sort.Sort(CnfBySeq(normalConfs))
	sort.Sort(CnfBySeq(defaultConfs))
	sort.Sort(CnfBySeq(hideConfs))
	confs = append(confs, showConfs...)
	confs = append(confs, normalConfs...)
	confs = append(confs, defaultConfs...)
	confs = append(confs, hideConfs...)

	//goazure.Error("Conf Sorted ",confs)
	//for _, c := range confs {
	//	goazure.Warning(c.Parameter.Id, "| Def:", c.IsDefault, "| Status:", c.Status, ".", c.SequenceNumber, "Channel:", c.ChannelType, ".", c.ChannelNumber)
	//}

	//time.Sleep(time.Minute)

	//ctl.Data["json"] = dfConfs
	//ctl.ServeJSON()

	return confs
}

type CnfBySeq []*models.RuntimeParameterChannelConfig

func (a CnfBySeq) Len() int      { return len(a) }
func (a CnfBySeq) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CnfBySeq) Less(i, j int) bool {
	if a[i].Status == a[j].Status && a[i].Status == models.CHANNEL_STATUS_SHOW && a[i].SequenceNumber != a[j].SequenceNumber {
		return a[i].SequenceNumber < a[j].SequenceNumber
	} else {
		return a[i].Parameter.Id < a[j].Parameter.Id
	}
}

func (ctl *ParameterController) ParamQueueWithBoiler(boiler *models.Boiler) []int {
	var paramQueue,
	rangeQueue,
	switchQueue []int
	var bConfs []*models.RuntimeParameterChannelConfig
	for _, cb := range boiler.TerminalsCombined {
		bConfs = append(bConfs, ctl.ChannelConfigList(strconv.FormatInt(cb.TerminalCode, 10))...)
	}

	ctmNum := 0
	for _, c := range bConfs {
		if !c.IsDefault {
			ctmNum++
		}

		if c.Status == models.CHANNEL_STATUS_HIDE {
			continue
		}

		if c.IsDefault && ctmNum > 0 {
			continue
		}

		if c.Parameter.Category.Id == models.RUNTIME_PARAMETER_CATEGORY_SWITCH {
			switchQueue = append(switchQueue, int(c.Parameter.Id))
			continue
		}

		if c.Parameter.Category.Id == models.RUNTIME_PARAMETER_CATEGORY_RANGE {
			rangeQueue = append(rangeQueue, int(c.Parameter.Id))
			continue
		}

		paramQueue = append(paramQueue, int(c.Parameter.Id))

		/*
		switch c.Status {
		case models.CHANNEL_STATUS_SHOW:
			b.RuntimeQueue = append([]int{int(c.Parameter.Id)}, b.RuntimeQueue...)
		case models.CHANNEL_STATUS_DEFAULT:
			b.RuntimeQueue = append(b.RuntimeQueue, int(c.Parameter.Id))
		case models.CHANNEL_STATUS_HIDE:
		default:

		}
		*/
	}

	paramQueue = append(paramQueue, rangeQueue...)
	paramQueue = append(paramQueue, switchQueue...)

	if boiler.Medium.Id == 1 && ctmNum <= 0 {
		paramQueue = append(paramQueue, 1201)
		paramQueue = append(paramQueue, 1202)
	}

	return paramQueue
}

func (ctl *ParameterController) BoilerHasChannelCustom() bool {
	var tid string
	if ctl.Input()["terminal"] == nil || len(ctl.Input()["terminal"]) == 0 {
		e := fmt.Sprintln("there is no boiler!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return false
	}

	tid = ctl.Input()["terminal"][0]

	hasCustom := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		Filter("Terminal__Uid", tid).
		Filter("IsDefault", false).Filter("IsDeleted", false).Exist()

	goazure.Warning("Boiler Has Custom Channel Config:", tid, hasCustom)

	type bJson struct {
		HasCustom bool
	}

	var js bJson
	js.HasCustom = hasCustom

	ctl.Data["json"] = js
	ctl.ServeJSON()

	return hasCustom
}

type ChannelConfig struct {
	TerminalCode string `json:"terminal_code"`
	ParameterId  int    `json:"parameter_id"`

	ChannelType   int `json:"channel_type"`
	ChannelNumber int `json:"channel_number"`

	Status         int32 `json:"status"`
	SequenceNumber int32 `json:"sequence_number"`

	SwitchValue int32 `json:"switch_status"`

	Scale float32 `json:"scale"`

	Ranges []bRange `json:"ranges"`

	IsDeleted bool `json:"is_deleted"`
}

type bRange struct {
	Name  string `json:"name"`
	Min   int64  `json:"min"`
	Max   int64  `json:"max"`
	Value int64  `json:"value"`
}

type ByRMin []bRange

func (a ByRMin) Len() int           { return len(a) }
func (a ByRMin) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRMin) Less(i, j int) bool { return a[i].Min < a[j].Min }

func (ctl *ParameterController) ChannelConfigMatrix() {
	c := ChannelConfig{}
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &c); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	bConfs := ctl.ChannelConfigList(c.TerminalCode)
	fmt.Println("bConfs:",bConfs)
	matrix := make([][]models.RuntimeParameterChannelConfigIssued, 16)
	for i := range matrix {
		matrix[i] = make([]models.RuntimeParameterChannelConfigIssued, 6)
	}
	cnfIss := models.RuntimeParameterChannelConfigIssued{}
	for _, cnf := range bConfs {
		analogSwitchIssued := models.IssuedAnalogueSwitch{}
		cnfIss.RuntimeParameterChannelConfig = cnf
		if cnf.IsDefault == false {
			if cnf.ChannelType == models.CHANNEL_TYPE_SWITCH  {
				if err := dba.BoilerOrm.QueryTable("issued_analogue_switch").RelatedSel("Function").Filter("channel_id", &cnf.Uid).One(&analogSwitchIssued); err != nil {
					goazure.Error("Query issued_switch Error:", err)
				}
			} else  {
				if err := dba.BoilerOrm.QueryTable("issued_analogue_switch").RelatedSel("Function").RelatedSel("Byte").Filter("channel_id", &cnf.Uid).One(&analogSwitchIssued); err != nil {
					goazure.Error("Query issued_analogue Error:", err)
				}
			}
		}
		cnfIss.AnalogueSwitch = analogSwitchIssued
		switch cnf.ChannelType {
		case models.CHANNEL_TYPE_TEMPERATURE:
			fallthrough
			//模拟量
		case models.CHANNEL_TYPE_ANALOG:
			matrix[cnf.ChannelNumber-1][cnf.ChannelType-1] = cnfIss
			//开关量
		case models.CHANNEL_TYPE_SWITCH:
			num := (cnf.ChannelNumber - 1) % 16
			ty := (cnf.ChannelNumber-1)/16 + models.CHANNEL_TYPE_SWITCH - 1
			matrix[num][ty] = cnfIss
			//计算量
		case models.CHANNEL_TYPE_CALCULATE:
			matrix[cnf.ChannelNumber-1][cnf.ChannelType+1] = cnfIss
			//状态量
		case models.CHANNEL_TYPE_RANGE:
			matrix[cnf.ChannelNumber-1][cnf.ChannelType] = cnfIss
		default:

		}
	}
	issuedSwitchDefault := []models.IssuedSwitchDefault{}
	if _,err := dba.BoilerOrm.QueryTable("issued_switch_default").RelatedSel("Function").Filter("Terminal__TerminalCode", c.TerminalCode).All(&issuedSwitchDefault); err != nil {
		goazure.Error("Query issued_switch_burn Error", err)
	} else {
		for _,c := range issuedSwitchDefault {
			if c.ChannelNumber == 1{
				matrix[0][2].AnalogueSwitch.Function = c.Function
				matrix[0][2].AnalogueSwitch.Modbus = c.Modbus
				matrix[0][2].AnalogueSwitch.BitAddress = c.BitAddress
			} else if c.ChannelNumber == 2 {
				matrix[1][2].AnalogueSwitch.Function = c.Function
				matrix[1][2].AnalogueSwitch.Modbus = c.Modbus
				matrix[1][2].AnalogueSwitch.BitAddress = c.BitAddress
			}
		}
	}
	goazure.Info("ChannelConfig Matrix:", matrix)
	ctl.Data["json"] = matrix
	ctl.ServeJSON()
}

func (ctl *ParameterController) ChannelConfigUpdate() {
	var aCnfs []*ChannelConfig
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &aCnfs); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Config Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	for i, c := range aCnfs {
		goazure.Warn("[", i, "]ChannelConfigUpdate: ", c)
		var cnf models.RuntimeParameterChannelConfig
		var ter models.Terminal
		code, _ := strconv.ParseInt(c.TerminalCode, 10, 32)
		if err := dba.BoilerOrm.QueryTable("Terminal").Filter("TerminalCode", code).Filter("IsDeleted", false).One(&ter); err != nil {
			goazure.Error("Get Channel's Terminal Error:", err)
		}
		cnf.Terminal = &ter
		cnf.ChannelType = int32(c.ChannelType)
		cnf.ChannelNumber = int32(c.ChannelNumber)
		cnf.IsDefault = false

		//TODO: Default Signed & Threshold
		cnf.Signed = true
		cnf.NegativeThreshold = 32768

		if c.ParameterId <= 0 {
			var aCnf models.RuntimeParameterChannelConfig
			if err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
				Filter("Terminal__Uid", ter.Uid).Filter("ChannelType", c.ChannelType).Filter("ChannelNumber", c.ChannelNumber).Filter("IsDefault", false).
				One(&aCnf); err != nil {
				goazure.Error("Get ChannelConfig To Delete Error:", err)
				continue
			}

			ctl.ChannelConfigDelete(&aCnf)
			continue
		} else {
			for _, p := range ParamCtrl.Parameters {
				if p.Id == int64(c.ParameterId) {
					cnf.Parameter = p
				}
			}
		}

		cnf.Status = c.Status
		if cnf.Status == models.CHANNEL_STATUS_SHOW {
			cnf.SequenceNumber = c.SequenceNumber
		} else {
			cnf.SequenceNumber = -1
		}

		cnf.Name = cnf.Parameter.Name
		cnf.Length = cnf.Parameter.Length

		if cnf.ChannelType == models.CHANNEL_TYPE_SWITCH {
			if c.SwitchValue == 0 {
				c.SwitchValue = 1
			}
			cnf.SwitchStatus = c.SwitchValue
		}

		if err := DataCtl.AddData(&cnf, true, "Terminal", "ChannelType", "ChannelNumber", "IsDefault"); err != nil {
			e := fmt.Sprintln("Channel Config Update Error:", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))

			continue
		}

		if cnf.Parameter.Category.Id == models.RUNTIME_PARAMETER_CATEGORY_RANGE {
			if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config_range").
				Filter("ChannelConfig__Uid", cnf.Uid).Delete(); err != nil {
				goazure.Error("Delete Old Ranges Error:", err, num)
			}

			var aRanges []*models.RuntimeParameterChannelConfigRange
			sort.Sort(ByRMin(c.Ranges))
			for i, r := range c.Ranges {
				//goazure.Info(i, "| range:", r)
				var rg models.RuntimeParameterChannelConfigRange
				if r.Max < r.Min {
					e := fmt.Sprintln("状态值范围设定有误！")
					goazure.Error(e)
					ctl.Ctx.Output.SetStatus(400)
					ctl.Ctx.Output.Body([]byte(e))
					return
				}

				for _, ar := range aRanges {
					if ar.Max >= r.Min {
						e := fmt.Sprintln("状态值范围设定有误，请不要设定重复的范围区间！")
						goazure.Error(e)
						ctl.Ctx.Output.SetStatus(400)
						ctl.Ctx.Output.Body([]byte(e))
						return
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
}

func (ctl *ParameterController) ChannelConfigDelete(cnf *models.RuntimeParameterChannelConfig) error {
	if cnf.ChannelType == models.CHANNEL_TYPE_RANGE {
		if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config_range").
			Filter("ChannelConfig__Uid", cnf.Uid).
			Delete(); err != nil {
			goazure.Error("Delete Ranges Error:", err, num)
		}
	}

	if num, err := dba.BoilerOrm.Delete(cnf); err != nil {
		goazure.Error("Delete Channel Config Error:", err, num)
		return err
	}

	return nil
}

func (ctl *ParameterController) RuntimeParameterUpdate() {
	//var p models.RuntimeParameter
	var param models.RuntimeParameter
	var bParam bRuntimeParameter
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &bParam); err != nil {
		e := fmt.Sprintln("Unmarshal Parameter JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	//goazure.Warn("Parameter:", p)

	if err := dba.BoilerOrm.QueryTable("runtime_parameter").Filter("Id", bParam.Id).One(&param); err != nil {
		e := fmt.Sprintln("Read Parameter Error", err)
		goazure.Warn(e)
		param.Id= bParam.Id
		param.ParamId = bParam.ParamId
		param.Length = bParam.Length
		param.Fix = bParam.Fix

		param.Category = runtimeParameterCategory(bParam.CategoryId)

		param.Medium = runtimeParameterMedium(0)
		param.AddBoilerMedium(0)
		param.CreatedBy = ctl.GetCurrentUser()
	}

	param.Name = bParam.Name
	param.Scale = bParam.Scale
	param.Unit = bParam.Unit
	param.Remark = bParam.Remark

	param.UpdatedBy = ctl.GetCurrentUser()
	param.IsDeleted = false

	if err := DataCtl.AddData(&param, true); err != nil {
		e := fmt.Sprintln("Add/Update Parameter Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
	}
	sql:="insert into issued_parameter_organization(parameter_id,create_time,update_time,organization_id) values(?,now(),now(),?) on duplicate key update update_time=now(),is_deleted=false, organization_id=?"
	if _,err:=dba.BoilerOrm.Raw(sql,bParam.Id,bParam.OrganizationId,bParam.OrganizationId).Exec();err!=nil{
		goazure.Error("Insert issued_parameter_organization Error",err)
	}

	go ParamCtrl.RefreshParameters()
}

func (ctl *ParameterController) RuntimeParameterDelete() {
	var p models.RuntimeParameter

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &p); err != nil {
		e := fmt.Sprintln("Unmarshal Parameter JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var cnfs []*models.RuntimeParameterChannelConfig
	if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").Filter("Parameter__Id", p.Id).All(&cnfs); err != nil {
		goazure.Warn("Get ChannelConfig For This Parameter Error:", err, num)
	}

	for _, c := range cnfs {
		ctl.ChannelConfigDelete(c)
	}

	goazure.Warn("Delete Parameter:", p)

	if err := DataCtl.DeleteData(&p); err != nil {
		e := fmt.Sprintln("Delete Parameter Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))

		return
	}
	sql:="update issued_parameter_organization set is_deleted=true where parameter_id=?"
	if _,err:=dba.BoilerOrm.Raw(sql,p.Id).Exec();err!=nil{
		goazure.Error("Update issued_parameter_organization Error",err)
	}
	go ParamCtrl.RefreshParameters()
}


func (ctl *ParameterController) RuntimeParameter(pid int) *models.RuntimeParameter {
	var param models.RuntimeParameter

	if err := dba.BoilerOrm.QueryTable("runtime_parameter").
		RelatedSel("Category").
		Filter("Id", pid).Filter("IsDeleted", false).
		One(&param); err != nil {
		goazure.Error("Read Parameter Error: ", err)
		return nil
	}

	return &param
}

func runtimeParameterCategory(cateId int64) *models.RuntimeParameterCategory {
	var category models.RuntimeParameterCategory
	if err := dba.BoilerOrm.QueryTable("runtime_parameter_category").
		Filter("Id", cateId).Filter("IsDeleted", false).
		One(&category); err != nil {
		e := fmt.Sprintln("Read ParameterCategory Error", err, cateId)
		goazure.Error(e)

		return nil
	}

	return &category
}

func runtimeParameterMedium(medId int64) *models.RuntimeParameterMedium {
	med := models.RuntimeParameterMedium{}
	med.Id = medId
	err := DataCtl.ReadData(&med)
	if err != nil {
		fmt.Println("Read Error: ", err)
	}

	return &med
}

func generateDefaultChannelConfig() error {

	const fieldRowNo int = 0
	const paramIdCol int = 0

	var ac models.RuntimeParameterChannelConfig
	aType := reflect.TypeOf(ac)

	err := filepath.Walk(rtmDefaultsPath, func(aPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(aPath, "channel_config") {
			return err
		}

		goazure.Info("FilePath", aPath)

		file, errFile := os.Open(aPath)
		if errFile != nil {
			fmt.Println("Read File Error:", errFile)
		}
		reader := csv.NewReader(file)

		records, errRead := reader.ReadAll()
		if errRead != nil {
			log.Fatal(errRead)
		}
		goazure.Info("Records: ", records)

		var fieldNames []string
		for i, row := range records {
			if i == fieldRowNo {
				fieldNames = row
			}
			if i > fieldRowNo {
				d := reflect.New(aType)
				da := d.Elem()
				in := d.Interface().(models.DataInterface)

				for j, field := range row {
					if fieldNames[j] == "" || field == "" {
						continue
					}

					var value interface{}
					var name string

					switch j {
					case paramIdCol:
						pid, _ := strconv.ParseInt(field, 10, 64)
						value = ParamCtrl.RuntimeParameter(int(pid))
						name = value.(*models.RuntimeParameter).Name
						if strings.Index(name, "室内温度") >= 0 {
							name = "室内温度"
						}
						name = "默认(" + name + ")"
						da.FieldByName("Name").Set(reflect.ValueOf(name))
					default:
						switch da.FieldByName(fieldNames[j]).Kind() {
						case reflect.Int32:
							vi64, _ := strconv.ParseInt(field, 10, 32)
							value = int32(vi64)
						case reflect.Int64:
							vi64, _ := strconv.ParseInt(field, 10, 32)
							value = vi64
						case reflect.Float32:
							vf32, _ := strconv.ParseFloat(field, 32)
							value = float32(vf32)
						case reflect.Float64:
							vf64, _ := strconv.ParseFloat(field, 64)
							value = vf64
						case reflect.Bool:
							vfb, e := strconv.ParseBool(field)
							if e != nil {
								goazure.Error("ParseBool Error:", field, e)
							}
							value = vfb
						default:
							value = field
						}
					}

					goazure.Info("Field(", fieldNames[j], ":", da.FieldByName(fieldNames[j]).Kind(), "): ", field, value)
					da.FieldByName(fieldNames[j]).Set(reflect.ValueOf(value))
				}

				//idx := i - fieldRowNo - 1
				DataCtl.AddData(in, true,
					"Parameter",
					"IsDefault",
					"ChannelType",
					"ChannelNumber")
			}
		}

		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (ctl *ParameterController) InitParameterChannelConfig() {
	// generateDefaultChannelConfig()

	interval := time.Second * 10
	ticker := time.NewTicker(interval)
	tick := func() {
		for t := range ticker.C {
			fmt.Println("处理ChannelDataReload:",t)
			ParamCtrl.ChannelDataReload(t)
		}
	}

	go tick()
	go ParamCtrl.ChannelDataReload(time.Now())
}
