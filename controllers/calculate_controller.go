package controllers

import (
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"

	"github.com/AzureTech/goazure"

	"encoding/json"
	"fmt"
	"time"
)

type CalculateController struct {
	MainController
}

var CalcCtl *CalculateController = &CalculateController{}

type CalculateParameter struct {
	models.BoilerCalculateResult

	BoilerId			string			`json:"boiler_id"`
	ParameterId			string			`json:"parameter_id"`
	FuelTypeId			int				`json:"fuel_type_id"`

	SmokeTemperature	float64			`json:"smoke_temper"`
	WindTemperature		float64			`json:"wind_temper"`
	SmokeO2				float64			`json:"smoke_o2"`
}


func init() {
	go CalcCtl.InitBoilerCalculateParameter([]*models.Boiler{})
}

func (ctl *CalculateController) InitBoilerCalculateParameter(boilers []*models.Boiler) {
	if len(boilers) == 0 {
		raw := "SELECT	DISTINCT `boiler`.* " +
			"FROM	`boiler` " +
			"WHERE 	`boiler`.`uid` NOT IN " +
			"(SELECT `boiler`.`uid` FROM `boiler`, `boiler_calculate_parameter` AS `calc` WHERE `boiler`.`uid` = `calc`.`boiler_id` AND `calc`.`is_deleted` = FALSE)" +
			"AND `boiler`.`is_deleted` = FALSE;"

		if num, err := dba.BoilerOrm.Raw(raw).QueryRows(&boilers); err != nil || num == 0 {
			goazure.Warn("Read Boilers Without CalcParam Error:", err, num)
			// MainCtrl.bWaitGroup.Wait()
			// boilers = MainCtrl.Boilers
		}
	}

	for i, b := range boilers {
		var calc models.BoilerCalculateParameter
		if err := dba.BoilerOrm.QueryTable("boiler_calculate_parameter").Filter("Boiler__Uid", b.Uid).Filter("IsDeleted", false).One(&calc); err != nil {
			goazure.Warn("Read Boiler CalcParam Error:", err, i, b.Name)

			calc.Boiler = b
			calc.Name = b.Name

			calc.Reserved2 = time.Now()
		}

		if  calc.CoalQnetvar <= 0 &&
			calc.CoalAar <= 0 &&
			calc.CoalClz <= 0 &&
			calc.CoalClm <= 0 &&
			calc.CoalCfh <= 0 &&
			calc.CoalAlz <= 0 &&
			calc.CoalAlm <= 0 &&
			calc.CoalAfh <= 0 &&
			calc.CoalQ3 <= 0 &&
			calc.CoalM <= 0 &&
			calc.CoalN <= 0 &&
			calc.CoalCtLz <= 0 {
			calc.CoalQnetvar = models.DefaultCalcParam.CoalQnetvar
			calc.CoalAar = models.DefaultCalcParam.CoalAar
			calc.CoalClz = models.DefaultCalcParam.CoalClz
			calc.CoalClm = models.DefaultCalcParam.CoalClm
			calc.CoalCfh = models.DefaultCalcParam.CoalCfh
			calc.CoalAlz = models.DefaultCalcParam.CoalAlz
			calc.CoalAlm = models.DefaultCalcParam.CoalAlm
			calc.CoalAfh = models.DefaultCalcParam.CoalAfh
			calc.CoalQ3 = models.DefaultCalcParam.CoalQ3
			calc.CoalM = models.DefaultCalcParam.CoalM
			calc.CoalN = models.DefaultCalcParam.CoalN
			calc.CoalCtLz = models.DefaultCalcParam.CoalCtLz
		}

		if	calc.GasQ3 <= 0 &&
			calc.GasM <= 0 &&
			calc.GasN <= 0 {
			calc.GasQ3 = models.DefaultCalcParam.GasQ3
			calc.GasM = models.DefaultCalcParam.GasM
			calc.GasN = models.DefaultCalcParam.GasN
		}

		if  calc.ConfParam1 <= 0 { calc.ConfParam1 = models.DefaultCalcParam.ConfParam1 }

		goazure.Info(calc)

		if err := DataCtl.AddData(&calc, true, "Boiler"); err != nil {
			goazure.Error("Added & Updated Boiler CalcParam Error:", err)
		}
	}
}

