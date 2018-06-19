package controllers

import (
	"fmt"

	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"
	"github.com/pborman/uuid"

	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/models/caches"
	"github.com/AzureRelease/boiler-server/conf"
	"github.com/AzureRelease/boiler-server/models/logs"

	"strconv"
	"io/ioutil"
	"os"
	"encoding/csv"
	"encoding/json"
	"log"
	"strings"
	"time"
	"math/rand"

	"errors"
	"math"
)

type RuntimeController struct {
	MainController

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

/*type IssuedInformationPush struct {
	BoilerUid string
	BoilerName string
	UserUid string
	AlarmCount int
	ParameterId int
	ParameterAlarmCount int
	Description string
	BoilerRuntime float64
	status bool
}*/

type BoilerPush struct {
	Uid string
	Name string
}

var RtmCtl *RuntimeController = &RuntimeController{}

const rtmDefaultsPath string = "models/properties/runtime_defaults/"

func init() {
	RtmCtl.MainController = *MainCtrl
}

func (ctl *RuntimeController) InitInformationPush() {
	ticker := time.NewTicker(time.Hour*1)
	tick := func() {
		for t := range ticker.C {
			if (t.Weekday().String() == "Monday")&& t.Hour() == 8{
				ctl.WeekInformation(t)
			}
		}
	}
	go tick()
	if (time.Now().Weekday().String() == "Monday")&& time.Now().Hour() == 8{
		go ctl.WeekInformation(time.Now())
	}
}
func (ctl *RuntimeController) WeekInformation(t time.Time) {
	var boiler []BoilerPush
	sql:="select uid ,name from boiler"
	if _,err:= dba.BoilerOrm.Raw(sql).QueryRows(&boiler);err!=nil{
		goazure.Error("Query boiler Error",err)
	}
	for _,b := range boiler {
		bid := b.Uid
		bName := b.Name
		var alarmCount int
		var parameterId int
		var parameterAlarmCount int
		var description string
		var runtime float64
		sql := "select count(*) from boiler_alarm_history where boiler_id=? and date(date_sub(?,interval 7 day)) <= date(start_date) and date(start_date) < date(?)"
		if err := dba.BoilerOrm.Raw(sql, bid, t,t).QueryRow(&alarmCount); err != nil {
			goazure.Error("Query boiler_alarm_history Error",err)
		}
		if alarmCount == 0{
			parameterId = 0
			parameterAlarmCount = 0
			description = "无"
		} else {
			var param []orm.Params
			maxsql:= "select parameter_id,description, count(*) as count  from boiler_alarm_history where boiler_id=? and date(date_sub(?,interval 7 day)) <= date(start_date) and date(start_date) < date(?) " +
				"group by parameter_id order by count desc limit 1"
			if _,err:=dba.BoilerOrm.Raw(maxsql,bid,t,t).Values(&param);err!=nil{
				goazure.Error("Query boiler_alarm_history Error",err)
			}
			if pid,err:=strconv.Atoi(param[0]["parameter_id"].(string));err!=nil{
				goazure.Error("Parse Int Error",err)
			} else {
				parameterId = pid
			}
			if pCount,err:= strconv.Atoi(param[0]["count"].(string));err!=nil{
				goazure.Error("Parse Int Error",err)
			} else {
				parameterAlarmCount = pCount
			}
			description = param[0]["description"].(string)
		}
		runtimeSql := "select Fire_duration from boiler_terminal_combined btc , issued_boiler_fire_duration ibf where btc.terminal_code = ibf.Boiler_term_id " +
			"and btc.boiler_id=? and btc.terminal_set_id = 1"
		if err:=dba.BoilerOrm.Raw(runtimeSql,bid).QueryRow(&runtime);err!=nil{
			goazure.Error("Query Fire_duration Error",err)
		}
		infoSql:="insert into issued_week_information_log(uid,boiler_id,boiler_name,start_date,end_date,create_time,alarm_count,parameter_id,parameter_alarm_count,description,boiler_runtime) "+
			"values(uuid(),?,?,date(date_sub(?,interval 7 day)),date(?),now(),?,?,?,?,?)"
		if _,err:=dba.BoilerOrm.Raw(infoSql,bid,bName,t,t,alarmCount,parameterId,parameterAlarmCount,description,runtime).Exec();err!=nil{
			goazure.Error("Insert issued_information_push_log Error",err)
		}
	}
	durationSql:="update issued_boiler_fire_duration " +
		"set Start_statistic_time = now(), End_statistic_time=now(),Fire_duration=0"
	if _,err:=dba.BoilerOrm.Raw(durationSql).Exec();err!=nil{
		goazure.Error("Update issued_boiler_fire_duration Error",err)
	}
}

/*func (ctl *RuntimeController) PushInformation(t time.Time) {
	var boiler []BoilerPush
	sql:="select distinct bms.boiler_id , b.name  from boiler_message_subscriber bms,boiler b where b.uid= bms.boiler_id"
	if _,err:=dba.BoilerOrm.Raw(sql).QueryRows(&boiler);err!=nil{
		goazure.Error("Query boiler_message_subscriber Error",err)
	}
	for _,b:= range boiler {
		bid := b.BoilerId
		bName := b.Name
		var alarmCount int
		var parameterId int
		var parameterAlarmCount int
		var description string
		var runtime float64
		sql := "select count(*) from boiler_alarm_history where boiler_id=? and date(date_sub(?,interval 7 day)) <= date(start_date) and date(start_date) < date(?)"
		if err := dba.BoilerOrm.Raw(sql, bid, t,t).QueryRow(&alarmCount); err != nil {
			goazure.Error("Query boiler_alarm_history Error",err)
		}
		if alarmCount == 0{
			parameterId = 0
			parameterAlarmCount = 0
			description = ""
		} else {
			var param []orm.Params
			maxsql:= "select parameter_id,description, count(*) as count  from boiler_alarm_history where boiler_id=? and date(date_sub(?,interval 7 day)) <= date(start_date) and date(start_date) < date(?) " +
				"group by parameter_id order by count desc limit 1"
			if _,err:=dba.BoilerOrm.Raw(maxsql,bid,t,t).Values(&param);err!=nil{
				goazure.Error("Query boiler_alarm_history Error",err)
			}
			if pid,err:=strconv.Atoi(param[0]["parameter_id"].(string));err!=nil{
				goazure.Error("Parse Int Error",err)
			} else {
				parameterId = pid
			}
			if pCount,err:= strconv.Atoi(param[0]["count"].(string));err!=nil{
				goazure.Error("Parse Int Error",err)
			} else {
				parameterAlarmCount = pCount
			}
			description = param[0]["description"].(string)
		}
		runtimeSql := "select Fire_duration from boiler_terminal_combined btc , issued_boiler_fire_duration ibf where btc.terminal_code = ibf.Boiler_term_id " +
			"and btc.boiler_id=? and btc.terminal_set_id = 1"
		if err:=dba.BoilerOrm.Raw(runtimeSql,bid).QueryRow(&runtime);err!=nil{
			fmt.Println("err:",err)
		}
		var users []*models.User
		raw := "SELECT 	`user`.* "
		raw += "FROM	`boiler`, `user`, `boiler_message_subscriber` AS `sub` "
		raw += "WHERE	`boiler`.`uid` = `sub`.`boiler_id` AND `user`.`uid` = `sub`.`user_id` "
		raw += fmt.Sprintf("AND	`boiler`.`uid` = '%s';", bid)
		if num, err := dba.BoilerOrm.Raw(raw).QueryRows(&users); err != nil {
			fmt.Println("Get Boiler Subscribers Error:", err, num)
		}
		for _,u := range users {
			var info IssuedInformationPush
			info.BoilerUid = bid
			info.BoilerName = bName
			info.UserUid =u.Uid
			info.AlarmCount =alarmCount
			info.ParameterId = parameterId
			info.ParameterAlarmCount = parameterAlarmCount
			info.Description =description
			info.BoilerRuntime =runtime
			var su models.UserThird
			qu := dba.BoilerOrm.QueryTable("user_third")
			qu = qu.Filter("User__Uid", u.Uid).Filter("App", "service").Filter("IsDeleted", false)
			if err := qu.One(&su); err != nil {
				goazure.Error("User", u.Name, "Is NOT Subscribed.")
				continue
			}
			tempMsg, _ := WxCtl.TemplateInformationPush(&info)
			err:=WxCtl.SendTemplateInformation(su.OpenId, tempMsg)
			if err!=nil {
				info.status = false
			} else {
				info.status = true
			}
			infoSql:="insert into issued_information_push_log(uid,boiler_uid,user_uid,create_time,alarm_count,parameter_id,parameter_alarm_count,description,boiler_runtime,status) "+
				"values(uuid(),?,?,now(),?,?,?,?,?,?)"
			if _,err:=dba.BoilerOrm.Raw(infoSql,info.BoilerUid,info.UserUid,info.AlarmCount,info.ParameterId,info.ParameterAlarmCount,info.Description,info.BoilerRuntime,info.status).Exec();err!=nil{
				goazure.Error("Insert issued_information_push_log Error",err)
			}
		}
	}
	durationSql:="update issued_boiler_fire_duration " +
		"set Start_statistic_time = now(), End_statistic_time=now(),Fire_duration=0"
		if _,err:=dba.BoilerOrm.Raw(durationSql).Exec();err!=nil{
			goazure.Error("Update issued_boiler_fire_duration Error",err)
		}
}*/

func (ctl *RuntimeController) GetBoilerRank() {
	ticker := time.NewTicker(time.Minute * 15)
	tick := func() {
		for t := range ticker.C {
			RtmCtl.RefreshBoilerHeatRank(t)
		}
	}

	go tick()
	go RtmCtl.RefreshBoilerHeatRank(time.Now())
}

func (ctl *RuntimeController) GetRunningDuration() {
	ticker := time.NewTicker(time.Minute * 30)
	tick := func() {
		for t := range ticker.C {
			RtmCtl.RefreshStatusRunningDuration(t)
		}
	}

	go tick()
	go RtmCtl.RefreshStatusRunningDuration(time.Now())
}

func (ctl *RuntimeController) RuntimeReload() {
	ticker := time.NewTicker(time.Second * 30)
	tick := func() {
		for t := range ticker.C {
			go runtimeReload(t)
		}
	}

	go tick()
}

func (ctl *RuntimeController) RuntimeDataReload(rtm *models.BoilerRuntime, due float64) {
	defer func() {
		if r := recover(); r != nil {
			goazure.Error("Recovered in RtmCtl.RuntimeDataReload()", r)
			//time.Sleep(time.Second * 10)
		}
	}()

	startTime := time.Now()

	var lgr logs.BoilerRuntimeLog
	lgr.Name = "RuntimeDataReload()"
	lgr.Runtime = rtm
	lgr.TableName = "boiler_runtime"
	lgr.Query = "UPDATED"
	lgr.CreatedDate = startTime
	lgr.DurationTotal = due
	lgr.Status = logs.BOILER_RUNTIME_LOG_STATUS_READY
	LogCtl.AddReloadLog(&lgr)

	v := float64(rtm.Value) * float64(rtm.Parameter.Scale)
	var val interface{}
	pow10_n := math.Pow10(int(rtm.Parameter.Fix))
	val = math.Trunc(v * pow10_n + 0.5) / pow10_n

	if rtm.Parameter.Category.Id == models.RUNTIME_PARAMETER_CATEGORY_RANGE {
		val = rtm.Remark
	}

	/* ALARM */
	ctl.ReloadAlarmWithRuntime(rtm, val)
	/* CACHE */
	ctl.ReloadCacheWithRuntime(rtm, val)
	/* HISTORY */
	//ctl.ReloadHistoryWithRuntime(rtm, val)

	rtm.Status = models.RUNTIME_STATUS_NEEDRELOAD

	if err := DataCtl.UpdateData(rtm); err != nil {
		goazure.Error("Updated Runtime With Alarm Error:", err, "\n", rtm)
	}

	var lgd logs.BoilerRuntimeLog
	lgd.Name = "RuntimeDataReload()"
	lgd.Runtime = rtm
	lgd.TableName = "boiler_runtime"
	lgd.Query = "UPDATED"
	lgd.CreatedDate = time.Now()
	lgd.Duration = float64(lgd.CreatedDate.Sub(startTime)) / float64(time.Second)
	lgd.DurationTotal = lgr.DurationTotal + lgd.Duration
	lgd.Status = logs.BOILER_RUNTIME_LOG_STATUS_READY
	LogCtl.AddReloadLog(&lgr)
}

func (ctl *RuntimeController) ReloadAlarmWithRuntime(rtm *models.BoilerRuntime, val interface{}) (*models.BoilerAlarm, error) {
	var alarm 	models.BoilerAlarm
	var rule 	models.RuntimeAlarmRule
	goazure.Error("处理的运行参数：",rtm.Parameter)
	boiler := BlrCtl.Boiler(rtm.Boiler.Uid)
	goazure.Error(fmt.Sprintf("运行参数的锅炉：%s,锅炉使用企业：,%s",rtm.Boiler.Name,boiler.Enterprise.Name))
	if boiler == nil {
		return nil, errors.New("boiler can not be nil")
	}
	if boiler.Enterprise == nil{
		return nil,errors.New("boiler Enterprise can not be nil")
	}
	qr := dba.BoilerOrm.QueryTable("runtime_alarm_rule").RelatedSel("Organization")
	cond := orm.NewCondition()
	if boiler.Form != nil {
		condFm := orm.NewCondition().Or("BoilerForm__Id", boiler.Form.Id).Or("BoilerForm__Id", 0)
		cond = cond.AndCond(condFm)
	}
	if boiler.Medium != nil {
		condMed := orm.NewCondition().Or("BoilerMedium__Id", boiler.Medium.Id).Or("BoilerMedium__Id", 0)
		cond = cond.AndCond(condMed)
	}
	if boiler.Fuel != nil && boiler.Fuel.Type != nil {
		condFt := orm.NewCondition().Or("BoilerFuelType__Id", boiler.Fuel.Type.Id).Or("BoilerFuelType__Id", 0)
		cond = cond.AndCond(condFt)
	}
	condEvaDummy := orm.NewCondition().And("BoilerCapacityMin", 0).And("BoilerCapacityMax", 0)
	condEvaValid := orm.NewCondition().And("BoilerCapacityMin__lte", boiler.EvaporatingCapacity).And("BoilerCapacityMax__gte", boiler.EvaporatingCapacity)
	condEva := orm.NewCondition().OrCond(condEvaDummy).OrCond(condEvaValid)
	cond = cond.AndCond(condEva)

	qr = qr.SetCond(cond)
	qr = qr.Filter("Parameter__Id", rtm.Parameter.Id).Filter("IsDeleted", false).Filter("Organization",boiler.Enterprise.Uid)
	if err := qr.One(&rule); err != nil {
		goazure.Warning("Get AlarmRule Error:", err, "\n", rtm.Boiler)
		return nil, err
	}
	goazure.Error("查询出来的告警规则:",rule)
	var alarmDesc string
	var alarmLevel int32 = models.RUNTIME_ALARM_LEVEL_UNDEFINED
	if rule.Warning > rule.Normal && val.(float64) > float64(rule.Warning) {
		alarmDesc = rule.Description
		alarmLevel = rule.Priority
	}
	if rule.Warning < rule.Normal && val.(float64) < float64(rule.Warning) {
		alarmDesc = rule.Description
		alarmLevel = rule.Priority
	}

	if alarmLevel <= models.RUNTIME_ALARM_LEVEL_NORMAL {
		return nil, errors.New("no alarm occurs")
	}

	qa := dba.BoilerOrm.QueryTable("boiler_alarm").
		Filter("TriggerRule__Uid", rule.Uid).Filter("Boiler__Uid", rtm.Boiler.Uid).
		Filter("EndDate__lte", rtm.CreatedDate.Add(time.Hour * 1)).Filter("EndDate__gte", rtm.CreatedDate.Add(time.Hour * -1)).
		Filter("IsDeleted", false)
	if err := qa.One(&alarm); err != nil {
		alarm.Uid = uuid.New()
		alarm.Boiler = boiler
		alarm.Parameter = rule.Parameter
		alarm.TriggerRule = &rule
		alarm.Description = alarmDesc
		alarm.AlarmLevel = alarmLevel
		alarm.Priority = rule.Priority
		alarm.State = models.BOILER_ALARM_STATE_NEW
		alarm.NeedSend = rule.NeedSend
		alarm.StartDate = rtm.CreatedDate
		alarm.EndDate = rtm.CreatedDate

		goazure.Warning("Get Exist Alarm Error:", alarm)
	} else {
		if !alarm.EndDate.After(rtm.CreatedDate) {
			alarm.EndDate = rtm.CreatedDate
		}

		if alarm.StartDate.After(rtm.CreatedDate) {
			alarm.StartDate = rtm.CreatedDate
		}
	}

	goazure.Info("Alarm Occurs:", alarm)

	if err := DataCtl.AddData(&alarm, true); err != nil {
		return nil, err
	}

	rtm.Alarm = &alarm

	return &alarm, nil
}

func (ctl *RuntimeController) ReloadCacheWithRuntime(rtm *models.BoilerRuntime, val interface{}) {
	/* CACHES */
	tableNamePrefix := "boiler_runtime_cache_"
	tableNameInst	:= "instant"
	tableNameList 	:= []string{"day", "week", "month"}

	if rtm.Parameter.Category.Id == models.RUNTIME_PARAMETER_CATEGORY_RANGE {
		val = rtm.Value
	}

	alarmId := ""
	alarmLv := int32(0)
	alarmDesc := ""
	if rtm.Alarm != nil {
		alarmId = rtm.Alarm.Uid
		alarmLv = rtm.Alarm.AlarmLevel
		alarmDesc = rtm.Alarm.Description
	}

	rawInst :=
		"INSERT INTO `" + tableNamePrefix + tableNameInst + "` " +
		"( " +
		"`runtime_id` , `boiler_id` , `parameter_id` , `alarm_id` , " +
		"`created_date` , `updated_date` , `name` , `value` , " +
		"`parameter_name` , `unit` , `alarm_level` , `alarm_description`, `remark` " +
		") " +
		"VALUES " +
		"( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"`runtime_id` = ? , " +
		"`alarm_id` = ? , " +
		"`updated_date` = IF(`updated_date` > ?, `updated_date`, ?) , " +
		"`value` = ? , " +
		"`name` = ? , " +
		"`parameter_name` = ? , " +
		"`unit` = ? , " +
		"`alarm_level` = ? , " +
		"`alarm_description` = ? , " +
		"`remark` = ?, " +
		"`is_deleted` = FALSE;"

	if 	res, err := dba.BoilerOrm.Raw(rawInst,
		rtm.Id, rtm.Boiler.Uid, rtm.Parameter.Id, alarmId,
		rtm.CreatedDate, rtm.CreatedDate, rtm.Boiler.Name, val,
		rtm.Parameter.Name, rtm.Parameter.Unit, alarmLv, alarmDesc, rtm.Remark,

		rtm.Id,
		alarmId,
		rtm.CreatedDate, rtm.CreatedDate,
		val,
		rtm.Boiler.Name,
		rtm.Parameter.Name,
		rtm.Parameter.Unit,
		alarmLv,
		alarmDesc,
		rtm.Remark).
		Exec(); err != nil {
		goazure.Error("Insert Instant Error:", err, res)
	}
	/*else {
		row, _ := res.RowsAffected()
		id, _ := res.LastInsertId()
		goazure.Info("Insert Instant Done!", rtm)
		goazure.Info("RowsAffected", row)
		goazure.Info("LastInsertId", id)
		//panic(0)
	}*/

	for _, name := range tableNameList {
		raw :=
			"INSERT IGNORE `" + tableNamePrefix + name + "` " +
			"( " +
			"`runtime_id` , `boiler_id` , `parameter_id` , `alarm_id` , " +
			"`created_date` , `name` , `value` , " +
			"`parameter_name` , `unit` , `alarm_level` , `alarm_description`, `remark` " +
			") " +
			"VALUES " +
			"( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?); "

		if	res, err := dba.BoilerOrm.
			Raw(raw,
			rtm.Id, rtm.Boiler.Uid, rtm.Parameter.Id, alarmId,
			rtm.CreatedDate, rtm.Boiler.Name, val,
			rtm.Parameter.Name, rtm.Parameter.Unit, alarmLv, alarmDesc, rtm.Remark).
			Exec(); err != nil {
			goazure.Error("Insert Runtime Cache Error:", err, res)
		} /*else {
			row, _ := res.RowsAffected()
			id, _ := res.LastInsertId()
			goazure.Info("Insert Runtime Cache Done!", rtm)
			goazure.Info("RowsAffected", row)
			goazure.Info("LastInsertId", id)
			//panic(0)
		} */

	}

	//if !isDefault {
	//	if err := DataCtl.AddData(cache.(models.DataInterface), true); err != nil {
	//		goazure.Error("Added/Updated Cache Failed:", err)
	//	}
	//}
}

func (ctl *RuntimeController) ReloadHistory() {
	interval := time.Minute * 5

	ticker := time.NewTicker(interval)
	tick := func() {
		for t := range ticker.C {
			ctl.ReloadHistoryWithArchived(t.Add(time.Minute * -10), t.Add(time.Minute * 2))
		}
	}

	go tick()
	go ctl.ReloadHistoryWithArchived(time.Now().Add(time.Minute * -10), time.Now().Add(time.Minute * 2))
}

func (ctl *RuntimeController) ReloadHistoryWithArchived(startDate time.Time, endDate time.Time) {
	defer func() {
		if r := recover(); r != nil {
			goazure.Error("Recovered in RtmCtl.ReloadHistoryWithArchived()", r)
			//time.Sleep(time.Second * 10)
		}
	}()

	BlrCtl.WaitGroup.Wait()
	ParamCtrl.WaitGroup.Wait()

	var archives 	[]*caches.BoilerRuntimeHistoryArchived
	//t := time.Now()
	//rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if  num, err := dba.BoilerOrm.QueryTable("boiler_runtime_history_archived").
		Filter("CreatedDate__gte", startDate).Filter("CreatedDate__lte", endDate).
		Limit(-1).
		All(&archives);
		err != nil {
		goazure.Error("History Archives Read Error:", err, num)

		return
	}

	goazure.Error("Get Archives:", archives)

	var histories []*caches.BoilerRuntimeHistory

	for _, arch := range archives {
		var history *caches.BoilerRuntimeHistory
		var isMatched bool = false
		for _, h := range histories {
			if h.Boiler != nil && h.Boiler.Uid == arch.Boiler.Uid && h.CreatedDate.Equal(arch.CreatedDate) {
				history = h
				isMatched = true
				break
			}
		}

		if !isMatched {
			boiler := BlrCtl.Boiler(arch.Boiler.Uid)
			if boiler == nil {
				continue
			}

			history = &caches.BoilerRuntimeHistory{}
			history.Boiler = boiler
			history.Name = boiler.Name
			history.CreatedDate = arch.CreatedDate
			history.Histories = []*caches.History{}
		}

		//goazure.Error("History:", history)
		his := &caches.History{}
		param := ParamCtrl.RuntimeParameter(int(arch.Parameter.Id))
		var val interface{}
		if param == nil {
			continue
		} else if param.Category == nil {
			continue
		}
		if param.Category.Id == models.RUNTIME_PARAMETER_CATEGORY_SWITCH {
			val = float64(arch.Value)
		} else {
			v := float64(arch.Value) * float64(param.Scale)
			pow10_n := math.Pow10(int(param.Fix))
			val = math.Trunc(v * pow10_n + 0.5) / pow10_n
		}

		isDataMatched := false
		for _, d := range history.Histories {
			if 	d.ParameterId == param.Id {
				his = d
				isDataMatched = true
				break
			}
		}

		his.Value = val

		if !isDataMatched {
			his.ParameterId = param.Id
			history.Histories = append(history.Histories, his)
		}

		if !isMatched {
			histories = append(histories, history)
		}

	}

	for _, history := range histories {
		history.Marshal()
		if  num, err := dba.BoilerOrm.InsertOrUpdate(history); err != nil {
			goazure.Error("History Add Failed:", err, num)
		} else {
			goazure.Info("History Added:", history)
		}
	}
}

func runtimeReload(t time.Time) {
	var runtime []*models.BoilerRuntime

	if num, err := dba.BoilerOrm.QueryTable("boiler_runtime").
		RelatedSel("Boiler__Fuel").RelatedSel("Parameter").
		Filter("Status", models.RUNTIME_STATUS_NEW).Limit(1000).All(&runtime); err != nil {
		goazure.Error("Get Runtime Data NEW Error:", err, num)
	}

	for _, rtm := range runtime {
		RtmCtl.RuntimeDataReload(rtm, float64(0))
	}

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
			orgCond := orm.NewCondition().Or("Enterprise__Uid", usr.Organization.Uid).Or("Factory__Uid", usr.Organization.Uid).Or("Maintainer__Uid", usr.Organization.Uid)
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

	boiler := BlrCtl.Boiler(b.Uid)

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

		switch b.Range {
		case RUNTIME_RANGE_DEFAULT:
			fallthrough
		case RUNTIME_RANGE_TODAY:
			tableName = "boiler_runtime_cache_day"
		case RUNTIME_RANGE_THERE_DAY:
			fallthrough
		case RUNTIME_RANGE_WEEK:
			tableName = "boiler_runtime_cache_week"
		}

		qr := dba.BoilerOrm.QueryTable(tableName)
		if 	tableName == "boiler_runtime" {
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
		fmt.Println("boilerForm:",boiler.Form)
		q := dba.BoilerOrm.QueryTable("runtime_alarm_rule")
		q = q.RelatedSel("Parameter__Category").RelatedSel("BoilerForm").RelatedSel("BoilerMedium").RelatedSel("BoilerFuelType")
		orCond := orm.NewCondition().Or("BoilerForm__Id", 0).
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

	if len(b.RuntimeQueue) == 0 {
		b.RuntimeQueue = ParamCtrl.ParamQueueWithBoiler(boiler)
	}

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
	boilers, err := BlrCtl.CurrentBoilerList(usr)
	if err != nil {
		goazure.Error("get current boiler list error!")
	}

	for _, b := range boilers {
		duration[b.Uid] = RtmCtl.StatusRunningDuration[b.Uid]
	}

	ctl.Data["json"] = duration
	ctl.ServeJSON()
}

func (ctl *RuntimeController) BoilerRuntimeHistory() {
	//goazure.Info("=== Ready to get Runtime History ===")
	usr := ctl.GetCurrentUser()

	b := bBoiler{}
	resStatus := 200
	resBody := "Success"
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	var isMatched bool = false
	var boiler models.Boiler
	boilers, _ := BlrCtl.CurrentBoilerList(usr)
	for _, cb := range boilers {
		if cb.Uid == b.Uid {
			boiler = *cb
			isMatched = true
		}
	}
	if !isMatched {
		e := fmt.Sprintln("Permission Error")
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(e))
		goazure.Error(e)

		return
	}

	if len(b.RuntimeQueue) == 0 {
		b.RuntimeQueue = ParamCtrl.ParamQueueWithBoiler(&boiler)
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
	var boiler *models.Boiler
	var scope string
	//goazure.Error("Instants Param:", ctl.Input())
	if ctl.Input()["scope"] != nil && len(ctl.Input()["scope"]) > 0 {
		scope = ctl.Input()["scope"][0]
	}

	var b bBoiler
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	boiler = BlrCtl.Boiler(b.Uid)
	var combined models.BoilerTerminalCombined
	if err:=dba.BoilerOrm.QueryTable("boiler_terminal_combined").RelatedSel("Terminal").Filter("Boiler__Uid",b.Uid).Filter("TerminalSetId",1).One(&combined);err!=nil{
		goazure.Error("Query boiler_terminal_combined Error",err)
		return
	}
	fmt.Println("terminal:",combined.TerminalCode)
	if len(b.RuntimeQueue) <= 0 {
		b.RuntimeQueue = ParamCtrl.ParamQueueWithBoiler(boiler)
		switch scope {
		case "thumb":
			if len(b.RuntimeQueue) >= 4 {
				b.RuntimeQueue = b.RuntimeQueue[:4]
			}

			for i, v := range b.RuntimeQueue {
				ids := strconv.FormatInt(int64(v), 10)
				if strings.HasPrefix(ids, strconv.FormatInt(int64(models.RUNTIME_PARAMETER_CATEGORY_RANGE), 10)) {
					b.RuntimeQueue = b.RuntimeQueue[:i]
					break
				}
				if strings.HasPrefix(ids, strconv.FormatInt(int64(models.RUNTIME_PARAMETER_CATEGORY_SWITCH), 10)) {
					b.RuntimeQueue = b.RuntimeQueue[:i]
					break
				}
			}

		/*	if len(b.RuntimeQueue) <= 0 {
				b.RuntimeQueue = []int{1201, 1015, 1002, 1202}
			}*/

			//goazure.Error("RuntimeQueue:", b.RuntimeQueue)
		case "wxapp":
			if len(b.RuntimeQueue) <= 0 {
				b.RuntimeQueue = []int{1201, 1002, 1015, 1202, 1016, 1001, 1004, 1003}
			}
		default:
			b.RuntimeQueue = ParamCtrl.ParamQueueWithBoiler(boiler)
		}
	}

	runtimes := func() []orm.Params {
		var ins []orm.Params
		q := dba.BoilerOrm.QueryTable("boiler_runtime_cache_instant")
		q = q.RelatedSel("Parameter")
		q = q.Filter("Boiler__Uid", b.Uid)
		if len(b.RuntimeQueue) > 0 {
			q = q.Filter("Parameter__Id__in", b.RuntimeQueue)
		}
		//q = q.Filter("UpdatedDate__gt", time.Now().Add(time.Hour * -6))
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
		param := ParamCtrl.RuntimeParameter(paramId)
		if param == nil {
			continue
		}
		for _, r := range rs {
			if  r["Parameter"] == int64(paramId) {
				r["ParameterCategory"] = param.Category.Id
				var channel models.RuntimeParameterChannelConfig
				if err:=dba.BoilerOrm.QueryTable("runtime_parameter_channel_config").Filter("Terminal__Uid",combined.Terminal.Uid).Filter("Parameter__Id",param.Id).One(&channel);err!=nil{
					goazure.Error("Query channel_config by param and terminal Error")
				}
				r["SwitchStatus"] = channel.SwitchStatus
				rtms = append(rtms, r)
				num++
				break
			}
		}

		if num == 0 {
			in := orm.Params{}
			in["Parameter"] = param.Id
			in["ParameterName"] = param.Name
			in["ParameterCategory"] = param.Category.Id

			rtms = append(rtms, in)
		}
	}

	if len(rtms) > 0 {
		r := rtms[0]
		var date string
		if r["UpdatedDate"] != nil {
			date = r["UpdatedDate"].(time.Time).Format("2006-01-02 15:04:05")
		}

		r["Date"] = date
	}

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
	//limit := 7
	//type DailyTotal struct {
	//	Boilers		[]*models.Boiler
	//	Parameters	[]*models.RuntimeParameter
	//	Flows		[][]orm.Params
	//	Heats		[][]orm.Params
	//}

	//dt := DailyTotal{}
	//params := parameters(1003, 1201)
	//boilers, _ := ctl.boilers()

	//dt.Boilers = boilers
	//dt.Parameters = params

	//TODO: Index Of Date Not Boiler
	//var aFlows [][]orm.Params
	//var aHeats [][]orm.Params

	raw := 	"SELECT 	`flows`.`date`, AVG(`flows`.`value`) AS `flow`, AVG(`heats`.`value`) AS `heat` " +
			"FROM		`boiler_runtime_cache_flow_daily` AS `flows`, `boiler_runtime_cache_heat_daily` AS `heats` " +
			"WHERE		`flows`.`date` = `heats`.`date` AND `flows`.`boiler_id` = `heats`.`boiler_id` " +
			"GROUP BY 	`heats`.`date`, `heats`.`boiler_id` " +
			"ORDER BY 	`heats`.`date`;"

	type total struct {
		Date		time.Time
		Flow		float64
		Heat		float64
	}

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

	for _, ht := range RtmCtl.BoilerHeats {
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
	
	ranks := RtmCtl.BoilerRanks

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
		//goazure.Info("Read SQL:", sqlCoal)
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

	RtmCtl.BoilerRanks = ranks
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

	RtmCtl.BoilerHeats = heats
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

	if boilers, err = BlrCtl.CurrentBoilerList(usr); err != nil {
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

//******************************* Init

func (ctl *RuntimeController) InitRuntimeParameters() {
	generateDefaultRuntimeParameterCategories()
	//generateDefaultRuntimeParameterMediums()
	generateDefaultRuntimeParameters()
}

func generateDefaultRuntimeParameterCategories() {
	var cate10, cate11, cate12, cate13 models.RuntimeParameterCategory
	cate10.Id = 10
	cate10.Name = "模拟量采集参数"
	//cate10.NameEn = "Analogue Collection Parameters"

	cate11.Id = 11
	cate11.Name = "开关量采集参数"
	//cate11.NameEn = "Switch Signal Collection Parameters"

	cate12.Id = 12
	cate12.Name = "计算参数"
	//cate12.NameEn = "Calculation Parameters"

	cate13.Id = 13
	cate13.Name = "状态量采集参数"
	//cate13.NameEn = "Status Parameters"

	addRuntimeParameterCategory(cate10)
	addRuntimeParameterCategory(cate11)
	addRuntimeParameterCategory(cate12)
	addRuntimeParameterCategory(cate13)
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
				//nameEn := row[2]
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
				//param.NameEn = nameEn
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
				param := ParamCtrl.RuntimeParameter(1107)
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

/*
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
*/

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
				if !BlrCtl.IsGenerateData(b.Uid) {
					continue
				}

				value := int64((i + t.Minute()) % 3)
				//goazure.Error("go generateBoilerStatus:", stat, t, b.Name, value)
				generateBoilerStatus(b, t, value, stat)
			}
		}

		go statGen(time.Now())
		for t := range statTicker.C {
			go statGen(t)
		}
	}

	BlrCtl.RefreshGlobalBoilerList()
	for _, st := range status {
		go generateStatus(time.Now(), st)
	}

	ticker := time.NewTicker(time.Minute * 30)
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
			if !BlrCtl.IsGenerateData(b.Uid) {
				continue
			}

			if !BlrCtl.IsBurning(b.Uid) {
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
			BlrCtl.RefreshGlobalBoilerList()
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
	param := ParamCtrl.RuntimeParameter(int(run.Pid))

	if param.Category.Id == 11 && value > 0 {
		value = 1
	}

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

	go RtmCtl.RuntimeDataReload(&rtm, float64(0))

	return value
}

func generateBoilerStatus(boiler *models.Boiler, date time.Time, value int64, run Runtime) {
	if value < 0 {
		value = int64(float32(run.BaseValue) + rand.Float32() * float32(run.DynamicValue))
	}

	param := ParamCtrl.RuntimeParameter(int(run.Pid))

	if param == nil {
		goazure.Error("Get Param", run.Pid, "Error!")

		return
	}

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

func (ctl *RuntimeController)ImportExistCache() {
	var rtms []*models.BoilerRuntime

	if	num, err := dba.BoilerOrm.QueryTable("boiler_runtime").
		RelatedSel("Boiler").RelatedSel("Parameter").RelatedSel("Alarm").
		Limit(-1).
		All(&rtms); err != nil {
		goazure.Error("Read Boiler Runtime Error:", err, num)
	}

	/* CACHES */
	tableNamePrefix := "boiler_runtime_cache_"
	tableNameList 	:= []string{"day", "week", "month"}

	for _, rtm := range rtms {
		v := float64(rtm.Value) * float64(rtm.Parameter.Scale)
		var val interface{}
		pow10_n := math.Pow10(int(rtm.Parameter.Fix))
		val = math.Trunc(v * pow10_n + 0.5) / pow10_n

		alarmId := ""
		alarmLv := int32(0)
		alarmDesc := ""
		if 	rtm.Alarm != nil {
			alarmId = rtm.Alarm.Uid
			alarmLv = rtm.Alarm.AlarmLevel
			alarmDesc = rtm.Alarm.Description
		}

		for _, name := range tableNameList {
			raw :=
				"INSERT IGNORE `" + tableNamePrefix + name + "` " +
					"( " +
					"`runtime_id` , `boiler_id` , `parameter_id` , `alarm_id` , " +
					"`created_date` , `name` , `value` , " +
					"`parameter_name` , `unit` , `alarm_level` , `alarm_description`, `remark` " +
					") " +
					"VALUES " +
					"( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?); "

			if	res, err := dba.BoilerOrm.
				Raw(raw,
				rtm.Id, rtm.Boiler.Uid, rtm.Parameter.Id, alarmId,
				rtm.CreatedDate, rtm.Boiler.Name, val,
				rtm.Parameter.Name, rtm.Parameter.Unit, alarmLv, alarmDesc, rtm.Remark).
				Exec(); err != nil {
				goazure.Error("Insert Runtime Cache Error:", err, res)
			} /*else {
			row, _ := res.RowsAffected()
			id, _ := res.LastInsertId()
			goazure.Info("Insert Runtime Cache Done!", rtm)
			goazure.Info("RowsAffected", row)
			goazure.Info("LastInsertId", id)
			//panic(0)
		} */

		}
	}
}