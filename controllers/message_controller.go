package controllers

import (
	"github.com/AzureTech/goazure"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/AzureRelease/boiler-server/models"
	"strconv"
	"os"
	"encoding/csv"
	"log"
	"io/ioutil"
	"strings"
)

type MessageController struct {
	goazure.Controller
}

var MsgCtl *MessageController

const msgDefautlsPath string = "models/properties/message_defaults/"

func (msgCtl *MessageController) InitMessageFormatter() {

}

func (msgCtl *MessageController) InitMessageTags() {
	generateMessageTags()
	generateMessageFormatters()
}

func generateMessageFormatters() {

	const msgFmtPath string = msgDefautlsPath + "formatters/"

	files,_ := ioutil.ReadDir(msgFmtPath)
	fmt.Println("MessageFormatters FileCount: ", len(files))
	for i, v := range files {
		fmt.Println("File: ", i, v)
	}

	const headerRowNo int = 0
	const fieldRowNo int = 1

	msgType := func(typeId int32) *models.MessageType {
		t := models.MessageType{ TypeId: typeId }
		err := DataCtl.ReadData(&t, "TypeId")
		if err != nil {
			fmt.Println("Read Error: ", err)
		}
		return &t
	}

	msgTag := func(tagId int32) *models.MessageTag {
		tag := models.MessageTag{ TagId: tagId }
		err := DataCtl.ReadData(&tag, "TagId")
		if err != nil {
			fmt.Println("Read Error: ", err)
		}
		return &tag
	}

	lenErrCount := 0
	var lenErrFmts []models.MessageFormatter

	for page := 1; page <= len(files); page++ {
		filename := msgFmtPath + "message_formatter-" + strconv.Itoa(page) + ".csv"
		file, errFile := os.Open(filename)
		if errFile != nil {
			fmt.Println("Read File Error:", errFile)
		}
		reader := csv.NewReader(file)

		records, errRead := reader.ReadAll()
		if errRead != nil {
			log.Fatal(errRead)
		}
		fmt.Println("Records: ", records)

		//build msg types
		const namePrefix string = "锅炉数据采集监测终端"

		headerRow := records[headerRowNo]
		for i, v := range headerRow {
			fmt.Println("HeaderRow: ", i, v)
		}

		typeId, _ := strconv.ParseInt(headerRow[0], 10, 32);
		typeStatus := headerRow[len(headerRow) - 1]
		typeRemark := headerRow[1]
		typeName := typeRemark
		if strings.Contains(typeName, namePrefix) {
			typeName = typeName[len(namePrefix) : strings.IndexRune(typeRemark, '（')]
		}
		messageType := models.MessageType{}
		messageType.TypeId = int32(typeId)
		messageType.Name = typeName
		messageType.Remark = typeRemark
		messageType.From = headerRow[len(headerRow) - 2]

		addMessageType(messageType)

		for i, row := range records {
			if i <= fieldRowNo {

			} else {
				tagId, _ := strconv.ParseInt(row[1], 10, 32)
				length, _ := strconv.ParseInt(row[4], 10, 32)
				start, _ := strconv.ParseInt(row[5], 10, 32)
				//end, _ := strconv.ParseInt(row[6], 10, 32)

				var defaults string
				switch tagId {
				case 1:
					defaults = "0xAC 0xEB"
				case 4:
					defaults = typeStatus
				case 999:
					defaults = "0xAF 0xED"
				}

				idx := i - fieldRowNo - 1
				msgFmt := models.MessageFormatter{}
				msgFmt.FormatterId = int32(page * 100 + idx)
				msgFmt.Type = msgType(int32(page))
				msgFmt.Tag = msgTag(int32(tagId))	//包头
				msgFmt.SequenceNumber = int32(idx)
				msgFmt.StartPoint = int32(start)

				msgFmt.Name = strconv.Itoa(page) + ". " + row[2]
				//msgFmt.NameEn = strconv.Itoa(page) + ". " +msgTag(int(tagId)).NameEn
				msgFmt.Length = int32(length)
				msgFmt.Default = defaults

				addMessageFormatter(msgFmt)

				if length != int64(msgTag(int32(tagId)).Length) {
					log.Println("Message Formatter Length Error!")
					lenErrCount++
					lenErrFmts = append(lenErrFmts, msgFmt)
				}
			}
		}
	}

	fmt.Println("\n===================\n, Length Error Count: ", lenErrCount)
	fmt.Println(lenErrFmts)
}


