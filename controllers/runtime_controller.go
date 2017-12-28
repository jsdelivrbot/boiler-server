package controllers

import (
	"fmt"

	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"
	//"github.com/pborman/uuid"

	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/models/caches"
	"github.com/AzureRelease/boiler-server/conf"

	"strconv"
	"io/ioutil"
	"os"
	"encoding/csv"
	"encoding/json"
	"log"
	"strings"
	"time"
	"math/rand"

	"reflect"

	"errors"
	"math"
)

type RuntimeController struct {
	MainController

	Parameters 				[]*models.RuntimeParameter

	StatusRunningDuration	map[string]time.Duration

	BoilerRanks				[]orm.Params
	BoilerHeats				[]*BoilerHeatAvg
}

type Status struct {
	Running		struct {

	}

	Alarm		struct {

	}
}

var RtmCtrl *RuntimeController = &RuntimeController{}

const rtmDefaultsPath string = "models/properties/runtime_defaults/"

func init() {
	RtmCtrl.MainController = *MainCtrl

	go RtmCtrl.RefreshParameters()
}

func (ctl *RuntimeController) GetBoilerRank() {
	ticker := time.NewTicker(time.Minute * 15)
	tick := func() {
		for t := range ticker.C {
			RtmCtrl.RefreshBoilerHeatRank(t)
		}
	}

	go tick()
	go RtmCtrl.RefreshBoilerHeatRank(time.Now())
}

func (ctl *RuntimeController) GetRunningDuration() {
	ticker := time.NewTicker(time.Minute * 30)
	tick := func() {
		for t := range ticker.C {
			RtmCtrl.RefreshStatusRunningDuration(t)
		}
	}

	go tick()
	go RtmCtrl.RefreshStatusRunningDuration(time.Now())
}

