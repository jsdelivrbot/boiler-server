package controllers

import (
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/conf"

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
	"errors"
	"github.com/AzureRelease/boiler-server/models/logs"
	"math/rand"
)

type ParameterController struct {
	MainController

	ReloadLimit				int64

	Parameters 				[]*models.RuntimeParameter
}

var ParamCtrl *ParameterController = &ParameterController{}

func init() {
	ParamCtrl.MainController = *MainCtrl

	go ParamCtrl.RefreshParameters()
}

func (ctl *ParameterController) RefreshParameters() {
	ParamCtrl.bWaitGroup.Add(1)

	var params []*models.RuntimeParameter
	qs := dba.BoilerOrm.QueryTable("runtime_parameter")
	if 	num, err := qs.RelatedSel("Category").
		Filter("IsDeleted", false).OrderBy("Id").
		All(&params); err != nil || num == 0 {
		goazure.Error("Get RuntimeParameterList Error:", num, err)
	}

	for i, v := range params {
		if num, err := dba.BoilerOrm.LoadRelated(v, "BoilerMediums"); err != nil && num == 0 {
			goazure.Error("[",i,"]", v, num, err)
		}

		for _, b := range v.BoilerMediums {
			b.Name = strings.TrimSuffix(b.Name, "锅炉")
		}
	}

	ParamCtrl.Parameters = params

	ParamCtrl.bWaitGroup.Done()
}

func (ctl *ParameterController) GetRuntimeParameter(typeId int32) (*models.RuntimeParameter, error) {
	for _, param := range ParamCtrl.Parameters {
		if param.ParamId == typeId {
			return param, nil
		}
	}

	err := errors.New("can not find any parameter with this Id")

	return nil, err
}

func (ctl *ParameterController) RuntimeParameterList() {
	goazure.Warning("Ready to Get RuntimeParameterList")
	ParamCtrl.bWaitGroup.Wait()

	ctl.Data["json"] = ParamCtrl.Parameters
	ctl.ServeJSON()
}

