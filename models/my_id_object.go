package models

type MyIdObject struct {
	Id		int64		`orm:"auto;pk;index"`

	MyObject
}

func (obj *MyIdObject) GetObject() interface{} {
	return &obj.MyObject
}

func (obj *MyIdObject) GetKey() interface{} {
	return obj.Id
}

func (obj *MyIdObject) SetKey(key interface{}) {
	obj.Id = key.(int64)
}