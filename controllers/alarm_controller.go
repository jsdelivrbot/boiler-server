package controllers

import (
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"os"
	"encoding/csv"
	"path/filepath"
	"log"
	"time"
	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"
	"encoding/json"
	"github.com/AzureRelease/boiler-server/conf"
)

var AlmCtl *AlarmController = &AlarmController{}

type AlarmController struct {
	MainController
}

const ALARM_INTERVAL = time.Hour * 4

func (ctl *AlarmController) InitAlarmSendService() {
	interval := time.Minute * 5
	if !conf.IsRelease {
		interval = time.Second * 5
	}

	ticker := time.NewTicker(interval)
	tick := func() {
		for t := range ticker.C {
			AlmCtl.SendAlarmMessage(t)
		}
	}

	go tick()
}

func (ctl *AlarmController) AlarmRuleList() {
	usr := ctl.GetCurrentUser()

	var alarmRules []models.RuntimeAlarmRule
	qs := dba.BoilerOrm.QueryTable("runtime_alarm_rule")
	qs = qs.RelatedSel("Parameter__Category").RelatedSel("BoilerForm").RelatedSel("BoilerMedium").RelatedSel("BoilerFuelType")
	if usr.IsCommonUser() ||
		usr.Status == models.USER_STATUS_INACTIVE || usr.Status == models.USER_STATUS_NEW {
		qs = qs.Filter("IsDemo", true)
	} else if usr.IsOrganizationUser() {
		//qs = qs.Filter("Scope", models.RUNTIME_ALARM_SCOPE_ENTERPRISE)
	}
	num, err := qs.Filter("IsDeleted", false).OrderBy("Parameter").All(&alarmRules)
	fmt.Printf("Returned Rows Num: %d, %s", num, err)

	//se := ctl.GetSession(SESSION_CURRENT_USER)
	//fmt.Println("\nBoiler Get CurrentUser: ", usr, " | ", se);

	ctl.Data["json"] = alarmRules
	ctl.ServeJSON()
}

func (ctl *AlarmController) BoilerAlarmList() {
	usr := ctl.GetCurrentUser()
	var boilerUid string
	goazure.Warn("Alarm GET URL:", ctl.Input())
	if ctl.Input()["boiler"] != nil &&
		len(ctl.Input()["boiler"]) > 0 &&
		ctl.Input()["boiler"][0] != "undefined" {
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

		if num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").All(&boilers); num == 0 || err != nil {
			goazure.Error("Read BoilerList For Alarm Error: ", err, num)
		}
	}

	if !usr.IsAdmin() && len(boilers) <= 0 {
		e := fmt.Sprintln("No Admin && No Boiler Accessed")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		goazure.Error(e)
		return
	}

	var objs []orm.Params
	qa := dba.BoilerOrm.QueryTable("boiler_alarm")
	qa = qa.RelatedSel("Boiler__Enterprise").RelatedSel("Parameter")
	qa = qa.Filter("EndDate__gt", time.Now().Add(ALARM_INTERVAL * -1))
	if !usr.IsAdmin() || len(boilers) > 0 {
		qa = qa.Filter("Boiler__in", boilers)
	}
	if num, err := qa.Filter("IsDeleted", false).
		OrderBy("-EndDate").
		Values(&objs, "Uid", "StartDate", "EndDate", "Priority", "State", "Description",
		"Boiler__Uid", "Boiler__Name", "Boiler__Enterprise__Name",
		"Parameter__Name"); num == 0 || err != nil {
		goazure.Error("Read AlarmList Error:", num, err)
	}

	for _, a := range objs {
		//dba.BoilerOrm.LoadRelated(&a, "Runtime")
		//goazure.Warn("RTM:", a)
		//startDateStr := a.StartDate.Format("2006-01-02 15:04")
		//endDateStr := a.EndDate.Format("2006-01-02 15:04")

		dateText := func(date time.Time) string {
			var str string
			nYear, nMonth, nDay := time.Now().Date()
			dYear, dMonth, dDay := date.Date()

			if nYear == dYear {
				if nMonth == dMonth && nDay == dDay {
					str = date.Format("15:04")
				} else {
					str = date.Format("01/02 15:04")
				}
			} else {
				str = date.Format("2006/01/02 15:04")
			}

			return str
		}

		start := a["StartDate"].(time.Time)
		end := a["EndDate"].(time.Time)
		duration := end.Sub(start)

		dueText := ""

		days := duration / (time.Hour * 24)
		duration -= days * time.Hour * 24
		hours := duration / time.Hour
		duration -= hours * time.Hour
		minutes := duration / time.Minute

		if days > 0 {
			dueText += fmt.Sprintf("%d天", days)
		}
		if hours > 0 {
			dueText += fmt.Sprintf("%d时", hours)
		}
		dueText += fmt.Sprintf("%d分", minutes)

		startText := dateText(start)
		endText := dateText(end)

		a["Duration"] = duration
		a["DueText"] = dueText
		a["StartText"] = startText
		a["EndText"] = endText
	}

	ctl.Data["json"] = objs
	ctl.ServeJSON()
}