func (ctl *CalculateController) BoilerCalculate() {
	//usr := ctl.GetCurrentUser()
	var calc CalculateParameter
	var result models.BoilerCalculateResult

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &calc); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Calculate Parameter Error", err)
		return
	}

	goazure.Warn("Calc:", calc)

	var q2, q3, q4, q5, q6, apy, heat float64

	switch calc.FuelTypeId {
	case 1:
		fallthrough
	case 4:
		// COAL
		q3 = calc.Q3
		q5 = calc.Q5

		q4 = ctl.CalculateQ4(calc.Qnetvar, calc.Aar, calc.Clz, calc.Clm, calc.Cfh, calc.Alz, calc.Alm, calc.Afh)
		apy = ctl.CalculateApy(calc.SmokeO2)
		q2 = ctl.CalculateQ2(calc.M, calc.N, apy, calc.SmokeTemperature, calc.WindTemperature, q4)
		q6 = ctl.CalculateQ6(calc.Afh, calc.Aar, calc.CtLz, calc.Qnetvar, calc.Clz)

		heat = ctl.CalculateHeatEfficiency(q2, q3, q4, q5, q6)
	case 2:
		fallthrough
	case 3:
		fallthrough
	default:
		// GAS
		q3 = calc.Q3
		q5 = calc.Q5

		apy = ctl.CalculateApy(calc.SmokeO2)
		q2 = ctl.CalculateQ2(calc.M, calc.N, apy, calc.SmokeTemperature, calc.WindTemperature, q4)

		heat = ctl.CalculateHeatEfficiency(q2, q3, q4, q5, q6)
	}

	result = calc.BoilerCalculateResult
	result.Boiler = &models.Boiler{}
	result.BasedParameter = &models.BoilerCalculateParameter{}

	result.Q2 = q2
	result.Q3 = q3
	result.Q4 = q4
	result.Q5 = q5
	result.Q6 = q6
	result.ExcessAir = apy
	result.Heat = heat

	goazure.Warn("CalcResult:", result)

	if err := dba.BoilerOrm.QueryTable("boiler").Filter("Uid", calc.BoilerId).One(result.Boiler); err != nil {
		e := fmt.Sprintln("Read CalculateBoiler Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if err := dba.BoilerOrm.QueryTable("boiler_calculate_parameter").Filter("Uid", calc.ParameterId).One(result.BasedParameter); err != nil {
		e := fmt.Sprintln("Read CalculateParameter Error:", err)
		goazure.Warn(e)
	}

	result.Name = result.Boiler.Name
	result.Fuel = result.Boiler.Fuel

	if err := DataCtl.AddData(&result, false); err != nil {
		e := fmt.Sprintln("Added Calculate Result Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	ctl.Data["json"] = result
	ctl.ServeJSON()
}

func (ctl *CalculateController) BoilerCalculateParameterUpdate() {
	var calc CalculateParameter

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &calc); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Calculate Parameter Error", err)
		return
	}

	goazure.Warn("Calc:", calc)

	ctl.CalculateParameterUpdate(&calc)
}

