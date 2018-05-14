package tsssdk

/*
#include "tss_sdk.h"
#include "tss_sdk_anti.h"
*/
import "C"
import (
	"unsafe"
	"zeus/iserver"
)

//export onSendDataToClient
func onSendDataToClient(roleid C.TSS_UINT64, data *C.uchar, datalen C.int) C.TssSdkProcResult {
	antiData := C.GoBytes(unsafe.Pointer(data), datalen)

	user := iserver.GetSrvInst().GetEntityByDBID("Player", uint64(roleid))
	if user == nil {
		return C.TSS_SDK_PROC_INTERNAL_ERR
	}

	if err := user.RPC(iserver.ServerTypeClient, "TssData", antiData); err != nil {
		return C.TSS_SDK_PROC_INTERNAL_ERR
	}

	return C.TSS_SDK_PROC_OK
}
