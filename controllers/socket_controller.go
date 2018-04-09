package controllers

import (
	"github.com/AzureTech/goazure"
	"github.com/AzureTech/goazure/orm"

	"strconv"
	"strings"
	"fmt"
	"os"
	"bytes"
	"encoding/binary"
	"net"
	"github.com/AzureRelease/boiler-server/dba"
	"time"
)

type SocketController struct {
	MainController
}

var SocketCtrl *SocketController = &SocketController{}
func c9(Code string)([]byte) {
	words_1:="\xac\xeb\x00\x09\x00\x00\xc9"+Code+"\x00\x00\xaf\xed"
	buf :=[]byte(words_1)
	copy(buf[13:15],CRC16(buf[4:13]))
	return buf
}

func (ctl *SocketController)c0(b []byte,Code string,ver int32)([]byte) {
	words_1:="\xac\xeb\x01\x62\x00\x00\xc0"+Code
	buf:=append([]byte(words_1),IntToByte(ver)...)
	fmt.Println(buf)
	buf=append(buf,b...)
	fmt.Println(buf)
	words_2:="\x00\x00\xaf\xed"
	buf=append(buf,[]byte(words_2)...)
	copy(buf[358:360],CRC16(buf[4:358]))
	return buf
}

func Send(conn net.Conn,code string) {
	buf:=c9(code)
	n, err := conn.Write(buf)
	if err != nil {
		goazure.Error("%s%s","Write error:", err)
	} else {
		goazure.Info(fmt.Sprintf("Write %d bytes, content is %x\n", n, string(buf[:n])))
	}
	//conn.Write(buffer)
	fmt.Println("send over")
}
func Receive(conn net.Conn) ([]byte) {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		goazure.Error("Receive error:", err)
	} else {
		goazure.Info(fmt.Sprintf("Receive %d bytes, content is %s\n", n, string(buf[:n])))
	}
	fmt.Println(buf[:n])
	return buf[:n]
}
type Info struct {
	Sn string
	CurrMessage string
}
func (ctl *SocketController)Message(Code string)(string) {
	var info Info
	info.Sn=Code
	fmt.Println("infoSn",info.Sn)
	sql:="select curr_message from issued_message where sn=?"
	if err:=dba.BoilerOrm.Raw(sql,Code).QueryRow(&info);err!=nil{
		goazure.Error("Query curr_message Error",err)
	}
	fmt.Println("下发的报文:",info.CurrMessage)
	return info.CurrMessage
}

func SendConfig(conn net.Conn,Code string) {
		b:=SocketCtrl.Message(Code)
		buf:=[]byte(b)
		n,err:=conn.Write(buf)
		if err!=nil {
			goazure.Error("%s%s","Write error:", err)
		} else {
			goazure.Info(fmt.Sprintf("Write %d bytes, content is %x\n", n, string(buf[:n])))
		}
	fmt.Println("send over")
}

//下发配置
func (ctl *SocketController)SocketConfigSend(Code string)([]byte) {
	server := "47.100.0.27:18887"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return nil
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return nil
	}
	goazure.Info("connect success")
	SendConfig(conn,Code)
	buf:=Receive(conn)
	conn.Close()
	return buf
}

//重启
func (ctl *SocketController)SocketTerminalRestart(code string)([]byte) {
	server := "47.100.0.27:18887"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return nil
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return nil
	}
	goazure.Info("connect success")
	Send(conn,code)
	buf:=Receive(conn)
	conn.Close()
	return buf
}
func b4()([]byte){
	words_2 := "\xac\xeb\x00\x0b\x00\x00\xb4\x30\x33\x30\x30\x39\x30\x30\x31\x00\x00\xaf\xed"
	buf     := []byte(words_2)
	copy(buf[15:17], CRC16(buf[4:15]))
	goazure.Info(fmt.Sprintf("%x", buf[4:15]))
	goazure.Info(fmt.Sprintf("%x", CRC16(buf[4:15])))
	goazure.Info(fmt.Sprintf("%x", buf))

	return buf
}

