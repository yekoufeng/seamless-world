package tlog

import (
	"fmt"
	"net"
	"reflect"

	log "github.com/cihub/seelog"
)

type iName interface {
	Name() string
}

var tlogger log.LoggerInterface
var udpConn net.Conn
var sendBuf chan []byte

func init() {
	var err error
	tlogger, err = log.LoggerFromConfigAsFile("../res/config/tlog.xml")
	if err != nil {
		panic(err)
	}
}

// Flush 退出时清空缓冲区
func Flush() {
	tlogger.Flush()
}

// ConfigRemoteAddr 配置远程入库地址
func ConfigRemoteAddr(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return err
	}

	udpConn, err = net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		return err
	}

	sendBuf = make(chan []byte, 1000)
	go sendLoop()

	log.Info("TLOG UDP Connect to ", addr)
	return nil
}

// Format 按照Tlog日志格式打印日志
func Format(content iName) {
	str := fmt.Sprintf("%s%s", content.Name(), toFormatStr(content))
	tlogger.Infof("%s", str)
	if udpConn != nil {
		select {
		case sendBuf <- []byte(str + "\n"):
		default:
		}
		// udpConn.Write([]byte(str + "\n"))
	}
}

func toFormatStr(content interface{}) string {
	var str string
	v := reflect.ValueOf(content)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface()
		str += fmt.Sprintf("|%v", field)
	}
	return str
}

func sendLoop() {
	for {
		select {
		case data := <-sendBuf:
			udpConn.Write(data)
		}
	}
}
