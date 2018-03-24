package controllers

import (
	"github.com/AzureTech/goazure"
	"fmt"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"errors"
	"strings"
	"sync"
	"time"
)

type MainController struct {
	goazure.Controller

	Boilers				[]*models.Boiler
	bWaitGroup			sync.WaitGroup
}

var IsRelease bool = (goazure.AppConfig.String("runmode") == "prod");
var MainCtrl *MainController = &MainController{}

func init() {
	//now := time.Now();
	time.LoadLocation("PRC")

	initApplications()
}

func initApplications() {
	/*======================== HOLDER ==========================*/
	wechatWebHolder := models.Application{}
	wechatWebHolder.Name = "微信网站 厚德能源2025"
	//wechatWebHolder.NameEn = "Weixin Web HolderBoiler"
	wechatWebHolder.Platform = "weixin"
	wechatWebHolder.App = "website"
	wechatWebHolder.Identity = "holder"
	wechatWebHolder.Domain = "www.holderboiler.com"
	wechatWebHolder.AppId = "wxcd2d43cd41ef6912"
	wechatWebHolder.AppSecret = "eb64cbee22c12c65f25af46761a3b98e"

	wechatServiceHolder := models.Application{}
	wechatServiceHolder.Name = "微信服务号 炉管家"
	//wechatServiceHolder.NameEn = "Weixin Service HolderBoiler"
	wechatServiceHolder.Platform = "weixin"
	wechatServiceHolder.App = "service"
	wechatServiceHolder.Identity = "holder"
	wechatServiceHolder.Domain = "www.holderboiler.com"
	wechatServiceHolder.Path = "wechat-server"
	wechatServiceHolder.AppId = "wx7057a9dd005c2bd9"
	wechatServiceHolder.AppSecret = "7cf71c5f2989a56c6cca5c323c351097"
	wechatServiceHolder.OriginId = "gh_51f45f1ce9fa"
	wechatServiceHolder.ApiToken = "kiaN8akaSl4ana"
	wechatServiceHolder.AesKey = "Wp1tbvGfIBADQkFAyTMNCTavQgF4n4woOhxrZDkEBjf"

	wechatMiniHolder := models.Application{}
	wechatMiniHolder.Name = "微信小程序 锅炉在线节能平台"
	//wechatMiniHolder.NameEn = "Weixin Mini Program HolderBoiler"
	wechatMiniHolder.Platform = "weixin"
	wechatMiniHolder.App = "mini"
	wechatMiniHolder.Identity = "holder"
	wechatMiniHolder.Domain = "www.holderboiler.com"
	wechatMiniHolder.Path = "wechat-mini-server"
	wechatMiniHolder.AppId = "wxe74ea02c8906c425"
	wechatMiniHolder.AppSecret = "5af961b73ea272eceae3d28c23eb30c4"
	wechatMiniHolder.ApiToken = "excalibur"
	wechatMiniHolder.AesKey = "BbZ69KBORjgvgdTSkDbTFxDCOURpxuBJq5Jb7sXP0b7"

	DataCtl.AddData(&wechatWebHolder, true, "AppId")
	DataCtl.AddData(&wechatServiceHolder, true, "AppId")
	DataCtl.AddData(&wechatMiniHolder, true, "AppId")
}

func (ctl *MainController) Get() {
	tplName := "index.html"
	domain := ctl.Ctx.Input.Domain()
	u := ctl.Ctx.Input.URL()
	goazure.Info("Get Domain:", domain)
	if strings.HasPrefix(domain, "login") || strings.Index(u, "login") > -1 {
		tplName = "login.html"
	}

	ctl.TplName = tplName
	ctl.Render()
}

func (ctl *MainController) GetCurrentUser() *models.User {
	if ctl == nil {
		fmt.Println("UserCtl is Nil!")
		//return ctl.getSysUser()
		return nil
	}

	//goazure.Error("\n===========", ctl, "\n", reflect.TypeOf(ctl) ,"===========\n")

	usrSession := ctl.GetSession(SESSION_CURRENT_USER)
	if usrSession == nil {
		//ctl.SetSession(SESSION_CURRENT_USER, ctl.getSysUser())
		return nil
	}

	//usr := ctl.GetSession(SESSION_CURRENT_USER)
	//CurrentUser = usr.(*models.User)
	//fmt.Println("Try to Get CurrentUser:", usrSession)
	return usrSession.(*models.User)
}

func (ctl *MainController) GetCurrentUserWithToken(token string) (*models.User, error) {
	if len(token) <= 0 {
		return nil, errors.New("Token Can Not Be Empty!");
	}

	session := models.UserSession{ Token: token }
	if err := DataCtl.ReadData(&session, "Token"); err != nil {
		goazure.Error("Read UserSession Error:", err, session.IsActived)
		return nil, errors.New("Read UserSession Error:" + err.Error());
	}

	ctl.UpdateCurrentUser(session.User)
	usr := ctl.GetCurrentUser()

	return usr, nil
}

func (ctl *MainController) UpdateCurrentUser(usr *models.User) {
	qs := dba.BoilerOrm.QueryTable("user")
	qs = qs.RelatedSel("Role").RelatedSel("Supervisor").
		RelatedSel("Organization__Type").RelatedSel("Address__Location").
		Filter("Uid", usr.Uid)
	//for i, v := range params {
	//	num, err := dba.BoilerOrm.LoadRelated(&v, "BoilerMediums")
	//	aParams = append(aParams, v)
	//	fmt.Println("[",i,"]", v, num, err)
	//}
	if err := qs.One(usr); err != nil {
		fmt.Println("Get CurrentUser Details Error: ", err)
	}

	var logins []*models.UserLogin
	if num, err := dba.BoilerOrm.QueryTable("user_login").Filter("User__Uid", usr.Uid).OrderBy("-CreatedDate").Limit(10).All(&logins); err != nil || num == 0 {
		fmt.Println("Read Logins(", num, ") Error:", err)
	}
	usr.Logins = append(usr.Logins, logins...)

	var thirds []*models.UserThird
	if num, err := dba.BoilerOrm.QueryTable("user_third").
		Filter("User__Uid", usr.Uid).Filter("IsDeleted", false).
		OrderBy("-CreatedDate").All(&thirds); err != nil || num == 0 {
		fmt.Println("Read Thirds(", num, ") Error:", err)
	}
	usr.Thirds = append(usr.Thirds, thirds...)

	session := &models.UserSession{}
	if err := dba.BoilerOrm.QueryTable("user_session").
		Filter("User__Uid", usr.Uid).Filter("IsDeleted", false).Filter("IsActived", true).
		OrderBy("-CreatedDate").One(session); err != nil {
		goazure.Error("Read Session Error:", err)
	}
	usr.Sessions = append(usr.Sessions, session)

	//ctl.setCookiesWithUser(usr)
	ctl.SetSession(SESSION_CURRENT_USER, usr)
}