func (ctl *AlarmController) BoilerAlarmHistoryList() {
	usr := ctl.GetCurrentUser()
	var boilerUid string
	goazure.Warn("Alarm GET URL:", ctl.Input())
	if ctl.Input()["boiler"] != nil &&
		len(ctl.Input()["boiler"]) > 0 &&
		ctl.Input()["boiler"][0] != "undefined" {
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

		if num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").All(&boilers); num == 0 || err != nil {
			goazure.Error("Read BoilerList For Alarm Error: ", err, num)
		}
	}

	//goazure.Warn("Boilers:", boilers)
	if !usr.IsAdmin() && len(boilers) <= 0 {
		e := fmt.Sprintln("No Admin && No Boiler Accessed")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		goazure.Error(e)
		return
	}

	var objs []orm.Params
	qa := dba.BoilerOrm.QueryTable("boiler_alarm_history")
	qa = qa.RelatedSel("Boiler__Enterprise").RelatedSel("Parameter")
	if !usr.IsAdmin() || len(boilers) > 0 {
		qa = qa.Filter("Boiler__in", boilers)
	}
	if num, err := qa.Filter("IsDeleted", false).OrderBy("-EndDate", "-IsDemo").Values(&objs, "Uid", "StartDate", "EndDate", "Priority", "Description", "IsDemo",
		"Boiler__Uid", "Boiler__Name", "Boiler__Enterprise__Name",
		"Parameter__Name"); num == 0 || err != nil {
		goazure.Error("Read Alarm History List Error:", num, err)
	}

	for _, a := range objs {
		//dba.BoilerOrm.LoadRelated(&a, "Runtime")
		//goazure.Warn("RTM:", a)
		//startDateStr := a.StartDate.Format("2006-01-02 15:04")
		//endDateStr := a.EndDate.Format("2006-01-02 15:04")

		dateText := func(date time.Time) string {
			var str string
			nYear, nMonth, nDay := time.Now().Date()
			dYear, dMonth, dDay := date.Date()

			if nYear == dYear {
				if nMonth == dMonth && nDay == dDay {
					str = date.Format("15:04")
				} else {
					str = date.Format("01/02 15:04")
				}
			} else {
				str = date.Format("2006/01/02 15:04")
			}

			return str
		}

		start := a["StartDate"].(time.Time)
		end := a["EndDate"].(time.Time)
		duration := end.Sub(start)

		dueText := ""

		days := duration / (time.Hour * 24)
		duration -= days * time.Hour * 24
		hours := duration / time.Hour
		duration -= hours * time.Hour
		minutes := duration / time.Minute

		if days > 0 {
			dueText += fmt.Sprintf("%d天", days)
		}
		if hours > 0 {
			dueText += fmt.Sprintf("%d时", hours)
		}
		dueText += fmt.Sprintf("%d分", minutes)

		startText := dateText(start)
		endText := dateText(end)

		a["Duration"] = duration
		a["DueText"] = dueText
		a["StartText"] = startText
		a["EndText"] = endText
	}

	ctl.Data["json"] = objs
	ctl.ServeJSON()
}

func (ctl *AlarmController) BoilerAlarmCount() {
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

		if num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").All(&boilers); num == 0 || err != nil {
			fmt.Printf("Read BoilerList Error: %d, %s", num, err)
		}
	}

	goazure.Warn("Boilers:", boilers)

	qa := dba.BoilerOrm.QueryTable("boiler_alarm")
	qa = qa.Filter("EndDate__gt", time.Now().Add(ALARM_INTERVAL * -1))
	if len(boilers) > 0 {
		qa = qa.Filter("Boiler__in", boilers)
	}

	num, err := qa.Filter("IsDeleted", false).Count()
	if num == 0 || err != nil {
		fmt.Printf("Read AlarmCount Error: %d, %v", num, err)
	}

	ctl.Data["json"] = num
	ctl.ServeJSON()
}

