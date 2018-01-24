package log

import (
	"github.com/AzureTech/goazure/logs"
)

func init() {
	logs.Async()
	logs.EnableFuncCallDepth(true)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"/var/log/boiler/main.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
}
