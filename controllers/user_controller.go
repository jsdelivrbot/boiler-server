package controllers

import (
	"fmt"

	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"

	"github.com/AzureRelease/boiler-server/models"

	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/AzureRelease/boiler-server/dba"
	"strconv"
	"time"

	"errors"
)

type UserController struct {
	MainController
}

const (
	SESSION_CURRENT_USER 	= "current_user"
	SESSION_THIRD_USER 	= "third_user"
)

const userDefautlsPath string = "models/properties/user_defaults/"

var UsrCtl *UserController
//var CurrentUser *models.User

func init()  {
	// var role models.UserRole
	// DataCtl.GenerateDefaultData(reflect.TypeOf(role), userDefautlsPath, "role", nil)
}

func (ctl *UserController) Get() {
	usr := ctl.GetCurrentUser()
	fmt.Println("Get User!", usr)
	//ctl.setCookiesWithUser(usr)

	ctl.Data["json"] = usr
	ctl.ServeJSON()
}

func (ctl *UserController) Post() {

}

func (ctl *UserController) GetCurrentUser() *models.User {
	if ctl == nil {
		//goazure.Warning("UserCtl is Nil!")
		//return ctl.getSysUser()
		return nil
	}
	usrSession := ctl.GetSession(SESSION_CURRENT_USER)
	if usrSession == nil {
		//ctl.SetSession(SESSION_CURRENT_USER, ctl.getSysUser())
		return nil
	}

	usr := ctl.GetSession(SESSION_CURRENT_USER)
	//CurrentUser = usr.(*models.User)
	fmt.Println("Try to Get CurrentUser:", usr)
	return usr.(*models.User)
}

type Usr struct {
	Uid			string		`json:"uid"`
	OpenId		string		`json:"openid"`
	UnionId		string		`json:"unionid"`

	Username 	string		`json:"username"`
	Password 	string		`json:"password"`
	PasswordNew	string		`json:"password_new"`

	Fullname	string		`json:"fullname"`
	Mobile 		string 		`json:"mobile"`
	Email 		string 		`json:"email"`
	Name 		string 		`json:"name"`

	RoleId		int			`json:"role"`
	OrgId		string		`json:"org"`
	Status		int			`json:"stat"`
	IpAddress	string		`json:"ip"`
}

func (ctl *UserController) UserList() {
	usr := ctl.GetCurrentUser()

	if usr.IsCommonUser() {
		return
	}

	qs := dba.BoilerOrm.QueryTable("user")
	qs = qs.RelatedSel()
	if usr.IsOrganizationUser() {
		qs = qs.Filter("Organization__Uid", usr.Organization.Uid)
	}
	//if usr != nil {
	//	uid := usr.Uid
	//	if err := qs.RelatedSel().Filter("uid", uid).One(usr); err != nil {
	//		fmt.Println("Get Current User Error: ", err)
	//	}
	//}

	var usrs []models.User

	if _, err := qs.Filter("IsDeleted", false).OrderBy("Status", "Username").All(&usrs); err != nil {
		fmt.Println("Fetch Users List Error: ", err)
	}

	ctl.Data["json"] = usrs
	ctl.ServeJSON()
}

func (ctl *UserController) UserRoleList() {
	var roles []models.UserRole

	if 	_, err := dba.BoilerOrm.QueryTable("user_role").
		Filter("RoleId__in", []int32{0, 1, 2, 10, 11}).
		OrderBy("RoleId").All(&roles); err != nil {
		goazure.Error("Fetch UserRoles List Error: ", err)
	}

	ctl.Data["json"] = roles
	ctl.ServeJSON()
}

func (ctl *UserController) UserActive() {
	us := ctl.GetCurrentUser()
	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	usr := models.User{
		MyUidObject: models.MyUidObject{
			Uid: u.Uid,
		},
	}
	if err := DataCtl.ReadData(&usr); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Find Active User Uid Error!"))
		return
	}

	usr.Status = models.USER_STATUS_NORMAL
	usr.UpdatedBy = us

	if err := DataCtl.AddData(&usr, true); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Active User Error!"))
		return
	}

	fmt.Println("\nActive User:", usr, u)
}