func (ctl *RuntimeController) RefreshParameters() {
	RtmCtrl.bWaitGroup.Add(1)

	var params []*models.RuntimeParameter
	qs := dba.BoilerOrm.QueryTable("runtime_parameter")
	if num, err := qs.RelatedSel("Category").OrderBy("Id").All(&params); err != nil || num == 0 {
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

	RtmCtrl.Parameters = params

	RtmCtrl.bWaitGroup.Done()
}

func (ctl *RuntimeController) RuntimeParameterList() {
	goazure.Warning("Ready to Get RuntimeParameterList")
	RtmCtrl.bWaitGroup.Wait()

	ctl.Data["json"] = RtmCtrl.Parameters
	ctl.ServeJSON()
}

func (ctl *RuntimeController) BoilerRuntimeCount() {
	usr := ctl.GetCurrentUser()
	var boilerUid string
	if ctl.Input()["boiler"] != nil && len(ctl.Input()["boiler"]) > 0 {
		boilerUid = ctl.Input()["boiler"][0]
	}

	var boilers []models.Boiler
	if !usr.IsAdmin() || len(boilerUid) > 0 {
		qs := dba.BoilerOrm.QueryTable("boiler")
		if usr.IsCommonUser() ||
			usr.Status == models.USER_STATUS_INACTIVE ||
			usr.Status == models.USER_STATUS_NEW {
			qs = qs.Filter("IsDemo", true)
		} else if usr.IsOrganizationUser() {
			orgCond := orm.NewCondition().Or("Enterprise__Uid", usr.Organization.Uid).Or("Factory__Uid", usr.Organization.Uid).Or("Installed__Uid", usr.Organization.Uid)
			cond := orm.NewCondition().AndCond(orgCond)
			qs = qs.SetCond(cond).Filter("IsDemo", false)
		}
		if len(boilerUid) > 0 {
			qs = qs.Filter("Uid", boilerUid)
		}

		if num, err := qs.Filter("IsDeleted", false).All(&boilers); num == 0 || err != nil {
			goazure.Error("Read BoilerList Error:", num, err)
		}
	}

	timeStr := time.Now().Format("2006-01-02")
	fmt.Println("timeStr:", timeStr)
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timeNumber := t.Unix()
	fmt.Println("timeNumber:", t, timeNumber)

	qr := dba.BoilerOrm.QueryTable("boiler_runtime")
	qr = qr.Filter("CreatedDate__gte", t)
	if !usr.IsAdmin() || len(boilerUid) > 0 {
		qr = qr.Filter("Boiler__in", boilers)
	}

	num, err := qr.Filter("IsDeleted", false).Count()
	if num == 0 || err != nil {
		goazure.Error("Read RuntimeCount Error: ", num, err)
	}

	ctl.Data["json"] = num
	ctl.ServeJSON()
}

const (
	RUNTIME_RANGE_DEFAULT = 0
	RUNTIME_RANGE_TODAY = 1
	RUNTIME_RANGE_THERE_DAY = 2
	RUNTIME_RANGE_WEEK = 3
	RUNTIME_RANGE_CUSTOM = 4
)

func (ctl *RuntimeController) BoilerRuntimeList() {
	usr := ctl.GetCurrentUser()
	goazure.Info("=== Ready to get Runtime ===\n", usr)

	b := bBoiler{}
	resStatus := 200
	resBody := "Success"
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	//goazure.Warning("Uid:", uid)
	var boiler *models.Boiler
	for _, ab := range MainCtrl.Boilers {
		if ab.Uid == b.Uid {
			boiler = ab
			break
		}
	}

	runtimes := func(param *models.RuntimeParameter, limit int) []orm.Params {
		var rtm []orm.Params

		tableName := "boiler_runtime"
		switch param.Id {
		case 1001:
			tableName = "boiler_runtime_cache_steam_temperature"
		case 1002:
			tableName = "boiler_runtime_cache_steam_pressure"
		case 1003:
			tableName = "boiler_runtime_cache_flow"

		case 1005:
			fallthrough
		case 1006:
			tableName = "boiler_runtime_cache_water_temperature"

		case 1014:
			fallthrough
		case 1015:
			tableName = "boiler_runtime_cache_smoke_temperature"

		case 1016:
			fallthrough
		case 1017:
			fallthrough
		case 1018:
			fallthrough
		case 1019:
			tableName = "boiler_runtime_cache_smoke_component"

		case 1021:
			fallthrough
		case 1022:
			tableName = "boiler_runtime_cache_environment_temperature"

		case 1201:
			tableName = "boiler_runtime_cache_heat"
		case 1202:
			tableName = "boiler_runtime_cache_excess_air"
		}

		if b.Range > RUNTIME_RANGE_TODAY {
			tableName = "boiler_runtime"
		}

		qr := dba.BoilerOrm.QueryTable(tableName)
		if tableName == "boiler_runtime" {
			qr = qr.RelatedSel("Alarm")
		}
		qr = qr.Filter("Boiler__Uid", boiler.Uid).Filter("Parameter__Id", param.Id)

		timeStr := time.Now().Format("2006-01-02")
		t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)

		switch b.Range {
		case RUNTIME_RANGE_DEFAULT:
			if limit > 0 {
				qr = qr.Limit(limit)
			}
		case RUNTIME_RANGE_TODAY:
			qr = qr.Filter("CreatedDate__gte", t).Limit(-1)
		case RUNTIME_RANGE_THERE_DAY:
			qr = qr.Filter("CreatedDate__gte", t.Add(time.Hour * 24 * -2)).Limit(-1)
		case RUNTIME_RANGE_WEEK:
			qr = qr.Filter("CreatedDate__gte", t.Add(time.Hour * 24 * -6)).Limit(-1)
		}

		//raw := fmt.Sprintf("SELECT boiler_runtime_view.*, boiler_alarm.alarm_level FROM boiler_runtime_view LEFT OUTER JOIN boiler_alarm ON boiler_runtime_view.alarm_id = boiler_alarm.uid WHERE boiler_runtime_view.boiler_id = '%s' && boiler_runtime_view.parameter_id = %d ORDER BY created_date %s", boiler.Uid, param.Id, lmt)
		//goazure.Warn("Runtime Raw:", raw)
		if num, err := qr.OrderBy("-CreatedDate").Values(&rtm, "CreatedDate", "Value"); err != nil || num == 0 {
			goazure.Error("Read BoilerRuntime Error", err, num)
			resStatus = 404
			resBody = "No Runtime Found!"
		} else {
			if tableName == "boiler_runtime" {
				for _, r := range rtm {
					value := float64(r["Value"].(int64)) * float64(param.Scale)
					pow10_n := math.Pow10(int(param.Fix))
					value = math.Trunc(value * pow10_n + 0.5) / pow10_n
					r["Value"] = value
				}
			}

			resStatus = 200
			resBody = "Success"
		}

		return rtm
	}

	alarmRules := func(boiler *models.Boiler) []*models.RuntimeAlarmRule {
		var rules []*models.RuntimeAlarmRule
		q := dba.BoilerOrm.QueryTable("runtime_alarm_rule")
		q = q.RelatedSel("Parameter__Category").RelatedSel("BoilerForm").RelatedSel("BoilerMedium").RelatedSel("BoilerFuelType")
		orCond := orm.NewCondition().Or("BoilerForm", boiler.Form).Or("BoilerForm__Id", 0).
			Or("BoilerMedium", boiler.Medium).Or("BoilerMedium__Id", 0).
			Or("BoilerFuelType", boiler.Fuel.Type).Or("BoilerFuelType__Id", 0)
		cond := orm.NewCondition().AndCond(orCond)
		q = q.SetCond(cond)
		if num, err := q.Filter("IsDeleted", false).All(&rules); err != nil || num == 0 {
			fmt.Println("Read RuntimeRule Error", err)
			resStatus = 404
			resBody = "No Runtime Found!"
		}

		return rules
	}

	var aRuntimes 	[][]orm.Params

	params := parameters(b.RuntimeQueue...)

	//fmt.Println("\n\nPARAMs:", params, "\nRTMs: \n", rtms)

	for _, p := range params {
		rtms := runtimes(p, b.Limit)
		aRuntimes = append(aRuntimes, rtms)
	}

	ab := aBoiler {
		Boiler: *boiler,
		Runtimes: aRuntimes,
		Parameters: params,
		Rules: alarmRules(boiler),
	}

	ctl.Data["json"] = ab
	ctl.ServeJSON()

	if resStatus != 200 {
		ctl.Ctx.Output.SetStatus(resStatus)
		ctl.Ctx.Output.Body([]byte(resBody))
	}
}

//TODO: Running Time Duration
func (ctl *RuntimeController) RefreshStatusRunningDuration(t time.Time) {
	dMap := map[string]time.Duration{}

	//goazure.Warn("\n=====", ctl, reflect.TypeOf(ctl), "\n", MainCtrl.Boilers ,"=====\n")
	for _, b := range MainCtrl.Boilers {
		//goazure.Warn("Ticker Duration Boiler:", b)

		var status []*caches.BoilerRuntimeCacheStatusRunning
		var duration time.Duration = time.Duration(0)
		qs := dba.BoilerOrm.QueryTable("boiler_runtime_cache_status_running")
		qs = qs.Filter("Boiler__Uid", b.Uid).Filter("Value", 1).
			OrderBy("CreatedDate")
		if num, err := qs.Filter("IsDeleted", false).All(&status); err != nil {
			goazure.Error("Get BoilerRuntimeCacheStatus Error:", err, num)
		}

		for i, st := range status {
			duration += time.Duration(st.Duration) * time.Microsecond
			if i == len(status) - 1 {
				duration += time.Now().Sub(st.CreatedDate)
			}
		}

		dMap[b.Uid] = duration
	}
	//goazure.Info("Durations:", dMap)

	ctl.StatusRunningDuration = dMap
}

