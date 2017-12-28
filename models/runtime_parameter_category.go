package models

type RuntimeParameterCategory struct {
	MyIdObject

	//CategoryId		int32				`orm:"index"`
	Parameters		[]*RuntimeParameter		`orm:"reverse(many)"`
}