func (ctl *ParameterController) ChannelDataReload(t time.Time) {
	var lgIn	logs.BoilerRuntimeLog
	nonce := rand.Intn(254)
	lgIn.Name = "ChannelDataReload()"
	lgIn.CreatedDate = t
	lgIn.Status = logs.BOILER_RUNTIME_LOG_STATUS_INIT
	go DataCtl.AddData(&lgIn, false)

	data := ctl.DataListNeedReload(nonce)

	var lgRd	logs.BoilerRuntimeLog
	lgRd.Name = "DataListNeedReload()"
	lgRd.TableName = "boiler_m163"
	lgRd.Query = "SELECT"
	lgRd.CreatedDate = time.Now()
	lgRd.Duration = float64(lgRd.CreatedDate.Sub(t)) / float64(time.Second)
	lgRd.DurationTotal = lgRd.Duration
	lgRd.Status = logs.BOILER_RUNTIME_LOG_STATUS_DONE
	go DataCtl.AddData(&lgRd, false)

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
		go DataCtl.AddData(&lgr, false)

		runtimeReload := func(cnf *models.RuntimeParameterChannelConfig) {
			var rtm models.BoilerRuntime
			var value int64
			chName := models.ChannelName(cnf.ChannelType, cnf.ChannelNumber)
			goazure.Info("Channel_Name", chName, d[chName], reflect.ValueOf(d[chName]).Kind())
			switch reflect.ValueOf(d[chName]).Kind() {
			case reflect.String:
				value, _ = strconv.ParseInt(d[chName].(string), 10, 64)
				if value == 0 {
					v, _ := strconv.ParseFloat(d[chName].(string), 10)
					value = int64(v)
				}
			case reflect.Int32:
				value = int64(d[chName].(int32))
			case reflect.Int64:
				value = d[chName].(int64)
			case reflect.Float32:
				value = int64(d[chName].(float32))
			case reflect.Float64:
				value = int64(d[chName].(float64))
			}

			if value > 65535 {
				go dba.BoilerOrm.Raw(rawDisUid).Exec()
				return
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
				idx := (cnf.ChannelNumber - 1) % 16 + 1
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
				var v 	int64 = -1
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
			rtm.Boiler = combined.Boiler
			rtm.Value = value
			rtm.Remark = cnf.Remark

			rtm.Status = models.RUNTIME_STATUS_NEW

			if tm, err := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local); err != nil {
				goazure.Error("Parse Time Error:", err)
			} else {
				//goazure.Info("Parse Time:", t, "||", tm)
				rtm.CreatedDate = tm
			}

			if err := DataCtl.AddData(&rtm, true); err != nil {
				goazure.Error("Reload Runtime Error:", err)
				//isSuccess = false
				return
			} else {
				var lgd logs.BoilerRuntimeLog
				lgd.Name = "runtimeReload()"
				lgd.Runtime = &rtm
				lgd.TableName = "boiler_m163 -> boiler_runtime"
				lgd.Query = "INSERT"
				lgd.CreatedDate = time.Now()
				lgd.Duration = float64(lgd.CreatedDate.Sub(startTime)) / float64(time.Second)
				lgd.DurationTotal = lgr.DurationTotal + lgd.Duration
				lgd.Status = logs.BOILER_RUNTIME_LOG_STATUS_DONE
				go DataCtl.AddData(&lgd, false)

				go RtmCtl.RuntimeDataReload(&rtm, lgd.DurationTotal)
			}

			go dba.BoilerOrm.Raw(rawDisUid).Exec()
		}

		if 	combined.Boiler != nil &&
			len(combined.Boiler.Uid) > 0 {

			bConfs := ctl.ChannelConfigList(code)
			for _, cnf := range bConfs {
				go runtimeReload(cnf)
			}
		}
	}

	rawDis :=
		"UPDATE	`boiler_m163` " +
		"SET	`need_reload` = 0 " +
		"WHERE	`need_reload` = ?;"
	go dba.BoilerOrm.Raw(rawDis, nonce).Exec()
}

