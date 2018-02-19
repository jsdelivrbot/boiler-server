package controllers

import (
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/conf"

	"io/ioutil"
	"fmt"
	"encoding/json"
	"net/url"
	"net/http"
	"encoding/xml"
	"sort"
	"crypto/sha1"
	"github.com/AzureTech/goazure"
	"time"
	"strings"
	"log"

	"github.com/AzureTech/wechat/mp/base"
	"github.com/AzureTech/wechat/mp/core"
	"github.com/AzureTech/wechat/mp/menu"
	"github.com/AzureTech/wechat/mp/user"
	"github.com/AzureTech/wechat/mp/message/callback/request"
	"github.com/AzureTech/wechat/mp/message/callback/response"
	"github.com/AzureTech/wechat/mp/message/custom"
	"github.com/AzureTech/wechat/mp/message/template"
)

type Request struct {
	XMLName					xml.Name 	`xml:"xml"`
	ToUserName   			string		//`xml:",cdata"`
	FromUserName 			string		//`xml:",cdata"`
	CreateTime   			time.Duration	//`xml:"CreateTime"`
	MsgType      			string		//`xml:",cdata"`
	Content      			string		//`xml:",cdata"`
	Location_X, Location_Y	float32
	Scale					int
	Label					string
	PicUrl					string
	MsgId					int
}

type Response struct {
	XMLName 		xml.Name 	`xml:"xml"`
	ToUserName   		string		`xml:",cdata"`
	FromUserName 		string		`xml:",cdata"`
	CreateTime   		time.Duration	//`xml:"CreateTime"`
	MsgType      		string		`xml:",cdata"`
	Content      		string		`xml:",cdata"`
	ArticleCount 		int     	`xml:",omitempty"`
	Articles     		[]*item 	`xml:"Articles>item,omitempty"`
	FuncFlag     		int
}

type item struct {
	XMLName     	xml.Name `xml:"item"`
	Title       	string
	Description 	string
	PicUrl      	string
	Url         	string
}

type WechatController struct {
	MainController
}

type WeixinMessage struct {
	XMLName			xml.Name	`xml:"xml"`
	ToUserName		string		`xml:"ToUserName"`
	FromUserName	string		`xml:"FromUserName"`
	CreateTime		int			`xml:"CreateTime"`
	MsgType			string		`xml:"MsgType"`
	Event			string		`xml:"Event"`
	EventKey		string		`xml:"EventKey"`
	Content			string		`xml:"Content"`
	MsgId			string		`xml:"MsgId"`
}

//var app *models.Application = &models.Application{}
var WxCtl *WechatController = &WechatController{}

var (
	identity		string
	app, mini		*models.Application = &models.Application{}, &models.Application{}

	msgHandler 		core.Handler
	msgServer  		*core.Server

	accessTokenServer 	core.AccessTokenServer
	wechatClient      	*core.Client
)

func (ctl *WechatController) InitWechatService() {
	//app.AppId = WEIXIN_PUBLIC_LUCID_APPID
	if conf.IsRelease {
		identity = "holder"
	} else {
		identity = "lucid"
	}

	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)

	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
	mux.EventHandleFunc(request.EventTypeSubscribe, eventSubscribeHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)

	msgHandler = mux

	if err := dba.BoilerOrm.QueryTable("application").Filter("Identity", identity).
		Filter("App", "service").One(app); err != nil {
		goazure.Error("Read App Error:", err)
	}

	if err := dba.BoilerOrm.QueryTable("application").Filter("Identity", identity).
		Filter("App", "mini").One(mini); err != nil {
		goazure.Error("Read Mini Error:", err)
	}

	msgServer = core.NewServer(app.OriginId, app.AppId, app.ApiToken, app.AesKey, msgHandler, nil)

	accessTokenServer = core.NewDefaultAccessTokenServer(app.AppId, app.AppSecret, nil)
	wechatClient = core.NewClient(accessTokenServer, nil)

	fmt.Println(base.GetCallbackIP(wechatClient))

	WxCtl.SyncUserList()
	WxCtl.SyncMenu()
}