func generateMessageTags() {
	const msgTagsPath string = msgDefautlsPath + "tags/"
	const fieldRow int = 0

	filename := msgTagsPath + "message_tag.csv"
	file, errFile := os.Open(filename)
	if errFile != nil {
		fmt.Println("Read File Error:", errFile)
	}
	reader := csv.NewReader(file)

	records, errRead := reader.ReadAll()
	if errRead != nil {
		log.Fatal(errRead)
	}
	fmt.Println("Message Tags Records: ", records)

	for i, row := range records {
		if i <= fieldRow {

		} else {
			tagId, _ := strconv.ParseInt(row[0], 10, 32)
			name := row[1]
			nameEn := row[2]
			column := row[3]
			length, _ := strconv.ParseInt(row[4], 10, 32)
			dataType := row[5]
			remark := row[6]

			msgTag := models.MessageTag{}

			msgTag.TagId = int32(tagId)
			msgTag.Name = name
			msgTag.NameEn = nameEn
			msgTag.Remark = remark
			msgTag.Column = column
			msgTag.DataType = dataType
			msgTag.Length = int32(length)

			addMessageTag(msgTag)
		}
	}

	/*
	addMessageTag(models.MessageTag{
		TagId: 1,
		MyUidObject: models.MyUidObject{
			Name: "包头",
			NameEn: "Header",
			Remark: "用2个字节来表示协议头，0xAC 0xEB是信息帧开头的关键字",
		},
		Column: "Header",
	})

	addMessageTag(models.MessageTag{
		TagId: 2,
		MyUidObject: models.MyUidObject{
			Name: "有效数据长度",
			NameEn: "Content Length",
			Remark: "代表当前数据包中的有效数据的长度，字节序使用BIG_ENDIAN",
		},
		Column: "ContentLength",
	})

	addMessageTag(models.MessageTag{
		TagId: 3,
		MyUidObject: models.MyUidObject{
			Name: "流水号",
			NameEn: "Serial Number",
			Remark: "从零开始递增的流水号，逐条计数到0xFFFF后归为0，上传下发数据包分别计数",
		},
		Column: "SerialNumber",
	})

	addMessageTag(models.MessageTag{
		TagId: 4,
		MyUidObject: models.MyUidObject{
			Name: "状态消息码",
			NameEn: "Status",
			Remark: "表示本数据包的功能，16进制，为0xA0",
		},
		Column: "Status",
		Length: 1,
	})

	addMessageTag(models.MessageTag{
		TagId: 5,
		MyUidObject: models.MyUidObject{
			Name: "终端身份识别码",
			NameEn: "Identifier",
			Remark: "终端身份识别码由平台产生",
		},
		Column: "Identifier",
		DataType: "CHAR",
		Length: 6,
	})

	addMessageTag(models.MessageTag{
		TagId: 6,
		MyUidObject: models.MyUidObject{
			Name: "响应数据",
			NameEn: "ResponseCode",
			Remark: "10：成功，00：失败原因A，01：失败原因B，02：失败原因C",
		},
		Column: "ResponseCode",
		Length: 1,
	})

	addMessageTag(models.MessageTag{
		TagId: 7,
		MyUidObject: models.MyUidObject{
			Name: "登录密码",
			NameEn: "Password",
			Remark: "在注册成功的情况下由平台分配登陆密码",
		},
		Column: "Password",
	})

	addMessageTag(models.MessageTag{
		TagId: 8,
		MyUidObject: models.MyUidObject{
			Name: "平台系统时间",
			NameEn: "Server Date",
			Remark: "在登陆成功的情况下，将平台时间下发至终端，每次登陆时进行一次时间校准",
		},
		Column: "ServerDate",
		DataType: "CHAR",
		Length: 14,
	})

	addMessageTag(models.MessageTag{
		TagId: 9,
		MyUidObject: models.MyUidObject{
			Name: "数据格式版本号",
			NameEn: "Version",
			Remark: "V1.0：0x10，V2.0：0x20，V3.0:0x30",
		},
		Column: "Version",
		Length: 1,
	})

	addMessageTag(models.MessageTag{
		TagId: 10,
		MyUidObject: models.MyUidObject{
			Name: "锅炉编号",
			NameEn: "Boiler Number",
			Remark: "被监测锅炉的编号，采用ASCII码，用16进制表示",
		},
		Column: "BoilerNumber",
		DataType: "CHAR",
		Length: 2,
	})

	addMessageTag(models.MessageTag{
		TagId: 11,
		MyUidObject: models.MyUidObject{
			Name: "上报时间",
			NameEn: "Upload Date",
			Remark: "终端设备在上报数据时的系统时间，采用ASCII码，并用16进制表示，如32 30 31 32 2D 30 37 2D 32 36 20 31 31 3A 33 35 3A 30 30代表的时间是2012-07-26 11:35:00",
		},
		Column: "UploadDate",
		DataType: "CHAR",
		Length: 14,
	})


	// Temperature Channels
	tempChanDefaults := []string{
		"蒸汽温度",
		"给水温度（冷水）",
		"给水温度（热水）",
		"排烟温度（节能器后）",
		"排烟温度（节能器前）",
		"环境温度",
	}
	for i := 1; i <= 12 ; i++ {
		remark := func() string {
			def := ""
			if i <= len(tempChanDefaults) {
				def = "，缺省配置：" + tempChanDefaults[i - 1]
			}

			return "具体参数由通道配置文件确定" + def
		}

		addMessageTag(models.MessageTag{
			TagId: 100 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: "温度通道" + strconv.Itoa(i),
				NameEn: "Temperature Channel " + strconv.Itoa(i),
				Remark: remark(),
			},
			Column: "TemperatureChannel" + strconv.Itoa(i),
		})
	}

	// Analogue Channels
	analogChanDefaults := []string{
		"蒸汽压力",
		"排烟氧量",
	}
	for i := 1; i <= 12 ; i++ {
		remark := func() string {
			def := ""
			if i <= len(analogChanDefaults) {
				def = "，缺省配置：" + analogChanDefaults[i - 1]
			}

			return "具体参数由通道配置文件确定" + def
		}

		addMessageTag(models.MessageTag{
			TagId: 200 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: "模拟量通道" + strconv.Itoa(i),
				NameEn: "Analogue Channel " + strconv.Itoa(i),
				Remark: remark(),
			},
			Column: "AnalogueChannel" + strconv.Itoa(i),
		})
	}

	// Switch Signal I/O Channels
	switchDefaults := []string{
		"锅炉运行",
		"点火信号",
		"极低水位报警",
		"超高水位报警",
		"超压报警",
		"综合报警",
		"熄火故障",
		"软水硬度报警",
	}
	inputRange := 32
	outputRange := 16
	for i := 1; i <= inputRange + outputRange ; i++ {
		ioSign := func() []string{
			if i <= inputRange {
				return []string{"输入", "Input"}
			} else {
				return []string{"输出", "Output"}
			}
		}
		remark := func() string {
			def := ""
			if i <= len(switchDefaults) {
				def = "，缺省配置：" + switchDefaults[i - 1]
			}

			return "具体参数由通道配置文件确定" + def
		}

		s_no := func() int{
			if i > inputRange {
				return i - inputRange
			} else {
				return i
			}
		}

		addMessageTag(models.MessageTag{
			TagId: 300 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: "开关量" + ioSign()[0] + "通道" + strconv.Itoa(s_no()),
				NameEn: "Switch Signal " + ioSign()[1] + " Channel " + strconv.Itoa(s_no()),
				Remark: remark(),
			},
			Column: "SwitchSignal" + ioSign()[1] + "Channel" + strconv.Itoa(s_no()),
		})
	}

	// Calculation Parameters
	calcDefaults := []string{
		"热效率（终端计算）",
		"q2（终端计算）",
		"q3（终端计算）",
		"q4（终端计算）",
		"q5（终端计算）",
		"q6（终端计算）",
		"排烟过量空气系数（终端计算）",
	}
	for i := 1; i <= 12 ; i++ {
		remark := func() string {
			def := ""
			if i <= len(calcDefaults) {
				def = "，缺省配置：" + calcDefaults[i - 1]
			}

			return "具体参数由通道配置文件确定" + def
		}

		addMessageTag(models.MessageTag{
			TagId: 400 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: "计算参数" + strconv.Itoa(i),
				NameEn: "Calculation Parameter " + strconv.Itoa(i),
				Remark: remark(),
			},
			Column: "CalculationParameter" + strconv.Itoa(i),
		})
	}

	//Energy Efficiency Factor
	energeDefaults := [][]string{
		[]string{"Coal_QNETVAR燃料收到基低位发热量", "", "", ""},
		[]string{"Coal_AAR燃料收到基灰分", "", "", ""},
		[]string{"Coal_MAR燃料收到基水分", "", "", ""},
		[]string{"Coal_VDAF干燥无灰基挥发分", "", "", ""},
		[]string{"Coal_CLZ炉渣可燃物含量", "", "", ""},
		[]string{"Coal_CLM漏煤可燃物含量", "", "", ""},
		[]string{"Coal_CFH飞灰可燃物含量", "", "", ""},
		[]string{"Coal_DED锅炉额定负荷", "", "", ""},
		[]string{"Coal_DSC锅炉实测负荷", "", "", ""},
		[]string{"Coal_ALZ炉渣含灰量占入炉煤总灰量百分比", "", "", ""},
		[]string{"Coal_ALM漏煤含灰量占入炉煤总灰量百分比", "", "", ""},
		[]string{"Coal_AFH飞灰含灰量占入炉煤总灰量百分比", "", "", ""},
		[]string{"Coal_Q3气体未完全燃烧热损失", "", "", ""},
		[]string{"Coal_M燃料计算系数", "", "", ""},
		[]string{"Coal_N燃料计算系数", "", "", ""},
		[]string{"Coal_TLZ炉渣温度", "", "", ""},
		[]string{"Coal_CT_LZ炉渣焓", "", "", ""},
		[]string{"Gas_DED锅炉额定负荷", "", "", ""},
		[]string{"Gas_DSC锅炉实测负荷", "", "", ""},
		[]string{"Gas_APY排烟处过量空气系数", "", "", ""},
		[]string{"Gas_Q3气体未完全燃烧热损失", "", "", ""},
		[]string{"Gas_M燃料计算系数", "", "", ""},
		[]string{"Gas_N燃料计算系数", "", "", ""},
	}

	// Setting Parameters
	settingDefaults := [][]string{
		[]string{"上报启闭", "Upload On-Off", "UploadOnOff", "0x00：关闭主动上报，0x01：开启主动上报"},
		[]string{"上报周期", "Upload Period", "UploadPeriod", "0x0000~0xFFFF，单位为s"},
	}
	for i := 1; i <= 12 ; i++ {
		name := func() string {
			if i <= len(settingDefaults) {
				return settingDefaults[i - 1][0]
			} else {
				return "配置参数预留" + strconv.Itoa(i - len(settingDefaults))
			}
		}

		nameEn := func() string {
			if i <= len(settingDefaults) {
				return settingDefaults[i - 1][1]
			} else {
				return "Reserved Setting Parameter " + strconv.Itoa(i - len(settingDefaults))
			}
		}

		remark := func() string {
			if i <= len(settingDefaults) {
				return settingDefaults[i - 1][3]
			} else {
				return "具体参数由通道配置文件确定"
			}
		}

		col := func() string {
			if i <= len(settingDefaults) {
				return settingDefaults[i - 1][2]
			} else {
				return "ReservedSettingParameter" + strconv.Itoa(i - len(settingDefaults))
			}
		}

		length := func() int32 {
			if i == 1 {
				return 1
			} else {
				return 2
			}
		}


		addMessageTag(models.MessageTag{
			TagId: 500 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: name(),
				NameEn: nameEn(),
				Remark: remark(),
			},
			Column: col(),
			Length: length(),
		})
	}


	// Terminals
	terminalInfos := [][]string{
		[]string{"终端故障", "Terminal Error", "TerminalError", "0x00：无故障，0x01：故障A，0x02：故障B，0x03：故障"},
		[]string{"终端重启", "Terminal Restart", "TerminalRestart", "0x00：不重启，0x01：重启"},
	}

	for i, v := range terminalInfos {
		addMessageTag(models.MessageTag{
			TagId: 600 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: v[0],
				NameEn: v[1],
				Remark: v[3],
			},
			Column: v[2],
			Length: 1,
		})
	}

	// Reserve
	for i := 1; i <= 4 ; i++ {
		addMessageTag(models.MessageTag{
			TagId: 800 + int32(i),
			MyUidObject: models.MyUidObject{
				Name: "保留字" + strconv.Itoa(i),
				NameEn: "Reserve " + strconv.Itoa(i),
				Remark: "保留字段",
			},
			Column: "Reserve" + strconv.Itoa(i),
		})
	}

	addMessageTag(models.MessageTag{
		TagId: 998,
		MyUidObject: models.MyUidObject{
			Name: "CRC校验",
			NameEn: "CRC",
			Remark: "对有效数据的CRC校验",
		},
		Column: "Crc",
	})

	addMessageTag(models.MessageTag{
		TagId: 999,
		MyUidObject: models.MyUidObject{
			Name: "包尾",
			NameEn: "Tailer",
			Remark: "用2个字节来表示协议尾，0xAF 0xED是信息帧结束的关键字",
		},
		Column: "Tailer",
	})
	*/
}

func addMessageType(msgType models.MessageType) error {

	err := DataCtl.AddData(&msgType, true, "TypeId")

	return err
}

func addMessageFormatter(formatter models.MessageFormatter) error {

	err := DataCtl.AddData(&formatter, true, "FormatterId")

	return err
}

func addMessageTag(tag models.MessageTag) error {

	if tag.DataType == "" {
		tag.DataType = "HEX"
	}

	if tag.Length == 0 {
		tag.Length = 2
	}

	err := DataCtl.AddData(&tag, true, "TagId")

	return err
}

func ByteToInt(Bytes []byte) (int8) {
	b_buf := bytes.NewBuffer(Bytes)
	var x int8
	err := binary.Read(b_buf, binary.BigEndian, &x)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(x)

	return x
}