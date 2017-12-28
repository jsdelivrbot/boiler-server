package controllers

import (
	"github.com/AzureRelease/boiler-server/dba"

	"fmt"
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"
)

type DeveloperController struct {
	MainController
}

var DevCtrl *DeveloperController = &DeveloperController{}

func (ctl *DeveloperController) TerminalOriginMessageList() {
	usr := ctl.GetCurrentUser()
	if !usr.IsAdmin() {
		e := fmt.Sprintln("没有权限查看此信息！")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if ctl.Input()["dev"] == nil || len(ctl.Input()["dev"]) == 0 || ctl.Input()["dev"][0] != "origin" {
		e := fmt.Sprintln("调试参数不正确！")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if ctl.Input()["terminal"] == nil || len(ctl.Input()["terminal"]) == 0 {
		e := fmt.Sprintln("there is no TerminalCode!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	code := ctl.Input()["terminal"][0]
	if len(code) < 6 {
		for i := 0; i < 6 - len(code); i++ {
			code = fmt.Sprintf("%d%s", 0, code)
		}
	}

	var messages []orm.Params
	raw :=  "SELECT	* " +
			"FROM	`boiler_m163` " +
			"WHERE	`Boiler_term_id` = " + code + " " +
			"ORDER BY `TS` DESC LIMIT 100;"
	if num, err := dba.BoilerOrm.Raw(raw).Values(&messages); err != nil {
		e := fmt.Sprintf("Get Origin Message With[%s] Error: %v | %d", code, err, num)
		goazure.Error(e)
	}

	goazure.Warn("Get Origin Messages:", messages)

	ctl.Data["json"] = messages
	ctl.ServeJSON()
}