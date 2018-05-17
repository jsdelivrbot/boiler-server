package models

type Boiler struct {
	MyUidObject

	//锅炉类型* 卧式-双筒 etc.
	Form				*BoilerTypeForm			`orm:"rel(fk);null"`
	//锅炉介质：蒸汽锅炉、热水锅炉、有机热载体锅炉
	Medium				*BoilerMedium			`orm:"rel(fk)"`
	//锅炉用途：热电锅炉、工业锅炉
	Usage				*BoilerUsage			`orm:"rel(fk)"`
	//锅炉燃料：燃煤、燃油（气）、生物质 燃料类型* 烟煤Ⅱ（17700≤Qnet.v.ar≤21000）
	Fuel				*Fuel					`orm:"rel(fk)"`
	Template			*BoilerTemplate			`orm:"rel(fk);null"`

	Factory				*Organization			`orm:"rel(fk);null"`
	Enterprise			*Organization			`orm:"rel(fk);null"`
	Maintainer			*Organization			`orm:"rel(fk);null"`
	Supervisor			*Organization			`orm:"rel(fk);null"`

	OrganizationsLinked	[]*Organization			`orm:"rel(m2m);null;index;rel_through(github.com/AzureRelease/boiler-server/models.BoilerOrganizationLinked)"`
	//使用单位地址
	//所在区域: 所在区域* 浙江省 杭州市 桐庐县	<-REGION
	//经度，纬度
	Address				*Address				`orm:"rel(fk);null"`	//出厂编号	<-FACTORY_Code
	FactoryNumber		string					`orm:"size(30)"`
	RegisterCode		string					`orm:"size(60)"`		//注册码		<-REGISTER_CODE
	RegisterOrg			*Organization			`orm:"rel(fk);null"`	//登记机构 	<-REGISTER_ORG
	CertificateNumber	string					`orm:"size(60)"`		//使用证号	<-USE_CODE
	DeviceCode			string					`orm:"size(60)"`		//设备代码	<-BOILER_CODE
	ModelCode			string					`orm:"size(60)"`		//锅炉型号	<-MAIN_PARAM

	EvaporatingCapacity	 float64							//蒸发量
	Contact				*Contact				`orm:"rel(fk);index;null"`	//联系人&联系电话

	Terminal			*Terminal				`orm:"rel(fk);null;index"`
	TerminalCode		string					`orm:"index"`
	TerminalSetId		int32					`orm:"index"`
	TerminalsCombined	[]*Terminal				`orm:"rel(m2m);null;index;rel_through(github.com/AzureRelease/boiler-server/models.BoilerTerminalCombined)"`

	BoilerMaintainSchedule

	//TODO: FOR the bug of goazure
	Status				map[string]int32					`orm:"-"`
	Runtime				[]*BoilerRuntime					`orm:"reverse(many)"`
	Calculate			[]*BoilerCalculateParameter			`orm:"reverse(many)"`
	Maintenance			[]*BoilerMaintenance				`orm:"reverse(many)"`
	Subscribers			[]*User								`orm:"rel(m2m);rel_through(github.com/AzureRelease/boiler-server/models.BoilerMessageSubscriber)"`

	/*
	是否派单 是   否
	流程图模式* 模板模式
	流程图模板 燃煤双锅筒  󰆦
	自定义流程图
	自定义分类
	*/

}

type BoilerMessageSubscriber struct{
	Id			int64		`orm:"auto"`
	Boiler		*Boiler		`orm:"rel(fk);index"`
	User		*User		`orm:"rel(fk);index"`
}

type BoilerTerminalCombined struct{
	Id				int64		`orm:"auto"`
	Boiler			*Boiler		`orm:"rel(fk);index"`
	Terminal		*Terminal	`orm:"rel(fk);index"`
	TerminalCode	int64		`orm:"index"`
	TerminalSetId	int32		`orm:"index"`
}

func (cb *BoilerTerminalCombined) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Terminal"},
		[]string{"TerminalCode", "TerminalSetId"},
	}
}

type BoilerTermStatus struct{
	BoilerTermId		string				`orm:"pk;column(Boiler_term_id)"`
	BoilerTermIp		string				`orm:"column(Boiler_term_ip)"`
	BoilerTermPwd		int32				`orm:"column(Boiler_term_pwd)"`
	BoilerTermStatus	bool				`orm:"column(Boiler_term_status)"`
}

type BoilerOrganizationLinked struct{
	Id					int64				`orm:"auto"`
	Boiler				*Boiler				`orm:"rel(fk);index"`
	OrganizationType	*OrganizationType	`orm:"rel(fk);index"`
	Organization		*Organization		`orm:"rel(fk);index"`
}

func (ob *BoilerOrganizationLinked) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "Organization"},
	}
}

/*** ***/
func (boiler *Boiler) IsBelongToOrganization(org *Organization) bool {
	if 	(boiler.Enterprise != nil && (boiler.Enterprise.Uid == org.Uid || (boiler.Enterprise.SuperOrganization != nil && boiler.Enterprise.SuperOrganization.Uid == org.Uid))) ||
		(boiler.Factory != nil && (boiler.Factory.Uid == org.Uid || (boiler.Factory.SuperOrganization != nil && boiler.Factory.SuperOrganization.Uid == org.Uid))) ||
		(boiler.Maintainer != nil && (boiler.Maintainer.Uid == org.Uid || (boiler.Maintainer.SuperOrganization != nil && boiler.Maintainer.SuperOrganization.Uid == org.Uid))) ||
		(boiler.Supervisor != nil && (boiler.Supervisor.Uid == org.Uid || (boiler.Supervisor.SuperOrganization != nil && boiler.Supervisor.SuperOrganization.Uid == org.Uid))) {
		return true
	}

	for _, li := range boiler.OrganizationsLinked {
		if 	li.Uid == org.Uid ||
			(li.SuperOrganization != nil && li.SuperOrganization.Uid == org.Uid) {
			return true
		}
	}

	return false
}