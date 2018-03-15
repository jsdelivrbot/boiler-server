package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"
	"strconv"
	"database/sql"
	"time"
	"github.com/AzureRelease/boiler-server/models"
	"github.com/AzureRelease/boiler-server/dba"
	"flag"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"errors"
)

type DBController struct {
	MainController

	db		*sql.DB
}

var DBCtl *DBController = &DBController{}

var (
	debug	  = flag.Bool("debug", true, "enable debugging")
	server    = flag.String("server", "101.231.74.206", "the database server")
	port *int = flag.Int("port", 1433, "the database port")
	user      = flag.String("user", "boiler", "the database user")
	password  = flag.String("password", "boiler123", "the database password")
	database  = flag.String("database", "boiler", "the database name")
)

func init() {
	err := DBCtl.InitDbConnect()
	if err != nil {
		return
	}

	//defer db.Close()
}

func (ctl *DBController) GetStringFromMap(m orm.Params, defaults string, col string, cols ...string) string {
	if m[col] != nil {
		return m[col].(string)
	} else {
		for _, c := range cols {
			if m[c] != nil {
				return m[c].(string)
			}
		}
		return defaults
	}
}

func (ctl *DBController)InitMSSQLData() {
	ctl.ImportMSSQLData(0, time.Time{})
}

func (ctl *DBController)LoadMSSQLData() {
	ticker := time.NewTicker(time.Minute * 5)
	tick := func() {
		for t := range ticker.C {
			DBCtl.ImportMSSQLData(1, t)
		}
	}

	go tick()
}

func (ctl *DBController)ImportMSSQLData(offset int, tm time.Time) error {
	if ctl.db == nil {
		return errors.New("db conn can not be nil!")
	}

	query :=
		"SELECT * " +
		"FROM BoilerData_310101C027 " +
		"ORDER BY timestamp ASC; "

	if offset > 0 {
		query =
			"SELECT * " +
			"FROM BoilerData_310101C027 " +
			"WHERE [timestamp] > DATEADD(HOUR, -" + strconv.FormatInt(int64(offset), 10) + " , GETDATE()) " +
			"ORDER BY [timestamp] ASC; "
	}

	stmt, err := ctl.db.Prepare(query)
	if err != nil {
		goazure.Error("Prepare failed:", err.Error())
		return err
	}

	rows, err := stmt.Query()
	var id 				int64
	var parameterName	string
	var timestamp		time.Time
	var value			float64

	BlrCtl.bWaitGroup.Wait()
	var boiler *models.Boiler
	for _, b := range MainCtrl.Boilers {
		if b.Uid == "8a3e47d0-759f-474d-b2e2-a89692f7c496" {
			boiler = b
			break
		}
	}

	var sRtm models.BoilerRuntime
	sParam := runtimeParameter(1107)
	sRtm.Boiler = boiler
	sRtm.Parameter = sParam
	sRtm.CreatedDate = time.Now()
	sRtm.Value = 1
	sRtm.Remark = "MSSQL Trigger"

	if er := DataCtl.AddData(&sRtm, true); er != nil {
		goazure.Error("Added MSSQL Runtime Error:", er)
	}

	for true {
		if rows.Next() {
			err = rows.Scan(&id, &parameterName, &timestamp, &value)
			if err != nil {
				goazure.Error("Scan failed:", err.Error())
			}

			var param	models.RuntimeParameter
			var rtm 	models.BoilerRuntime

			if er := dba.BoilerOrm.QueryTable("runtime_parameter").Filter("Name", parameterName).One(&param); er != nil {
				goazure.Error("Read Parameter By Name Error:", er, parameterName)
			}

			rtm.Boiler = boiler
			rtm.Parameter = &param
			rtm.CreatedDate = timestamp.Local()
			rtm.Value = int64(value * 1000)
			rtm.Remark = "MSSQL"
			rtm.Name = strconv.FormatInt(id, 10)

			if er := DataCtl.AddData(&rtm, false); er != nil {
				goazure.Error("Added MSSQL Runtime Error:", er)
			}

			ReloadRuntimeHistory(&rtm, value)

			goazure.Info("obj:", id, parameterName, timestamp, value)
		} else {
			break
		}
	}

	defer stmt.Close()

	return nil
}

func (ctl *DBController)InitDbConnect() error {
	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", *server, *user, *password, *port, *database)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}

	db, err := sql.Open("mssql", connString)
	if err != nil {
		goazure.Error("Open connection failed:", err.Error())
	}

	ctl.db = db

	return err
}