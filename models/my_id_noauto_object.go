package models

type MyIdNoAutoObject struct {
	Id		int64		`orm:"pk;index"`

	MyObject
}

func (obj *MyIdNoAutoObject) GetObject() interface{} {
	return &obj.MyObject
}

func (obj *MyIdNoAutoObject) GetKey() interface{} {
	return obj.Id
}

func (obj *MyIdNoAutoObject) SetKey(key interface{}) {
	obj.Id = key.(int64)
}