func b7()([]byte){
	words_2 :=
			"\xac\xeb\x00\x5e\x00\x00\xb7\x30\x33\x30\x30\x39\x30\x10\x30\x31" +
			"\x06\x0d\x02\x19\x00\x00\x00\x00\x00\xcb\x04\xda\x05\x99\x00\x00" +
			"\x00\x00\x1d\x4c\x01\xf4\x07\xd0\x00\x14\x00\x3c\x01\x7c\x00\x00" +
			"\xde\xbc\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x10\x30\xaf\xed"
	buf := []byte(words_2)
	copy(buf[98:100], CRC16(buf[4:98]))
	goazure.Info(fmt.Sprintf("%x", buf[4:98]))
	goazure.Info(fmt.Sprintf("%x", CRC16(buf[4:98])))
	goazure.Info(fmt.Sprintf("%x", buf))
	return buf
}

func b5(code string, isOn bool, period int)([]byte){
	var buf			[]byte
	var msgFmt 		[]orm.Params
	var msgBuilder 	[]orm.Params
	var msgStatus	[]orm.Params
	dummy := "\x00\x00"

	rawFmt := "SELECT a.id,a.name, a.elem,a.attr, a.data_desc, a.elem_sseq, b.atom_len, b.atom_type, b.dis_len, b.dis_type,b.dec_len " +
		"FROM pf_message_data AS a, pf_atomic_data AS b " +
		"WHERE (a.name = 'header_msg' " +
		fmt.Sprintf("OR a.name = (select name from pf_message_list where id = %d) ", 181) +
		"OR a.name = 'tailer_msg') " +
		"AND a.elem = b.name " +
		"ORDER BY a.id,elem_sseq;"
	if num, err := dba.BoilerOrm.Raw(rawFmt).Values(&msgFmt); err != nil || num == 0 {
		goazure.Error("Get Message Format Error:", err, num)
	}

	rawBuilder := "SELECT * FROM `pf_message_list` " +
		fmt.Sprintf("WHERE	`id` = %d;", 181)
	if num, err := dba.BoilerOrm.Raw(rawBuilder).Values(&msgBuilder); err != nil || num == 0 {
		goazure.Error("Get Message Builder Error:", err, num)
	}

	rawSecret := fmt.Sprintf("SELECT * FROM `boiler_term_status` WHERE `Boiler_term_id` = %s", code)
	if num, err := dba.BoilerOrm.Raw(rawSecret).Values(&msgStatus); err != nil || num == 0 {
		goazure.Error("Get Message Status Error:", err, num)
		//return buf
	}

	var status orm.Params
	if len(msgStatus) > 0 { status = msgStatus[0] }
	pwd := int64(0)
	if status["Boiler_term_pwd"] != nil { pwd, _ = strconv.ParseInt(status["Boiler_term_pwd"].(string), 10, 64) }
	sysTime := time.Now().Format("2006-01-02 15:04:05")
	goazure.Warn("sysTime:", sysTime)
	on := "\x01"
	if !isOn {
		on = "\x00"
	}

	sn		:= IntToByte(1)
	//"%2x%2x%2x%1x%6s%1x%2x%1x%2s%26s%1x%2x%2x%2x%2x%2x%2x"
	body 	:= fmt.Sprintf("%2s%1s", sn, "\xb5")
	body 	+= fmt.Sprintf("%6s%1s%2x%1s%2x", code, "\x10", pwd, on, period)
	body	+= fmt.Sprintf("%19s%1s%2s%2s%2s%2s", sysTime, "0x00", dummy, dummy, dummy, dummy)
	bodyBuf := []byte(body)

	availLen := IntToByte(int32(len(body)))
	header	:= fmt.Sprintf("%2s%2s", "\xac\xeb", availLen)
	crc     := CRC16(bodyBuf)
	tailer 	:= fmt.Sprintf("%2s%2s", crc, "\xaf\xed")
	msg		:= header + body + tailer

	buf 	= []byte(msg)
	//copy(buf[98:100], CRC16(buf[4:98]))
	goazure.Info("length:", availLen, len(body))
	goazure.Info(fmt.Sprintf("%x\n", []byte(body)))
	goazure.Info(fmt.Sprintf("%s\n", []byte(body)))
	goazure.Info(fmt.Sprintf("%x\n", CRC16(crc)))
	goazure.Info(fmt.Sprintf("%x\n", buf))

	return buf
}

