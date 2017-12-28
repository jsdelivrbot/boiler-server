package models

type MyUidObject struct {
	//Uid		string		`orm:"pk;type(uuid);size(36);index"`
	Uid 		string		`orm:"pk;size(36);index"`

	MyObject
}

func (obj *MyUidObject) GetObject() interface{} {
	return &obj.MyObject
}

func (obj *MyUidObject) GetKey() interface{} {
	return obj.Uid
}

func (obj *MyUidObject) SetKey(key interface{}) {
	obj.Uid = key.(string)
}