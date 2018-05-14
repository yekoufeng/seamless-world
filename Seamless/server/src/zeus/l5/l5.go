package l5

/*
#cgo CFLAGS: -I ./
#cgo LDFLAGS: -L ./ -lqostrans -lqos_client
#include <stdlib.h>
#include "qos_c.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"time"
	"unsafe"
)

// DataForUpdate 需要调用更新接口的数据
type DataForUpdate struct {
	Modid    int
	Cmd      int
	Ret      int
	Duration int
}

var updateChan = make(chan DataForUpdate, 10240)

func init() {
	go func() {
		for data := range updateChan {
			APIRouteResultUpdate(data.Modid, data.Cmd, data.Ret, data.Duration)
		}
	}()
}

// CallService 发起请求
func CallService(modid, cmd int, fn func(string, uint) error) error {
	now := time.Now()
	ip, port, ret, err := APIGetRoute(modid, cmd)
	defer func() {
		data := DataForUpdate{
			Modid:    modid,
			Cmd:      cmd,
			Ret:      ret,
			Duration: int(time.Since(now).Nanoseconds() / 1000000),
		}
		select {
		case updateChan <- data:
			return
		default:
			log.Println("更新L5结果的channel满了")
		}
	}()
	if err != nil {
		return err
	}
	if len(ip) == 0 {
		return errors.New("l5 get route fail")
	}
	return fn(ip, port)
}

// APIGetRoute 获取ip:port
func APIGetRoute(modid, cmd int) (string, uint, int, error) {
	hostIP := C.CString("")
	port := C.ushort(0)

	defer C.free(unsafe.Pointer(hostIP))

	ret := C.ApiGetRoute_C(
		C.int(modid),
		C.int(cmd),
		C.float(1000),
		hostIP,
		C.uint(32),
		&port,
	)

	if ret == 0 {
		ip := C.GoString(hostIP)
		port := uint(port)
		return ip, port, 0, nil
	}

	return "", 0, int(ret), fmt.Errorf("l5 APIGetRoute fail: %d", ret)
}

// APIRouteResultUpdate 上报服务调用情况, duration毫秒
func APIRouteResultUpdate(modid, cmd, ret, duration int) {
	C.ApiRouteResultUpdate_C(
		C.int(modid),
		C.int(cmd),
		C.int(ret),
		C.int(duration),
	)
}