func (ctl *ParameterController) DataListNeedReload(nonce int) []orm.Params {
	var data 	[]orm.Params

	limit := ParamCtrl.ReloadLimit
	if limit <= 0 { limit = 600 }

	rawReady :=
		"UPDATE	`boiler_m163` " +
		"SET	`need_reload` = ? " +
		"WHERE	`need_reload` = 1 " +
		"LIMIT	" + strconv.FormatInt(limit, 10) + ";"

	raw :=
		/*
		"SELECT	`rtm`.* " +
		"FROM	`boiler_m163` AS `rtm`, `boiler_terminal_combined` AS `boiler` "  +
		//"WHERE	`boiler`.`terminal_code` = CAST(`rtm`.`Boiler_term_id` AS SIGNED) " +
		//"AND 	`boiler`.`terminal_set_id` = CAST(`rtm`.`Boiler_boiler_id` AS SIGNED) " +
		"WHERE	`boiler`.`terminal_code` = `rtm`.`Boiler_term_id` " +
		"  AND 	`boiler`.`terminal_set_id` = `rtm`.`Boiler_boiler_id` " +
		"  AND	`rtm`.`need_reload` = TRUE "
		*/

		"SELECT	* " +
		"FROM	`boiler_m163` "  +
		"WHERE	`need_reload` = ?;"

	var lg		logs.BoilerRuntimeLog
	lg.Name = "DataListNeedReload()"
	lg.TableName = "boiler_163m"
	lg.Query = raw
	lg.CreatedDate = time.Now()
	lg.Status = logs.BOILER_RUNTIME_LOG_STATUS_READY
	go DataCtl.AddData(&lg, false)

	if res, err := dba.BoilerOrm.Raw(rawReady, nonce).Exec(); err != nil {
		goazure.Error("Ready DataListNeedReload Error", err, res)
		return data
	}

	if num, err := dba.BoilerOrm.Raw(raw, nonce).Values(&data); err != nil {
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

	if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		RelatedSel("Parameter").
		Filter("IsDefault", true).Filter("IsDeleted", false).All(&dfConfs); err != nil || num == 0 {
		goazure.Error("Default Channel Config is Missing!", err, num)
	} else {
		goazure.Info("Get Boiler Channel Config:", num, "\n", bConfs)
	}

	if num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
		RelatedSel("Parameter").
		RelatedSel("Terminal").
		Filter("Terminal__TerminalCode", co).Filter("IsDeleted", false).All(&bConfs); err != nil || num == 0 {
		goazure.Error("Get Boiler Channel Config Error:", err, num)
	}
	/*else {
		goazure.Info("Get Boiler Channel Config:", num, "\n", bConfs)
	}*/

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
	for _, dc := range dfConfs {
		var isMatched bool = false
		for _, c := range bConfs {
			//goazure.Warning("ChannelType:", c.ChannelType, dc.ChannelType)
			//goazure.Warning("ChannelNumber:", c.ChannelNumber, dc.ChannelNumber)
			//goazure.Warning("Parameter.Id:", c.Parameter.Id, dc.Parameter.Id)
			if (c.ChannelType 	== dc.ChannelType &&
				c.ChannelNumber == dc.ChannelNumber) ||
				c.Parameter.Id == dc.Parameter.Id {

				isMatched = true

				matchIds = append(matchIds, c.Parameter.Id)
				break
			}
		}

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

	var confs		[]*models.RuntimeParameterChannelConfig
	var showConfs, hideConfs, normalConfs, defaultConfs 	[]*models.RuntimeParameterChannelConfig

	for _, c := range bConfs {
		if c.IsDefault {
			defaultConfs = append(defaultConfs, c)
		} else {
			switch c.Status {
			case models.CHANNEL_STATUS_SHOW:
				showConfs = append(showConfs, c)
			case models.CHANNEL_STATUS_HIDE:
				hideConfs = append(hideConfs, c)
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

type CnfBySeq 		[]*models.RuntimeParameterChannelConfig
func (a CnfBySeq) Len() int 			{ return len(a) }
func (a CnfBySeq) Swap(i, j int)		{ a[i], a[j] = a[j], a[i]}
func (a CnfBySeq) Less(i, j int) bool	{
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
	var bConfs		[]*models.RuntimeParameterChannelConfig
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
		HasCustom	bool
	}

	var js bJson
	js.HasCustom = hasCustom

	ctl.Data["json"] = js
	ctl.ServeJSON()

	return hasCustom
}

type ChannelConfig struct {
	TerminalCode	string		`json:"terminal_code"`
	ParameterId		int			`json:"parameter_id"`

	ChannelType		int			`json:"channel_type"`
	ChannelNumber	int			`json:"channel_number"`
	
	Status			int32		`json:"status"`
	SequenceNumber	int32		`json:"sequence_number"`

	SwitchValue		int32		`json:"switch_status"`
	
	Scale			float32		`json:"scale"`

	Ranges			[]bRange 	`json:"ranges"`
	
	IsDeleted		bool		`json:"is_deleted"`
}

type bRange struct {
	Name		string		`json:"name"`
	Min			int64		`json:"min"`
	Max			int64		`json:"max"`
	Value		int64		`json:"value"`
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

	matrix := make([][]*models.RuntimeParameterChannelConfig, 16)
	for i := range matrix {
		matrix[i] = make([]*models.RuntimeParameterChannelConfig, 6)
	}
	for _, cnf := range bConfs {
		switch cnf.ChannelType {
		case models.CHANNEL_TYPE_TEMPERATURE:
			fallthrough
		case models.CHANNEL_TYPE_ANALOG:
			matrix[cnf.ChannelNumber - 1][cnf.ChannelType - 1] = cnf
		case models.CHANNEL_TYPE_SWITCH:
			num := (cnf.ChannelNumber - 1) % 16
			ty := (cnf.ChannelNumber - 1) / 16 + models.CHANNEL_TYPE_SWITCH - 1
			matrix[num][ty] = cnf
		case models.CHANNEL_TYPE_CALCULATE:
			matrix[cnf.ChannelNumber - 1][cnf.ChannelType + 1] = cnf
		case models.CHANNEL_TYPE_RANGE:
			matrix[cnf.ChannelNumber - 1][cnf.ChannelType] = cnf
		default:

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
		goazure.Warn("[", i ,"]ChannelConfigUpdate: ", c)
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
			if	err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").
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
			if c.SwitchValue == 0 { c.SwitchValue = 1}
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
			if 	num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config_range").
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
		if 	num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config_range").
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
	var p 		models.RuntimeParameter
	var param 	models.RuntimeParameter

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &p); err != nil {
		e := fmt.Sprintln("Unmarshal Parameter JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	goazure.Warn("Parameter:", p)

	if err := dba.BoilerOrm.QueryTable("runtime_parameter").Filter("Id", p.Id).One(&param); err != nil {
		e := fmt.Sprintln("Read Parameter Error", err)
		goazure.Warn(e)

		param = p
		param.Length = p.Length
		param.Fix = p.Fix

		param.Category = runtimeParameterCategory(p.Category.Id)

		param.Medium = runtimeParameterMedium(0)
		param.AddBoilerMedium(0)
		param.CreatedBy = ctl.GetCurrentUser()
	}

	param.Name = p.Name
	param.Scale = p.Scale
	param.Unit = p.Unit
	param.Remark = p.Remark

	param.UpdatedBy = ctl.GetCurrentUser()
	param.IsDeleted = false

	if err := DataCtl.AddData(&param, true); err != nil {
		e := fmt.Sprintln("Add/Update Parameter Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
	}

	go ParamCtrl.RefreshParameters()
}

func (ctl *ParameterController) RuntimeParameterDelete() {
	var p models.RuntimeParameter

	if 	err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &p); err != nil {
		e := fmt.Sprintln("Unmarshal Parameter JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var cnfs []*models.RuntimeParameterChannelConfig
	if 	num, err := dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").Filter("Parameter__Id", p.Id).All(&cnfs); err != nil {
		goazure.Warn("Get ChannelConfig For This Parameter Error:", err, num)
	}

	for _, c := range cnfs {
		ctl.ChannelConfigDelete(c)
	}

	goazure.Warn("Delete Parameter:", p)

	if 	err := DataCtl.DeleteData(&p); err != nil {
		e := fmt.Sprintln("Delete Parameter Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))

		return
	}

	go ParamCtrl.RefreshParameters()
}


func runtimeParameter(pid int) *models.RuntimeParameter {
	param := models.RuntimeParameter{}

	if  err := dba.BoilerOrm.QueryTable("runtime_parameter").
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
	if	err := dba.BoilerOrm.QueryTable("runtime_parameter_category").
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
		fmt.Println("Records: ", records)

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
						value = runtimeParameter(int(pid))
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

					goazure.Info("Field(", fieldNames[j],":", da.FieldByName(fieldNames[j]).Kind(),"): ", field, value)
					da.FieldByName(fieldNames[j]).Set(reflect.ValueOf(value))
				}

				//fmt.Println("Da: ", da)
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

func (ctl *ParameterController) InitParameterChannelConfig(limit int64) {
	// generateDefaultChannelConfig()
	ParamCtrl.ReloadLimit = limit

	interval := time.Second * 5
	if !conf.IsRelease {
		interval = time.Minute * 1
	}
	ticker := time.NewTicker(interval)
	tick := func() {
		for t := range ticker.C {
			ParamCtrl.ChannelDataReload(t)
		}
	}

	go tick()
	go ParamCtrl.ChannelDataReload(time.Now())
}