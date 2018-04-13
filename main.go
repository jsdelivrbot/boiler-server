package main

import (
	"github.com/AzureTech/goazure"
	_ "github.com/AzureTech/goazure/session/mysql"

	"github.com/AzureRelease/boiler-server/controllers"
	_ "github.com/AzureRelease/boiler-server/routers"
	_ "github.com/AzureRelease/boiler-server/dba"
	_ "github.com/AzureRelease/boiler-server/models"
	_ "github.com/AzureRelease/boiler-server/controllers"
	_ "github.com/AzureRelease/boiler-server/log"
)

var wechatServerEnabled = true

func main() {
	//go controllers.BlrCtl.InitBoilerDefaults()
	//go controllers.RtmCtl.InitRuntimeParameters()

	//goazure.EnableAdmin = true
	//go controllers.RtmCtl.:RefreshRuntimeStatusRunning()

	//go generateRandomData(true)

	go initWechatServer()

	//屏蔽错误页详细信息
	goazure.ErrorController(&controllers.ErrorController{})

	//controllers.OrgCtrl.InitOrganizationDefaults()
	//go initDefaultData()
	//go controllers.CalcCtl.ImportBoilerCalculateFromHSEI()

	//go controllers.ParamCtrl.InitParameterChannelConfig()
	//go controllers.RtmCtl.ImportExistCache()

	//t := time.Now()
	//rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	//go controllers.RtmCtl.ReloadHistoryWithArchived(rounded, time.Now().Add(time.Minute * 2))
	//goazure.Warn(fmt.Sprintf("%2x", 17867))
	//goazure.Warn(fmt.Sprintf("%x", 17867))

	//go controllers.RtmCtl.GetBoilerRank()
	//go controllers.RtmCtl.GetRunningDuration()

	goazure.Run()
}

func initWechatServer() {
	if !wechatServerEnabled {
		return
	}

	go controllers.AlmCtl.InitAlarmSendService()
	go controllers.WxCtl.InitWechatService()
}

func initDefaultData() {
	controllers.BlrCtl.InitBoilerDefaults()
	controllers.RtmCtl.InitRuntimeParameters()
	controllers.MsgCtl.InitMessageTags()
	controllers.OrgCtrl.InitOrganizationDefaults()
	controllers.AlmCtl.InitAlarmRules()

	trimBoilerData()
	//go generateRandomData(true)
}

func generateRandomData(isOn bool) {
	go controllers.RtmCtl.GenerateBoilerStatus(isOn)
	go controllers.RtmCtl.GenerateBoilerRuntime(isOn)

}

func trimBoilerData() {
	controllers.UsrCtl.UpdateUserName()
}