package controllers

import (
	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/models/caches"

	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"

	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"reflect"
	"errors"
	"io/ioutil"
)

type BoilerController struct {
	MainController
}

var BlrCtl *BoilerController = &BoilerController{}

const boilerDefautlPath string = "models/properties/boiler_defaults/"

func init() {
	BlrCtl.MainController = *MainCtrl

	go BlrCtl.RefreshGlobalBoilerList()

	initConfig()
	initCalculateParameter()
}

func initConfig() {
	var fileName string = "dba/sql/init_boiler_config.sql"
	var sql string
	if buf, err := ioutil.ReadFile(fileName); err != nil {
		goazure.Error("read SQL File", fileName, "Error", err)
		return
	} else {
		sql = string(buf)
		goazure.Info("Read SQL:", sql)
	}

	if res, err := dba.BoilerOrm.Raw(sql).Exec(); err != nil {
		goazure.Error("Init Boiler Config Error:", err, res)
	}
}

func initCalculateParameter() {

}

func (ctl *BoilerController) BoilerCount() {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		if ctl.Input()["token"] == nil || len(ctl.Input()["token"]) == 0 {
			return
		}
		token := ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
			return
		}
	}

	qs := dba.BoilerOrm.QueryTable("boiler")
	if usr.IsCommonUser() ||
		usr.Status == models.USER_STATUS_INACTIVE || usr.Status == models.USER_STATUS_NEW {
		qs = qs.Filter("IsDemo", true)
	} else {
		qs = qs.Filter("IsDemo", false)
		if usr.IsOrganizationUser() {
			orgCond := orm.NewCondition().Or("Enterprise__Uid", usr.Organization.Uid).Or("Factory__Uid", usr.Organization.Uid).Or("Maintainer__Uid", usr.Organization.Uid).Or("Supervisor__Uid", usr.Organization.Uid)
			cond := orm.NewCondition().AndCond(orgCond)
			qs = qs.SetCond(cond)
		}
	}

	if num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").Count(); err != nil || num == 0 {
		goazure.Error("Get BoilerList Error:", num, err)
	} else {
		goazure.Info("Returned Rows Num:", num, err)

		ctl.Data["json"] = num
		ctl.ServeJSON()
	}
}

func (ctl *BoilerController) RefreshGlobalBoilerList() {
	BlrCtl.bWaitGroup.Add(1)
	//goazure.Error("")
	var boilers []*models.Boiler
	//var bMap []orm.Params
	qs := dba.BoilerOrm.QueryTable("boiler")
	qs = qs.RelatedSel("Form__Type").RelatedSel("Medium").RelatedSel("Usage").
		RelatedSel("Fuel__Type").RelatedSel("Template").
		RelatedSel("Factory").RelatedSel("Enterprise").RelatedSel("Maintainer").RelatedSel("Supervisor").
		RelatedSel("RegisterOrg").
		RelatedSel("Terminal").
		RelatedSel("Contact").
		RelatedSel("Address__Location")

	if  num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").
		All(&boilers);
		//All(&boilers,
		//	"Uid", "Name", "CreatedDate", "UpdatedDate", "IsDemo",
		//	"Form__Id", "Form__Name", "Form__Type__Id", "Form__Type__Name",
		//	"Medium__Id", "Medium__Name",
		//	"Usage__Id", "Usage__Name",
		//	"Fuel__Uid", "Fuel__Name", "Fuel__FuelId", "Fuel__Type__Id", "Fuel__Type__Name", "Fuel__QNetVArMin", "Fuel__QNetVArMax",
		//	"Template",
		//	"Factory__Uid", "Factory__Name", "Factory__IsSupervisor", "Factory__SuperOrganization",
		//	"Enterprise__Uid", "Enterprise__Name", "Enterprise__IsSupervisor", "Enterprise__SuperOrganization",
		//	"Installed__Uid", "Installed__Name", "Installed__IsSupervisor", "Installed__SuperOrganization");
		err != nil || num == 0 {
		goazure.Error("Get BoilerList Error:", num, err)
	} else {
		goazure.Info("Returned Boilers RowNum:", num)
	}

	var combines []*models.BoilerTerminalCombined
	if  num, err := dba.BoilerOrm.QueryTable("boiler_terminal_combined").RelatedSel("Terminal").OrderBy("TerminalSetId").
		All(&combines); err != nil {
		goazure.Error("Get TerminalCombined Error:", err, num)
	}

	var links []*models.BoilerOrganizationLinked
	if  num, err := dba.BoilerOrm.QueryTable("boiler_organization_linked").
		RelatedSel("Organization").
		OrderBy("Id").
		All(&links); err != nil {
		goazure.Error("Get OrganizationLinked Error:", err, num)
	}

	var bList []string
	for _, b := range boilers {
		bList = append(bList, b.Uid)
		//calc, _ := BlrCtl.GetBoilerCalculateParameter(b)
		//b.Calculate = append(b.Calculate, calc)
	}

	//var bCalc []*models.BoilerCalculateParameter
	//if num, err := dba.BoilerOrm.QueryTable("boiler_calculate_parameter").
	//	Filter("Boiler__Uid__in", bList).
	//	All(&bCalc); err != nil || num == 0 {
	//	e := fmt.Sprintln("Read Boilers CalcParam Error:", err, num)
	//	goazure.Error(e)
	//}

	var bStatus []orm.Params
	qi := dba.BoilerOrm.QueryTable("boiler_runtime_cache_instant")
	if len(bList) > 0 {
		qi = qi.Filter("Boiler__Uid__in", bList)
	}
	if 	num, err := qi.Filter("Parameter__Id", 1107).Filter("IsDeleted", false).
		Values(&bStatus, "Boiler__Uid", "Value"); err != nil {
		goazure.Warning("Read Boiler Burning Status Error!", err, num)
	}

	for _, b := range boilers {
		//for _, c := range bCalc {
		//	if c.Boiler.Uid == b.Uid {
		//		b.Calculate = append(b.Calculate, c)
		//		break
		//	}
		//}

		b.Status = make(map[string]int32)
		b.Status["IsBurning"] = 0
		for _, st := range bStatus {
			//goazure.Error("st:", st)
			if st["Boiler__Uid"].(string) == b.Uid && st["Value"].(float64) > 0 {
				b.Status["IsBurning"] = 1
				break
			}
		}

		for _, cb := range combines {
			if 	cb.Boiler.Uid == b.Uid {
				cb.Terminal.Remark = strconv.FormatInt(int64(cb.TerminalSetId), 10)
				cb.Terminal.TerminalSetId = cb.TerminalSetId
				b.TerminalsCombined = append(b.TerminalsCombined, cb.Terminal)
			}
		}

		for _, li := range links {
			if 	li.Boiler.Uid == b.Uid {
				b.OrganizationsLinked = append(b.OrganizationsLinked, li.Organization)
				//goazure.Warn("Boiler Linked", b, "\n", li.Organization)
				//panic(0)
			}
		}

		if b.InspectInnerDateNext.IsZero() { b.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30) }
		if b.InspectOuterDateNext.IsZero() { b.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30) }
		if b.InspectValveDateNext.IsZero() { b.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30) }
		if b.InspectGaugeDateNext.IsZero() { b.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30) }
		//goazure.Info("[",i,"]", b)
	}

	goazure.Info("MainCtrl.Boilers = boilers")

	MainCtrl.Boilers = boilers

	BlrCtl.bWaitGroup.Done()

	go RtmCtl.RefreshStatusRunningDuration(time.Now())
	go RtmCtl.RefreshBoilerRank(time.Now())

	go TermCtl.TerminalBindReload()
}