func (ctl *RuntimeController) BoilerStatusRunningDuration() {
	usr := ctl.GetCurrentUser()
	duration := map[string]time.Duration{}
	boilers, err := BoilerCtrl.CurrentBoilerList(usr)
	if err != nil {
		goazure.Error("get current boiler list error!")
	}

	for _, b := range boilers {
		duration[b.Uid] = RtmCtrl.StatusRunningDuration[b.Uid]
	}

	ctl.Data["json"] = duration
	ctl.ServeJSON()
}

func (ctl *RuntimeController) BoilerRuntimeHistory() {
	//goazure.Info("=== Ready to get Runtime History ===")
	b := bBoiler{}
	resStatus := 200
	resBody := "Success"
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	if len(b.RuntimeQueue) == 0 {
		b.RuntimeQueue = []int{1001, 1002, 1003, 1005, 1006, 1015, 1014, 1016, 1021, 1202, 1201}
	}

	var	 histories	[]*caches.BoilerRuntimeHistory
	type hData		struct {
		Date		time.Time					`json:"date"`
		Data		[]*caches.History			`json:"data"`
	}
	type historyData struct {
		History		[]hData						`json:"history"`
		Parameters	[]*models.RuntimeParameter	`json:"parameter"`
	}

	var hisData historyData

	params := parameters(b.RuntimeQueue...)
	hisData.Parameters = params

	qh := dba.BoilerOrm.QueryTable("boiler_runtime_history")
	qh = qh.Filter("Boiler__Uid", b.Uid)
	qh = qh.Filter("CreatedDate__gte", b.StartDate).Filter("CreatedDate__lte", b.EndDate)
	qh = qh.OrderBy("-CreatedDate")
	if b.Limit != 0 {
		qh = qh.Limit(b.Limit)
	}
	if num, err := qh.All(&histories); err != nil || num == 0 {
		goazure.Error("Read BoilerRuntimeHistory Error", err, num)
		resStatus = 404
		resBody = "No Runtime Found!"
	} else {
		for _, his := range histories {
			his.Unmarshal()

			var h hData
			h.Date = his.CreatedDate
			h.Data = his.Histories

			hisData.History = append(hisData.History, h)
		}

		resStatus = 200
		resBody = "Success"
	}

	goazure.Info("Final HistoryData:", hisData)

	ctl.Data["json"] = hisData
	ctl.ServeJSON()

	if resStatus != 200 {
		ctl.Ctx.Output.SetStatus(resStatus)
		ctl.Ctx.Output.Body([]byte(resBody))
	}
}

func (ctl *RuntimeController) BoilerRuntimeInstants() {
	//goazure.Error("Instants Param:", ctl.Input())
	var b bBoiler
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	if len(b.RuntimeQueue) <= 0 {
		b.RuntimeQueue = []int{1201, 1015, 1002, 1202}
	}

	runtimes := func() []orm.Params {
		var ins []orm.Params
		q := dba.BoilerOrm.QueryTable("boiler_runtime_cache_instant")
		q = q.RelatedSel("Parameter")
		q = q.Filter("Boiler__Uid", b.Uid)
		if len(b.RuntimeQueue) > 0 {
			q = q.Filter("Parameter__Id__in", b.RuntimeQueue)
		}
		q = q.Filter("UpdatedDate__gt", time.Now().Add(time.Hour * -12))
		if num, err := q.Filter("IsDeleted", false).OrderBy("Parameter__Id").Values(&ins);
			/* "Runtime", "Boiler", "Parameter", "Alarm",
			"Value", //"CreatedDate",
			"ParameterName", "Unit",
			"AlarmLevel", "AlarmDescription",
			"IsDeleted", "IsDemo") */
			err != nil || num == 0 {
			goazure.Error("Read BoilerRuntime Instants Error", err, num)
		}
		//for _, in := range ins {
		//	goazure.Warn("Instants: ", in)
		//}

		return ins
	}

	rs := runtimes()
	rtms := []orm.Params{}
	for _, paramId := range b.RuntimeQueue {
		num := 0
		for _, r := range rs {
			if r["Parameter"] == int64(paramId) {
				rtms = append(rtms, r)
				num++
				break
			}
		}

		if num == 0 {
			in := orm.Params{}
			param := runtimeParameter(paramId)
			in["Parameter"] = param.Id
			in["ParameterName"] = param.Name

			rtms = append(rtms, in)
		}
	}

	r := rtms[0]
	var date string
	if r["UpdatedDate"] != nil {
		date = r["UpdatedDate"].(time.Time).Format("2006-01-02 15:04:05")
	}

	r["Date"] = date

	//goazure.Warning("==========>>>>>>>>> RTMs: ", rtms)

	ctl.Data["json"] = rtms
	ctl.ServeJSON()
}

func (ctl *RuntimeController) BoilerRuntimeDaily() {
	b := bBoiler{}
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	type Daily struct {
		Parameters	[]*models.RuntimeParameter
		Flows		[]orm.Params
		Heats		[]orm.Params
	}

	params := parameters(1003, 1201)
	flows := runtimeDaily(b.Uid, b.Limit, 0, 1003)
	heats := runtimeDaily(b.Uid, b.Limit, 0, 1201)

	da := Daily {
		Parameters: params,
		Flows: flows,
		Heats: heats,
	}

	ctl.Data["json"] = da
	ctl.ServeJSON()
}

func (ctl *RuntimeController) BoilerRuntimeDailyTotal() {
	//TODO: Index Of Date Not Boiler


	raw := "SELECT 	`flows`.`date`, AVG(`flows`.`value`) AS `flow`, AVG(`heats`.`value`) AS `heat` "
	raw += "FROM	`boiler_runtime_cache_flow_daily` AS `flows`, `boiler_runtime_cache_heat_daily` AS `heats` "
	raw += "WHERE	`flows`.`date` = `heats`.`date` AND `flows`.`boiler_id` = `heats`.`boiler_id` "
	raw += "GROUP BY `heats`.`date`, `heats`.`boiler_id` "
	raw += "ORDER BY `heats`.`date`;"

	type total struct {
		Date		time.Time
		Flow		float64
		Heat		float64
	}	//var aFlows [][]orm.Params
	//var aHeats [][]orm.Params

	var ttls []*total
	if num, err := dba.BoilerOrm.Raw(raw).QueryRows(&ttls); err != nil {
		goazure.Error("Ready Month Total Error:", err, num)
	}

	ctl.Data["json"] = ttls
	ctl.ServeJSON()
}

