package models

import "time"

type BoilerCalculateParameter struct {
	MyUidObject

	Boiler				*Boiler			`orm:"rel(fk);index"`

	CoalQnetvar			float64
	CoalAar				float64
	CoalMar				float64
	CoalVdaf			float64
	CoalClz				float64
	CoalClm				float64
	CoalCfh				float64
	CoalDed				float64
	CoalDsc				float64
	CoalAlz				float64
	CoalAlm				float64
	CoalAfh				float64
	CoalQ3				float64
	CoalM				float64
	CoalN				float64
	CoalTlz				float64
	CoalCtLz			float64

	GasDed				float64
	GasDsc				float64
	GasApy				float64
	GasQ3				float64
	GasM				float64
	GasN				float64

	ConfParam1			float64			//Q5
	ConfParam2			float64
	ConfParam3			float64
	ConfParam4			float64
	ConfParam5			float64
	ConfParam6			float64

	AlarmThreshold1		float64
	AlarmThreshold2		float64
	AlarmThreshold3		float64
	AlarmThreshold4		float64
	AlarmThreshold5		float64
	AlarmThreshold6		float64
	AlarmThreshold7		float64
	AlarmThreshold8		float64

	Reserved1			float64
	Reserved2			time.Time		`orm:"type(datetime)"`
	Reserved3			float64
	Reserved4			float64
}

var DefaultCalcParam BoilerCalculateParameter

func init() {
	DefaultCalcParam.CoalQnetvar = 20710
	DefaultCalcParam.CoalAar = 7.92
	DefaultCalcParam.CoalClz = 43.5
DefaultCalcParam.CoalClm = 31.66
DefaultCalcParam.CoalCfh = 16.08
DefaultCalcParam.CoalAlz = 75
DefaultCalcParam.CoalAlm = 5
DefaultCalcParam.CoalAfh = 20
DefaultCalcParam.CoalQ3 = 0.5
DefaultCalcParam.CoalM = 0.4
DefaultCalcParam.CoalN = 3.6
DefaultCalcParam.CoalCtLz = 570.2

DefaultCalcParam.GasQ3 = 0.5
DefaultCalcParam.GasM = 0.4
DefaultCalcParam.GasN = 3.6

DefaultCalcParam.ConfParam1 = 2.9
}