func (ctl *BoilerController) CurrentBoilerList(usr *models.User) ([]*models.Boiler, error) {
	BlrCtl.bWaitGroup.Wait()

	var boilers []*models.Boiler
	var err error
	if usr == nil || len(usr.Uid) <= 0 {
		err = errors.New("user can not be nil")
		return boilers, err
	}

	for _, b := range MainCtrl.Boilers {
		if ctl.IsBoilerBelongToUser(b, usr) {
			boilers = append(boilers, b)
		}
	}

	return boilers, nil
}

func (ctl *BoilerController) BoilerCalculateParameter() {
	usr := ctl.GetCurrentUser()
	if ctl.Input()["boiler"] == nil || len(ctl.Input()["boiler"]) == 0 {
		e := fmt.Sprintln("Unknown Boiler Id!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	uid := ctl.Input()["boiler"][0]

	var boiler models.Boiler
	boiler.Uid = uid
	if isBelong := ctl.IsBoilerBelongToUser(&boiler, usr); !isBelong {
		e := fmt.Sprintln("Boiler Permission Denied: ", uid)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	calc, _ := BlrCtl.GetBoilerCalculateParameter(&boiler)

	ctl.Data["json"] = calc
	ctl.ServeJSON()
}

func (ctl *BoilerController) GetBoilerCalculateParameter(boiler *models.Boiler) (*models.BoilerCalculateParameter, error) {
	var calc models.BoilerCalculateParameter
	if err := dba.BoilerOrm.QueryTable("boiler_calculate_parameter").Filter("Boiler__Uid", boiler.Uid).One(&calc); err != nil {
		e := fmt.Sprintln("Read Boiler CalcParam Error:", boiler, err)
		goazure.Error(e)

		return nil, errors.New(e)
	}

	return &calc, nil
}

//TODO: Failed
func (ctl *BoilerController) GetBoilerStatus(boiler *models.Boiler) (*caches.BoilerRuntimeCacheInstant, error) {
	var status caches.BoilerRuntimeCacheInstant
	if err := dba.BoilerOrm.QueryTable("boiler_runtime_cache_instant").Filter("Boiler__Uid", boiler.Uid).Filter("Parameter__Id", 1107).Filter("IsDeleted", false).One(&status); err != nil {
		e := fmt.Sprintln("Read Boiler Runtime Status Error:", boiler, err)
		goazure.Error(e)

		return nil, errors.New(e)
	}

	return &status, nil
}

func (ctl *BoilerController) IsBoilerBelongToUser(boiler *models.Boiler, usr *models.User) bool {
	if usr.IsAdmin() {
		return true
	}

	if usr.IsOrganizationUser() {
		return boiler.IsBelongToOrganization(usr.Organization)
	}

	if usr.IsCommonUser() {
		if boiler.IsDemo {
			return true
		} else {
			return false
		}
	}

	return true
}

func (ctl *BoilerController) BoilerList() ([]*models.Boiler, error) {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		if ctl.Input()["token"] == nil || len(ctl.Input()["token"]) == 0 {
			return nil, errors.New("there is no token")
		}
		token := ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))

			return nil, err
		}
	}

	var aBoilers, boilers []*models.Boiler
	var err error
	if aBoilers, err = ctl.CurrentBoilerList(usr); err != nil {
		goazure.Error("Get CurrentBoilerList Error:", err)
	}

	boilers = aBoilers

	//goazure.Warning("Get CurrentBoilers:", boilers)

	if ctl.Input()["boiler"] != nil && len(ctl.Input()["boiler"]) > 0 {
		//goazure.Warn("BoilerUid:", ctl.Input()["boiler"][0], len(ctl.Input()["boiler"][0]))
		if len(ctl.Input()["boiler"][0]) == 36 {
			//goazure.Warn("Select Boiler by Uid:", ctl.Input()["boiler"][0], len(ctl.Input()["boiler"][0]))
			boilers = []*models.Boiler{}
			for _, ab := range aBoilers {
				if ab.Uid == ctl.Input()["boiler"][0] {
					boilers = append(boilers, ab)
					break
				}
			}
		}
	} else {
		boilers = aBoilers
	}

	ctl.Data["json"] = boilers
	ctl.ServeJSON()

	return boilers, nil
}

func (ctl *BoilerController) BoilerFormList() {
	var forms []orm.Params

	qs := dba.BoilerOrm.QueryTable("boiler_type_form")
	num, err := qs.RelatedSel("Type").Filter("IsDeleted", false).OrderBy("Id").Values(&forms, "Id", "Name", "NameEn")
	fmt.Printf("Returned Rows Num: %d, %s", num, err)

	ctl.Data["json"] = forms
	ctl.ServeJSON()
}

func (ctl *BoilerController) BoilerMediumList() {
	var mediums []orm.Params
	qs := dba.BoilerOrm.QueryTable("boiler_medium")

	num, err := qs.Filter("IsDeleted", false).OrderBy("Id").Values(&mediums, "Id", "Name", "NameEn")
	fmt.Printf("Returned Rows Num: %d, %s", num, err)

	ctl.Data["json"] = mediums
	ctl.ServeJSON()
}

func (ctl *BoilerController) BoilerListWeixin() {
	usr := ctl.GetCurrentUser()
	if usr == nil {
		goazure.Info("Params:", ctl.Input())
		token := ctl.Input()["token"][0]

		var err error
		usr, err = ctl.GetCurrentUserWithToken(token)
		if err != nil {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))

			return
		}
	}

	var boilers []orm.Params
	qs := dba.BoilerOrm.QueryTable("boiler")
	qs = qs.RelatedSel("Form__Type").RelatedSel("Medium").RelatedSel("Usage").
		RelatedSel("Fuel__Type").RelatedSel("Template").
		RelatedSel("Factory").RelatedSel("Enterprise").RelatedSel("Maintainer").RelatedSel("Supervisor").
		RelatedSel("Address__Location")
	if usr.IsCommonUser() ||
		usr.Status == models.USER_STATUS_INACTIVE || usr.Status == models.USER_STATUS_NEW {
		qs = qs.Filter("IsDemo", true)
	} else {
		qs = qs.Filter("IsDemo", false)
		if usr.IsOrganizationUser() {
			orgCond := orm.NewCondition().Or("Enterprise__Uid", usr.Organization.Uid).Or("Factory__Uid", usr.Organization.Uid).Or("Maintainer__Uid", usr.Organization.Uid)
			cond := orm.NewCondition().AndCond(orgCond)
			qs = qs.SetCond(cond)
		}
	}
	num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").
		Values(&boilers,
		"Uid", "Name",
		"Enterprise__Name", "Factory__Name", "Maintainer__Name", "Supervisor__Name",
		"CertificateNumber", "DeviceCode", "ModelCode", "RegisterCode",
		"Form__Name", "Form__Id",
		"Fuel__Name", "Fuel__Type__Name", "Fuel__Type__Id",
		"Medium__Id", "Medium__Name",
		"Address__Location__LocationId", "Address__Location__Name", "Address__Location__LocationName",
		"EvaporatingCapacity",
		"TerminalCode")

	for _, b := range boilers {
		b["IsBurning"] = ctl.IsBurning(b["Uid"].(string))
	}

	goazure.Info("Returned Rows Num:", num, err)

	ctl.Data["json"] = boilers
	ctl.ServeJSON()
}

