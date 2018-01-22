package controllers

import (
	"fmt"
	"bytes"
	"encoding/json"
	"time"
	"github.com/AzureRelease/boiler-server/models"
	"net/http"
	"net/url"
	"github.com/AzureTech/goazure"
	"github.com/pborman/uuid"
	"github.com/AzureRelease/boiler-server/dba"
	"errors"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type UserThirdController struct {
	UserController
}

//&{200 OK 200 HTTP/1.1 1 1 map[Connection:[keep-alive] Content-Type:[text/plain] Date:[Sat, 04 Mar 2017 09:13:16 GMT] Content-Length:[78]] 0xc4203547c0 78 [] false false map[] 0xc4200d31d0 0xc4203c5d90}
//&{0xc4207c6040 {0 0} false <nil> 0x161870 0x161800}

/*
{
	"access_token":"E3U6nI1QnDReW2NVtGdZRg-jEIW9NEOmTHzAY92Z5IaXysNKViMZKDL8zdRv1oKCsARHHI9q3PTvt2fmimqjNfne5Qw4LNjSU6TjkK0WpMs",
	"expires_in":7200,
	"refresh_token":"XbrhPvVXakFg3DzpFZtB-q6yfLDru4-nP1qWUJZOCZpelA-okbXe11NcBm7_h4ahMBxO9W8M0YPRH6t1S6bxRH5MXmvIyg5PInXoNVGrJHg",
	"openid":"oMIu4w_7XwePp2bruFuGvusQjGHg",
	"scope":"snsapi_login",
	"unionid":"ojfIXwH4p2pS_SAiX2SDLrif-10I"
}
*/

type WeixinAccess struct {
	OpenId			string		`json:"openid"`
	UnionId			string		`json:"unionid"`
	
	SessionKey		string		`json:"session_key"`
	
	AccessToken		string		`json:"access_token"`
	RefreshToken	string		`json:"refresh_token"`
	ExpiresIn		int			`json:"expires_in"`
	Scope			string		`json:"scope"`

	ErrCode			int			`json:"errcode"`
	ErrMsg			string		`json:"errmsg"`
}

/*
	openid		普通用户的标识，对当前开发者帐号唯一
	nickname	普通用户昵称
	sex		普通用户性别，1为男性，2为女性
	province	普通用户个人资料填写的省份
	city		普通用户个人资料填写的城市
	country		国家，如中国为CN
	headimgurl	用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
	privilege	用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	unionid		用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
*/

type WeixinUserInfo struct {
	OpenId			string		`json:"openid"`
	UnionId			string		`json:"unionid"`
	Code			string		`json:"code"`

	EncryptData		string		`json:"encryptData"`
	EncryptedData	string		`json:"encryptedData"`
	Iv				string		`json:"iv"`
	Signature		string		`json:"signature"`

	Username		string		`json:"username"`
	Password		string		`json:"password"`
	
	Language		string		`json:"language"`
	GroupId       	int 		`json:"groupid"`

	Nickname		string		`json:"nickname"`
	Sex				int			`json:"sex"`
	Province		string		`json:"province"`
	City			string		`json:"city"`
	Country			string		`json:"country"`
	HeadImageUrl	string		`json:"headimgurl"`
	Privilege		[]string	`json:"privilege"`

	Subscribe		int			`json:"subscribe"`
	SubscribeTime	int 		`json:"subscribe_time"`

	Remark			string		`json:"remark"`

	ErrCode			int			`json:"errcode"`
	ErrMsg			string		`json:"errmsg"`
}

type WeixinSecretInfo struct {
	OpenId		string		`json:"openId"`
	UnionId		string		`json:"unionId"`
	
	NickName	string		`json:"nickName"`
	Gender		int		`json:"gender"`
	Language	string		`json:"language"`
	City		string		`json:"city"`
	Province	string		`json:"province"`
	Country		string		`json:"country"`
	AvatarUrl	string		`json:"avatarUrl"`
	
	Watermark	struct {
		Timestamp	int64	`json:"timestamp"`
		AppId		string	`json:"appId"`
			 }		`json:"watermark"`
}

type WeixinTemplateMessage struct {
	ToUser		string		`json:"touser"`		// *接收者（用户）的 openid
	TemplateId	string		`json:"template_id"`	// *所需下发的模板消息的id
	Page		string		`json:"page"`		// 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	FormId		string		`json:"form_id"`	// *表单提交场景下，为 submit 事件带上的 formId；支付场景下，为本次支付的 prepay_id
	Data		struct{}	`json:"data"`		// *模板内容，不填则下发空模板
	Color		string		`json:"color"`		// 模板内容字体的颜色，不填默认黑色
	EmphasisKeyword string	`json:"emphasis_keyword"`	// 模板需要放大的关键词，不填则默认无放大
}

func (ctl *UserThirdController) UserLoginWeixinWeb() {
	fmt.Println("Ready to Login With Weixin!!!")
	fmt.Println(ctl.Ctx.Input.URI(), ctl.Ctx.Input.URL(), ctl.Ctx.Input.Domain(), ctl.Ctx.Input.Params())

	domain := ctl.Ctx.Input.Domain()
	app := models.Application{ Domain: domain, App: "website" }
	if err := DataCtl.ReadData(&app, "Domain", "App"); err != nil {
		fmt.Println("Read AppInfo Error:", app)
	}

	if err := ctl.LoginWeixinWeb(&app, "user_login_weixin"); err != nil {
		fmt.Println("Login With Weixin Web Failed")
	}
}

func (ctl *UserThirdController) UserLoginWeixinWebCallback() {
	fmt.Println("Ready to Login Callback With Weixin!!!")
	fmt.Println(ctl.Ctx.Input.IP(), ctl.Ctx.Input.URI(), ctl.Ctx.Input.URL(), ctl.Ctx.Input.Domain())
	fmt.Println(ctl.Input())

	code := ctl.Input()["code"][0]
	//state := ctl.Input()["state"]
	domain := ctl.Ctx.Input.Domain()
	app := models.Application{ Domain: domain, App: "website" }
	if err := DataCtl.ReadData(&app, "Domain", "App"); err != nil {
		fmt.Println("Read AppInfo Error:", app)
	}

	wa, err := GetWeixinAccess(&app, code, false)
	if err != nil || wa.ErrCode != 0 {
		fmt.Println("GetWeixinAccess Error:", err)
		ctl.Ctx.Redirect(302, "/")
		return
	}

	info, err := GetWeixinUserInfo(wa)
	if err != nil || info.ErrCode != 0 {
		fmt.Println("GetWeixinUserInfo Error:", err)
		ctl.Ctx.Redirect(302, "/")
		return
	}

	var third models.UserThird

	if err := dba.BoilerOrm.QueryTable("user_third").Filter("OpenId", wa.OpenId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&third); err != nil {
		goazure.Warn("Read UserThird(OpenId) Error", err)

		third.Application = &app
		third.Platform = app.Platform
		third.App = app.App
		third.Identity = app.Identity
		third.OpenId = wa.OpenId
		third.UnionId = wa.UnionId

		third.Name = info.Nickname
		third.Sex = info.Sex
		third.Province = info.Province
		third.City = info.City
		third.Country = info.Country
		third.HeadImageUrl = info.HeadImageUrl

		var ath models.UserThird
		if err := dba.BoilerOrm.QueryTable("user_third").Filter("UnionId", wa.UnionId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&ath); err != nil {
			user := models.User{}

			goazure.Warn("Read UserThird(UnionId) Error", err)
			user.Status = models.USER_STATUS_THIRD
			user.Role = role(models.USER_ROLE_USER)
			user.Thirds = append(user.Thirds, &third)

			user.Name = third.Name
			user.Picture = third.HeadImageUrl

			ctl.SetSession(SESSION_CURRENT_USER, &user)
		} else {
			goazure.Info("Get UserThird(UnionId):", ath)

			third.User = ath.User

			third.AccessToken = wa.AccessToken
			third.RefreshToken = wa.RefreshToken
			third.ExpiresIn = time.Now().Add(time.Second * time.Duration(wa.ExpiresIn))

			if err := DataCtl.AddData(&third, true, "OpenId"); err != nil {
				goazure.Error("Third Record Add Error:", err, third)
			}

			ctl.LoginWeixinWebDone(&third)
		}
	} else {
		third.Name = info.Nickname

		third.AccessToken = wa.AccessToken
		third.RefreshToken = wa.RefreshToken
		third.ExpiresIn = time.Now().Add(time.Second * time.Duration(wa.ExpiresIn))

		DataCtl.UpdateData(&third)

		ctl.LoginWeixinWebDone(&third)
	}

	ctl.Ctx.Redirect(302, "/admin#!/")
}

func (ctl *UserThirdController) LoginWeixinWebDone(third *models.UserThird) error {
	if third.Application == nil {
		return errors.New("UserThird Application Can NOT be nil!")
	}

	if third.User == nil {
		return errors.New("UserThird User Can NOT be nil!")
	}

	login := models.UserLogin{}
	login.User = third.User

	login.Name = third.Name
	login.Remark = "微信登录成功"
	login.CreatedBy = third.User

	login.LoginPassword = "Weixin QRCode Scaned"
	login.IsSuccess = true
	login.LoginMethod = "Weixin Web"

	if err := DataCtl.InsertData(&login); err != nil {
		goazure.Error("LoginLog Error", err)
	}

	var sessions []*models.UserSession
	if num, err := dba.BoilerOrm.QueryTable("user_session").
		Filter("User__Uid", third.User.Uid).Filter("Application__Uid", third.Application.Uid).
		Filter("IsDeleted", false).Filter("IsActived", true).
		OrderBy("-CreatedDate").All(&sessions); err != nil || num == 0 {
		goazure.Warning("Read Actived Session Error:", err)
	} else {
		for _, se := range sessions {
			se.IsActived = false
			DataCtl.UpdateData(se)
		}
	}

	session := models.UserSession {
		User: third.User,
		Login: &login,
		Application: third.Application,
		Token: uuid.New(),
		IsActived: true,
	}
	if err := DataCtl.InsertData(&session); err != nil {
		goazure.Error("Session Error", err)
	}

	ctl.UpdateCurrentUser(third.User)

	return nil
}

func (ctl *UserThirdController) UserBindWeixin() {
	fmt.Println("Ready to Bind With Weixin!!!")
	domain := ctl.Ctx.Input.Domain()
	app := models.Application{ Domain: domain, App: "website" }
	if err := DataCtl.ReadData(&app, "Domain", "App"); err != nil {
		fmt.Println("Read AppInfo Error:", app)
	}

	if err := ctl.LoginWeixinWeb(&app, "user_bind_weixin"); err != nil {
		fmt.Println("Login With Weixin Web Failed")
	}
}

func (ctl *UserThirdController) UserBindWeixinCallback() {
	fmt.Println("Ready to Bind Callback With Weixin!!!")

	code := ctl.Input()["code"][0]
	//state := ctl.Input()["state"]
	domain := ctl.Ctx.Input.Domain()
	app := models.Application{ Domain: domain, App: "website" }
	if err := DataCtl.ReadData(&app, "Domain", "App"); err != nil {
		fmt.Println("Read AppInfo Error:", app)
	}

	wa, err := GetWeixinAccess(&app, code, false)
	if err != nil || wa.ErrCode != 0 {
		fmt.Println("GetWeixinAccess Error:", err)
		ctl.Ctx.Redirect(302, "/")
		return
	}

	info, err := GetWeixinUserInfo(wa)
	if err != nil || info.ErrCode != 0 {
		fmt.Println("GetWeixinUserInfo Error:", err)
		ctl.Ctx.Redirect(302, "/")
		return
	}

	usr := ctl.GetCurrentUser()
	third := models.UserThird{}
	third.Application = &app
	third.Platform = app.Platform
	third.OpenId = wa.OpenId
	third.UnionId = wa.UnionId
	third.AccessToken = wa.AccessToken
	third.RefreshToken = wa.RefreshToken
	third.ExpiresIn = time.Now().Add(time.Second * time.Duration(wa.ExpiresIn))

	third.Name = info.Nickname
	if len(usr.Name) == 0 || usr.Name == usr.Username {
		usr.Name = info.Nickname
	}
	if usr.Gender == models.USER_GENDER_UNKNOWN {
		usr.Gender = info.Sex
	}

	third.User = usr
	third.IsDeleted = false

	if err := DataCtl.AddData(&third, true, "OpenId"); err != nil {
		goazure.Error("Bind UserThird Error", err, third.User, third.IsDeleted)
	}

	DataCtl.UpdateData(usr)
	ctl.UpdateCurrentUser(usr)

	ctl.Ctx.Redirect(302, "/admin#!/profile/account")
}

func (ctl *UserThirdController) UserLoginBindThird() {
	au := ctl.GetSession(SESSION_CURRENT_USER).(*models.User)
	third := au.Thirds[0]

	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}

	usr, err := ctl.Login(&u)
	if err != nil {
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	third.User = usr
	third.IsDeleted = false

	if len(usr.Name) == 0 || usr.Name == usr.Username {
		usr.Name = third.Name
	}
	if err := DataCtl.AddData(third, true, "OpenId"); err != nil {
		goazure.Warn("Login And Bind User Error: ", err)
	}
	DataCtl.UpdateData(usr)

	var aThirds []*models.UserThird
	if num, err := dba.BoilerOrm.QueryTable("user_third").
		Filter("UnionId", third.UnionId).All(&aThirds); err != nil {
		goazure.Error("Update UnionIds:", err, num)
	}
	for _, th := range aThirds {
		th.User = third.User
		th.IsDeleted = false

		DataCtl.UpdateData(th)
	}

	fmt.Println("User: ", usr)
	fmt.Println("Third: ", third)
	ctl.UpdateCurrentUser(usr)

	ctl.Ctx.Output.SetStatus(200)
	ctl.Ctx.Output.Body([]byte("绑定成功，您可以通过微信扫码直接登录平台"))
}

func (ctl *UserThirdController) UserRegisterBindThird() {
	au := ctl.GetSession(SESSION_CURRENT_USER).(*models.User)
	third := au.Thirds[0]

	usr, err := ctl.UserRegister()
	if err != nil {
		goazure.Error("UserRegisterBindThird Error:", err)
		return
	}

	if len(third.HeadImageUrl) > 0 {
		usr.Picture = third.HeadImageUrl
	}

	if len(third.Name) > 0 {
		usr.Name = third.Name
	}

	third.User = usr
	third.IsDeleted = false

	fmt.Println("User Name:", usr.Name)
	if len(usr.Name) == 0 || usr.Name == usr.Username {
		usr.Name = third.Name
		goazure.Info("User Name Updated:", usr.Name)
	}
	if err := DataCtl.AddData(third, true, "OpenId", "UnionId"); err != nil {
		goazure.Error("Login And Bind User Error: ", err)
	}
	DataCtl.UpdateData(usr)

	fmt.Println("Register User: ", usr)
	fmt.Println("Register Third: ", third)

	ctl.UpdateCurrentUser(usr)

	ctl.Ctx.Output.SetStatus(200)
	ctl.Ctx.Output.Body([]byte("绑定成功，您可以通过微信扫码直接登录平台"))
}

func (ctl *UserThirdController) UserUnbindWeixin() {
	var token string
	var openid string
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		if ctl.Input()["token"] == nil || len(ctl.Input()["token"]) == 0 {
			return
		}
		token = ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))

			return
		}
	}

	if ctl.Input()["openid"] != nil && len(ctl.Input()["openid"]) > 0 {
		openid = ctl.Input()["openid"][0]
	}

	var oThirds, uThirds []*models.UserThird
	qs := dba.BoilerOrm.QueryTable("user_third").Filter("Platform", "weixin").Filter("IsDeleted", false)
	if len(openid) <= 0 {
		qs = qs.Filter("User__Uid", usr.Uid)
	} else {
		qs = qs.Filter("OpenId", openid)
	}

	if num, err := qs.All(&oThirds); err != nil || num == 0 {
		e := fmt.Sprintln("Get UserThird By OpenId For Delete Error:", err, num)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))

		return
	}

	var th *models.UserThird
	for _, t := range oThirds {
		if len(t.UnionId) > 0 {
			th = t
		}
	}

	if th == nil || len(th.UnionId) <= 0 {
		e := fmt.Sprintln("Get UnionId for Delete Error!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))

		return
	}

	if num, err := dba.BoilerOrm.QueryTable("user_third").Filter("UnionId", th.UnionId).All(&uThirds); err != nil || num == 0 {
		e := fmt.Sprintln("Get UserThird By UnionId For Delete Error:", err, num)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))

		return
	}

	for _, t := range uThirds {
			t.User = nil
			t.IsDeleted = true
			DataCtl.UpdateData(t)
	}

	ctl.UpdateCurrentUser(usr)
}