func (ctl *AlarmController) BoilerAlarmDetail() {
	type bAlarm struct {
		Uid		string		`json:"uid"`
	}

	var al bAlarm
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &al); err != nil {
		e := fmt.Sprintln("Unmarshal AlarmUid Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		goazure.Error(e)
		return
	}

	var alarms []orm.Params
	qa := dba.BoilerOrm.QueryTable("boiler_alarm")
	qa = qa.RelatedSel("Boiler__Enterprise").RelatedSel("Parameter").RelatedSel("TriggerRule")
	qa = qa.Filter("Uid", al.Uid)
	if num, err := qa.Filter("IsDeleted", false).Values(&alarms,
		"Uid", "StartDate", "EndDate", "Priority", "State", "Description",
		"Boiler__Uid", "Boiler__Name", "Boiler__Enterprise__Name",
		"Parameter__Id", "Parameter__Name", "Parameter__Unit", "Parameter__Scale", "Parameter__Fix",
		"TriggerRule__Normal", "TriggerRule__Warning");
		err != nil || num == 0 {
		goazure.Error("Read Alarm Error:", err, num)
	}

	alarm := alarms[0]
	var runtime []orm.Params
	qs := dba.BoilerOrm.QueryTable("boiler_runtime")
	qs = qs.Filter("Boiler__Uid", alarm["Boiler__Uid"]).Filter("Parameter__Id", alarm["Parameter__Id"]).
		Filter("CreatedDate__gte", alarm["StartDate"]).Filter("CreatedDate__lte", alarm["EndDate"])
	if num, err := qs.Filter("IsDeleted", false).OrderBy("CreatedDate").Values(&runtime, "CreatedDate", "Value"); err != nil || num == 0 {
		goazure.Error("Read AlarmRuntimeList Error:", err, num)
	}

	alarm["runtime"] = runtime

	ctl.Data["json"] = alarm
	ctl.ServeJSON()
}

func (ctl *AlarmController) BoilerAlarmFeedbackList() {
	type bAlarm struct {
		Uid		string		`json:"uid"`
	}

	var al bAlarm
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &al); err != nil {
		e := fmt.Sprintln("Unmarshal AlarmUid Error", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		goazure.Error(e)
		return
	}

	var feedbacks []*models.BoilerAlarmFeedback
	qf := dba.BoilerOrm.QueryTable("boiler_alarm_feedback").RelatedSel("CreatedBy__Organization").RelatedSel("CreatedBy__Role")
	qf = qf.Filter("Alarm__Uid", al.Uid)
	if num, err := qf.Filter("IsDeleted", false).All(&feedbacks); err != nil {
		goazure.Error("Read Alarm Error:", err, num)
	}

	ctl.Data["json"] = feedbacks
	ctl.ServeJSON()
}

type bAlarmRule struct {
	Uid					string		`json:"uid"`
	ParamId 			int64		`json:"paramId"`
	BoilerFormId 		int64		`json:"boilerFormId"`
	BoilerMediumId		int64		`json:"boilerMediumId"`
	BoilerFuelTypeId 	int64 		`json:"boilerFuelTypeId"`
	BoilerCapacityMin 	int 		`json:"boilerCapacityMin"`
	BoilerCapacityMax 	int 		`json:"boilerCapacityMax"`

	NormalValue			float32		`json:"normalValue"`
	WarningValue		float32		`json:"warningValue"`
	Delay				int64		`json:"delay"`
	Priority			int			`json:"priority"`
	Description			string		`json:"description"`
}

