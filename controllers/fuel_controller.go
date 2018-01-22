package controllers

import (
	"github.com/AzureRelease/boiler-server/dba"
	"fmt"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"
	"encoding/json"
	"time"
)

type FuelController struct {
	MainController
}

var fuelCtl *FuelController = &FuelController{}

func (ctl *FuelController) FuelList() {
	var fuels []models.Fuel
	qs := dba.BoilerOrm.QueryTable("fuel")
	num, err := qs.RelatedSel().OrderBy("Type__Id", "FuelId").All(&fuels)
	fmt.Printf("Returned Rows Num: %d, %s", num, err)

	ctl.Data["json"] = fuels
	ctl.ServeJSON()
}

func (ctl *FuelController) FuelTypeList() {
	var types []models.FuelType
	qs := dba.BoilerOrm.QueryTable("fuel_type")
	num, err := qs.OrderBy("Id").All(&types)
	fmt.Printf("Returned Rows Num: %d, %s", num, err)

	ctl.Data["json"] = types
	ctl.ServeJSON()
}

func (ctl *FuelController) FuelListWeixin() {
	var types []orm.Params
	qt := dba.BoilerOrm.QueryTable("fuel_type")
	if num, err := qt.OrderBy("Id").Values(&types, "Id", "Name", "NameEn"); err != nil || num == 0 {
		goazure.Error("Read Fuel-Type Error:", err, num)
	}

	for _, t := range types {
		var fuels []orm.Params
		qf := dba.BoilerOrm.QueryTable("fuel")
		qf = qf.Filter("Type__Id", t["Id"])
		if num, err := qf.OrderBy("FuelId").Values(&fuels, "Uid", "FuelId", "Name", "NameEn"); err != nil || num == 0 {
			goazure.Error("Read Fuel Error:", err, num)
		}
		t["Fuels"] = fuels
	}

	ctl.Data["json"] = types
	ctl.ServeJSON()
}

func (ctl *FuelController) FuelRecordList() {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		var token string
		if ctl.Input()["token"] != nil && len(ctl.Input()["token"]) > 0 {
			token = ctl.Input()["token"][0]
		}

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
			return
		}
	}

	var records []*models.BoilerFuelRecord
	qs := dba.BoilerOrm.QueryTable("boiler_fuel_record")
	qs = qs.RelatedSel("Fuel__Type").
		RelatedSel("Boiler__Factory").RelatedSel("Boiler__Enterprise").RelatedSel("Boiler__Maintainer")
	if usr.IsCommonUser() ||
		usr.Status == models.USER_STATUS_INACTIVE || usr.Status == models.USER_STATUS_NEW {
		qs = qs.Filter("Boiler__IsDemo", true)
	} else {
		qs = qs.Filter("Boiler__IsDemo", false)
		if usr.IsOrganizationUser() {
			orgCond := orm.NewCondition().Or("Boiler__Enterprise__Uid", usr.Organization.Uid).Or("Boiler__Factory__Uid", usr.Organization.Uid).Or("Boiler__Maintainer__Uid", usr.Organization.Uid)
			cond := orm.NewCondition().AndCond(orgCond)
			qs = qs.SetCond(cond)
		}
	}
	num, err := qs.Filter("IsDeleted", false).OrderBy("-UpdatedDate").All(&records)

	goazure.Info("Returned Rows Num:", num, err)
	//se := ctl.GetSession(SESSION_CURRENT_USER)
	//fmt.Println("\nBoiler Get CurrentUser: ", usr, " | ", se);

	ctl.Data["json"] = records
	ctl.ServeJSON()
}

type Record struct {
	Uid		string		`json:"uid"`
	BoilerUid	string		`json:"boiler_uid"`
	Fuel		string		`json:"fuel"`

	StartDate	time.Time	`json:"start_date"`
	EndDate		time.Time	`json:"end_date"`

	FuelAmount	float64		`json:"fuel_amount"`
	
	Remark		string		`json:"remark"`
}

func (ctl *FuelController) FuelRecordUpdate() {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		token := ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
			return
		}
	}

	//resStatus := 200;
	//resBody := "Success";
	rec := Record{}
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &rec); err != nil || len(rec.BoilerUid) <= 0 || rec.FuelAmount <= 0.0 {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	var rtms []orm.Params
	total := 0.0
	qr := dba.BoilerOrm.QueryTable("boiler_runtime_cache_flow")
	qr = qr.Filter("Boiler__Uid", rec.BoilerUid).Filter("Parameter__Id", 1003).
		Filter("CreatedDate__gte", rec.StartDate).Filter("CreatedDate__lte", rec.EndDate)
	if num, err := qr.Filter("IsDeleted", false).OrderBy("-CreatedDate").Values(&rtms, "CreatedDate", "Value"); (err != nil || num == 0) {
		goazure.Warning("Read BoilerRuntime Error", err)
	} else {
		sum := float64(0)
		avg := 0.0
		due := float64(rec.EndDate.Sub(rec.StartDate)) / float64(time.Hour)
		for _, r := range rtms {
			sum += r["Value"].(float64)
		}
		avg = float64(sum) / float64(len(rtms))
		total = avg * due // * float64(param.Scale)
	}
	rate := total / rec.FuelAmount
	goazure.Warn("Fuel Flows", total, rate);

	record := &models.BoilerFuelRecord{}
	boiler := boilerWithUid(rec.BoilerUid)
	fuel := &models.Fuel{ FuelId: int32(401) }
	if err := DataCtl.ReadData(fuel, "FuelId"); err != nil {
		goazure.Error("Read Biomess Fuel Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Read Biomess Fuel Error:" + err.Error()))
		return
	}

	record.Boiler = boiler
	record.Fuel = fuel
	record.StartDate = rec.StartDate
	record.EndDate = rec.EndDate

	record.TotalFlow = total
	record.FuelAmount = rec.FuelAmount
	record.Rate = rate
	record.Remark = rec.Remark

	if err := DataCtl.AddData(record, true); err != nil {
		goazure.Error("Added Fuel Record Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Added Fuel Record Error:" + err.Error()))
	} else {
		goazure.Info("Fuel Record Added:", record)
	}
}

func (ctl *FuelController) FuelRecordDelete() {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		token := ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
			return
		}
	}

	rec := Record{}
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &rec); err != nil {
		e := fmt.Sprintln("Unmarshal Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var record models.BoilerFuelRecord
	record.Uid = rec.Uid

	if err := DataCtl.DeleteData(&record); err != nil {
		e := fmt.Sprintln("Delete FuelRecord Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
}
