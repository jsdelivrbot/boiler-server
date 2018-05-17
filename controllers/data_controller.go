package controllers

import (
	"reflect"
	"errors"
	"fmt"

	"github.com/AzureTech/goazure/orm"
	"github.com/AzureTech/goazure"
	"github.com/pborman/uuid"

	"github.com/AzureRelease/boiler-server/dba"
	"github.com/AzureRelease/boiler-server/models"
	"os"
	"encoding/csv"
	"path/filepath"
	"strings"
	"strconv"
	"log"
)

type DataController struct {
	MainController
}

var DataCtl *DataController = &DataController{}

func (ctl *DataController) ReadData(data models.DataInterface, cols ...string) (error) {
	var err error

	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("input param is not a struct")
	}
	m := fmt.Sprintf("Read Type %v:%v\n", t, data)
	goazure.Info(m)

	err = dba.BoilerOrm.Read(data, cols...)

	switch err {
	case orm.ErrNoRows:
	case orm.ErrMissPK:
		fallthrough
	default:
		//fmt.Println("Get Data: ", data)
	}

	return err
}

func (ctl *DataController) AddData(data models.DataInterface, needUpdate bool, cols ...string) (error) {
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("input param is not a struct")
	}
	//fmt.Printf("Output Type %v:\n", t)
	d := reflect.New(t)
	eData := d.Elem()

	for _, col := range cols {
		value := reflect.ValueOf(data).Elem().FieldByName(col)
		goazure.Info("Col:", col, "Value:", value)
		eData.FieldByName(col).Set(value)
	}

	aData := d.Interface().(models.DataInterface)

	//fmt.Println("Data[",reflect.TypeOf(data),"]:", data, "\naData[",reflect.TypeOf(aData),"]:", aData)
	aData.SetKey(data.GetKey())
	//goazure.Debug("Ready to Read aData:", aData, cols)
	err := dba.BoilerOrm.Read(aData, cols...)
	if err != nil {
		warn := fmt.Sprintf("ReadData Error: %T, %v, %s", aData, aData, err)
		goazure.Warn(warn)
	}

	obj := data.GetObject().(*models.MyObject)
	aObj := aData.GetObject().(*models.MyObject)

	//goazure.Info("obj:", obj, reflect.TypeOf(data.GetObject()))
	//goazure.Info("abj:", aObj, reflect.TypeOf(aData.GetObject()))

	var crtUsr *models.User

	isSysUser := t == reflect.TypeOf(models.User{}) &&
		obj.Name == "system" //&& obj.NameEn == "system"
	isSysRole := t == reflect.TypeOf(models.UserRole{}) &&
		reflect.ValueOf(data).Elem().FieldByName("RoleId").Int() == 0

	if isSysUser || isSysRole {
		crtUsr = nil
	} else {
		crtUsr = UsrCtl.GetCurrentUser()
	}
	//fmt.Println("isSysUser: ", isSysUser, "\nisSysRole: ", isSysRole)
	switch err {
	case orm.ErrNoRows:
		fallthrough
	case orm.ErrMissPK:
		err = ctl.InsertData(data)
	default:
		//fmt.Println("Get A Data: ", data)
		data.SetKey(aData.GetKey())

		obj.CreatedDate = aObj.CreatedDate

		if isSysUser || isSysRole {
			obj.CreatedBy = nil
		} else  {
			obj.CreatedBy = aObj.CreatedBy
		}
	}

	/** Unknown Issue
	/*  panic: reflect: call of reflect.Value.Int on string Value
	isCreated, id, err := dba.BoilerOrm.ReadOrCreate(data, cols[0], cols[1:]...);
	fmt.Println("Data:", data, "\naData:", aData)
	if err != nil {
		fmt.Println("***Error:", err)
	}

	if err == nil {
		if isCreated {
			//obj.CreatedDate = time.Now()
			//obj.CreatedBy = UsrCtl.CurrentUser()
			fmt.Println("Created Data:", data)
		} else {
			fmt.Println("Get ", id, " Data:", data)
		}
	}
	*/

	if !needUpdate {
		return err
	}

	if crtUsr != nil {
		obj.UpdatedBy = crtUsr
	}
	err = DataCtl.UpdateData(data)

	return err
}

func (ctl *DataController) InsertData(data models.DataInterface) error {
	obj := data.GetObject().(*models.MyObject)
	key := data.GetKey()
	switch reflect.TypeOf(key).Kind() {
	case reflect.Int32:
		//goazure.Info("Auto Id")
	case reflect.String:
		//goazure.Info("uuid")
		if key.(string) == "" {
			if dba.BoilerOrm.Driver().Type() == orm.DRPostgres {

			} else {

			}
			data.SetKey(uuid.New())

			//quoteRune := rune('\'')
			//obj.Uid = fmt.Sprintf("%s%s%s::uuid", string(quoteRune), uuid.New(), string(quoteRune))
			//obj.Uid = "uuid_generate_v4()"
			//fmt.Println("Create New UUID:", data.GetKey())
		}
	}

	if obj.CreatedBy == nil {
		obj.CreatedBy = UsrCtl.GetCurrentUser()
	}
	if obj.UpdatedBy == nil {
		obj.UpdatedBy = obj.CreatedBy
	}

	//if obj.CreatedDate.IsZero() {
	//	obj.CreatedDate = time.Now();
	//}
	//obj.UpdatedDate = time.Now();

	var err error
	if _, err = dba.BoilerOrm.Insert(data); err != nil {
		msg := fmt.Sprintf("Insert Error: %T, %s", data, err)
		goazure.Error(msg)
	} else {
		msg := fmt.Sprintf("Inserted[ %T ]: %v", data, data)
		goazure.Info(msg)
	}

	return err
}