func (ctl *AlarmController) AlarmRuleUpdate() {
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	var al bAlarmRule

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &al); err != nil {
		e := fmt.Sprintln("Unmarshal AlarmRule Update Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	rule := models.RuntimeAlarmRule{}
	if len(al.Uid) > 0 {
		rule.Uid = al.Uid

		if err := DataCtl.ReadData(&rule); err != nil {
			fmt.Println("Read Exist Alarm Rule Error:", err)
		}
	}

	param := models.RuntimeParameter{}
	if rule.Parameter == nil {
		if al.ParamId <= 0 {
			//panic("ParamId on AlarmRule Error!")
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("ParamId on AlarmRule Error!"))
			return
		}

		param.Id = al.ParamId
		if err := DataCtl.ReadData(&param); err != nil {
			fmt.Println("Read Exist AlarmRule Parameter Error:", err)
			//panic(err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("Read Exist AlarmRule Parameter Error:"))
			return
		}

		rule.Parameter = &param
	}

	form := models.BoilerTypeForm{}
	form.Id = al.BoilerFormId
	if err := DataCtl.ReadData(&form); err != nil {
		fmt.Println("Read AlarmRule BoilerForm Error:", err)
	} else {
		rule.BoilerForm = &form
	}

	medium := models.BoilerMedium{}
	medium.Id = al.BoilerMediumId
	if err := DataCtl.ReadData(&medium); err != nil {
		fmt.Println("Read AlarmRule BoilerMedium Error:", err)
	} else {
		rule.BoilerMedium = &medium
	}

	fuelType := models.FuelType{}
	fuelType.Id = al.BoilerFuelTypeId
	if err := DataCtl.ReadData(&fuelType); err != nil {
		fmt.Println("Read AlarmRule BoilerForm Error:", err)
	} else {
		rule.BoilerFuelType = &fuelType
	}

	//capacityMin, _ := strconv.ParseInt(al.BoilerCapacityMin, 10, 64)
	//capacityMax, _ := strconv.ParseInt(al.BoilerCapacityMax, 10, 64)
	capacityMin := al.BoilerCapacityMin
	capacityMax := al.BoilerCapacityMax

	if capacityMin < 0 {
		capacityMin = 0
	}
	if capacityMax < 0 {
		capacityMax = 0
	}

	if capacityMax <= capacityMin {
		capacityMin = capacityMax
		capacityMax = 0
	}
	rule.BoilerCapacityMin = int32(capacityMin)
	rule.BoilerCapacityMax = int32(capacityMax)

	rule.Normal = float32(al.NormalValue)
	rule.Warning = float32(al.WarningValue)
	//TODO: AlarmRule Danger

	rule.Delay = al.Delay
	rule.Priority = int32(al.Priority)

	rule.Description = al.Description

	if err := DataCtl.AddData(&rule, true); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update AlarmRule Error!"))
		return
	}

	goazure.Info("\nUpdate AlarmRule:", rule, al)
}

