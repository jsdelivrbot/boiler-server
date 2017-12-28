package models

type Boiler struct {
	MyUidObject

	//锅炉类型* 卧式-双筒
	Form				*BoilerTypeForm			`orm:"rel(fk)"`
	//锅炉介质：蒸汽锅炉、热水锅炉、有机热载体锅炉
	Medium				*BoilerMedium			`orm:"rel(fk)"`
	//锅炉用途：热电锅炉、工业锅炉
	Usage				*BoilerUsage			`orm:"rel(fk)"`
	//锅炉燃料：燃煤、燃油（气）、生物质 燃料类型* 烟煤Ⅱ（17700≤Qnet.v.ar≤21000）
	Fuel				*Fuel					`orm:"rel(fk)"`
	Template			*BoilerTemplate			`orm:"rel(fk);null"`

	Factory				*Organization			`orm:"rel(fk);null"`
	Enterprise			*Organization			`orm:"rel(fk);null"`
	Installed			*Organization			`orm:"rel(fk);null"`

	//使用单位地址
	//所在区域: 所在区域* 浙江省 杭州市 桐庐县	<-REGION
	//经度，纬度
	Address				*Address				`orm:"rel(fk);null"`
	//出厂编号	<-FACTORY_Code
	FactoryNumber		string					`orm:"size(30)"`
	//注册码		<-REGISTER_CODE
	RegisterCode		string					`orm:"size(60)"`
	//登记机构 	<-REGISTER_ORG
	//RegisterOrg		string					`orm:"size(60)"`
	RegisterOrg			*Organization			`orm:"rel(fk);null"`
	//使用证号	<-USE_CODE
	CertificateNumber	string					`orm:"size(60)"`
	//设备代码	<-BOILER_CODE
	DeviceCode			string					`orm:"size(60)"`
	//锅炉型号	<-MAIN_PARAM
	ModelCode			string					`orm:"size(60)"`

	//蒸发量
	EvaporatingCapacity	int64
	//联系人&联系电话
	Contact				*Contact				`orm:"rel(fk);index;null"`

	Terminal			*Terminal				`orm:"rel(fk);null;index"`
	TerminalCode		string					`orm:"index"`
	TerminalSetId		int32					`orm:"index"`

	BoilerMaintainSchedule

	//TODO: FOR the bug of goazure
	Status				map[string]int32					`orm:"-"`
	Runtime				[]*BoilerRuntime					`orm:"reverse(many)"`
	Calculate			[]*BoilerCalculateParameter			`orm:"reverse(many)"`
	Maintenance			[]*BoilerMaintenance				`orm:"reverse(many)"`
	Subscribers			[]*User								`orm:"rel(m2m);rel_through(github.com/AzureRelease/boiler-server/models.BoilerMessageSubscriber)"`
	//`orm:"rel(m2m);rel_table(boiler_subscriber_users)"`

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
	Boiler		*Boiler		`orm:"rel(fk)"`
	User		*User		`orm:"rel(fk)"`
}