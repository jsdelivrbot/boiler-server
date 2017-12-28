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

var TerminalCtrl *TerminalController

func (ctl *TerminalController) TerminalList() {
	usr := ctl.GetCurrentUser()

	var terminals []*models.Terminal

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

	var uList 		[]string
	var tBoilers 	[]*models.Boiler
	for _, ter := range terminals {
		uList = append(uList, ter.Uid)
		//aTerminals = append(aTerminals, ter)
	}

	if num, err := dba.BoilerOrm.QueryTable("boiler").Filter("Terminal__Uid__in", uList).OrderBy("TerminalSetId").All(&tBoilers); err != nil || num == 0 {
		goazure.Error("Read Terminal Boilers Error:", err, num)
	}

	for _, ter := range terminals {
		for _, b := range tBoilers {
			if b.Terminal.Uid == ter.Uid {
				ter.Boilers = append(ter.Boilers, b)
			}
		}
	}

	//se := ctl.GetSession(SESSION_CURRENT_USER)
	//fmt.Println("\nBoiler Get CurrentUser: ", usr, " | ", se);

	ctl.Data["json"] = terminals
	ctl.ServeJSON()
}

type aTerminal struct {
	Uid				string		`json:"uid"`
	TerminalCode 	string		`json:"code"`
	Name			string		`json:"name"`
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
}

func (ctl *TerminalController) TerminalConfig() {

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

	goazure.Debug("TerminalCode:", ter.TerminalCode)
	if code, err := strconv.ParseInt(ter.TerminalCode, 10, 32); err == nil && code > 0  {
		terminal.TerminalCode = code
	}

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
		fmt.Println("Read Exist Comment Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Read Exist AlarmRule Parameter Error:"))
		return
	}

	terminal.IsDeleted = true
	terminal.TerminalCode = terminal.TerminalCode * 10000 + 99

	if err := DataCtl.UpdateData(&terminal); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Delete Terminal Error!"))
		return
	}

	fmt.Println("\nDeleted Terminal:", terminal, ter)
}