type Alarm struct {
	Uid     string    `json:"uid"`
}
func (ctl *AlarmController) AlarmRuleDelete() {
	usr := ctl.GetCurrentUser()
	var a Alarm
	var alarmRule models.RuntimeAlarmRule

	if !usr.IsAdmin() {
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte("Permission Denied!"))
		goazure.Error("Permission Denied!")
		return
	}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &a); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}
	if err := dba.BoilerOrm.QueryTable("runtime_alarm_rule").RelatedSel("Parameter__Category").RelatedSel("BoilerForm").RelatedSel("BoilerMedium").RelatedSel("BoilerFuelType").Filter("Uid", a.Uid).One(&alarmRule); err != nil {
		e := fmt.Sprintf("Read runtime_alarm_rule for Delete Error: %v", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	if err := DataCtl.DeleteData(&alarmRule); err != nil {
		e := fmt.Sprintln("Delete runtime_alarm_rule Error!", alarmRule, err)
		goazure.Error(e)
		ctl.Ctx.Output.Body([]byte(e))
	}
}

type bAlarmFeedback struct {
	Uid			string		`json:"uid"`
	AlarmId 		string		`json:"alarm_id"`
	State	 		int		`json:"state"`
	Topic			string		`json:"topic"`
	Content		 	string 		`json:"content"`
}

func (ctl *AlarmController) BoilerAlarmUpdate() {
	usr := ctl.GetCurrentUser()
	//if !usr.IsAdmin() {
	//	e := fmt.Sprintln("Permission Denied!")
	//	goazure.Error(e)
	//	ctl.Ctx.Output.SetStatus(403)
	//	ctl.Ctx.Output.Body([]byte(e))
	//	return
	//}
	var cmt bAlarmFeedback

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cmt); err != nil {
		e := fmt.Sprintln("Unmarshal AlarmFeedback Update Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var fb models.BoilerAlarmFeedback
	var alarm models.BoilerAlarm
	var fbExist bool = false
	if len(cmt.Uid) > 0 {
		fb.Uid = cmt.Uid

		if err := DataCtl.ReadData(&fb); err != nil {
			goazure.Warn("Read Exist AlarmFeedback Error:", err)
		} else {
			fbExist = true
		}
	}

	alarm.Uid = cmt.AlarmId
	if err := DataCtl.ReadData(&alarm); err != nil {
		e := fmt.Sprintln("Read Alarm Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if !fbExist {
		fb.Alarm = &alarm
	}

	if len(cmt.Topic) <= 0 {
		fb.Name = cmt.Topic
	}
	fb.Content = cmt.Content
	fb.State = int32(cmt.State)
	fb.CreatedBy = usr
	fb.UpdatedBy = usr

	switch fb.State {
	case models.ALARM_FEEDBACK_TYPE_CONFIRM:
		alarm.State = models.BOILER_ALARM_STATE_CONFIRMED
	case models.ALARM_FEEDBACK_TYPE_REJECT:
		alarm.State = models.BOILER_ALARM_STATE_REJECTED
	case models.ALARM_FEEDBACK_TYPE_VERIFIED:
		alarm.State = models.ALARM_FEEDBACK_TYPE_VERIFIED
	default:

	}

	if err := DataCtl.AddData(&fb, true); err != nil {
		e := fmt.Sprintln("Update AlarmFeedback Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if err := DataCtl.UpdateData(&alarm); err != nil {
		e := fmt.Sprintln("Update Alarm State Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	//goazure.Info("Update AlarmFeedback:", alarm, fb)
}

func (ctl *AlarmController) SendAlarmMessage(t time.Time) {
	var alarms []*models.BoilerAlarm
	if 	num, err := dba.BoilerOrm.QueryTable("boiler_alarm").
		RelatedSel("Boiler__Address").
		Filter("State", models.BOILER_ALARM_STATE_NEW).Filter("StartDate__lte", t.Add(time.Minute * -10)).Filter("IsDeleted", false).
		All(&alarms); err != nil {
		goazure.Error("Fetch New Alarm Error:", num, err)
	}

	for i, al := range alarms {
		if al.Priority > 0 {
			var users []*models.User
			raw := 	"SELECT `user`.* " +
					"FROM	`boiler`, `user`, `boiler_message_subscriber` AS `sub` " +
					"WHERE	`boiler`.`uid` = `sub`.`boiler_id` AND `user`.`uid` = `sub`.`user_id` " +
					fmt.Sprintf("AND	`boiler`.`uid` = '%s';", al.Boiler.Uid)
			if 	num, err := dba.BoilerOrm.Raw(raw).QueryRows(&users); err != nil {
				goazure.Error("Get Boiler Subscribers Error:", err, num)
			}
			al.Boiler.Subscribers = users

			for _, u := range al.Boiler.Subscribers {
				var tds []*models.UserThird
				if 	num, err := dba.BoilerOrm.QueryTable("user_third").
					Filter("User__Uid", u.Uid).Filter("App", "service").Filter("IsDeleted", false).
					All(&tds); err != nil {
					goazure.Error("User", u.Name, "Is NOT Subscribed.", num)
					continue
				}

				tempMsg, _ := WxCtl.TemplateMessageAlarm(al)
				for _, su := range tds {
					WxCtl.SendTemplateMessage(su.OpenId, tempMsg)
				}
			}
		}

		al.State = models.BOILER_ALARM_STATE_PENDING
		if err := DataCtl.UpdateData(al); err != nil {
			goazure.Error("Update Alarm State Error:", err, al, i)
		}
	}
}

func generateDefaultAlarmRules() (error) {
	const fieldRowNo int = 0
	const paramIdCol int = 0
	const boilerFormCol int = 1
	const boilerMediumCol int = 2
	const boilerFuelTypeCol int = 3

	var ac models.RuntimeAlarmRule
	aType := reflect.TypeOf(ac)

	err := filepath.Walk(rtmDefaultsPath, func(aPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(aPath, "alarm_config") {
			return err
		}

		log.Println("FilePath", aPath)

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

					switch j {
					case paramIdCol:
						pid, _ := strconv.ParseInt(field, 10, 64)
						value = runtimeParameter(int(pid))
					case boilerFormCol:
						id, _ := strconv.ParseInt(field, 10, 64)
						value = boilerForm(int(id))
					case boilerMediumCol:
						id, _ := strconv.ParseInt(field, 10, 64)
						value = boilerMedium(int(id))
					case boilerFuelTypeCol:
						id, _ := strconv.ParseInt(field, 10, 64)
						value = boilerFuelType(int(id))
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

						default:
							value = field
						}
					}

					fmt.Println("Field(", fieldNames[j],":", da.FieldByName(fieldNames[j]).Kind(),"): ", field, value)
					da.FieldByName(fieldNames[j]).Set(reflect.ValueOf(value))
				}

				//fmt.Println("Da: ", da)
				//fmt.Println("In:", in)
				//idx := i - fieldRowNo - 1
				DataCtl.AddData(in, true,
					"Parameter",
					"BoilerForm",
					"BoilerMedium",
					"BoilerFuelType",
					"BoilerCapacityMin", "BoilerCapacityMax",
					"Scope")
			}
		}

		return err
	})

	return err
}


func (ctl *AlarmController) InitAlarmRules() {
	generateDefaultAlarmRules()
}