func (ctl *RuntimeController) BoilerHeatRank() {
	boilers, _ := ctl.boilers()
	heats := []*BoilerHeatAvg{}

	for _, ht := range RtmCtrl.BoilerHeats {
		for _, b := range boilers {
			if ht.Uid == b.Uid {
				heats = append(heats, ht)
				break
			}
		}
	}

	var c0, c1, c2, c3, g0, g1 BoilerHeatRank

	c0.FuelCate = "coal"
	c0.EvaporateCate = "D≤1"
	c0.EvaporateId = "c0"
	c0.Rank = true
	c1.FuelCate = "coal"
	c1.EvaporateCate = "1＜D≤2"
	c1.EvaporateId = "c1"
	c1.Rank = true
	c2.FuelCate = "coal"
	c2.EvaporateCate = "2＜D≤8"
	c2.EvaporateId = "c2"
	c2.Rank = true
	c3.FuelCate = "coal"
	c3.EvaporateCate = "8＜D≤20"
	c3.EvaporateId = "c3"
	c3.Rank = true

	g0.FuelCate = "gas"
	g0.EvaporateCate = "D≤2"
	g0.EvaporateId = "g0"
	g0.Rank = true
	g1.FuelCate = "gas"
	g1.EvaporateCate = "D＞2"
	g1.EvaporateId = "g1"
	g1.Rank = true

	for _, h := range heats {
		goazure.Warning(h)
	}
	
	ranks := RtmCtrl.BoilerRanks

	ctl.Data["json"] = ranks
	ctl.ServeJSON()
}

func (ctl *RuntimeController) RefreshBoilerRank(t time.Time) {

	var fileNameCoal string = "dba/sql/select_boiler_evaporate_rank_coal.sql"
	var fileNameGas string = "dba/sql/select_boiler_evaporate_rank_gas.sql"
	var sqlCoal, sqlGas string
	if buf, err := ioutil.ReadFile(fileNameCoal); err != nil {
		goazure.Error("read SQL File", fileNameCoal, "Error", err)
		return
	} else {
		sqlCoal = string(buf)
		goazure.Info("Read SQL:", sqlCoal)
	}

	if buf, err := ioutil.ReadFile(fileNameGas); err != nil {
		goazure.Error("read SQL File", fileNameGas, "Error", err)
		return
	} else {
		sqlGas = string(buf)
	}

	var ranksCoal []orm.Params
	var ranksGas  []orm.Params
	if num, err := dba.BoilerOrm.Raw(sqlCoal).Values(&ranksCoal); err != nil {
		goazure.Error("Select Boiler Evaporate Rank Coal Error:", err, num)
	}
	if num, err := dba.BoilerOrm.Raw(sqlGas).Values(&ranksGas); err != nil {
		goazure.Error("Select Boiler Evaporate Rank Gas Error:", err, num)
	}

	ranks := append(ranksCoal, ranksGas...)

	RtmCtrl.BoilerRanks = ranks
}

func (ctl *RuntimeController) RefreshBoilerHeatRank(t time.Time) {
	var fileNameSQL string = "dba/sql/select_boiler_runtime_heat_avg_week.sql"

	var sql string
	if buf, err := ioutil.ReadFile(fileNameSQL); err != nil {
		goazure.Error("read SQL File", fileNameSQL, "Error", err)
		return
	} else {
		sql = string(buf)
		goazure.Info("Read SQL:", sql)
	}

	var heats []*BoilerHeatAvg
	if num, err := dba.BoilerOrm.Raw(sql).QueryRows(&heats); err != nil {
		goazure.Error("Select Boiler Heat Rank Error:", err, num)
	}

	RtmCtrl.BoilerHeats = heats
}

func (ctl *RuntimeController) boilers() ([]*models.Boiler, error) {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		if ctl.Input()["token"] == nil || len(ctl.Input()["token"]) == 0 {
			return nil, errors.New("no current user, no availabel token")
		}
		token := ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
			return nil, err
		}
	}

	var boilers []*models.Boiler
	var err error

	if boilers, err = BoilerCtrl.CurrentBoilerList(usr); err != nil {
		goazure.Error("Get CurrentBoilerList Error:", err)
	}

	return boilers, nil
}

func paramSelect(ps []*models.RuntimeParameter, pid int) (*models.RuntimeParameter) {
	var ret *models.RuntimeParameter
	for _, p := range ps {
		if p.Id == int64(pid) {
			ret = p
			break
		}
	}

	return ret
}

func parameters(ids... int) []*models.RuntimeParameter {
	var params []*models.RuntimeParameter
	q := dba.BoilerOrm.QueryTable("runtime_parameter")

	q = q.Filter("Id__in", ids)
	if num, err := q.OrderBy("Id").All(&params); err != nil || num == 0 {
		goazure.Error("Read Param Error", err)
	}

	var aParams []*models.RuntimeParameter
	for _, id := range ids {
		param := paramSelect(params, id)
		aParams = append(aParams, param)
	}

	return aParams
}