func (ctl *UserController) UserDelete() {
	us := ctl.GetCurrentUser()
	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	usr := models.User{
		MyUidObject: models.MyUidObject{
			Uid: u.Uid,
		},
	}

	DataCtl.ReadData(&usr)

	usr.Username = usr.Username + "_deleted"
	usr.UpdatedBy = us
	DataCtl.UpdateData(&usr)

	if err := DataCtl.DeleteData(&usr); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Delete User Error!"))
		return
	}

	fmt.Println("\nDelete User:", usr, u)
}

func (ctl *UserController) UserUpdate() {
	us := ctl.GetCurrentUser()

	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal Error", err)
		return
	}

	org := OrgCtrl.organizationWithUid(u.OrgId)
	fmt.Println("\nRoleId: ", u.RoleId, "\nOrg:", org)

	if (u.RoleId == models.USER_ROLE_ORG_USER || u.RoleId == models.USER_ROLE_ORG_ADMIN) && org == nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("机构用户必须选择一个有效的企业或机构"))
		goazure.Error("Org User With No Org!")
		return
	}

	if u.RoleId <= models.USER_ROLE_SUPERADMIN || u.RoleId >= models.USER_ROLE_USER {
		org = nil
	}

	usr := models.User{}
	if len(u.Uid) > 0 {
		usr.Uid = u.Uid
		if err := DataCtl.ReadData(&usr); err != nil {
			e := fmt.Sprintln("Find Update User Uid Error!", err)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}

		if len(u.PasswordNew) > 0 {
			usr.Password = ctl.hashedPassword(u.PasswordNew)
		}
		usr.Status = int(u.Status)
	} else {
		usr.Username = u.Username
		if 	err := DataCtl.ReadData(&usr, "Username"); err != nil || usr.IsDeleted {
			e := fmt.Sprintln("Find Update Username Error!", err)
			goazure.Warn(e)
			usr.IsDeleted = false
			usr.Password = ctl.hashedPassword(u.Password)
			usr.Status = models.USER_STATUS_NORMAL
			usr.CreatedBy = us
		} else {
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("用户名已经存在，请重新输入"))
			return
		}

	}

	//roleId, _ := strconv.ParseInt(u.RoleId, 10, 64);
	//stat, _ := strconv.ParseInt(u.Status, 10, 64);
	//usr.Username = u.Username
	usr.Name = u.Fullname
	usr.Role = role(int(u.RoleId))
	usr.Organization = org
	usr.UpdatedBy = us
	if len(usr.Picture) == 0 {
		usr.Picture  = "avatar0.png"
	}

	if err := DataCtl.AddData(&usr, true); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("更新用户信息失败!"))
		return
	}

	fmt.Println("\nUpdate User:", usr, u)
}

func (ctl *UserController) UserUpdateAdmin() {

}

func (ctl *UserController) UserProfileUpdate() {
	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal UserProfile Error", err)
		return
	}

	usr := ctl.GetCurrentUser()
	usr.Name = u.Fullname
	usr.MobileNumber = u.Mobile
	usr.Email = u.Email

	if err := DataCtl.UpdateData(usr); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Update UserProfile Error!"))
		return
	}

	fmt.Println("\nUpdate UserProfile:", usr, u)
}

func (ctl *UserController) UserPasswordUpdate() {
	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Updated Json Error!"))
		fmt.Println("Unmarshal UserProfile Error", err)
		return
	}

	usr := ctl.GetCurrentUser()

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(u.Password)); err != nil {
		fmt.Println("Hashed Pwd Error:", err)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("原密码错误!"))
		fmt.Println("Password Error", err)
		return
	}

	usr.Password = ctl.hashedPassword(u.PasswordNew)

	if err := DataCtl.UpdateData(usr); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("修改密码错误!"))
		return
	}

	fmt.Println("\nUpdate UserPassword:", usr, u)
}

func (ctl *UserController) UserLogin() {
	var u Usr
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("Login Json Error!"))
		goazure.Info("Unmarshal Error", err)
		return
	}

	goazure.Info("User Login! ", u)

	if _, err := ctl.Login(&u); err != nil {
		ctl.Ctx.Output.SetStatus(403)
		ctl.Ctx.Output.Body([]byte(err.Error()))
	}
}

