package conf

import (
	"github.com/AzureTech/goazure"
	"net/url"
)

var IsRelease bool = goazure.AppConfig.String("runmode") == "prod"

var Version string = "e884cb938479205f9e25c91152bf359db1c35d6f"
var BinPath string =    "/home/apps/bin/"
var TermNoRegist = []byte("Term not be registed")
var TermTimeout = []byte("Term Response Timeout")
var BoilerStart string = "\xAA"
var BoilerShut string = "\x55"
var BoilerReset string = "\x99"
var TermConfig string = "0xC0(终端配置)"
var TermRestart string = "0xC9(终端重启)"
var BoilerController string = "0xCF(锅炉控制)"
var TermOnline bool = true
var TermOffline bool = false
var BatchFlag bool = true
var ContentLogsFlag bool = true
var IsReloadLogEnabled = true
var Server string
var DbConnection string

func init() {
	if IsRelease {
		DbConnection = "holder2025:hold+123456789@tcp(rm-uf6s78595q8r68it7vo.mysql.rds.aliyuncs.com:3306)/boiler_main?charset=utf8&loc=" + url.QueryEscape("PRC")
		Server = "47.100.0.27:18887"
		goazure.Warning("数据库连接到压力服")
	} else {
		DbConnection = "azureadmin:azure%2016@tcp(rm-a0z2ur23e09te04c8h4n.mysql.rds.aliyuncs.com:3306)/boiler_main?charset=utf8&loc=" + url.QueryEscape("PRC")
		Server = "139.196.152.127:12000"
		goazure.Warning("数据库连接到正式服")
	}
}



