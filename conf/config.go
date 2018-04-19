package conf

import "github.com/AzureTech/goazure"

var IsRelease bool = goazure.AppConfig.String("runmode") == "prod"
var Version string = "e884cb938479205f9e25c91152bf359db1c35d6f"
var BinPath string = "/home/apps/bin/"
var TermNoRegist = []byte("Term not be registed")
var TermTimeout = []byte("Term Response Timeout")
var BoilerStart string = "\xAA"
var BoilerShut string = "\x55"
var BoilerReset string = "\x99"
var TermOnline string = "1000000000"
var TermOffline string = "0000000000"