func textMsgHandler(ctx *core.Context) {
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)

	msg := request.GetText(ctx.MixedMsg)
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	//ctx.RawResponse(resp) // 明文回复
	//WXCtrl.SendText(msg.FromUserName, "锅炉好严重告警啊", "")

	ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	//ctx.RawResponse(resp) // 明文回复
	ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func eventSubscribeHandler(ctx *core.Context) {
	msg := fmt.Sprintf("收到订阅事件:\n%s\n", ctx.MsgPlaintext)
	openid := fmt.Sprintf("OpenId: %s\n", ctx.MixedMsg.FromUserName)
	goazure.Warn(msg)
	goazure.Info(openid)

	info, err := user.Get(wechatClient, string(ctx.MixedMsg.FromUserName), user.LanguageZhCN)
	var third models.UserThird
	if err := dba.BoilerOrm.QueryTable("user_third").RelatedSel("User__Role").
		Filter("OpenId", info.OpenId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&third); err != nil {
		goazure.Warn("Read UserThird(OpenId) Error", err)

		third.Application = app
		third.Platform = app.Platform
		third.App = app.App
		third.Identity = app.Identity
		third.OpenId = info.OpenId
		third.UnionId = info.UnionId

		third.Name = info.Nickname
		third.Sex = info.Sex
		third.Province = info.Province
		third.City = info.City
		third.Country = info.Country
		third.HeadImageUrl = info.HeadImageURL

		var ath models.UserThird
		if err := dba.BoilerOrm.QueryTable("user_third").RelatedSel("User__Role").
			Filter("UnionId", info.UnionId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&ath); err != nil {
			goazure.Warn("Read UserThird(UnionId) Error", err)
		} else {
			goazure.Info("Get UserThird(UnionId):", ath)

			third.User = ath.User

			//ctl.LoginWeixinWebDone(&third)
		}

		if err := DataCtl.AddData(&third, true, "OpenId"); err != nil {
			goazure.Error("Third Record Add Error:", err, third)
		}
	} else {
		third.Name = info.Nickname

		DataCtl.UpdateData(&third)
		//ctl.LoginWeixinWebDone(&third)
	}

	var content string = "欢迎关注锅炉在线节能平台，"
	if third.User != nil {
		content += fmt.Sprintf("%s %s！\n您可以通过本平台接收锅炉相关告警信息，或者通过本平台的在线小程序来查看锅炉的实时运行情况。", third.User.Role.Name, third.User.Name)
	} else {
		content += fmt.Sprintf("新用户 %s！\n如果您已经是本平台用户，请点击下方菜单中“在线小程序 —— 账号·绑定”进行绑定，\n如果您还未注册本平台，请访问%s进行注册。", info.Nickname, app.Domain)
	}

	event := request.GetSubscribeEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, content)
	//ctx.RawResponse(resp) // 明文回复
	ctx.AESResponse(resp, 0, "", nil) // aes密文回复

	//text := custom.NewText(third.OpenId, "锅锅锅炉发生严重告警！", "")
	//custom.Send(wechatClient, text)
	goazure.Error(info, "|UnionId:", info.UnionId, "|Errors:", err)
}


func (ctl *WechatController) SendText(openId, content, customerService string) {
	goazure.Error("ready to send:", content)
	text := custom.NewText(openId, content, customerService)
	if err := custom.Send(wechatClient, text); err != nil {
		goazure.Error("Wechat Mesg Send Error:", err)
	}
}

