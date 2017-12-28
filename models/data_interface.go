package models

type DataInterface interface {
	GetObject()		interface{}

	GetKey() 		interface{}
	SetKey(key interface{})
}