//Common Login
func (ctl *UserController) Login(u *Usr) (*models.User, error) {
	fmt.Println("Ready to Login:", u)
	if u == nil || u.Username == "" || u.Password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}

	var usr models.User
	err := dba.BoilerOrm.QueryTable("user").Filter("Username", u.Username).Filter("IsDeleted", false).One(&usr)

	var resBody string
	var isSuccess bool = false
	var savePwd bool = false
	switch err {
	case orm.ErrNoRows:
		resBody = "用户名或密码错误"//"没有这个用户，请检查你的用户名是否正确"
	case orm.ErrMissPK:
		resBody = "未知错误"
	default:
		if usr.Status == 2 {
			return nil,errors.New("该用户已经被管理员禁用")
		}
		err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(u.Password))
		if err != nil {
			fmt.Println("Hashed Pwd Error:", err)
			resBody = "用户名或密码错误"
			savePwd = true
		} else {
			resBody = "登录成功"
			isSuccess = true
		}
	}

	loginPwd := func() string {
		if savePwd {
			return u.Password
		} else {
			return "Accepted"
		}
	}
	var login models.UserLogin

	login.User = &usr

	login.Name = u.Username
	login.Remark = resBody
	login.CreatedBy = &usr

	login.IsLogin = true
	login.IsSuccess = isSuccess

	login.LoginPassword = loginPwd()
	login.LoginIp = u.IpAddress

	if err := DataCtl.InsertData(&login); err != nil {
		fmt.Println("LoginLog Error", err)
	}

	goazure.Info("User Has Login", usr)

	if isSuccess {
		ctl.UpdateCurrentUser(&usr)
	} else {
		//defer ctl.Ctx.Output.SetStatus(resStat)
		//defer ctl.Ctx.Output.Body([]byte(resBody))
		fmt.Println("Login Error", err)

		return nil, errors.New(resBody)
	}

	return &usr, nil
}

func (ctl *UserController) UserRegister() (*models.User, error) {
	u := Usr{}
	if err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &u); err != nil {
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte("注册信息非法"))
		fmt.Println("User Reg Unmarshal Error", err)
		return nil, errors.New("注册信息非法")
	}

	var usr models.User

	var resStat int
	var resBody string
	var isSuccess bool = false

	switch err := DataCtl.ReadData(&models.User{ Username: u.Username }, "Username"); err {
	case orm.ErrNoRows:
		org := OrgCtrl.organizationWithUid(u.OrgId)
		fmt.Println("\nRoleId: ", u.RoleId, "\nOrg:", org)

		if (u.RoleId == models.USER_ROLE_ORG_USER || u.RoleId == models.USER_ROLE_ORG_ADMIN) && org == nil {
			fmt.Println("Org Error Get!")
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte("机构用户必须选择一个有效的企业或机构"))
			fmt.Println("Org User With No Org!")
			return nil, errors.New("机构用户必须选择一个有效的企业或机构")
		}

		usr = models.User{}

		usr.Name = u.Username
		usr.Role = role(u.RoleId)
		usr.Username = u.Username
		usr.Password = ctl.hashedPassword(u.Password)
		usr.Picture = "avatar0.png"

		usr.MobileNumber = u.Mobile
		//Email = u.Email

		usr.Organization = org
		usr.RegisterIp = u.IpAddress

		usr.IsDeleted = false

		if err := DataCtl.AddData(&usr, true, "Username"); err != nil {
			fmt.Println("Add User Error", err)
			resStat = 400
			resBody = "注册用户库出错"
		} else {
			resStat = 200
			resBody = "注册成功"
			isSuccess = true
		}
	case orm.ErrMissPK:
		resStat = 403
		resBody = "未知错误"
	default:
		resStat = 403
		resBody = "该用户名已经注册"
	}
	//fmt.Println("Register Error", err, isSuccess, resStat, resBody)
	fmt.Println("User Register!\n", u)

	if isSuccess {
		if _, err := ctl.Login(&u); err != nil {
			ctl.Ctx.Output.SetStatus(403)
			ctl.Ctx.Output.Body([]byte(err.Error()))

			return nil, err
		}
	} else {
		//defer ctl.Ctx.Output.SetStatus(resStat)
		//defer ctl.Ctx.Output.Body([]byte(resBody))
		ctl.Ctx.Output.SetStatus(resStat)
		ctl.Ctx.Output.Body([]byte(resBody))

		return nil, errors.New(resBody)
	}

	return &usr, nil
}