func (ctl *WechatController) SendTemplateMessage(openId string, tempMsg *template.TemplateMessage2) {
	tempMsg.ToUser = openId

	if id, err := template.Send(wechatClient, tempMsg); err != nil {
		goazure.Error("Wechat TemplateMsg Send Error:", err, id)
	}
	/*
		ToUser     string          `json:"touser"`             // 必须, 接受者OpenID
    	TemplateId string          `json:"template_id"`        // 必须, 模版ID
    	URL        string          `json:"url,omitempty"`      // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
    	TopColor   string          `json:"topcolor,omitempty"` // 可选, 整个消息的颜色, 可以不设置
    	Data       json.RawMessage `json:"data"`               // 必须, JSON 格式的 []byte, 满足特定的模板需求
	 */
	/*
		{{first.DATA}}
		设备地址：{{keyword1.DATA}}
		告警时间：{{keyword2.DATA}}
		告警级别：{{keyword3.DATA}}
		告警内容：{{keyword4.DATA}}
		{{remark.DATA}}
	 */
}

func (ctl *WechatController) TemplateMessageAlarm(alarm *models.BoilerAlarm) (*template.TemplateMessage2, error) {
	var tempId string
	if conf.IsRelease {
		tempId = "7SGwryjW4LQYj2b9Q6CY0yeVOdGIbzkDQqsGUZjewgs"
	} else {
		tempId = "_r4xPpT8qqZqSNrOBFKn8w8J92oxo1ayf4rxYlULEpI"
	}

	var tempMsg template.TemplateMessage2
	//tempMsg.ToUser = openId
	tempMsg.TemplateId = tempId
	tempMsg.TopColor = "#999"
	tempMsg.MiniProgram = template.TemplateMiniProgram{
		AppId: mini.AppId,
		PagePath: "page/monitor/index",
	}
	type tData struct {
		First		template.DataItem	`json:"first"`
		Keyword1	template.DataItem	`json:"keyword1"`
		Keyword2	template.DataItem	`json:"keyword2"`
		Keyword3	template.DataItem	`json:"keyword3"`
		Keyword4	template.DataItem	`json:"keyword4"`
		Keyword5	template.DataItem	`json:"keyword5"`
		Remark		template.DataItem	`json:"remark"`
	}
	data := tData{}
	data.First = template.DataItem{
		Value: "您有一台锅炉发生告警！",
		Color: "#333333",
	}
	data.Keyword1 = template.DataItem{
		Value: alarm.Boiler.Name,
		Color: "#333333",
	}
	address := ""
	if alarm.Boiler.Address != nil { address = alarm.Boiler.Address.Address }
	data.Keyword2 = template.DataItem{
		Value: address,
		Color: "#333333",
	}
	data.Keyword3 = template.DataItem{
		Value: alarm.Description,
		Color: "#333333",
	}
	levelStrs := []string{"低", "中", "高"}
	priority := int(alarm.Priority)
	if priority >= len(levelStrs) { priority = len(levelStrs) - 1 }
	data.Keyword4 = template.DataItem{
		Value: levelStrs[priority],
		Color: "#333333",
	}
	data.Keyword5 = template.DataItem{
		Value: alarm.StartDate.Format("2006-01-02 15:04"),
		Color: "#333333",
	}
	data.Remark = template.DataItem{
		Value: "请登录Web端平台或在线小程序查看详情。",
		Color: "#333333",
	}

	tempMsg.Data = data

	return &tempMsg, nil
}

func (ctl *WechatController) SyncUserList() {
	list, err := user.List(wechatClient, "")
	if err != nil {
		goazure.Error("Get Wechat User List Error", err)
		return
	}

	goazure.Warn("WechatUserList: ", list)
	for _, openid := range list.Data.OpenIdList {
		info, err := user.Get(wechatClient, openid, user.LanguageZhCN)
		var third models.UserThird

		if err := dba.BoilerOrm.QueryTable("user_third").Filter("OpenId", info.OpenId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&third); err != nil {
			goazure.Warn("Read UserThird(OpenId) Error", err)

			third.Application = app
			third.Platform = app.Platform
			third.App = app.App
			third.Identity = app.Identity
			third.OpenId = info.OpenId
			third.UnionId = info.UnionId

			third.Name = info.Nickname
			third.Sex = info.Sex
			third.Province = info.Province
			third.City = info.City
			third.Country = info.Country
			third.HeadImageUrl = info.HeadImageURL

			var ath models.UserThird
			if err := dba.BoilerOrm.QueryTable("user_third").Filter("UnionId", info.UnionId).Filter("User__isnull", false).Filter("IsDeleted", false).One(&ath); err != nil {
				goazure.Warn("Read UserThird(UnionId) Error", err)
			} else {
				goazure.Info("Get UserThird(UnionId):", ath)
				third.User = ath.User
				//ctl.LoginWeixinWebDone(&third)
			}

			if err := DataCtl.AddData(&third, true, "OpenId"); err != nil {
				goazure.Error("Third Record Add Error:", err, third)
			}

		} else {
			goazure.Warn("Get UserThird(OpenId):", err, third, info)
			third.Name = info.Nickname
			third.UnionId = info.UnionId

			DataCtl.UpdateData(&third)
			//ctl.LoginWeixinWebDone(&third)
		}

		if 	err != nil {
			goazure.Error(info, "|UnionId:", info.UnionId, "|Errors:", err)
		}
	}
}

