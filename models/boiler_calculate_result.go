package models

type BoilerCalculateResult struct {
	MyIdObject

	Boiler			*Boiler						`orm:"rel(fk);index"`
	Fuel			*Fuel						`orm:"rel(fk);index"`
	BasedParameter 	*BoilerCalculateParameter	`orm:"rel(fk);null;index"`

	Qnetvar			float64			`json:"qnetvar"`
	Aar				float64			`json:"aar"`
	Mar				float64			`json:"mar"`
	Vdaf			float64			`json:"vdaf"`
	Clz				float64			`json:"clz"`
	Clm				float64			`json:"clm"`
	Cfh				float64			`json:"cfh"`
	Ded				float64			`json:"ded"`
	Dsc				float64			`json:"dsc"`
	Alz				float64			`json:"alz"`
	Alm				float64			`json:"alm"`
	Afh				float64			`json:"afh"`
	Tlz				float64			`json:"tlz"`
	CtLz			float64			`json:"ct_lz"`

	//Apy				float64			`json:"apy"`
	M				float64			`json:"m"`
	N				float64			`json:"n"`

	Q2				float64			`json:"q2"`
	Q3				float64			`json:"q3"`
	Q4				float64			`json:"q4"`
	Q5				float64			`json:"q5"`
	Q6				float64			`json:"q6"`

	ExcessAir		float64			`json:"apy"`
	Heat			float64
}
