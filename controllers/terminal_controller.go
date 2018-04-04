package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"

	"strconv"
	"fmt"

	"encoding/json"
)

type TerminalController struct {
	MainController
}

var TermCtl *TerminalController = &TerminalController{}

func (ctl *TerminalController) TerminalList() {
	usr := ctl.GetCurrentUser()

	var terminals []*models.Terminal
	//查询终端
	qs := dba.BoilerOrm.QueryTable("terminal")
	qs = qs.RelatedSel("Organization__Type")
	if usr.IsCommonUser() ||
		usr.Status == models.USER_STATUS_INACTIVE || usr.Status == models.USER_STATUS_NEW {
		qs = qs.Filter("IsDemo", true)
	} else if usr.IsOrganizationUser() {
		qs = qs.Filter("Organization__Uid", usr.Organization.Uid)
	}
	if num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "TerminalCode").All(&terminals); err != nil || num == 0 {
		goazure.Error("Read Terminal List Error:", err, num)
	} else {
		goazure.Info("Returned Terminals RowNum:", num)
	}
	//查询绑定锅炉的终端
	var combines []*models.BoilerTerminalCombined
	if  num, err := dba.BoilerOrm.QueryTable("boiler_terminal_combined").RelatedSel("Boiler").
		OrderBy("TerminalSetId").
		All(&combines); err != nil {
		goazure.Error("Get TerminalCombined Error:", err, num)
	}
	//查询终端版本号
	verSql:="select sn,ver,update_time from issued_version"
	var vi []VersionIssued
	termIssued:=make([]models.TerminalIssued,len(terminals))
	if _,error:=dba.BoilerOrm.Raw(verSql).QueryRows(&vi);error!=nil{
		goazure.Error("Select issued_version Error")
	}
	fmt.Println("versionIssued:",vi)
    //有绑定的将绑定的锅炉加进去
	for t, ter := range terminals {
		fmt.Println("ttt:",t)
		termIssued[t].Terminal=ter
		for _, cb := range combines {
			if ter.Uid == cb.Terminal.Uid {
				cb.Boiler.TerminalSetId = cb.TerminalSetId
				ter.Boilers = append(ter.Boilers, cb.Boiler)
			}
		}
		for _,v := range vi {
			i,err:=strconv.ParseInt(v.Sn,10,64)
			if err!=nil{
				goazure.Error("ParseInt Error")
			}
			if i == ter.TerminalCode {
				termIssued[t].Ver = v.Ver
				termIssued[t].UpdateTime = v.UpdateTime
			}
		}
	}

	//se := ctl.GetSession(SESSION_CURRENT_USER)
	//fmt.Println("\nBoiler Get CurrentUser: ", usr, " | ", se);

	ctl.Data["json"] = termIssued
	ctl.ServeJSON()
}

type aTerminal struct {
	Uid				string		`json:"uid"`
	TerminalCode 	string		`json:"code"`
	Name			string		`json:"name"`
	OrganizationId	string		`json:"org_id"`
	
	SimNumber		string		`json:"sim_number"`
	LocalIp			string		`json:"ip"`
	UploadFlag		bool		`json:"upload_flag"`
	UploadPeriod	int64		`json:"upload_period"`
	
	Description		string		`json:"description"`
}