func (ctl *BoilerController) BoilerMaintainList() {
	usr := ctl.GetCurrentUser()
	var boilerUid string
	if ctl.Input()["boiler"] != nil &&
		len(ctl.Input()["boiler"]) > 0 &&
		ctl.Input()["boiler"][0] != "undefined" {
		boilerUid = ctl.Input()["boiler"][0]
	}

	var maintains []models.BoilerMaintenance
	var boilers []models.Boiler
	if !usr.IsAdmin() || len(boilerUid) > 0 {
		qs := dba.BoilerOrm.QueryTable("boiler")
		qs = qs.RelatedSel("Form__Type").RelatedSel("Medium").RelatedSel("Usage").
			RelatedSel("Fuel__Type").RelatedSel("Template").
			RelatedSel("Factory").RelatedSel("Enterprise").RelatedSel("Maintainer").
			RelatedSel("RegisterOrg").
			RelatedSel("Terminal").
			RelatedSel("Contact").
			RelatedSel("Address__Location")
		if usr.IsCommonUser() ||
			usr.Status == models.USER_STATUS_INACTIVE ||
			usr.Status == models.USER_STATUS_NEW {
			qs = qs.Filter("IsDemo", true)
		} else if usr.IsOrganizationUser() {
			orgCond := orm.NewCondition().Or("Enterprise__Uid", usr.Organization.Uid).Or("Factory__Uid", usr.Organization.Uid).Or("Maintainer__Uid", usr.Organization.Uid)
			cond := orm.NewCondition().AndCond(orgCond)
			qs = qs.SetCond(cond)
		}
		if len(boilerUid) > 0 {
			qs = qs.Filter("Uid", boilerUid)
		}

		if num, err := qs.Filter("IsDeleted", false).OrderBy("IsDemo", "Name").All(&boilers); num == 0 || err != nil {
			goazure.Error("Read BoilerList For Maintain Error: ", err, num)
		}
	}

	goazure.Warn("Boilers:", boilers)
	if !usr.IsAdmin() && len(boilers) <= 0 {
		e := fmt.Sprintln("No Admin && No Boiler Accessed")
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		goazure.Error(e)
		return
	}

	qm := dba.BoilerOrm.QueryTable("boiler_maintenance").RelatedSel()
	if !usr.IsAdmin() || len(boilers) > 0 {
		qm = qm.Filter("Boiler__in", boilers)
	}

	if num, err := qm.Filter("IsDeleted", false).OrderBy("-CreatedDate").All(&maintains); num == 0 || err != nil {
		goazure.Error("Read BoilerMaintainList Error:", num, err)
	}

	ctl.Data["json"] = maintains
	ctl.ServeJSON()
}


type Maintain struct {
	Uid           string		`json:"uid"`
	BoilerId      string		`json:"boiler_id"`
	InspectDate   time.Time		`json:"inspect_date"`
	Topic         string		`json:"topic"`
	Content       string		`json:"content"`
	Attachment    string 		`json:"attachment"`
	IsDone        bool 			`json:"is_done"`

	Burner        string		`json:"burner"`         // 燃烧器
	ImportGrate   string		`json:"import_grate"`   // 进料及炉排
	WaterSoftener string        `json:"water_softener"` // 软水器
	WaterPump     string        `json:"water_pump"`     // 水泵
	BoilerBody    string        `json:"boiler_body"`    // 锅炉本体
	EnergySaver   string        `json:"energy_saver"`   // 节能器
	AirPreHeater  string        `json:"air_pre_heater"` // 空预器
	DustCatcher   string        `json:"dust_catcher"`   // 除尘器
	DraughtFan    string        `json:"draught_fan"`    // 引风机
}

func (ctl *BoilerController) BoilerMaintainUpdate() {
	usr := ctl.GetCurrentUser()
	var mt Maintain
	maintain := models.BoilerMaintenance{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &mt); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Maintain Error", err)
		return
	}

	if len(mt.Uid) > 0 {
		fmt.Println(">>>>>>There is Maintain Uid:", mt.Uid)
		maintain.Uid = mt.Uid

		if err := DataCtl.ReadData(&maintain); err != nil {
			fmt.Println("Read Exist Maintain Error:", err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("Read Exist Maintain Error:"))
			return
		}
	}

	boiler := models.Boiler{}
	boiler.Uid = mt.BoilerId
	if err := DataCtl.ReadData(&boiler); err != nil {
		fmt.Println("Read Boiler Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update BoilerMaintain Error!"))
		return
	}

	maintain.Boiler = &boiler

	maintain.InspectDate = mt.InspectDate
	maintain.Name = mt.Topic
	maintain.Content = mt.Content
	maintain.Attachment = mt.Attachment
	maintain.IsDone = mt.IsDone

	maintain.Burner = mt.Burner
	maintain.ImportGrate = mt.ImportGrate
	maintain.WaterSoftener = mt.WaterSoftener
	maintain.WaterPump = mt.WaterPump
	maintain.BoilerBody = mt.BoilerBody
	maintain.EnergySaver = mt.EnergySaver
	maintain.AirPreHeater = mt.AirPreHeater
	maintain.DustCatcher = mt.DustCatcher
	maintain.DraughtFan = mt.DraughtFan

	if maintain.CreatedBy == nil {
		maintain.CreatedBy = usr
	}
	maintain.UpdatedBy = usr

	if err := DataCtl.AddData(&maintain, true); err != nil {
		fmt.Println("Add Maintain Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update Maintain Error!"))
		return
	}

	goazure.Info("Update Comment:", maintain, mt)
}