func Sender(conn net.Conn, code string, isOn bool, period int) {
	//words := "hello world!"
	//temp1   := fmt.Sprintf("%2x",255)
	/*
		Avail_Len := Common.IntToByte(11)
		Sn        := Common.IntToByte(0)
		CRC       := Common.IntToByte(17846)
		fmt.Printf("%x\n",CRC)
		//words_2   := fmt.Sprintf("%2s%2s%2s%1s%6s%2s%2s","\xac\xeb",Avail_Len,Sn,"\xa0","030090",CRC,"\xaf\xed")

		words_2   := fmt.Sprintf("%2s%2s%2s%1s%6s%2s%2s%2s","\xac\xeb",Avail_Len,Sn,"\xb4","030090","01","","\xaf\xed")
		//words_2   := fmt.Sprintf("%2s%2s%2s%1s%6s%2s%2s%2s","\xac\xeb",Avail_Len,Sn,"\xb8","030090","01","","\xaf\xed")
		//words_2   := fmt.Sprintf("%2s%2s%2s%1s%6s%2s%1s%2s%2s","\xac\xeb",Avail_Len,"\x00\x61","\xb3","030090","08","\x10","","\xaf\xed")
		//words_2   := fmt.Sprintf("%d","\x00")
		buf       := []byte(words_2)
		copy(buf[15:17],Common.CRC16(buf[4:15]))
		fmt.Printf("%x\n",buf[4:15])
		fmt.Printf("%x\n",Common.CRC16(buf[4:15]))
	*/
	buf := b5(code, isOn, period)
	//buf := b4()
	n, err := conn.Write(buf)
	if err != nil {
		goazure.Error("%s%s","Write error:", err)
	} else {
		goazure.Info(fmt.Sprintf("Write %d bytes, content is %x\n", n, string(buf[:n])))
	}
	//conn.Write(buffer)
	fmt.Println("send over")
	return
}

func Recive(conn net.Conn){
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		goazure.Error("Receive error:", err)
	} else {
		goazure.Info(fmt.Sprintf("Receive %d bytes, content is %s\n", n, string(buf[:n])))
	}
	// 注册请求测试
	// *报文头
	goazure.Info(fmt.Sprintf("%s,header:%x", conn.RemoteAddr().String(), buf[0:2]))
	// 有效数据长度
	goazure.Info(fmt.Sprintf("%s,len   :%d", conn.RemoteAddr().String(), ByteToInt(buf[2:4])))
	// 流水号
	goazure.Info(fmt.Sprintf("%s,sn    :%d", conn.RemoteAddr().String(), ByteToInt(buf[4:6])))
	// 状态消息码
	goazure.Info(fmt.Sprintf("%s,stauts:%x", conn.RemoteAddr().String(), buf[6]))
	// 终端身份识别码
	goazure.Info(fmt.Sprintf("%s,term  :%s", conn.RemoteAddr().String(), buf[7:13]))
	// 注册响应数据
	goazure.Info(fmt.Sprintf("%s,reply :%x", conn.RemoteAddr().String(), buf[13]))
	// 登陆密码
	goazure.Info(fmt.Sprintf("%s,pwd   :%x", conn.RemoteAddr().String(), ByteToInt(buf[14:16])))
	// CRC校验
	goazure.Info(fmt.Sprintf("%s,crc1  :%x", conn.RemoteAddr().String(), ByteToInt(buf[16:18])))
	crc := CRC16(buf[4:16])
	goazure.Info(fmt.Sprintf("%s,crc2  :%d", conn.RemoteAddr().String(), crc))
	// Log("%s,crc int16  :%x\n",conn.RemoteAddr().String(),int16(crc))
	// Log("%s,crc int32  :%x\n",conn.RemoteAddr().String(),int32(crc))
	// 包尾
	goazure.Info(fmt.Sprintf("%s,end   :%x\n", conn.RemoteAddr().String(), buf[18:20]))
	return
}

func BoilerSocketClient1(){
	b4()
	b7()
	Ip, _ := strconv.Atoi(strings.Replace("192.168.0.62",".","",-1))
	IpMax, _ := strconv.Atoi(strings.Replace("192.168.0.254",".","",-1))
	IpMin, _ := strconv.Atoi(strings.Replace("192.168.0.1",".","",-1))
	if Ip <= IpMax && Ip >= IpMin {
		fmt.Printf("%d\n", Ip)
	}
}