//TODO: need Reset Operations
func (ctl *TerminalController) TerminalReset() {
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		fmt.Println("Only Admin Access")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Comment Error!"))
		return
	}

	var ter 		aTerminal
	var terminal 	models.Terminal

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &ter); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}

	terminal.Uid = ter.Uid
	if err := DataCtl.ReadData(&terminal); err != nil {
		e := fmt.Sprintln("Read Terminal Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
}

//TODO: need Reset Operations
func (ctl *TerminalController) TerminalConfig() {
	/*
	usr := ctl.GetCurrentUser()
	if !usr.IsAdmin() {
		fmt.Println("Only Admin Access")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Comment Error!"))
		return
	}
	*/

	var ter aTerminal
	terminal := models.Terminal{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &ter); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}

	terminal.Uid = ter.Uid
	if err := DataCtl.ReadData(&terminal); err != nil {
		e := fmt.Sprintln("Read Terminal Error:", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	SocketCtrl.SocketClientMessageSend(ter.TerminalCode, ter.UploadFlag, int(ter.UploadPeriod))
}

func (ctl *TerminalController) TerminalUpdate() {
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		goazure.Info("Only Admin Access")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Comment Error!"))
		return
	}

	var ter aTerminal
	terminal := models.Terminal{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &ter); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Terminal Error", err)
		return
	}

	if len(ter.Uid) > 0 {
		terminal.Uid = ter.Uid
		if err := DataCtl.ReadData(&terminal); err != nil {
			goazure.Warn("Read Terminal Failed!", err)
		}
		//raw := fmt.Sprintf("SELECT * FROM terminal WHERE uid = '%s'", ter.Uid)
		//if err := dba.BoilerOrm.Raw(raw).QueryRow(&terminal); err != nil  {
		//	goazure.Warn("Raw Termial Error:", raw, "\n", err)
		//}
		//goazure.Info("Terminal:", terminal)
	}

	if len(ter.Name) > 0 { terminal.Name = ter.Name }

	if len(ter.TerminalCode) > 0 {
		goazure.Debug("TerminalCode:", ter.TerminalCode)
		if code, err := strconv.ParseInt(ter.TerminalCode, 10, 64); err == nil && code > 0 {
			if exist := dba.BoilerOrm.QueryTable("terminal").Filter("TerminalCode", ter.TerminalCode).Exist(); exist && terminal.TerminalCode != code {
				e := fmt.Sprintln("终端编码不可重复！ErrorCode:", "Dup Code")
				goazure.Error(e)
				ctl.Ctx.Output.SetStatus(400)
				ctl.Ctx.Output.Body([]byte(e))
				return
			}

			terminal.TerminalCode = code
		} else {
			e := fmt.Sprintln("终端编码错误！ErrorCode:", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
	}

	org := OrgCtrl.organizationWithUid(ter.OrganizationId)
	terminal.Organization = org

	terminal.SimNumber = ter.SimNumber
	terminal.LocalIp = ter.LocalIp
	terminal.UploadFlag = ter.UploadFlag
	terminal.UploadPeriod = ter.UploadPeriod

	terminal.IsDeleted = false

	if err := DataCtl.AddData(&terminal, true); err != nil {
		e := fmt.Sprintln("Insert/Update Terminal Error!", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	fmt.Println("\nUpdated Terminal:", terminal, ter)
}

func (ctl *TerminalController) TerminalDelete() {
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		fmt.Println("Only Admin Access")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Comment Error!"))
		return
	}

	var ter aTerminal
	terminal := models.Terminal{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &ter); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Terminal Error", err)
		return
	}

	terminal.Uid = ter.Uid
	if err := DataCtl.ReadData(&terminal); err != nil {
		e := fmt.Sprintln("Read Exist Terminal For Delete Error!", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	tCode := terminal.TerminalCode * 1000000
	var tm models.Terminal
	if  err := dba.BoilerOrm.QueryTable("terminal").Filter("TerminalCode__contains", terminal.TerminalCode).Filter("TerminalCode__gte", tCode).OrderBy("-TerminalCode").One(&tm); err != nil {
		goazure.Info("Not Get Exist Dead Code!")
	} else {
		goazure.Info("Get Exist Dead Code:", tm)
		tCode = tm.TerminalCode + 1
	}

	terminal.TerminalCode = tCode
	terminal.IsDeleted = true

	if err := DataCtl.UpdateData(&terminal); err != nil {
		e := fmt.Sprintln("Delete Terminal Error!", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	goazure.Info("Deleted Terminal:", terminal, ter)
}

func (ctl *TerminalController) TerminalBindReload() {
	for _, b := range MainCtrl.Boilers {
		if b.Terminal != nil {
			var combined models.BoilerTerminalCombined
			combined.Boiler = b
			combined.Terminal = b.Terminal

			code, _ := strconv.ParseInt(b.TerminalCode, 10, 64)
			combined.TerminalCode = code
			combined.TerminalSetId = b.TerminalSetId

			if num, err := dba.BoilerOrm.InsertOrUpdate(&combined); err != nil {
				goazure.Error("M2M Terminals Combined Add Error:", err, num)
			}
		}
	}
}

func (ctl *TerminalController) TerminalTrim() {
	var boilers []*models.Boiler
	qs := dba.BoilerOrm.QueryTable("boiler")
	if num, err := qs.All(&boilers); err != nil || num == 0 {
		goazure.Error("Read Boilers For Terminal Update Error.", err)
	}

	for _, b := range boilers {
		qt := dba.BoilerOrm.QueryTable("terminal")
		var tid int64
		var err error
		if tid, err = strconv.ParseInt(b.TerminalCode, 10, 32); err != nil {
			goazure.Error("Parse TerminalCode Error.", err, b.TerminalCode, tid)
		}

		var ter models.Terminal
		if err := qt.Filter("TerminalCode", tid).One(&ter); err != nil {
			goazure.Error("Read Terminal Error.", err, tid)
			break
		}

		b.Terminal = &ter

		if err := DataCtl.UpdateData(b); err != nil {
			goazure.Error("Update Boiler's Terminal Error.", err)
		}
	}
}