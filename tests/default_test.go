package test

import (
	"testing"
	"runtime"
	"path/filepath"

	//"net/http"
	//"net/http/httptest"
	"github.com/AzureTech/goazure"

	"github.com/AzureRelease/boiler-server/controllers"
	"github.com/AzureRelease/boiler-server/models"

	//. "github.com/smartystreets/goconvey/convey"



)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	goazure.TestAzureInit(apppath)
}

/*
func TestAzure(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	goazure.BeeApp.Handlers.ServeHTTP(w, r)
	goazure.Trace("testing", "TestAzure", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
*/

func TestAlarmReload(t *testing.T) {
	controllers.BlrCtl.WaitGroup.Wait()
	controllers.ParamCtrl.WaitGroup.Wait()

	boiler := controllers.BlrCtl.Boiler("7551c963-76ff-4fe0-a50a-c65aa40537e4")
	param := controllers.ParamCtrl.RuntimeParameter(10274)

	var rtm models.BoilerRuntime
	rtm.Parameter = param
	rtm.Boiler = boiler
	rtm.Value = 5555
	rtm.Remark = "temp"

	rtm.Status = models.RUNTIME_STATUS_NEW

	controllers.RtmCtl.RuntimeDataReload(&rtm, 0)
}
