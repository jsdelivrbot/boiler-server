package controllers

import (
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/conf"

	"github.com/AzureTech/goazure"
)

type AdminController struct {
	MainController
}

var AdminCtl *AdminController = &AdminController{}

func (ctl *AdminController) Get() {
	session := ctl.GetSession(SESSION_CURRENT_USER)
	version := conf.Version
	goazure.Error("Admin Get Session:", session)
	if session == nil {
		ctl.Ctx.Redirect(302, "/")
	} else {
		crtUsr := session.(*models.User)
		//goazure.Error("Admin Role:", crtUsr.Role.Name)
		ctl.Ctx.SetCookie("isLogin", "1")
		ctl.Ctx.SetCookie("username", crtUsr.Username)
		ctl.Ctx.SetCookie("user_picture", crtUsr.Picture)
		ctl.Ctx.SetCookie("user_role", crtUsr.Role.Name)
		ctl.Ctx.SetCookie("request_session", version)
		ctl.Data["sys_version"] = version
		ctl.Ctx.Output.Header("Cache-Control", "max-age=7200")
		ctl.TplName = "admin.html"
		ctl.Render()
	}
}