func (ctl *SocketController) SocketClientMessageSend(code string, isOn bool, period int) {
	server := "139.196.152.127:12000"
	//server := "127.0.0.1:12020"
	//server := "223.104.25.227:2577"
	fmt.Print(server[0:2])
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	goazure.Info("connect success")
	Sender(conn, code, isOn, period)
	Recive(conn)
	conn.Close()
}

func CRC(data []byte)(int,error){
	var crc int = 65535
	var table = []int{0, 49345, 49537, 320, 49921, 960, 640, 49729, 50689, 1728,
					  1920, 51009, 1280, 50625, 50305, 1088, 52225, 3264, 3456, 52545,
					  3840, 53185, 52865, 3648, 2560, 51905, 52097, 2880, 51457, 2496,
					  2176, 51265, 55297, 6336, 6528, 55617, 6912, 56257, 55937, 6720,
					  7680, 57025, 57217, 8000, 56577, 7616, 7296, 56385, 5120, 54465,
					  54657, 5440, 55041, 6080, 5760, 54849, 53761, 4800, 4992, 54081,
					  4352, 53697, 53377, 4160, 61441, 12480, 12672, 61761, 13056, 62401,
					  62081, 12864, 13824, 63169, 63361, 14144, 62721, 13760, 13440, 62529,
					  15360, 64705, 64897, 15680, 65281, 16320, 16000, 65089, 64001, 15040,
					  15232, 64321, 14592, 63937, 63617, 14400, 10240, 59585, 59777, 10560,
					  60161, 11200, 10880, 59969, 60929, 11968, 12160, 61249, 11520, 60865,
					  60545, 11328, 58369, 9408, 9600, 58689, 9984, 59329, 59009, 9792,
					  8704, 58049, 58241, 9024, 57601, 8640, 8320, 57409, 40961, 24768,
					  24960, 41281, 25344, 41921, 41601, 25152, 26112, 42689, 42881, 26432,
					  42241, 26048, 25728, 42049, 27648, 44225, 44417, 27968, 44801, 28608,
					  28288, 44609, 43521, 27328, 27520, 43841, 26880, 43457, 43137, 26688,
					  30720, 47297, 47489, 31040, 47873, 31680, 31360, 47681, 48641, 32448,
					  32640, 48961, 32000, 48577, 48257, 31808, 46081, 29888, 30080, 46401,
					  30464, 47041, 46721, 30272, 29184, 45761, 45953, 29504, 45313, 29120,
					  28800, 45121, 20480, 37057, 37249, 20800, 37633, 21440, 21120, 37441,
					  38401, 22208, 22400, 38721, 21760, 38337, 38017, 21568, 39937, 23744,
					  23936, 40257, 24320, 40897, 40577, 24128, 23040, 39617, 39809, 23360,
					  39169, 22976, 22656, 38977, 34817, 18624, 18816, 35137, 19200, 35777,
					  35457, 19008, 19968, 36545, 36737, 20288, 36097, 19904, 19584, 35905,
					  17408, 33985, 34177, 17728, 34561, 18368, 18048, 34369, 33281, 17088,
					  17280, 33601, 16640, 33217, 32897, 16448}
	l := len(data)
	for i := 0; i < l; i++{
		b := int(data[i])
		crc = crc >> 8 ^ table[(crc ^ b) & 0xff]
	}
	return crc,nil
}

func CRC16(data []byte)([]byte){
	crc, crcErr := CRC(data)
	CRCByte     := make([]byte,2)
	if crcErr != nil {
		fmt.Fprintf(os.Stdout,"Crc Check Err:%s\n", crcErr)
	}
	copy(CRCByte,IntToByte(int32(int16(crc))))
	return CRCByte
}

func IntToByte(Int int32)([]byte){
	b_buf := bytes.NewBuffer([]byte{})
	err := binary.Write(b_buf, binary.BigEndian, Int)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	r_buf := []byte{b_buf.Bytes()[2],b_buf.Bytes()[3]}
	return r_buf
}