func runtimeDaily(boilerId string, limit int, offset int, pid int) []orm.Params {
	var rtm []orm.Params
	var tableName string
	switch pid {
	case 1003:
		tableName = "boiler_runtime_cache_flow_daily"
	case 1201:
		tableName = "boiler_runtime_cache_heat_daily"
	default:
		return rtm

	}

	q := dba.BoilerOrm.QueryTable(tableName)
	q = q.Filter("Boiler__Uid", boilerId).
		Filter("Date__gte", time.Now().Add(time.Hour * -24 * time.Duration(limit + offset))).Limit(limit)
	if num, err := q.OrderBy("Date").Values(&rtm, "Date", "Value"); err != nil || num == 0 {
		goazure.Error("Read BoilerRuntime Error", err)
	}

	return rtm
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
	if err := dba.BoilerOrm.QueryTable("runtime_parameter_category").Filter("Id", cateId).One(&category); err != nil {
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

//******************************* Init

func (ctl *RuntimeController) InitRuntimeParameters() {
	generateDefaultRuntimeParameterCategories()
	//generateDefaultRuntimeParameterMediums()
	generateDefaultRuntimeParameters()
}

func generateDefaultRuntimeParameterCategories() {
	var cate10, cate11, cate12 models.RuntimeParameterCategory
	cate10.Id = 10
	cate10.Name = "模拟量采集参数"
	cate10.NameEn = "Analogue Collection Parameters"

	cate11.Id = 11
	cate11.Name = "开关量采集参数"
	cate11.NameEn = "Switch Signal Collection Parameters"

	cate12.Id = 12
	cate12.Name = "计算参数"
	cate12.NameEn = "Calculation Parameters"

	addRuntimeParameterCategory(cate10)
	addRuntimeParameterCategory(cate11)
	addRuntimeParameterCategory(cate12)
}

func generateDefaultRuntimeParameterMediums() {
	var medium1, medium2, medium3, medium4, medium5, medium6, medium7, medium8, medium0 models.RuntimeParameterMedium

	medium1.Id = 1
	medium1.Name = "蒸汽"
	medium1.NameEn = "Steam"

	medium2.Id = 2
	medium2.Name = "水"
	medium2.NameEn = "Water"

	medium3.Id = 3
	medium3.Name = "空气"
	medium3.NameEn = "Air"

	medium4.Id = 4
	medium4.Name = "电"
	medium4.NameEn = "Electricity"

	medium5.Id = 5
	medium5.Name = "燃料"
	medium5.NameEn = "Fuel"

	medium6.Id = 6
	medium6.Name = "烟气"
	medium6.NameEn = "Exhaust Gas"

	medium7.Id = 7
	medium7.Name = "炉膛"
	medium7.NameEn = "Furnace"

	medium8.Id = 8
	medium8.Name = "负荷"
	medium8.NameEn = "Load"

	medium0.Id = 0
	medium0.Name = "其他"
	medium0.NameEn = "Other"

	addRuntimeParameterMedium(medium1)
	addRuntimeParameterMedium(medium2)
	addRuntimeParameterMedium(medium3)
	addRuntimeParameterMedium(medium4)
	addRuntimeParameterMedium(medium5)
	addRuntimeParameterMedium(medium6)
	addRuntimeParameterMedium(medium7)
	addRuntimeParameterMedium(medium8)
}

func generateDefaultRuntimeParameters() {
	//const msgFmtPath string = msgDefautlsPath + "formatters/"

	files, _ := ioutil.ReadDir(rtmDefaultsPath)
	fmt.Println("RuntimeParameters FileCount: ", len(files))

	const headerRowNo int = 0
	const fieldRowNo int = 1

	for page := 0; page < len(files) - 1; page++ {
		filename := rtmDefaultsPath + "parameters-" + strconv.Itoa(10 + page) + ".csv"
		file, errFile := os.Open(filename)
		if errFile != nil {
			goazure.Error("Read File Error:", errFile)
			continue
		}
		reader := csv.NewReader(file)

		records, errRead := reader.ReadAll()
		if errRead != nil {
			log.Fatal(errRead)
		}
		fmt.Println("Records: ", records)

		for i, row := range records {
			if i <= fieldRowNo {

			} else {
				paramId, _ := strconv.ParseInt(row[0], 10, 32)
				name := row[1]
				nameEn := row[2]
				categoryId, _ := strconv.ParseInt(row[3], 10, 32)
				mediumId, _ := strconv.ParseInt(row[4], 10, 32)
				applicableBoilerTypes := row[5]
				unit := row[6]
				scale, _ := strconv.ParseFloat(row[7], 32)
				fix, _ := strconv.ParseInt(row[8], 10, 32)
				length, _ := strconv.ParseInt(row[9], 10, 32)

				remark := row[10]

				param := models.RuntimeParameter{}

				cateIdStr := strconv.FormatInt(categoryId, 10)
				paramIdStr := strconv.FormatInt(paramId, 10)
				param.Id, _ = strconv.ParseInt(cateIdStr + paramIdStr, 10, 64)
				param.ParamId = int32(paramId)

				param.Category = runtimeParameterCategory(categoryId)
				param.Medium = runtimeParameterMedium(mediumId)

				param.Name = name
				param.NameEn = nameEn
				param.Remark = remark

				param.Unit = unit
				param.Scale = float32(scale)
				param.Fix = int32(fix)
				param.Length = int32(length)

				addRuntimeParameter(&param)

				boilerTypes := strings.Split(applicableBoilerTypes, "")
				for _, boilerType := range boilerTypes {
					typeId, _ := strconv.ParseInt(boilerType, 10, 32)
					param.AddBoilerMedium(typeId)
				}
				if len(boilerTypes) == 0 {
					param.AddBoilerMedium(0)
				}
			}
		}
	}
}

func addRuntimeParameter(param *models.RuntimeParameter) (error) {
	var err error

	if param.Scale == 0 {
		param.Scale = 1.0
	}
	if param.Length == 0 {
		param.Length = 2
	}

	err = DataCtl.AddData(param, true, "ParamId", "Category")

	return err
}



func addRuntimeParameterCategory(category models.RuntimeParameterCategory) (error) {
	err := DataCtl.AddData(&category, true)

	return err
}

func addRuntimeParameterMedium(medium models.RuntimeParameterMedium) (error) {
	err := DataCtl.AddData(&medium, true)

	return err
}

//TODO: Ready to Discard.
func (ctl *RuntimeController) RefreshRuntimeStatusRunning() {
	var boilers []*models.Boiler
	qb := dba.BoilerOrm.QueryTable("boiler").Limit(-1)
	if num, err := qb.All(&boilers, "Uid", "Name"); err != nil {
		goazure.Error("Get Boilers Error:", err, num)
	}

	for _, b := range boilers {
		var status []*caches.BoilerRuntimeCacheStatus
		qs := dba.BoilerOrm.QueryTable("boiler_runtime_cache_status")
		qs = qs.Filter("Boiler__Uid", b.Uid).OrderBy("CreatedDate").Limit(-1)
		if num, err := qs.Filter("IsDeleted", false).All(&status); err != nil {
			goazure.Error("Get BoilerRuntimeCacheStatus Error:", err, num)
		}

		var isBurning bool = false
		var sinceDate time.Time
		var duration time.Duration
		var lastRun *caches.BoilerRuntimeCacheStatusRunning
		for i, st := range status {
			goazure.Warn("Status:", i, st)
			if (!isBurning && st.Value > 0) ||
				(isBurning && st.Value <= 0) {
				isBurning = st.Value > 0
				if !sinceDate.IsZero() && lastRun != nil {
					duration = st.CreatedDate.Sub(sinceDate)
					lastRun.Duration = int64(duration / time.Microsecond)
					if err := DataCtl.AddData(lastRun, true, "Boiler", "CreatedDate"); err != nil {
						goazure.Error("Update Running Duratiuon Error:", err, lastRun)
					}
				}
				sinceDate = st.CreatedDate

				timeStr := sinceDate.Format("2006-01-02")
				fmt.Println("timeStr:", timeStr)
				t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
				timeNumber := t.Unix()
				fmt.Println("timeNumber:", t, timeNumber)

				var run caches.BoilerRuntimeCacheStatusRunning
				param := runtimeParameter(1107)
				value := 0

				if isBurning {
					value = 1
				}

				run.Name = b.Name
				run.Boiler = b
				run.Runtime = st.Runtime
				run.Parameter = param
				run.ParameterName = param.Name
				run.Value = float64(value)

				run.AlarmLevel = 0
				run.AlarmDescription = ""

				run.CreatedDate = sinceDate
				run.Date = t

				if err := DataCtl.AddData(&run, true, "Boiler", "CreatedDate"); err != nil {
					goazure.Error("Add Running State Error:", err, run)
				}

				lastRun = &run
			}
		}
		if lastRun != nil {
			duration = time.Now().Sub(sinceDate)
			lastRun.Duration = int64(duration / time.Microsecond)
			if err := DataCtl.AddData(lastRun, true, "Boiler", "CreatedDate"); err != nil {
				goazure.Error("Update Running Duratiuon Error:", err, lastRun)
			}
		}

	}

}

//=========DEMO=========//

func (ctl *RuntimeController) UpdateRuntimeHistory(since time.Time, until time.Time) {
	//if since.IsZero() {
	//	since = time.Now().Add(time.Hour * -24)
	//}
	goazure.Warning("UpdateRuntimeHistory")

	qc := dba.BoilerOrm.QueryTable("boiler_runtime")
	count, err := qc.Count()
	goazure.Warn("boiler_runtime Count:", count, err)

	qs := dba.BoilerOrm.QueryTable("boiler_runtime")
	qs = qs.RelatedSel("Parameter").RelatedSel("Alarm")
	qs = qs.Filter("CreatedDate__gte", since)
	if !until.IsZero() {
		goazure.Info("Until:", until)
		qs = qs.Filter("CreatedDate__lte", until)
	}
	qs = qs.Filter("IsDeleted", false).OrderBy("-CreatedDate")

	for i := int64(0); i < count; i += 1000 {
		var runtimes []*models.BoilerRuntime
		if num, err := qs.Offset(i).All(&runtimes); err != nil || num == 0 {
			goazure.Error("Read Runtime for History Error!", err, num)
		} else {
			goazure.Info("Get Runtime for History:", num)
		}

		var baseTime time.Time
		//var value float32
		//var valueSum float32
		//var count int

		var alarmLevel int = models.RUNTIME_ALARM_LEVEL_NORMAL
		for _, r := range runtimes {
			tm := r.CreatedDate
			if baseTime.IsZero() ||
				tm.Sub(baseTime) > time.Minute * 5 ||
				tm.Sub(baseTime) < time.Minute * -5 {
				baseTime = tm
				//valueSum = 0.0
				//count = 0
			}

			history := caches.BoilerRuntimeCacheHistory{}
			history.Boiler = r.Boiler
			history.CreatedDate = baseTime
			if err := DataCtl.ReadData(&history, "Boiler", "CreatedDate"); err != nil {
				goazure.Warning("Read Runtime History Error:", err)
			}

			/*
			if r.Value > 0 {
				v := float32(r.Value) * r.Parameter.Scale
				valueSum += v
				count++
			}
			value = valueSum / float32(count)
			*/

			value := float32(r.Value) * r.Parameter.Scale
			if r.Alarm != nil && alarmLevel < int(r.Alarm.AlarmLevel) {
				alarmLevel = int(r.Alarm.AlarmLevel)
			}

			reflect.ValueOf(&history).Elem().FieldByName(fmt.Sprintf("P%d", r.Parameter.Id)).Set(reflect.ValueOf(value))
			reflect.ValueOf(&history).Elem().FieldByName(fmt.Sprintf("A%d", r.Parameter.Id)).Set(reflect.ValueOf(alarmLevel))
			//history.UpdatedDate = tm

			if err := DataCtl.AddData(&history, true, "Boiler", "CreatedDate"); err != nil {
				goazure.Error("Add RuntimeHistory, ", history.Boiler.Name, " Error:", err)
			}
		}
	}
}

func (ctl *RuntimeController) GenerateBoilerStatus(isOn bool) {
	var status []Runtime
	status = append(status, Runtime{
		Pid: 1107,
		BaseValue: 0,
		DynamicValue: 3,
	})

	generateStatus := func(date time.Time, stat Runtime) {
		statTicker := time.NewTicker(time.Minute * 5)
		statGen := func(t time.Time) {
			for i, b := range MainCtrl.Boilers {
				if !BoilerCtrl.IsGenerateData(b.Uid) {
					continue
				}

				value := int64((i + t.Hour()) % 3)
				//goazure.Error("go generateBoilerStatus:", stat, t, b.Name, value)
				generateBoilerStatus(b, t, value, stat)
			}
		}

		go statGen(time.Now())
		for t := range statTicker.C {
			go statGen(t)
		}
	}

	BoilerCtrl.RefreshGlobalBoilerList()
	for _, st := range status {
		go generateStatus(time.Now(), st)
	}

	ticker := time.NewTicker(time.Hour * 3)
	for t := range ticker.C {
		for _, st := range status {
			generateStatus(t, st)
		}
	}
}

func (ctl *RuntimeController) GenerateBoilerRuntime(isOn bool) {
	var runtimes []Runtime
	runtimes = append(runtimes, Runtime{
		Pid: 1001,			//蒸汽温度
		BaseValue: 1000,
		DynamicValue: 1200,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1002,			//蒸汽压力
		BaseValue: 400,
		DynamicValue: 800,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1003,			//瞬时流量
		BaseValue: 20,		//int64(float32(boiler.EvaporatingCapacity) / float32(2) * float32(1000)),
		DynamicValue: 40,	//int64(float32(boiler.EvaporatingCapacity) / float32(2) * float32(1000)),
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1005,			//给水温度(低)
		BaseValue: 100,
		DynamicValue: 500,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1006,			//给水温度(高)
		BaseValue: 260,
		DynamicValue: 440,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1014,			//排烟温度(高)
		BaseValue: 800,
		DynamicValue: 800,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1015,			//排烟温度(低)
		BaseValue: 500,
		DynamicValue: 500,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1016,			//烟气O2含量
		BaseValue: 0,
		DynamicValue: 64,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1017,			//烟气CO含量
		BaseValue: 0,
		DynamicValue: 42,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1018,			//烟气CO2含量
		BaseValue: 0,
		DynamicValue: 126,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1019,			//烟气NOX含量
		BaseValue: 0,
		DynamicValue: 40,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1021,			//室内温度(环境温度)
		BaseValue: -40,
		DynamicValue: 440,
	})

	runtimes = append(runtimes, Runtime{
		Pid: 1201,			//热效率
		BaseValue: 700,
		DynamicValue: 260,
	})
	runtimes = append(runtimes, Runtime{
		Pid: 1202,			//过量空气系数
		BaseValue: 100,
		DynamicValue: 60,
	})

	generateRuntime := func(date time.Time, rtm Runtime) {
		var currentValue int64 = rtm.BaseValue + int64(float64(rtm.DynamicValue) * 0.4)
		for _, b := range MainCtrl.Boilers {
			if !BoilerCtrl.IsGenerateData(b.Uid) {
				continue
			}

			if !BoilerCtrl.IsBurning(b.Uid) {
				continue
			}

			currentValue = generateBoilerRuntime(b, date, rtm, currentValue)
		}
	}

	interval := 5
	if conf.IsRelease {
		interval = 45
	}

	tickBoiler := time.NewTicker(time.Minute * 15)
	bTick := func() {
		for range tickBoiler.C {
			BoilerCtrl.RefreshGlobalBoilerList()
		}
 	}

	ticker := time.NewTicker(time.Second * time.Duration(interval))
	tick := func() {
		for t := range ticker.C {
			on1 := rand.Float32() * 3 > 2
			for _, r := range runtimes {
				on2 := rand.Float32() * 2 > 1
				if on1 && on2 {
					generateRuntime(t, r)
				} else {
					goazure.Warning("time skip", t)
				}
			}
		}
	}

	go bTick()
	go tick()
}

func triggerRules(boiler *models.Boiler, param *models.RuntimeParameter) []models.RuntimeAlarmRule {
	var rules []models.RuntimeAlarmRule

	q := dba.BoilerOrm.QueryTable("runtime_alarm_rule")
	q = q.RelatedSel("Parameter").
		RelatedSel("BoilerForm").
		RelatedSel("BoilerMedium").RelatedSel("BoilerFuelType")
	cond := orm.NewCondition()
	condForm := cond.Or("BoilerForm__Id", boiler.Form.Id).Or("BoilerForm__Id", 0)
	condMedium := cond.Or("BoilerMedium__Id", boiler.Medium.Id).Or("BoilerMedium__Id", 0)
	condFuel := cond.Or("BoilerFuelType__Id", boiler.Fuel.Type.Id).
		Or("BoilerFuelType__Id", 0)
	condition := cond.AndCond(condForm).AndCond(condMedium).
		AndCond(condFuel)
	q = q.SetCond(condition).Filter("Parameter", param)
	if num, err := q.Filter("IsDeleted", false).All(&rules); err != nil || num == 0 {
		fmt.Println("Read RuntimeRule Error", err, num)
	}

	return rules
}

func generateBoilerRuntime(boiler *models.Boiler, date time.Time, run Runtime, currentValue int64) int64 {
	if currentValue < 0 {
		currentValue = run.BaseValue
	}
	currentDynamic := currentValue - run.BaseValue
	rate := float64(currentDynamic) / float64(run.DynamicValue)
	step := int64(float64(run.DynamicValue) / 100.0 * rand.Float64() * 6)
	//fl := rand.Float32()
	//fmt.Printf("t.Round(%6s) = %s | %v | %v \n", time.Minute, date.Format("2006-01-02 15:04:05.999999999"), fl, date)
	//goazure.Info("Runtime: ", run)
	var value int64 = currentValue
	fl := rand.Float64()
	if fl > rate {
		value = currentValue + step
	} else if fl < rate {
		value = currentValue - step
	}

	//value := int64(float32(run.BaseValue) + fl * float32(run.DynamicValue))
	param := runtimeParameter(int(run.Pid))

	//if param.Category.Id == int32(11) && value > 0 {
	//	value = 1
	//}

	rtm := models.BoilerRuntime {}

	//rtm.Uid = uuid.New()
	rtm.Boiler = boiler
	rtm.Parameter = param
	rtm.Value = value
	rtm.CreatedDate = date
	rtm.UpdatedDate = date
	rtm.Status = models.RUNTIME_STATUS_NEW

	rtm.IsDemo = true

	if rtm.Parameter == nil {
		goazure.Error("Invalid ParameterId:", rtm)
		panic("ParameterId can not be zero")
	}

	if err := DataCtl.AddData(&rtm, false); err != nil {
		goazure.Error("Add Runtime", rtm.Parameter.Name, "Error:", err)
	}

	return value
}

func generateBoilerStatus(boiler *models.Boiler, date time.Time, value int64, run Runtime) {
	if value < 0 {
		value = int64(float32(run.BaseValue) + rand.Float32() * float32(run.DynamicValue))
	}

	param := runtimeParameter(int(run.Pid))

	rtm := models.BoilerRuntime {}

	//rtm.Uid = uuid.New()
	rtm.Boiler = boiler
	rtm.Parameter = param
	rtm.Value = value
	rtm.CreatedDate = date
	//rtm.UpdatedDate = date

	rtm.Status = models.RUNTIME_STATUS_NEW

	rtm.IsDemo = true

	if err := DataCtl.AddData(&rtm, false); err != nil {
		goazure.Error("Add Status", rtm.Parameter.Name, "Error:", err)
	}
}

type aBoiler struct {
	models.Boiler
	Runtimes		[][]orm.Params
	Parameters		[]*models.RuntimeParameter
	Rules			[]*models.RuntimeAlarmRule
}

type bBoiler struct {
	Uid				string		`json:"uid"`
	RuntimeQueue	[]int		`json:"runtimeQueue"`
	Range			int			`json:"range"`
	Limit			int			`json:"limit"`
	StartDate		time.Time	`json:"startDate"`
	EndDate			time.Time	`json:"endDate"`
}

type Runtime struct {
	Pid				int32
	BaseValue		int64
	DynamicValue	int64
}

type BoilerHeatAvg struct {
	Uid		string
	Heat	float64
	Week	int64
}

type BoilerHeatRank struct {
	Count			int64		`json:"count"`
	FuelCate		string		`json:"fuel_cate"`
	EvaporateCate	string		`json:"evaporate_cate"`
	EvaporateId		string		`json:"evaporate_id"`
	Rank			bool		`json:"rank"`
}

//TODO: -------HISTORY IMPORT---------
func (ctl *RuntimeController) ImportHistoryDataFromOldTable(tableName string) {
	interval := time.Minute * 3

	ticker := time.NewTicker(interval)
	tick := func() {
		for t := range ticker.C {
			importHistoryDataFromOldTable(tableName, t)
		}
	}

	go tick()

	go importHistoryDataFromOldTable(tableName, time.Now())
}

func importHistoryDataFromOldTable(tableName string, t time.Time) {
	var hMaps []orm.Params
	qh := dba.BoilerOrm.QueryTable(tableName)
	if num, err := qh.OrderBy("CreatedDate").Limit(1000).Values(&hMaps); err != nil {
		goazure.Error("Get HistoryData For Reload Error!", err, num)
	} else {
		goazure.Info("Get HistoryData For Reload:", num)
	}

	//var histories []*caches.BoilerRuntimeHistory

	for _, hMap := range hMaps {
		//goazure.Info("Get hMap:", hMap)
		var boiler models.Boiler
		var hData caches.BoilerRuntimeHistory

		if err := dba.BoilerOrm.QueryTable("boiler").Filter("Uid", hMap["Boiler"]).One(&boiler); err != nil {
			goazure.Error("Get Boiler For OLD History Error:", err, "\n", hMap)
			continue
		}

		hData.Boiler = &boiler
		hData.Name = boiler.Name
		hData.CreatedDate = hMap["CreatedDate"].(time.Time)
		hData.IsDemo = hMap["IsDemo"].(bool)
		hData.Remark = hMap["Remark"].(string)

		for k, v := range hMap {
			//goazure.Info("Get hMap:", k, "|", v)
			if strings.HasPrefix(k, "P") {
				key := strings.TrimPrefix(k, "P")
				aK := strings.Replace(key, "P", "A", 1)

				var h caches.History
				pid, er := strconv.ParseInt(key, 10, 64)
				if er != nil {
					goazure.Error("Value-Key Error:", er)
					continue
				}

				h.ParameterId = pid
				if v != nil && v.(float64) != 0.0 {
					h.Value = v.(float64)
					if hMap[aK] != nil {
						h.Alarm = hMap[aK].(int)
					}
					//goazure.Info("Get hData:", h)
				}

				hData.Histories = append(hData.Histories, &h)
			}
		}

		hData.Marshal()

		//histories = append(histories, &hData)

		if err := DataCtl.AddData(&hData, true); err != nil {
			goazure.Error("Insert TRANS History Error:", err, hData)
		} else {
			goazure.Info("Insert TRANS History Done.")
			rawDel := "DELETE FROM `" + tableName + "` WHERE `id` = ?;"
			if res, err := dba.BoilerOrm.Raw(rawDel, hMap["Id"]).Exec(); err != nil {
				goazure.Error("Delete OLD History Data Error:", err, "\n", rawDel)
			} else {
				rowNum, err := res.RowsAffected()
				goazure.Info("Deleted OLD hData:", rowNum, err)
			}
		}
	}

	/*
	if num, err := dba.BoilerOrm.InsertMulti(1000, histories); err != nil {
		goazure.Error("Insert TRANS History Error:", err, num)
	} else {
		goazure.Info("Insert TRANS History Done:", num)
		for i, d := range histories {
			if n, e := dba.BoilerOrm.Delete(d); e != nil {
				goazure.Error("Delete OLD HistoryData Error:", e, i, "/", n)
			}
		}
	}
	*/
}