func (ctl *CalculateController) CalculateParameterUpdate(calc *CalculateParameter) {
	var param models.BoilerCalculateParameter

	if err := dba.BoilerOrm.QueryTable("boiler_calculate_parameter").Filter("Boiler__Uid", calc.BoilerId).One(&param); err != nil {
		e := fmt.Sprintln("Read CalculateBoiler Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	switch calc.FuelTypeId {
	case 1:
		fallthrough
	case 4:
		// COAL
		if calc.Qnetvar > 0 { param.CoalQnetvar = calc.Qnetvar }
		if calc.Aar > 0 	{ param.CoalAar		= calc.Aar }
		if calc.Mar > 0 	{ param.CoalMar		= calc.Mar }
		if calc.Vdaf > 0 	{ param.CoalVdaf	= calc.Vdaf }
		if calc.Clz > 0 	{ param.CoalClz		= calc.Clz }
		if calc.Clm > 0 	{ param.CoalClm		= calc.Clm }
		if calc.Cfh > 0 	{ param.CoalCfh		= calc.Cfh }
		if calc.Ded > 0 	{ param.CoalDed		= calc.Ded }
		if calc.Dsc > 0 	{ param.CoalDsc		= calc.Dsc }
		if calc.Alz > 0 	{ param.CoalAlz		= calc.Alz }
		if calc.Alm > 0 	{ param.CoalAlm		= calc.Alm }
		if calc.Afh > 0 	{ param.CoalAfh		= calc.Afh }
		if calc.Tlz > 0 	{ param.CoalTlz		= calc.Tlz }
		if calc.CtLz > 0 	{ param.CoalCtLz	= calc.CtLz }
		if calc.Q3 > 0 		{ param.CoalQ3		= calc.Q3 }
		if calc.M > 0 		{ param.CoalM		= calc.M }
		if calc.N > 0 		{ param.CoalN		= calc.N }

	case 2:
		fallthrough
	case 3:
		fallthrough
	default:
		// GAS
		if calc.Qnetvar > 0 { param.CoalQnetvar = calc.Qnetvar }
		if calc.Ded > 0 	{ param.GasDed 		= calc.Ded }
		if calc.Dsc > 0 	{ param.GasDsc 		= calc.Dsc }
		if calc.ExcessAir > 0 { param.GasApy 	= calc.ExcessAir }
		if calc.Q3 > 0 		{ param.GasQ3		= calc.Q3 }
		if calc.M > 0 		{ param.GasM		= calc.M }
		if calc.N > 0 		{ param.GasN 		= calc.N }
	}

	param.ConfParam1 = calc.Q5

	if calc.Q2 > 0 { param.ConfParam2 = calc.Q2 }
	if calc.Q3 > 0 { param.ConfParam3 = calc.Q3 }
	if calc.Q4 > 0 { param.ConfParam4 = calc.Q4 }
	if calc.Q5 > 0 { param.ConfParam5 = calc.Q5 }
	if calc.Q6 > 0 { param.ConfParam6 = calc.Q6 }

	if err := DataCtl.AddData(&param, true, "Boiler"); err != nil {
		e := fmt.Sprintln("Add BoilerCalculate Parameter Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
	}
}

//*====================== Calculate Procedure ==========================*//

func (ctl *CalculateController) CalculateQ2(m, n, apy, tpy, tlk, q4 float64) float64 {
	q2 := (m + n * apy) * ((tpy - tlk) / 100.0) * (1.0 - q4 / 100.0)

	return q2
}

func (ctl *CalculateController) CalculateQ4(qNetVar, aar, clz, clm, cfh, alz, alm, afh float64) float64 {
	c := ctl.CalculateC(clz, clm, cfh, alz, alm, afh)
	q4 := (328.66 * aar * c) / qNetVar

	return q4
}

func (ctl *CalculateController) CalculateQ6(afh, aar, cthz, qNetVar, clz float64) float64 {
	q6 := (100.0 - afh * aar * cthz) / (qNetVar * (100.0 - clz))

	return q6
}

func (ctl *CalculateController) CalculateC(clz, clm, cfh, alz, alm, afh float64) float64 {
	c := (alz * clz) / (100.0 - clz) + (alm * clm) / (100.0 - clm) + (afh * cfh) / (100.0 - cfh)

	return c
}

func (ctl *CalculateController) CalculateApy(o2 float64) float64 {
	apy := 21.0 / (21.0 - o2)

	return apy
}

func (ctl *CalculateController) CalculateHeatEfficiency(q2, q3, q4, q5, q6 float64) float64 {
	heat := 100.0 - q2 - q3 - q4 - q5 - q6

	return heat
}

