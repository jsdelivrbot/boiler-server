package controllers

import (
	"fmt"
)

type DataUpdateController struct {
	MainController
}

var DataUpdateCtl *DataUpdateController = &DataUpdateController{}

func (ctl *DataUpdateController) Get() {
	session := ctl.GetSession(SESSION_CURRENT_USER)
	fmt.Println("Get Session:", session)
	if session == nil {
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(E403))
		return
	} else {
		//if NewData == true{
		//TODO:创建一个json数据格式
		ctl.Ctx.Output.Header("Cache-Control", "no-cache")
		ctl.Ctx.Output.Body([]byte("{}"))
		return
		//else {
		// 	ctl.Ctx.Output.SetStatus(204)
		//	return
	}
}