func (ctl *WechatController) SyncMenu() {
	m := menu.Menu{}
	/*
	btnMini := menu.Button{}
	btnMini.Name = "在线小程序"
	btnMini.SubButtons = append(btnMini.SubButtons, menu.Button{
		Type: menu.ButtonTypeMiniProgram,
		Name: "账号·绑定",
		URL: "https://mp.weixin.qq.com",
		AppId: mini.AppId,
		PagePath: "page/user/index",
	})
	btnMini.SubButtons = append(btnMini.SubButtons, menu.Button{
		Type: menu.ButtonTypeMiniProgram,
		Name: "主监控台",
		URL: "https://mp.weixin.qq.com",
		AppId: mini.AppId,
		PagePath: "page/monitor/index",
	})
	btnWeb := menu.Button{
		Type: menu.ButtonTypeView,
		Name: "节能平台",
		URL: "https://www.holderboiler.com/",
	}
	btnLocation := menu.Button{
		Type: menu.ButtonTypeLocationSelect,
		Name: "确认定位",
		Key: "rselfmenu_2_0",
	}

	m.Buttons = append(m.Buttons, btnMini)
	m.Buttons = append(m.Buttons, btnWeb)
	m.Buttons = append(m.Buttons, btnLocation)
	*/
	btnMonitor := menu.Button{
		Type: menu.ButtonTypeMiniProgram,
		Name: "主监控台",
		URL: "https://mp.weixin.qq.com",
		AppId: mini.AppId,
		PagePath: "page/monitor/index",
	}
	btnAccount := menu.Button{
		Type: menu.ButtonTypeMiniProgram,
		Name: "账号·绑定",
		URL: "https://mp.weixin.qq.com",
		AppId: mini.AppId,
		PagePath: "page/user/index",
	}
	btnManual := menu.Button{
		Type: menu.ButtonTypeView,
		Name: "使用手册",
		URL: "http://www.hold2025.com/h5/h.html",
	}

	m.Buttons = append(m.Buttons, btnMonitor)
	m.Buttons = append(m.Buttons, btnAccount)
	m.Buttons = append(m.Buttons, btnManual)

	if err := menu.Create(wechatClient, &m); err != nil {
		goazure.Error("Create Menu Error:", err)
	}
}

func (ctl *WechatController) WXCallbackHandler() {
	w := ctl.Ctx.ResponseWriter
	r := ctl.Ctx.Request

	//domain := ctl.Ctx.Input.Domain()
	//if app == nil || app.Domain != domain {
	//	app = &models.Application{ Domain: domain, App: "service" }
	//	if err := DataCtl.ReadData(app, "Domain", "App"); err != nil {
	//		goazure.Error("Read AppInfo Error:", app)
	//	}
	//
	//	msgServer = core.NewServer(app.OriginId, app.AppId, app.ApiToken, app.AesKey, msgHandler, nil)
	//
	//	accessTokenServer = core.NewDefaultAccessTokenServer(app.AppId, app.AppSecret, nil)
	//	wechatClient = core.NewClient(accessTokenServer, nil)
	//
	//	//fmt.Println(base.GetCallbackIP(wechatClient))
	//
	//	ctl.SyncUserList()
	//	ctl.SyncMenu()
	//}

	msgServer.ServeHTTP(w, r, nil)
}