func (ctl *BoilerController) BoilerMaintainDelete() {
	usr := ctl.GetCurrentUser()
	var m Maintain
	var maintain models.BoilerMaintenance

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &m); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		goazure.Error("Unmarshal Maintain Error", err)
		return
	}

	if err := dba.BoilerOrm.QueryTable("boiler_maintenance").RelatedSel("CreatedBy__Role").Filter("Uid", m.Uid).One(&maintain); err != nil {
		e := fmt.Sprintf("Read Maintain for Delete Error: %v", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	author := maintain.CreatedBy

	if author.Uid == usr.Uid ||
		(usr.IsAdmin() && usr.Role.RoleId < author.Role.RoleId) ||
		(usr.IsOrganizationUser() && usr.Organization.Uid == author.Organization.Uid && usr.Role.RoleId == models.USER_ROLE_ORG_ADMIN) {
		if err := DataCtl.DeleteData(&maintain); err != nil {
			e := fmt.Sprintf("Delete Maintain Error: %v", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}
	}
}

//*** CONFIG START

type bConfig struct {
	Uid			string		`json:"uid"`
	Config		string		`json:"config"`
	Value		string		`json:"value"`
	State		bool		`json:"state"`

	Cascade		bool		`json:"cascade"`
}

func (ctl *BoilerController) GetBoilerConfig() {
	goazure.Info("Ready to GetBoilerConfig!")

	var cnf bConfig
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cnf); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerConfig JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	value, err := ctl.BoilerConfig(cnf.Uid, cnf.Config)
	if err != nil {
		e := fmt.Sprintln("Get BoilerConfig Error", cnf, err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	cnf.Value = strconv.FormatBool(value.(bool))

	ctl.Data["json"] = cnf
	ctl.ServeJSON()
}

func (ctl *BoilerController) SetBoilerConfig() {
	goazure.Info("Ready to Updated Boiler!")
	usr := ctl.GetCurrentUser()

	if usr.Role.RoleId > 2  {
		e := fmt.Sprintln("Permission Denied, Only SuperAdmin Access!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var cnf bConfig
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cnf); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerConfig JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var conf models.BoilerConfig
	qs := dba.BoilerOrm.QueryTable("boiler_config")
	if err := qs.Filter("Boiler__Uid", cnf.Uid).One(&conf); err != nil {
		goazure.Error("Read BoilerConfig Error", err, conf)
		return
	}

	value, err := strconv.ParseBool(cnf.Value)
	if err != nil {
		goazure.Error("Parse BoilerConfig Value Error", err)
	}

	reflect.ValueOf(&conf).Elem().FieldByName(cnf.Config).SetBool(value)

	if cnf.Config == "IsGenerateData" && !value && cnf.Cascade {
		//goazure.Warn("Ready to Cascade:", cnf)
		raw := "DELETE FROM	`boiler_runtime` " +
			fmt.Sprintf("WHERE `is_demo` = TRUE AND `boiler_id` = '%s';", cnf.Uid)

		if res, err := dba.BoilerOrm.Raw(raw).Exec(); err != nil {
			goazure.Error("Cascade Failed:", err, res)
		}
		/*
		else {
			row, e := res.RowsAffected()
			goazure.Warn("Cascaded:", cnf, "\n", raw, "\n", res, row, e)
		}
		*/
	}

	if err := DataCtl.AddData(&conf, true); err != nil {
		e := fmt.Sprintln("Update BoilerConfig JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
}

func (ctl *BoilerController) IsGenerateData(boilerUid string) bool {
	value, err := ctl.BoilerConfig(boilerUid, "IsGenerateData")
	if err != nil {
		goazure.Warn("Get BoilerConfig IsGenerateData Error:", err)
		return false
	}

	return value.(bool)
}

func (ctl *BoilerController) BoilerConfig(boilerUid string, config string) (interface{}, error) {
	if len(boilerUid) == 0 || len(config) == 0 {
		e := errors.New("Boiler && Config Can Not be Nil!")
		return nil, e
	}

	conf := models.BoilerConfig{}
	qs := dba.BoilerOrm.QueryTable("boiler_config")
	if err := qs.Filter("Boiler__Uid", boilerUid).One(&conf); err != nil {
		goazure.Error("Read BoilerConfig Error", err)
		return nil, err
	}

	value := reflect.ValueOf(&conf).Elem().FieldByName(config).Interface()
	//goazure.Info("BoilerConfig[", config,"]:", value)

	return value, nil
}

//***************************************** CONFIG END
func (ctl *BoilerController) BoilerIsOnline(){
	if ctl.Input()["boiler"] == nil || len(ctl.Input()["boiler"]) == 0 {
		e := fmt.Sprintln("there is no boiler!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	uid := ctl.Input()["boiler"][0]
	type IsOnline struct {
		IsOnline bool `json:"IsOnline"`
	}
	var sql="select t.is_online from boiler_terminal_combined as bt,terminal as t"+
		" where bt.terminal_code=t.terminal_code and bt.boiler_id=?"
	var status IsOnline
	err:=dba.BoilerOrm.Raw(sql,uid).QueryRow(&status)
	if err!=nil {
		goazure.Warning("Read Boiler online Status Error!",err)
	}
	fmt.Println(status)
	ctl.Data["json"]=status
	ctl.ServeJSON()
}

func (ctl *BoilerController) BoilerIsBurning() {
	//goazure.Info("Ready to BoilerIsBurning!")
	//goazure.Info("Params:", ctl.Input())
	if ctl.Input()["boiler"] == nil || len(ctl.Input()["boiler"]) == 0 {
		e := fmt.Sprintln("there is no boiler!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	uid := ctl.Input()["boiler"][0]

	isBurning := ctl.IsBurning(uid)
	cnf := bConfig{}
	cnf.Uid = uid
	cnf.Config = "IsBurning"
	cnf.Value = strconv.FormatBool(isBurning)

	goazure.Warn("BoilerIsBurning:", cnf)
	// panic(0)

	ctl.Data["json"] = cnf
	ctl.ServeJSON()
}

func (ctl *BoilerController) IsBurning(boilerUid string) bool {
	var status []orm.Params
	q := dba.BoilerOrm.QueryTable("boiler_runtime_cache_instant")
	q = q.Filter("Boiler__Uid", boilerUid).Filter("Parameter__Id", 1107)//.Filter("UpdatedDate__gt", time.Now().Add(time.Minute * -30))
	if num, err := q.Filter("IsDeleted", false).Values(&status, "Value"); err != nil || num == 0 {
		goazure.Warning("Read Boiler Burning Status Error!", err, num)
		return false
	}
	goazure.Info("Boiler[", "|", status, "]", boilerUid)
	if status[0]["Value"].(float64) <= 0 {
		return false
	}

	return true
}

func (ctl *BoilerController) BoilerHasSubscribed() {
	//goazure.Info("Ready to BoilerIsBurning!")
	var boiler *models.Boiler = &models.Boiler{}
	var usr *models.User = &models.User{}
	if ctl.Input()["boiler"] == nil || len(ctl.Input()["boiler"]) == 0 {
		e := fmt.Sprintln("there is no boiler!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
	bid := ctl.Input()["boiler"][0]
	boiler.Uid = bid

	if ctl.Input()["uid"] == nil || len(ctl.Input()["uid"]) == 0 {
		e := fmt.Sprintln("there is no usr!")
		goazure.Error(e)
		usr = ctl.GetCurrentUser()
	} else {
		uid := ctl.Input()["uid"][0]
		usr.Uid = uid
	}

	m2m := dba.BoilerOrm.QueryM2M(boiler, "Subscribers")
	hasSubscribed := m2m.Exist(usr)

	cnf := bConfig{}
	cnf.Uid = boiler.Uid
	cnf.Config = "HasSubscribed"
	cnf.Value = strconv.FormatBool(hasSubscribed)

	ctl.Data["json"] = cnf
	ctl.ServeJSON()
}

func (ctl *BoilerController) BoilerSetSubscribe() {
	goazure.Info("Ready to Set Boiler Subscribe!")
	boiler := &models.Boiler{}
	usr := ctl.GetCurrentUser()

	var cnf bConfig
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cnf); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerConfig JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	boiler.Uid = cnf.Uid
	m2m := dba.BoilerOrm.QueryM2M(boiler, "Subscribers")

	if m2m.Exist(usr) {
		goazure.Warn("Subscibe Already Exist!")
	}

	var err error
	var num int64
	if cnf.State {
		num, err = m2m.Add(usr)
	} else {
		num, err = m2m.Remove(usr)
	}

	if err != nil || num == 0 {
		e := fmt.Sprintln("Subcribe Error:", boiler, "\n", usr, "\n", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}
}

func (ctl *BoilerController) BoilerMessageSend() {
	goazure.Info("Ready to Set Boiler Subscribe!")
	boiler := &models.Boiler{}
	usr := ctl.GetCurrentUser()
	var u *models.User

	var cnf bConfig
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &cnf); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerConfig JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	qs := dba.BoilerOrm.QueryTable("boiler")
	qs = qs.RelatedSel("Form__Type").RelatedSel("Medium").RelatedSel("Usage").
		RelatedSel("Fuel__Type").RelatedSel("Template").
		RelatedSel("Factory").RelatedSel("Enterprise").RelatedSel("Maintainer").
		RelatedSel("RegisterOrg").
		RelatedSel("Terminal").
		RelatedSel("Contact").
		RelatedSel("Address__Location")
	qs = qs.Filter("Uid", cnf.Uid)
	if err := qs.Filter("IsDeleted", false).One(boiler); err != nil {
		goazure.Warn("Get Boiler Info For Test Message: ", err, "\n", boiler, "\n", cnf)
	} else {
		var users []*models.User
		raw := "SELECT 	`user`.* "
		raw += "FROM	`user`, `boiler_message_subscriber` AS `sub` "
		raw += "WHERE	`user`.`uid` = `sub`.`user_id` "
		raw += fmt.Sprintf("AND	`sub`.`boiler_id` = '%s' ", boiler.Uid)
		raw += fmt.Sprintf("AND	`user`.`uid` = '%s';", usr.Uid)

		if num, err := dba.BoilerOrm.Raw(raw).QueryRows(&users); err != nil || num == 0 {
			goazure.Error("Get Boiler Subscribers Error:", err, num)
		} else {
			u = users[0]

			var su models.UserThird
			qu := dba.BoilerOrm.QueryTable("user_third")
			qu = qu.Filter("User__Uid", u.Uid).Filter("App", "service").Filter("IsDeleted", false)
			if err := qu.One(&su); err != nil {
				goazure.Error("User", u.Name, "Is NOT Subscribed.")
			} else {
				//content := "锅炉测试消息：\n"
				//content += "送达" + u.Name + "OpenId:" + su.OpenId + "\n"
				//content += "锅炉名称：" + boiler.Name + "\n"
				//content += "使用企业：" + boiler.Enterprise.Name + "\n"
				//content += "制造厂家：" + boiler.Factory.Name + "\n"
				//content += "测试完毕!\n"
				//WXCtrl.SendText(su.OpenId, content, "")
				var alarm models.BoilerAlarm
				qa := dba.BoilerOrm.QueryTable("boiler_alarm")
				qa = qa.RelatedSel("Boiler__Address").Filter("Boiler__Uid", boiler.Uid)
				if err := qa.One(&alarm); err != nil {
					goazure.Error("Get Sample Alarm Error:", err)
				}

				tempMsg, _ := WxCtl.TemplateMessageAlarm(&alarm)
				goazure.Info("WXCtrl.SendTemplateMessage(su.OpenId, tempMsg)", su.OpenId, tempMsg)
				WxCtl.SendTemplateMessage(su.OpenId, tempMsg)
			}
		}
	}
}

//*** INFO UPDATE START
type BoilerInfo struct {
	Uid					string		`json:"uid"`
	Name				string		`json:"name"`
	RegisterCode		string		`json:"registerCode"`
	DeviceCode			string		`json:"deviceCode"`
	ModelCode			string		`json:"modelCode"`
	FactoryNumber		string		`json:"factoryNumber"`
	CertificateNumber	string		`json:"certificateNumber"`

	EvaporatingCapacity	int64		`json:"evaporatingCapacity"`

	MediumId			int64		`json:"mediumId"`
	FuelId				string		`json:"fuelId"`
	FormId				int64		`json:"formId"`

	EnterpriseId		string		`json:"enterpriseId"`
	FactoryId			string		`json:"factoryId"`
	MaintainerId		string		`json:"maintainerId"`
	SupervisorId		string		`json:"supervisorId"`
	
	Address				string		`json:"address"`
	LocationId			int64		`json:"location_id"`
	Longitude			float64		`json:"longitude"`
	Latitude			float64		`json:"latitude"`
	
	Contact				string		`json:"contact"`
	PhoneNumber			string		`json:"phoneNumber"`
	MobileNumber		string		`json:"mobileNumber"`
	Email				string		`json:"email"`

	Links				[]bLink		`json:"links"`
	
	InspectDate			struct
	{
		Inner			time.Time	`json:"inner"`
		Outer			time.Time	`json:"outer"`
		Valve			time.Time	`json:"valve"`
		Gauge			time.Time	`json:"gauge"`
	}								`json:"inspectDate"`
}

type bLink struct {
	Num		int64		`json:num`
	Type 	int32		`json:"type"`
	Uid		string		`json:"uid"`
}

func (ctl *BoilerController) BoilerUpdate() {
	var boiler 	*models.Boiler
	var err		error
	var scope string = "basic"
	if ctl.Input()["scope"] == nil || len(ctl.Input()["scope"]) == 0 {
		e := fmt.Sprintln("there is no specify scope here!")
		goazure.Error(e)
	} else {
		scope = ctl.Input()["scope"][0]
	}

	switch scope {
	case "location":
		boiler, err = ctl.BoilerUpdateLocation()
		if err != nil {
			goazure.Error(err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
		}
	case "maintain":
		boiler, err = ctl.BoilerUpdateMaintain()
		if err != nil {
			goazure.Error(err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
		}
	case "basic":
		fallthrough
	default:
		boiler, err = ctl.BoilerUpdateBasic()
		if err != nil {
			goazure.Error(err)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(err.Error()))
		}
	}

	go BlrCtl.RefreshGlobalBoilerList()

	ctl.Data["json"] = boiler
	ctl.ServeJSON()
}

func (ctl *BoilerController) BoilerUpdateBasic() (*models.Boiler, error) {
	goazure.Info("Ready to Updated Boiler!")
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		return nil, errors.New(e)
	}

	var info 	BoilerInfo
	var boiler 	models.Boiler

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &info); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerInfo JSON Error", err)
		return nil, errors.New(e)
	}

	if len(info.Uid) > 0 {
		boiler.Uid = info.Uid
		if err := DataCtl.ReadData(&boiler); err != nil {
			goazure.Warn("Read BoilerInfo Failed!", err)
		}
	} else {
		boiler.CreatedBy = usr
	}

	if len(info.Name) > 0 { boiler.Name = info.Name }
	if len(info.RegisterCode) > 0 { boiler.RegisterCode = info.RegisterCode }
	if len(info.DeviceCode) > 0 { boiler.DeviceCode = info.DeviceCode }
	if len(info.ModelCode) > 0 { boiler.ModelCode = info.ModelCode }
	if info.EvaporatingCapacity > 0 { boiler.EvaporatingCapacity = info.EvaporatingCapacity }

	if len(info.FactoryNumber) > 0 { boiler.FactoryNumber = info.FactoryNumber }
	if len(info.CertificateNumber) > 0 { boiler.CertificateNumber = info.CertificateNumber }

	var usage models.BoilerUsage
	var med models.BoilerMedium
	var fuel models.Fuel
	var form models.BoilerTypeForm
	usage.Id = 1
	med.Id = info.MediumId
	fuel.Uid = info.FuelId
	form.Id = info.FormId
	if err := DataCtl.ReadData(&usage); err == nil { boiler.Usage = &usage }
	if err := DataCtl.ReadData(&med); err == nil { boiler.Medium = &med }
	if err := DataCtl.ReadData(&fuel); err == nil { boiler.Fuel = &fuel }
	if err := DataCtl.ReadData(&form); err == nil { boiler.Form = &form }

	var enterprise, factory, maintainer, supervisor models.Organization
	enterprise.Uid = info.EnterpriseId
	factory.Uid = info.FactoryId
	maintainer.Uid = info.MaintainerId
	supervisor.Uid = info.SupervisorId
	if err := DataCtl.ReadData(&enterprise); err == nil { boiler.Enterprise = &enterprise }
	if err := DataCtl.ReadData(&factory); err == nil { boiler.Factory = &factory }
	if err := DataCtl.ReadData(&maintainer); err == nil { boiler.Maintainer = &maintainer }
	if err := DataCtl.ReadData(&supervisor); err == nil { boiler.Supervisor = &supervisor }

	if	num, err := dba.BoilerOrm.QueryTable("boiler_organization_linked").
		Filter("Boiler__Uid", boiler.Uid).Delete(); err != nil {
		goazure.Warn("Deleted Old Links Error:", err, num)
	}

	for _, li := range info.Links {
		var linked 	models.BoilerOrganizationLinked
		var og 		models.Organization
		var ogType	models.OrganizationType

		og.Uid = li.Uid
		ogType.TypeId = li.Type

		linked.Boiler = &boiler
		if 	err := DataCtl.ReadData(&og, "Uid"); err == nil {
			linked.Organization = &og
		} else {
			goazure.Error("Org Uid Is Not Valid!", err)
			continue
		}

		if 	err := DataCtl.ReadData(&ogType, "TypeId"); err == nil {
			linked.OrganizationType = &ogType
		} else {
			goazure.Error("OrgType Is Not Valid!", err)
			continue
		}

		if num, err := dba.BoilerOrm.InsertOrUpdate(&linked); err != nil {
			goazure.Error("M2M OrganizationLinked Add Error:", err, num)
		}
	}

	if boiler.InspectInnerDateNext.IsZero() { boiler.InspectInnerDateNext = time.Now().Add(time.Hour * 24 * 30) }
	if boiler.InspectOuterDateNext.IsZero() { boiler.InspectOuterDateNext = time.Now().Add(time.Hour * 24 * 30) }
	if boiler.InspectValveDateNext.IsZero() { boiler.InspectValveDateNext = time.Now().Add(time.Hour * 24 * 30) }
	if boiler.InspectGaugeDateNext.IsZero() { boiler.InspectGaugeDateNext = time.Now().Add(time.Hour * 24 * 30) }

	boiler.UpdatedBy = usr

	if err := DataCtl.AddData(&boiler, true); err != nil {
		e := fmt.Sprintln("Insert/Update Boiler Error!", err)
		return nil, errors.New(e)
	}

	go CalcCtl.InitBoilerCalculateParameter([]*models.Boiler{&boiler})

	goazure.Info("Updated Boiler:", boiler, info)

	return &boiler, nil
}

func (ctl *BoilerController) BoilerUpdateLocation() (*models.Boiler, error) {
	goazure.Info("Ready to Updated Boiler!")
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		return nil, errors.New(e)
	}

	var info BoilerInfo
	var boiler models.Boiler
	var address *models.Address
	var location *models.Location = &models.Location{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &info); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerInfo JSON Error", err)
		return nil, errors.New(e)
	}

	if err := dba.BoilerOrm.QueryTable("boiler").RelatedSel("Address__Location").Filter("Uid", info.Uid).One(&boiler); err != nil {
		e := fmt.Sprintln("Uid Invalid!")
		return nil, errors.New(e)
	}

	address = boiler.Address
	location.LocationId = info.LocationId
	DataCtl.ReadData(location, "LocationId")

	if address == nil {
		address = &models.Address{}
	}

	if err := DataCtl.ReadData(location, "LocationId"); err == nil { address.Location = location }
	address.Address = info.Address
	address.Longitude = info.Longitude
	address.Latitude = info.Latitude

	boiler.Address = address

	boiler.UpdatedBy = usr

	if err := DataCtl.AddData(address, true); err != nil {
		e := fmt.Sprintln("Insert/Update Address Error!", err)
		return nil, errors.New(e)
	}

	if err := DataCtl.UpdateData(&boiler); err != nil {
		e := fmt.Sprintln("Update Boiler Address Error!", err)
		return nil, errors.New(e)
	}

	goazure.Info("\nUpdated Boiler:", boiler, info)

	return &boiler, nil
}

func (ctl *BoilerController) BoilerUpdateMaintain() (*models.Boiler, error) {
	goazure.Info("Ready to Updated Boiler Maintain!")
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		return nil, errors.New(e)
	}

	var info BoilerInfo
	var boiler models.Boiler
	var contact *models.Contact

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &info); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerInfo JSON Error", err)
		return nil, errors.New(e)
	}

	if err := dba.BoilerOrm.QueryTable("boiler").RelatedSel("Contact").Filter("Uid", info.Uid).One(&boiler); err != nil {
		e := fmt.Sprintln("Uid Invalid!")
		return nil, errors.New(e)
	}

	contact = boiler.Contact

	if contact == nil {
		contact = &models.Contact{}
	}

	contact.Name = info.Contact
	contact.PhoneNumber = info.PhoneNumber
	contact.MobileNumber = info.MobileNumber
	contact.Email = info.Email
	boiler.Contact = contact


	boiler.InspectInnerDateNext = info.InspectDate.Inner
	boiler.InspectOuterDateNext = info.InspectDate.Outer
	boiler.InspectValveDateNext = info.InspectDate.Valve
	boiler.InspectGaugeDateNext = info.InspectDate.Gauge

	boiler.UpdatedBy = usr

	if err := DataCtl.AddData(contact, true); err != nil {
		e := fmt.Sprintln("Insert/Update Contact Error!", err)
		return nil, errors.New(e)
	}

	if err := DataCtl.UpdateData(&boiler); err != nil {
		e := fmt.Sprintln("Update Boiler Address Error!", err)
		return nil, errors.New(e)
	}

	goazure.Info("\nUpdated Boiler:", boiler, info)

	return &boiler, nil
}

func (ctl *BoilerController) BoilerDelete()  {
	goazure.Info("Ready to Delete Boiler!")
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var info BoilerInfo
	boiler := models.Boiler{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &info); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerInfo JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if len(info.Uid) > 0 {
		boiler.Uid = info.Uid
		boiler.UpdatedBy = usr

		if  err := DataCtl.DeleteData(&boiler); err != nil {
			e := fmt.Sprintln("Delete BoilerInfo Failed!", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}

		if  num, err := dba.BoilerOrm.QueryTable("boiler_terminal_combined").
			Filter("Boiler__Uid", info.Uid).Delete(); err != nil {
			e := fmt.Sprintln("Delete Boiler Terminal Combined Error:", err, num)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
		}

	} else {
		e := fmt.Sprintln("Invalid Boiler Uid!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
	}

	for i, b := range MainCtrl.Boilers {
		if b.Uid == boiler.Uid {
			MainCtrl.Boilers = append(MainCtrl.Boilers[:i], MainCtrl.Boilers[i + 1:]...)
			break
		}
	}

	go BlrCtl.RefreshGlobalBoilerList()
}

type BoilerBind struct {
	BoilerId		string		`json:"boiler_id"`
	TerminalId		string		`json:"terminal_id"`
	TerminalSetId	int32		`json:"terminal_set_id"`
}

func (ctl *BoilerController) BoilerBind() {
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var bind BoilerBind
	boiler := models.Boiler{}
	terminal := models.Terminal{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &bind); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerBind JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	goazure.Info("Ready to Bind Boiler!", bind)

	boiler.Uid = bind.BoilerId
	terminal.Uid = bind.TerminalId

	errB := DataCtl.ReadData(&boiler)
	errT := DataCtl.ReadData(&terminal)

	if errB != nil || errT != nil {
		e := fmt.Sprintln("Read Bind Error:", errB, errT, "\nBoiler:", boiler, "\nTerminal", terminal)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	//qt := dba.BoilerOrm.QueryTable("boiler").Filter("Terminal__Uid", bind.TerminalId)
	//count, errC := qt.Count()
	//if errC != nil {
	//	goazure.Warn("Read Terminal Boiler Count Error", errC)
	//}
	//goazure.Info("TerminalBoilerCount:", count)

	setId := bind.TerminalSetId
	if setId <= 0 {
		var combines []*models.BoilerTerminalCombined
		if  num, err := dba.BoilerOrm.QueryTable("boiler_terminal_combined").
			Filter("Terminal__Uid", bind.TerminalId).Filter("Boiler__Uid", bind.BoilerId).OrderBy("TerminalSetId").All(&combines); err != nil {
				goazure.Warn("Get Exist Combined Error:", err, num)
		}

		for i := int32(1); i <= 8; i++ {
			isMatched := false
			for _, cb := range combines {
				if cb.TerminalSetId == i {
					isMatched = true
					break
				}
			}

			if !isMatched {
				setId = i
				break
			}
		}

	}

	if boiler.Terminal == nil {
		boiler.Terminal = &terminal

		code := strconv.FormatInt(terminal.TerminalCode, 10)
		if len(code) < 6 {
			for l := len(code); l < 6; l++ {
				code = "0" + code
			}
		}
		boiler.TerminalCode = code
		boiler.TerminalSetId = setId

		if err := DataCtl.UpdateData(&boiler); err != nil {
			e := fmt.Sprintln("Boiler Bind Error:", err, boiler, terminal)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
		}
	}

	var combined models.BoilerTerminalCombined
	combined.Boiler = &boiler
	combined.Terminal = &terminal
	combined.TerminalCode = terminal.TerminalCode
	combined.TerminalSetId = setId

	if num, err := dba.BoilerOrm.InsertOrUpdate(&combined); err != nil {
		goazure.Error("M2M TerminalCombined Add Error:", err, num)
	}

	//m2m := dba.BoilerOrm.QueryM2M(&boiler, "TerminalsCombined")
	//if num, err := m2m.Add(&terminal); err != nil {
	//	goazure.Error("M2M Terminals Combined Add Error:", err, num)
	//}

	go BlrCtl.RefreshGlobalBoilerList()
}

func (ctl *BoilerController) BoilerUnbind() {
	goazure.Info("Ready to Unbind Boiler!")
	usr := ctl.GetCurrentUser()

	if !usr.IsAdmin() {
		e := fmt.Sprintln("Permission Denied, Only Admin Access!")
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	var bind BoilerBind
	boiler := models.Boiler{}
	terminal := models.Terminal{}

	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &bind); err != nil {
		e := fmt.Sprintln("Unmarshal BoilerBind JSON Error", err)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	boiler.Uid = bind.BoilerId
	terminal.Uid = bind.TerminalId

	errB := DataCtl.ReadData(&boiler)
	errT := DataCtl.ReadData(&terminal)

	if errB != nil || errT != nil {
		e := fmt.Sprintln("Read Unbind Data Error:", errB, errT, boiler, terminal)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	}

	if boiler.Terminal.Uid == bind.TerminalId {
		boiler.Terminal = nil
		boiler.TerminalCode = ""
		boiler.TerminalSetId = 0

		if err := DataCtl.UpdateData(&boiler); err != nil {
			e := fmt.Sprintln("Boiler Unbind Error:", err, boiler, terminal)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
		}
	}

	if  num, err := dba.BoilerOrm.QueryTable("boiler_terminal_combined").
		Filter("Boiler__Uid", bind.BoilerId).Filter("Terminal__Uid", bind.TerminalId).
		Delete(); err != nil {
			goazure.Error("Unbind Boiler/Terminal Error:", err, num)
	}

	go BlrCtl.RefreshGlobalBoilerList()
}

//***************************************** INFO UPDATE END

func (ctl *BoilerController) GetBoilerRuntime() {
	type Boil struct {
		Uid		string		`json:"uid"`
		RuntimeQueue	[]int		`json:"runtimeQueue"`
		Limit		int		`json:"limit"`
	}
	b := Boil{}
	resStatus := 200
	resBody := "Success"
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &b); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	var boiler models.Boiler
	qs := dba.BoilerOrm.QueryTable("boiler")
	if err := qs.RelatedSel("Form__Type").RelatedSel("Medium").RelatedSel("Usage").
		RelatedSel("Fuel__Type").RelatedSel("Template").
		RelatedSel("Factory").RelatedSel("Enterprise").RelatedSel("Maintainer").
		RelatedSel("RegisterOrg").
		RelatedSel("Terminal").RelatedSel("Contact").RelatedSel("Address__Location").
		Filter("Uid", b.Uid).OrderBy("Name").One(&boiler);
		err != nil {
		fmt.Println("Get Boiler Info Error: ", err)
	}

	type aBoiler struct {
		//models.Boiler
		Runtimes		[][]models.BoilerRuntime
		Parameters		[]models.RuntimeParameter
		Rules			[]models.RuntimeAlarmRule
	}
	runtimes := func(params []models.RuntimeParameter, limit int) []models.BoilerRuntime {
		var rtm []models.BoilerRuntime
		q := dba.BoilerOrm.QueryTable("boiler_runtime")
		q = q.RelatedSel("Boiler").RelatedSel("Parameter__Category").RelatedSel("Alarm").
			Filter("Boiler__Uid", boiler.Uid)
		if len(params) > 0 {
			q = q.Filter("Parameter__in", params)
		}
		if num, err := q.Filter("CreatedDate__gt", time.Now().Add(time.Hour * -24)).
			OrderBy("-CreatedDate", "Parameter__Category__Id", "Parameter__ParamId").Limit(limit).All(&rtm); err != nil || num == 0 {
			fmt.Println("Read BoilerRuntime Error", err)
			resStatus = 404
			resBody = "No Runtime Found!"
		} else {
			resStatus = 200
			resBody = "Success"
		}

		return rtm
	}

	parameters := func(rtmIds []int) []models.RuntimeParameter {
		var params []models.RuntimeParameter
		q := dba.BoilerOrm.QueryTable("runtime_parameter")
		q = q.RelatedSel("Category")
		if len(rtmIds) > 0 {
			q = q.Filter("Id__in", rtmIds)
		}
		if num, err := q.OrderBy("Category__Id", "ParamId").All(&params); err != nil || num == 0 {
			fmt.Println("Read Params Error", err)
		}

		return params
	}

	alarmRules := func(boiler *models.Boiler) []models.RuntimeAlarmRule {
		var rules []models.RuntimeAlarmRule
		q := dba.BoilerOrm.QueryTable("runtime_alarm_rule")
		q = q.RelatedSel("Parameter__Category").RelatedSel("BoilerForm").RelatedSel("BoilerMedium").RelatedSel("BoilerFuelType")
		cond := orm.NewCondition().Or("BoilerForm", boiler.Form).Or("BoilerForm__Id", 0).
			Or("BoilerMedium", boiler.Medium).Or("BoilerMedium__Id", 0).
			Or("BoilerFuelType", boiler.Fuel.Type).Or("BoilerFuelType__Id", 0)
		qs = qs.SetCond(cond)
		if num, err := q.Filter("IsDeleted", false).All(&rules); err != nil || num == 0 {
			fmt.Println("Read RuntimeRule Error", err)
			resStatus = 404
			resBody = "No Runtime Found!"
		}

		return rules
	}

	var aRuntimes 	[][]models.BoilerRuntime
	var aParams 	[]models.RuntimeParameter

	params := parameters(b.RuntimeQueue)
	rtms := runtimes(params, b.Limit)

	fmt.Println("\n\nPARAMs:", params, "\nRTMs: \n", rtms)

	paramSelect := func(ps []models.RuntimeParameter, rtmId int) (models.RuntimeParameter) {
		var ret models.RuntimeParameter
		for _, p := range ps {
			if p.Id == int64(rtmId) {
				ret = p
				break
			}
		}

		return ret
	}

	rtmSelect := func(rs []models.BoilerRuntime, rtmId int) ([]models.BoilerRuntime) {
		var ret []models.BoilerRuntime = []models.BoilerRuntime{}

		for _, r := range rs {
			if r.Parameter.Id == int64(rtmId) {
				ret = append(ret, r)
			}
		}

		return ret
	}

	for _, rtmId := range b.RuntimeQueue {
		aParams = append(aParams, paramSelect(params, rtmId))
		aRuntimes = append(aRuntimes, rtmSelect(rtms, rtmId))
	}

	ab := aBoiler {
		//Boiler: boiler,
		Runtimes: aRuntimes,
		Parameters: aParams,
		Rules: alarmRules(&boiler),
	}

	//fmt.Println("\nA Boiler: ", ab)

	//raw := "SELECT *, boiler.name AS boiler_name, param.name AS param_name, cate.name AS cate_name " +
	//"FROM boiler_runtime AS rtm " +
	//"LEFT JOIN boiler ON boiler.uid = rtm.boiler_id " +
	//"LEFT JOIN runtime_parameter AS param ON rtm.runtime_parameter_id = param.uid " +
	//"LEFT JOIN runtime_parameter_category AS cate ON param.category_id = cate.uid "
	//"WHERE cate.category_id=10 AND param.param_id=2 " +
	//"ORDER BY rtm.created_date LIMIT 100"
	//num, err = dba.BoilerOrm.LoadRelated(&v, "Maintenance")
	//num, err = dba.BoilerOrm.LoadRelated(&v, "Runtime", 100)

	//fmt.Println("Uid:", b.Uid, "\nRuntimeQueue:", b.RuntimeQueue);

	ctl.Data["json"] = ab
	ctl.ServeJSON()

	if resStatus != 200 {
		ctl.Ctx.Output.SetStatus(resStatus)
		ctl.Ctx.Output.Body([]byte(resBody))
	}
}

//TODO: Init
func (ctl *BoilerController) InitBoilerDefaults() {
	var bt models.BoilerType
	var bts models.BoilerTypeForm
	DataCtl.GenerateDefaultData(reflect.TypeOf(bts), boilerDefautlPath, "type", reflect.TypeOf(bt))

	var bm models.BoilerMedium
	DataCtl.GenerateDefaultData(reflect.TypeOf(bm), boilerDefautlPath, "medium", nil)

	var bu models.BoilerUsage
	DataCtl.GenerateDefaultData(reflect.TypeOf(bu), boilerDefautlPath, "usage", nil)

	var btp models.BoilerTemplate
	DataCtl.GenerateDefaultData(reflect.TypeOf(btp), boilerDefautlPath, "template", nil)

	var bf models.Fuel
	var bft models.FuelType
	DataCtl.GenerateDefaultData(reflect.TypeOf(bf), boilerDefautlPath, "fuel", reflect.TypeOf(bft))
}

func boilerWithUid(uid string) *models.Boiler {
	if len(uid) <= 0 {
		return nil
	}
	boiler := &models.Boiler{}
	boiler.Uid = uid
	if err := DataCtl.ReadData(boiler); err != nil {
		goazure.Error("Read Boiler With Uid Error:", err.Error())
		return nil
	}

	return boiler
}

func boilerForm(formId int) *models.BoilerTypeForm {
	form := models.BoilerTypeForm{}
	form.Id = int64(formId)
	err := DataCtl.ReadData(&form)
	if err != nil {
		fmt.Println("Read Error: ", err)
	}
	return &form
}

func boilerMedium(mediumId int) *models.BoilerMedium {
	m := models.BoilerMedium{}
	m.Id = int64(mediumId)
	err := DataCtl.ReadData(&m)
	if err != nil {
		fmt.Println("Read Error: ", err)
	}
	return &m
}

func boilerFuelType(fuelTypeId int) *models.FuelType {
	fuelType := models.FuelType{}
	fuelType.Id = int64(fuelTypeId)
	err := DataCtl.ReadData(&fuelType)
	if err != nil {
		fmt.Println("Read Error: ", err)
	}
	return &fuelType
}

