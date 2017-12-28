package main

import (
	"github.com/AzureTech/goazure"
	_ "github.com/AzureTech/goazure/session/mysql"

	"github.com/AzureRelease/boiler-server/controllers"
	_ "github.com/AzureRelease/boiler-server/routers"
	_ "github.com/AzureRelease/boiler-server/dba"
	_ "github.com/AzureRelease/boiler-server/models"
	_ "github.com/AzureRelease/boiler-server/controllers"
	"time"
)

func main() {
	//go controllers.BoilerCtrl.InitBoilerDefaults()
	//go controllers.RtmCtrl.InitRuntimeParameters()

	//屏蔽错误页详细信息
	goazure.ErrorController(&controllers.ErrorController{})

	//go initDefaultData()

	goazure.Run()
}

func initDefaultData() {
	controllers.BoilerCtrl.InitBoilerDefaults()
	controllers.RtmCtrl.InitRuntimeParameters()
	controllers.MsgCtl.InitMessageTags()
	controllers.OrgCtrl.InitOrganizationDefaults()
	controllers.AlarmCtl.InitAlarmRules()

	trimBoilerData()
	//go generateRandomData(true)
}

func generateRandomData(isOn bool) {
	go controllers.RtmCtrl.GenerateBoilerStatus(isOn)
	go controllers.RtmCtrl.GenerateBoilerRuntime(isOn)

	go controllers.RtmCtrl.UpdateRuntimeHistory(time.Time{}, time.Time{})

	//controllers.RtmCtl.UpdateRuntimeHistory(time.Now().Add(time.Hour * time.Duration(-hours)), time.Time{})
}

func trimBoilerData() {
	controllers.UsrCtl.UpdateUserName()

	controllers.BoilerCtrl.UpdatedBoilerModel()
}