func (ctl *UserController) UserLogout() {
	usr := ctl.GetCurrentUser()
	if  usr == nil {
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
	goazure.Info("Ready to Logout")
	//ctl.setCookiesWithUser(nil)
	ctl.DestroySession()

	var login models.UserLogin
	login.User = usr
	login.Name = usr.Username
	login.Remark = "用户登出"
	login.CreatedBy = usr

	login.IsLogin = false
	login.IsSuccess = true

	login.LoginPassword = "Logout"

	if err := DataCtl.InsertData(&login); err != nil {
		goazure.Error("LoginLog Error", err)
	}

	//defer ctl.Ctx.Output.SetStatus(200)
	//defer ctl.Ctx.Output.Body([]byte("注销成功！"))
	ctl.Ctx.Output.SetStatus(200)
	ctl.Ctx.Output.Body([]byte("注销成功！"))

	ctl.Ctx.Redirect(302, "/")
}

func (ctl *UserController) setCookiesWithUser(usr *models.User) error {
	if usr == nil {
		fmt.Println("User IS NULL!!!!")
		ctl.Ctx.SetCookie("isLogin", "0")
		ctl.Ctx.SetCookie("username", "")
		ctl.Ctx.SetCookie("user_picture", "")
		ctl.Ctx.SetCookie("user_role_zh", "")

		return nil
	}

	uid := usr.Uid
	qs := dba.BoilerOrm.QueryTable("user")
	err := qs.RelatedSel().Filter("uid", uid).One(usr)

	textQuoted := strconv.QuoteToASCII(usr.Role.Name)
	//textQuoted = strings.Replace(textQuoted, "\\u", "", -1)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	ctl.Ctx.SetCookie("isLogin", "1")
	ctl.Ctx.SetCookie("username", usr.Username)
	ctl.Ctx.SetCookie("user_picture", usr.Picture)
	ctl.Ctx.SetCookie("user_role_zh", textUnquoted)

	if err != nil {
		fmt.Println("SetCookies Error:", err)
	}

	return err
}

func (ctl *UserController) hashedPassword(pwd string) string {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPwd)
}

func (ctl *UserController) UpdateUserName() {
	var users []models.User
	qs := dba.BoilerOrm.QueryTable("user")
	if num, err := qs.All(&users); err != nil {
		fmt.Println("Read Boilers Error:", err, num)
	}

	for i, u := range users {
		if len(u.Name) <= 0 {
			u.Name = u.Username
			if err := DataCtl.UpdateData(&u); err != nil {
				fmt.Println("Updated User Name Error: ", i, err)
			}
		}
	}
}


func (ctl *UserController) UserImageUpload() {
	usr := ctl.GetCurrentUser()

	if file, header, er := ctl.GetFile("file"); er != nil {
		e := fmt.Sprintln("Upload Image File Failed:", er)
		goazure.Error(e)
		ctl.Ctx.Output.SetStatus(400)
		ctl.Ctx.Output.Body([]byte(e))
		return
	} else {
		//goazure.Warn("Upload Image Success:", file, header)
		if file == nil {
			e := fmt.Sprintln("File Is Null!")
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}

		// get the filename
		fileName := "headImage_" + time.Now().Format("20060102150405") + ".png"
		basePath := "static/images/upload/"
		filePath := basePath + fileName
		// save to server
		if err := ctl.SaveToFile("file", filePath); err != nil {
			e := fmt.Sprintln("Save File Error:", err, fileName)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}

		usr.Picture = "/upload/" + fileName

		if 	err := DataCtl.UpdateData(usr); err != nil {
			e := fmt.Sprintln("Updated User Image Error:", err, fileName)
			goazure.Error(e)
			ctl.Ctx.Output.SetStatus(400)
			ctl.Ctx.Output.Body([]byte(e))
			return
		}

		//ctl.UpdateCurrentUser(usr)
		goazure.Info("Save Done:", header.Filename, fileName)
	}
}

func addUser(usr models.User) (error) {
	err := DataCtl.AddData(&usr, true, "Uid")
	if err != nil {
		fmt.Println("Add User Error:", err)
	}

	return err
}

func role(roleId int) *models.UserRole {
	r := models.UserRole{ RoleId: int32(roleId) }
	if err := DataCtl.ReadData(&r, "RoleId"); err != nil {
		fmt.Println("Read Role Error", err)
	}

	return &r
}