func (ctl *UserThirdController) UserLoginWeixinMini() {
	goazure.Error("Ready to UserLoginWeixinMini()")
	third := models.UserThird{}

	info := WeixinUserInfo{}
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &info); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	goazure.Info("WeixinMini User Info:", info)
	
	domain := ctl.Ctx.Input.Domain()
	app := models.Application{ Domain: domain, App: "mini" }
	if err := DataCtl.ReadData(&app, "Domain", "App"); err != nil {
		goazure.Error("Read AppInfo Error:", app)
	}

	access, err := GetWeixinAccess(&app, info.Code, false)
	if err != nil || access.ErrCode != 0 {
		e := fmt.Sprintln("GetWeixinMiniAccess Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	sif, err := ctl.UserDataDecryptBase64(access.SessionKey, info.EncryptedData, info.Iv)
	if err != nil {
		goazure.Error(err)
	}

	if err := dba.BoilerOrm.QueryTable("user_third").Filter("OpenId", access.OpenId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&third); err != nil {
		goazure.Error("Read UserThird Error", err)

		third.OpenId = access.OpenId
		third.UnionId = sif.UnionId

		var ath models.UserThird
		if err := dba.BoilerOrm.QueryTable("user_third").Filter("UnionId", sif.UnionId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&ath); err != nil {
			goazure.Error("Not Get UnionId!", ath)
			//usr.Status = models.USER_STATUS_THIRD
			//usr.Role = role(models.USER_ROLE_USER)
			//usr.Thirds = append(usr.Thirds, &third)
			//
			//ctl.SetSession(SESSION_CURRENT_USER, &usr)
			//ctl.Ctx.Output.SetStatus(400)
			//ctl.Ctx.Output.Body([]byte("Third Read Error!"))
			//return
		} else {
			goazure.Warn("Get Exist UnionId User:", ath)
			third.User = ath.User
			third.UnionId = ath.UnionId
			third.IsDeleted = false
		}
	}

	if len(third.UnionId) == 0 {
		third.UnionId = sif.UnionId
	}

	third.Application = &app
	third.Platform = app.Platform
	third.App = app.App
	third.Identity = app.Identity
	third.SessionKey = access.SessionKey
	third.ExpiresIn = time.Now().Add(time.Second * time.Duration(access.ExpiresIn))

	third.Name = info.Nickname
	third.Sex = info.Sex
	third.HeadImageUrl = info.HeadImageUrl
	third.Province = info.Province
	third.City = info.City
	third.Country = info.Country

	//third.IsDeleted = false
	//goazure.Info("Third: ", third.User)
	DataCtl.AddData(&third, true, "OpenId")
	//goazure.Info("Third Updated: ", third.User)

	if err := ctl.LoginWeixinMini(&third); err != nil {
		goazure.Error(err)
		ctl.Ctx.Output.SetStatus(302)
		ctl.Data["json"] = third
	} else {
		ctl.Data["json"] = ctl.GetCurrentUser()
	}

	ctl.ServeJSON()
}

func (ctl *UserThirdController) UserDataDecryptBase64(key, encryptedData, iv string) (*WeixinSecretInfo, error) {
	aKey, _ := base64.StdEncoding.DecodeString(key)
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedData)

	block, err := aes.NewCipher(aKey)
	if err != nil {
		e := errors.New(fmt.Sprintln("ciphertext Error:", ciphertext))
		return nil, e
	}
	aIv, _ := base64.StdEncoding.DecodeString(iv)

	// CBC mode always works in whole blocks.
	if len(ciphertext) % aes.BlockSize != 0 {
		e := errors.New("ciphertext is not a multiple of the block size")
		return nil, e
	}

	mode := cipher.NewCBCDecrypter(block, aIv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
	//goazure.Error("ciphertext", fmt.Sprintf("%s", ciphertext))
	//goazure.Warn(ciphertext)

	if idx := bytes.LastIndexByte(ciphertext, byte(125)); idx > -1 {
		//goazure.Warn("Find Last '}' Index:", idx)
		ciphertext = ciphertext[:idx + 1]
	}
	//goazure.Warn(ciphertext)
	goazure.Info(fmt.Sprintf("%s", ciphertext))

	var sinfo WeixinSecretInfo
	if err := json.Unmarshal(ciphertext, &sinfo); err != nil {
		e := errors.New(fmt.Sprintln("Weixin User Secret Ummarshal Error:", err))
		return nil, e
	}
	goazure.Info("WeixinMini Secret Info:", sinfo)

	return &sinfo, nil
}

func (ctl *UserThirdController) UserBindWeixinMini() {
	var third models.UserThird

	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Error("Unmarshal Error", err)
		return
	}
	goazure.Info("WeixinMini Bind:", u)

	//domain := ctl.Ctx.Input.Domain()
	//app := models.Application{ Domain: domain, App: "mini" }
	//if err := DataCtl.ReadData(&app, "Domain", "App"); err != nil {
	//	goazure.Error("Read AppInfo Error:", err)
	//}

	usr, err := ctl.Login(&u)
	if err != nil {
		e := fmt.Sprintln("Weixin Mini Login Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if err := dba.BoilerOrm.QueryTable("user_third").Filter("OpenId", u.OpenId).One(&third); err != nil {
		e := fmt.Sprintln("Read UserThird Error", err, third.User, third.IsDeleted)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	third.User = usr
	third.IsDeleted = false
	DataCtl.AddData(&third, true, "OpenId")

	var aThirds []*models.UserThird
	if num, err := dba.BoilerOrm.QueryTable("user_third").
		Filter("UnionId", third.UnionId).All(&aThirds); err != nil {
		goazure.Error("Update UnionIds:", err, num)
	}
	for _, th := range aThirds {
		th.User = third.User
		th.IsDeleted = false

		DataCtl.UpdateData(th)
	}

	if err := ctl.LoginWeixinMini(&third); err != nil {
		goazure.Error(err)
		ctl.Ctx.Output.SetStatus(302)
		ctl.Data["json"] = third
	} else {
		ctl.Data["json"] = ctl.GetCurrentUser()
	}

	ctl.ServeJSON()
}

func (ctl *UserThirdController) LoginWeixinMini(third *models.UserThird) error {
	goazure.Info("Weixin Mini Login: ", third, third.User)
	if third.Application == nil {
		return errors.New("UserThird Application Can NOT be nil!")
	}

	if third.User == nil {
		return errors.New("UserThird User Can NOT be nil!")
	}

	login := models.UserLogin {}
	login.User = third.User

	login.Name = third.Name
	login.Remark = "微信小程序登录成功"
	login.CreatedBy = third.User

	login.LoginPassword = "Weixin Mini Login"
	login.IsSuccess = true
	login.LoginMethod = "Weixin Mini"

	if err := DataCtl.InsertData(&login); err != nil {
		goazure.Error("LoginLog Error", err)
	}

	var sessions []*models.UserSession
	if num, err := dba.BoilerOrm.QueryTable("user_session").
		Filter("User__Uid", third.User.Uid).Filter("Application__Uid", third.Application.Uid).
		Filter("IsDeleted", false).Filter("IsActived", true).
		OrderBy("-CreatedDate").All(&sessions); err != nil || num == 0 {
		goazure.Warning("Read Actived Session Error:", err)
	} else {
		for _, se := range sessions {
			se.IsActived = false
			DataCtl.UpdateData(se)
		}
	}

	session := models.UserSession {
		User: third.User,
		Login: &login,
		Application: third.Application,
		Token: uuid.New(),
		IsActived: true,
	}
	if err := DataCtl.InsertData(&session); err != nil {
		goazure.Error("Session Error", err)
	}

	ctl.UpdateCurrentUser(third.User)

	return nil
}

func (ctl *UserThirdController) LoginWeixinWeb(app *models.Application, callbackPath string) error {
	u, _ := url.Parse(WEIXIN_OPEN_BASE_URL)
	u.Path = "/connect/qrconnect"
	q := u.Query()
	q.Set("appid", app.AppId)
	q.Set("redirect_uri", fmt.Sprintf("http://%s/%s/callback", app.Domain, callbackPath))
	q.Set("response_type", "code")
	q.Set("scope", "snsapi_login")
	q.Set("state", "boiler_login#wechat_redirect")
	u.RawQuery = q.Encode()

	//ctl.Redirect(url, 302)
	ctl.Ctx.Redirect(302, u.String())

	return nil;
}

func (ctl *UserThirdController) SendTemplateMessageMini(app *models.Application, tempMsg *WeixinTemplateMessage) {
	access, err := GetWeixinAccess(app, "", true)
	if err != nil {
		goazure.Error("Get WeixinAccess For TemplateMessageMini Error:", access)
	}

	//"https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=ACCESS_TOKEN"
}

func GetWeixinAccess(app *models.Application, code string, isTemplate bool) (*WeixinAccess, error) {
	var uri string
	switch app.App {
	case "website":
		uri = fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", app.AppId, app.AppSecret, code)
	case "mini":
		if isTemplate {
			uri = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", app.AppId, app.AppSecret)
		} else {
			uri = fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", app.AppId, app.AppSecret, code)
		}
	}

	goazure.Info("Weixin Access URL:", uri)
	resp, err := http.Get(uri)
	//req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		goazure.Error("Get Weixin LoginInfo Error:", err)
	}
	goazure.Info("Get Weixin LoginInfo: ", resp)
	if resp != nil && resp.StatusCode != 200 {
		goazure.Error("Get Weixin LoginInfo ErrorCode:", resp.Status, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println("WeixinAccess Body:", buf.String())
	//body := make([]byte, resp.ContentLength)
	//if n, err := resp.Body.Read(body); err != nil || n == 0 {
	//	fmt.Println("Read Weixin LoginInfo Body Error:", err, n)
	//}
	/* Do Not Work (by == []byte{} )
	by, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Weixin LoginInfo Body Error:", err)
	}
	fmt.Println("byte:", by)
	*/
	//var obj map[string]json.RawMessage   //Not Good for values still be []byte
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(buf.Bytes(), &obj); err != nil {
		goazure.Error("Unmarshal Error", err)
		return nil, err
	}

	var wa WeixinAccess
	if err := json.Unmarshal(buf.Bytes(), &wa); err != nil {
		goazure.Error("Unmarshal Error", err)
		return nil, err
	}

	goazure.Info("WeixinAccess:", wa)
	return &wa, err
}

func GetWeixinUserInfo(wa *WeixinAccess) (*WeixinUserInfo, error) {
	//https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
	uri := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", wa.AccessToken, wa.OpenId)
	resp, err := http.Get(uri)
	if err != nil {
		goazure.Error("Get Weixin LoginInfo Error:", err)
	}
	if resp != nil && resp.StatusCode != 200 {
		goazure.Error("Get Weixin LoginInfo ErrorCode:", resp.Status, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	goazure.Info("WeixinUserInfo Body:", buf.String())

	var info WeixinUserInfo
	if err := json.Unmarshal(buf.Bytes(), &info); err != nil {
		goazure.Error("Unmarshal Error", err)
		return nil, err
	}

	goazure.Info("WeixinUserInfo:", info)

	return &info, err
}