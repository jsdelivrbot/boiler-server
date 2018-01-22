package models

type RuntimeParameterCategory struct {
	MyIdObject

	//CategoryId		int32				`orm:"index"`
	Parameters		[]*RuntimeParameter		`orm:"reverse(many)"`
}

const (
	RUNTIME_PARAMETER_CATEGORY_ANALOG 		= 10
	RUNTIME_PARAMETER_CATEGORY_SWITCH 		= 11
	RUNTIME_PARAMETER_CATEGORY_CALCULATE 	= 12
	RUNTIME_PARAMETER_CATEGORY_RANGE 		= 13
)