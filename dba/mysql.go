package dba

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/AzureTech/goazure/orm"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/models/caches"
	"net/url"
	"github.com/AzureRelease/boiler-server/common"

	"github.com/AzureRelease/boiler-server/conf"
)

var MyORM 		orm.Ormer
var BoilerOrm 	orm.Ormer

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)


	var myConnection string




	if conf.IsRelease {
		myConnection = "azureadmin:azure%2016@tcp(rm-a0z2ur23e09te04c8h4.mysql.rds.aliyuncs.com:3306)/boiler_main?charset=utf8&loc=" + url.QueryEscape("PRC")
	} else {
		myConnection = "azureadmin:azure%2016@tcp(rm-a0z2ur23e09te04c8h4n.mysql.rds.aliyuncs.com:3306)/boiler_main?charset=utf8&loc=" + url.QueryEscape("PRC")
	}
	orm.RegisterDataBase("default", "mysql", myConnection)

	orm.RegisterModel(
		new(models.Application),

		new(models.User),
		new(models.UserRole),
		new(models.UserThird),
		new(models.UserLogin),
		new(models.UserSession),

		new(models.Organization),
		new(models.OrganizationType),

		new(models.Contact),
		new(models.Address),
		new(models.Location),

		new(models.Boiler),
		new(models.BoilerType),
		new(models.BoilerTypeForm),
		new(models.BoilerMedium),
		new(models.BoilerUsage),
		new(models.BoilerTemplate),
		new(models.BoilerMaintenance),

		new(models.BoilerMessageSubscriber),
		new(models.BoilerTerminalCombined),
		new(models.BoilerOrganizationLinked),

		new(models.Fuel),
		new(models.FuelType),

		new(models.BoilerConfig),
		new(models.BoilerRuntime),
		new(models.BoilerRuntimeArchived),

		new(caches.BoilerRuntimeCacheInstant),
		//new(caches.BoilerRuntimeCacheHistory),
		new(caches.BoilerRuntimeCacheDay),
		new(caches.BoilerRuntimeCacheWeek),
		new(caches.BoilerRuntimeCacheMonth),

		new(caches.BoilerRuntimeCacheFlow),
		new(caches.BoilerRuntimeCacheFlowDaily),
		new(caches.BoilerRuntimeCacheSteamTemperature),
		new(caches.BoilerRuntimeCacheSteamPressure),
		new(caches.BoilerRuntimeCacheSmokeTemperature),
		new(caches.BoilerRuntimeCacheSmokeComponent),
		new(caches.BoilerRuntimeCacheWaterTemperature),
		new(caches.BoilerRuntimeCacheEnvironmentTemperature),
		new(caches.BoilerRuntimeCacheHeat),
		new(caches.BoilerRuntimeCacheHeatDaily),
		new(caches.BoilerRuntimeCacheExcessAir),

		new(caches.BoilerRuntimeHistory),

		new(caches.BoilerRuntimeCacheStatus),
		new(caches.BoilerRuntimeCacheStatusRunning),

		new(models.BoilerCalculateParameter),
		new(models.BoilerCalculateResult),

		new(models.BoilerAlarm),
		new(models.BoilerAlarmHistory),
		new(models.BoilerAlarmFeedback),
		new(models.BoilerFuelRecord),

		new(models.RuntimeParameter),
		new(models.RuntimeParameterCategory),
		new(models.RuntimeParameterMedium),
		new(models.RuntimeParameterChannelConfig),
		new(models.RuntimeParameterChannelConfigRange),
		new(models.RuntimeAlarmRule),

		new(models.Message),
		new(models.MessageType),
		new(models.MessageFormatter),
		new(models.MessageTag),

		new(models.MessageLog),
		new(models.Message16Bit),
		new(models.Message32Bit),

		new(models.Terminal),

		new(models.Dialogue),
		new(models.DialogueComment),
	)

	orm.Debug = false//!conf.IsRelease

	MyORM = orm.NewOrm()
	MyORM.Using("default")
	// orm.RunSyncdb("default", false, true)

	BoilerOrm = MyORM
	common.BoilerOrm = BoilerOrm
	// orm.DefaultTimeLoc = time.UTC
}