func (ctl *DataController) UpdateData(data models.DataInterface) error {
	var e error

	//obj := data.GetObject()
	//obj.UpdatedDate = time.Now()

	if num, err := dba.BoilerOrm.Update(data); err != nil {
		msg := fmt.Sprintf("Update Error: %T, %v, %s", data, data, err)
		goazure.Error(msg)
		e = errors.New(msg)
	} else {
		msg := fmt.Sprintf("Updated Data: %T, %v, %d", data, data, num)
		goazure.Info(msg)
	}

	return e
}

func (ctl *DataController) DeleteData(data models.DataInterface) error {
	var err error
	if err = dba.BoilerOrm.Read(data); err != nil {
		goazure.Error("Delete Find Error: ", err)
	}

	obj := data.GetObject().(*models.MyObject)
	obj.IsDeleted = true

	if num, err := dba.BoilerOrm.Update(data); err == nil {
		goazure.Info("Delete Data: ", data, num)
	} else {
		goazure.Error("Delete Error: ", err)
	}

	return err
}

func (ctl *DataController) GenerateDefaultData(aType reflect.Type, path string, filenameKey string, superType reflect.Type) error {
	const headerRowNo int = 0
	const fieldRowNo int = 1

	if aType.Kind() == reflect.Ptr {
		aType = aType.Elem()
	}
	if aType.Kind() != reflect.Struct {
		return errors.New("input param is not a struct")
	}
	fmt.Printf("Read CSV Type %v:\n", aType)
	fmt.Printf("Read CSV SuperType %v:\n", superType)

	var superFieldNames []string
	var relName string
	if superType != nil {
		superPath := path + filenameKey + "-HEAD.csv"

		file, fileErr := os.Open(superPath)
		if fileErr != nil {
			fmt.Println("Read File Error:", fileErr)
		}
		reader := csv.NewReader(file)

		records, readErr := reader.ReadAll()
		if readErr != nil {
			log.Fatal(readErr)
		}
		superFieldNames = records[headerRowNo]
		relName = records[fieldRowNo][0]
		fmt.Println("Super Records: ", records)
	}

	err := filepath.Walk(path, func(aPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(aPath, filenameKey) {
			return err
		}

		log.Println("FilePath", aPath)

		file, errFile := os.Open(aPath)
		if errFile != nil {
			goazure.Error("Read File Error:", errFile)
		}
		reader := csv.NewReader(file)

		records, errRead := reader.ReadAll()
		if errRead != nil {
			log.Fatal(errRead)
		}
		goazure.Info("Records: ", records)

		var su reflect.Value

		if superType != nil {
			if strings.Contains(aPath, "HEAD") {
				return err
			}

			su = reflect.New(superType)
			sda := su.Elem()
			sin := su.Interface().(models.DataInterface)

			//Build Super Data
			headerRow := records[headerRowNo]
			for i, fieldName := range superFieldNames {
				if fieldName == "" {
					continue
				}

				var value reflect.Value
				switch sda.FieldByName(fieldName).Kind() {
				case reflect.Int32:
					vi64, _ := strconv.ParseInt(headerRow[i], 10, 32)
					v := int32(vi64)
					value = reflect.ValueOf(v)
				case reflect.Int64:
					vi64, _ := strconv.ParseInt(headerRow[i], 10, 32)
					value = reflect.ValueOf(vi64)
				default:
					value = reflect.ValueOf(headerRow[i])
				}

				goazure.Info(i, ". Super Field(", fieldName,":", sda.FieldByName(fieldName).Kind(),"): ", headerRow[i], value)

				sda.FieldByName(fieldName).Set(value)
			}

			ctl.AddData(sin, true, superFieldNames[0])
		}

		var fieldNames []string
		for i, row := range records {
			if strings.Contains(aPath, "HEAD") {
				continue
			}

			if i == fieldRowNo {
				fieldNames = row
			}
			if i > fieldRowNo {
				d := reflect.New(aType)
				da := d.Elem()
				in := d.Interface().(models.DataInterface)

				for j, field := range row {
					if fieldNames[j] == "" {
						continue
					}

					var value reflect.Value
					switch da.FieldByName(fieldNames[j]).Kind() {
					case reflect.Int32:
						vi64, _ := strconv.ParseInt(field, 10, 32)
						v := int32(vi64)
						value = reflect.ValueOf(v)
					case reflect.Int64:
						vi64, _ := strconv.ParseInt(field, 10, 32)
						value = reflect.ValueOf(vi64)
					case reflect.Float32:
						vf32, _ := strconv.ParseFloat(field, 32)
						v := float32(vf32)
						value = reflect.ValueOf(v)
					case reflect.Float64:
						vf64, _ := strconv.ParseFloat(field, 64)
						value = reflect.ValueOf(vf64)
					default:
						value = reflect.ValueOf(field)
					}

					//fmt.Println("Field(", fieldNames[j],":", da.FieldByName(fieldNames[j]).Kind(),"): ", field, value)

					da.FieldByName(fieldNames[j]).Set(value)
					if superType != nil {
						da.FieldByName(relName).Set(su)
					}

				}

				goazure.Info(fmt.Sprintf("Da: %T %v", da, da))
				goazure.Info(fmt.Sprintf("In: %T %v", in, in))

				ctl.AddData(in, true, fieldNames[0])
			}
		}

		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	return err
}