//TODO: ???
func (ctl *WechatController) Get() {
	fmt.Println(ctl.Ctx.Input.URI(), ctl.Ctx.Input.URL())
	domain := ctl.Ctx.Input.Domain()
	path := strings.Replace(ctl.Ctx.Input.URL(), "/", "", -1)
	//Register
	//map[signature:[3f7439732a0e4507b4200ad00b966aaf9a388db2] echostr:[6951275124616035966] timestamp:[1488739904] nonce:[2125553201]]
	app := models.Application{ Domain: domain, Path: path}
	if err := DataCtl.ReadData(&app, "Domain", "Path"); err != nil {
		goazure.Error("Read App Error:", err)
	}
	signature := ctl.Input().Get("signature")
	echostr := ctl.Input().Get("echostr")
	timestamp := ctl.Input().Get("timestamp")
	nonce := ctl.Input().Get("nonce")

	goazure.Info(signature)
	goazure.Info(timestamp)
	goazure.Info(nonce)
	goazure.Info(echostr)

	goazure.Info(Signature(&app, timestamp, nonce))
	if Signature(&app, timestamp, nonce) == signature {
		ctl.Ctx.WriteString(echostr)
	} else {
		ctl.Ctx.WriteString("")
	}

}

func (ctl *MainController) Post() {
	fmt.Println(ctl.Ctx.Input.URI(), ctl.Ctx.Input.URL(), ctl.Ctx.Input.UserAgent())

	body, err := ioutil.ReadAll(ctl.Ctx.Request.Body)
	if err != nil {
		goazure.Error(err)
		ctl.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	goazure.Info(string(body))
	var wreq *Request
	if wreq, err = DecodeRequest(body); err != nil {
		goazure.Error(err)
		ctl.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	goazure.Info(wreq.Content)
	wresp, err := dealwith(wreq)
	if err != nil {
		goazure.Error(err)
		ctl.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	data, err := wresp.Encode()
	if err != nil {
		goazure.Error(err)
		ctl.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	goazure.Warning(data)
	goazure.Warning(string(data))

	ctl.Ctx.Output.SetStatus(200)
	ctl.Ctx.Output.Body(data)
	//this.Data["xml"] = wresp
	//this.ServeXML()
	return
}

func (ctl *WechatController) WechatPublicRegister()  {

}

func dealwith(req *Request) (resp *Response, err error) {
	resp = NewResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text
	goazure.Info(req.MsgType)
	goazure.Info(req.Content)
	if req.MsgType == Text {
		if strings.Trim(strings.ToLower(req.Content), " ") == "help" || req.Content == "Hello2BizUser" || req.Content == "subscribe" {
			resp.Content = "目前支持包的使用说明及例子的说明，这些例子和说明来自于github.com/AzureTech/gopkg，例如如果你想查询strings有多少函数，你可以发送：strings，你想查询strings.ToLower函数，那么请发送：strings.ToLower"
			return resp, nil
		}
		strs := strings.Split(req.Content, ".")
		var resurl string
		var a item
		if len(strs) == 1 {
			resurl = "https://raw.github.com/AzureTech/gopkg/master/" + strings.Trim(strings.ToLower(strs[0]), " ") + "/README.md"
			a.Url = "https://github.com/AzureTech/gopkg/tree/master/" + strings.Trim(strings.ToLower(strs[0]), " ") + "/README.md"
		} else {
			var other []string
			for k, v := range strs {
				if k < (len(strs) - 1) {
					other = append(other, strings.Trim(strings.ToLower(v), " "))
				} else {
					other = append(other, strings.Trim(strings.Title(v), " "))
				}
			}
			resurl = "https://raw.github.com/AzureTech/gopkg/master/" + strings.Join(other, "/") + ".md"
			a.Url = "https://github.com/AzureTech/gopkg/tree/master/" + strings.Join(other, "/") + ".md"
		}
		goazure.Info(resurl)
		rsp, err := http.Get(resurl)
		if err != nil {
			resp.Content = "不存在该包内容"
			return nil, err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode == 404 {
			resp.Content = "找不到你要查询的包:" + req.Content
			return resp, nil
		}
		resp.MsgType = News
		resp.ArticleCount = 1
		body, err := ioutil.ReadAll(rsp.Body)
		goazure.Info(string(body))
		a.Description = string(body)
		a.Title = req.Content
		a.PicUrl = "http://bbs.gocn.im/static/image/common/logo.png"
		resp.Articles = append(resp.Articles, &a)
		resp.FuncFlag = 1
	} else {
		resp.Content = "暂时还不支持其他的类型"
	}
	return resp, nil
}

func Signature(app *models.Application, timestamp string, nonce string) string {
	strs := sort.StringSlice{app.ApiToken, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func DecodeRequest(data []byte) (req *Request, err error) {
	req = &Request{}
	if err = xml.Unmarshal(data, req); err != nil {
		goazure.Error(err)
		return
	}
	req.CreateTime *= time.Second
	return
}

func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func (resp Response) Encode() (data []byte, err error) {
	resp.CreateTime = time.Duration(time.Now().Unix())
	data, err = xml.Marshal(resp)
	return
}

type WechatServer struct{}
func (s WechatServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Wechat Server Starting...")
	if req.Method == "GET" {
		query := req.URL.Query()
		fmt.Printf("request GET Body: %v\n\n", query)
		//resp.Write([]byte(echostr))
	}

	if req.Method == "POST" {
		if req.Body == nil {}
		fmt.Printf("Request URL: %v\n\n", req.URL)
		fmt.Printf("Request Header: %v\n\n", req.Header)
		body, _ := ioutil.ReadAll(req.Body)
		req.Body.Close()
		fmt.Printf("Request POST Body: %v\n\n", string(body))

		var msg WeixinMessage
		xml.Unmarshal(body, &msg)

		switch msg.MsgType {
		case "event":
			if msg.Event == "subscribe" {
				urlStr := "https://www.jiananhuaxia.com/wechat/message/get"
				msgJSON, err := json.Marshal(msg)
				if err != nil {
					fmt.Printf("json.Marshal error: %v\n", err)
				}
				fmt.Printf("msgJSON: %v", string(msgJSON))
				request, err := http.NewRequest("POST", urlStr, nil)
				if err != nil {
					fmt.Printf("request error: %v\n", err)
				}

				client := &http.Client{}
				response, _ := client.Do(request)
				fmt.Printf(response.Status)
			}
		case "text":

		}

		fmt.Printf("xml.Unmarshal\n %v\n", msg)

	}

	/*
	<xml>
	<ToUserName><![CDATA[toUser]]></ToUserName>
	<FromUserName><![CDATA[fromUser]]></FromUserName>
	<CreateTime>1348831860</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	<Content><![CDATA[this is a test]]></Content>
	<MsgId>1234567890123456</MsgId>
	</xml>
	*/

	/*
	<xml>
	<ToUserName><![CDATA[gh_904658786b9c]]></ToUserName>
	<FromUserName><![CDATA[o-yvHjv2wnEi_jW_Xdy2O98dDUCs]]></FromUserName>
	<CreateTime>1447848738</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[subscribe]]></Event>
	<EventKey><![CDATA[]]></EventKey>
	<Encrypt><![CDATA[V5fLp3FGB6NQ/P3b5h+/+UY6Ef2pnARumyg6oh9xymKa1MRlTCfa57E+LnAZH2kLKeFBWlbfwBX8RfhEpuM54goJ8Q5gjySDaXtK2zdrMK/Qw4q1Y5bpOWa7x0E0V/L/QvLByGwGFpZyPZFMRNZPq/dDZ1hG6Iuk+Vau1mS4+AVJtknEjgq4BPrnOiSiMLSMUDiKWh7Kc9YRM7XdI6F8un7AR0NcYFugEb52/68zz2AqyhYxOSkKyrheMqw1e9fi/wQx2/7wOBE2ZyDd6iZZCIiJR8g4vQ5iq3tUPERfbUGaUwu0MjZWKdYrC7/0NxMYHwYF/xbH0XLhGAU0iF5CpoIFhzSA6I8OE1kgfKUa+2s+F9mbCEzcMWTRMTCi+dDwApg/6Fay6kX3P59rGSgkNnODjnz1M6Fran7xhSkdYNc=]]></Encrypt>
	</xml>
	*/

	fmt.Printf("request: %v\n\n", req)
	//resp.Write([]byte("<h1>Hello, 世界</h1>\n<p>Behold my Go web app.</p>"))
}

func GetAccessToken() (app *models.Application, acc_token string, err error) {
	u, _ := url.Parse("https://api.weixin.qq.com/")
	u.Path = token_path
	q := u.Query()
	q.Set("grant_type", "client_credential")
	q.Set("appid", app.AppId)
	q.Set("secret", app.AppSecret)
	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	resp, err_get := http.Get(u.String())
	if err_get != nil {
		err = err_get
	}
	body, err_read := ioutil.ReadAll(resp.Body)
	if err_read != nil {
		err = err_read
	}
	fmt.Printf("body: %v\n\n", string(body))

	type Access struct {
		AccessToken string `json:"access_token"`
		ExpiresIn int `json:"expires_in"`
	}
	var acc Access
	json.Unmarshal(body, &acc)

	acc_token = acc.AccessToken

	fmt.Printf("JSON: %v\n", acc)
	fmt.Printf("AccessToken: %v\n", access_token)

	return
}

func GetUserList() (err error) {
	u, _ := url.Parse(WEIXIN_API_BASE_URL)
	u.Path = userlist_path
	q := u.Query()
	q.Set("access_token", access_token)
	q.Set("next_openid", "")
	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	resp, err_get := http.Get(u.String())
	if err_get != nil {
		err = err_get
	}
	body, err_read := ioutil.ReadAll(resp.Body)
	if err_read != nil {
		err = err_read
	}
	fmt.Printf("body: %v\n\n", string(body))

	type UserList struct {
		Total      int `json:"total"`
		Count      int `json:"count"`
		NextOpenid string `json:"next_openid"`
		Data       struct {
				   Openid []string `json:"openid"`
			   } `json:"data"`
	}

	var usrlist UserList
	json.Unmarshal(body, &usrlist)

	infoURL, _ := url.Parse("https://api.weixin.qq.com/")
	infoURL.Path = userinfo_path

	go func() {
		for i, openid := range usrlist.Data.Openid {
			q := infoURL.Query()
			q.Set("access_token", access_token)
			q.Set("openid", openid)
			infoURL.RawQuery = q.Encode()
			//fmt.Println(u.String())
			//fmt.Printf("openid: %v\n", openid)
			respInfo, err_get := http.Get(infoURL.String())
			if err_get != nil {
				err = err_get
			}
			body, err_read := ioutil.ReadAll(respInfo.Body)
			if err_read != nil {
				err = err_read
			}
			fmt.Printf("userInfo[%d]: %v\n\n", i, string(body))

			var usrInfo WeixinUserInfo
			json.Unmarshal(body, &usrInfo)

		}
	}()

	//fmt.Printf("JSON: %v\n", acc)
	//fmt.Printf("AccessToken: %v\n", access_token)

	return
}

const (
	WEIXIN_API_BASE_URL string = "https://api.weixin.qq.com/"
	WEIXIN_OPEN_BASE_URL string = "https://open.weixin.qq.com/"
)

const (
	TOKEN    = "AzureTechfromqqweixin"
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"
)

var token_path string = "/cgi-bin/token"
var userlist_path string = "/cgi-bin/user/get"
var userinfo_path string = "/cgi-bin/user/info"

var access_token string