package tsssdk

import (
	"time"
)

/*
#include "tss_sdk.h"
#include "tss_sdk_anti.h"
#include "tss_sdk_uic.h"
#include <stdio.h>
#include <string.h>

const TssSdkBusiInterf *busi_interf_;
const TssSdkAntiInterfV3 *anti_interf_;
const TssSdkUicInterfV3 *uic_interf_;

extern TssSdkProcResult onSendDataToClient(const TSS_UINT64 role_id, const unsigned char *data, const int data_len);
// extern TssSdkProcResult onUICCallback(const TssSdkUicChatJudgeResultInfoV2 *result_info);

TssSdkProcResult
on_send_data_to_client(const TssSdkAntiSendDataInfoV3 *send_data_info)
{
	return onSendDataToClient(send_data_info->role_id_, send_data_info->anti_data_, send_data_info->anti_data_len_);
}

TssSdkProcResult
uic_callback_proc(const TssSdkUicChatJudgeResultInfoV3 *result_info)
{
	return TSS_SDK_PROC_OK;
}

int init_anti_interf(unsigned int instance_id)
{
	TssSdkInitInfo init_data;
    memset(&init_data, 0, sizeof(init_data));
    init_data.unique_instance_id_ = instance_id;
    init_data.tss_sdk_conf_ = "./config";
    busi_interf_ = tss_sdk_load(".", &init_data);

    if (busi_interf_ == NULL)
    {
        return -1;
    }

    TssSdkAntiInitInfoV3 init_info;
    init_info.send_data_to_client_ = on_send_data_to_client;
    anti_interf_ =  TSS_SDK_GET_ANTI_INTERF_V3(&init_info);
    if (anti_interf_ == NULL)
    {
        return -2;
	}

	TssSdkUicInitInfoV3 uic_init_info;
    memset(&uic_init_info, 0, sizeof(uic_init_info));
	uic_init_info.on_chat_judge_result_ = uic_callback_proc;
	uic_interf_ = TSS_SDK_GET_UIC_INTERF_V3(&uic_init_info);
	if (uic_interf_ == NULL)
	{
		return -2;
	}

    return 0;
}

int anti_add_user(unsigned char *open_id, unsigned char openid_len,
	unsigned char plat_id, unsigned int world_id,
	TSS_UINT64 rold_id, unsigned int client_ver,
	unsigned int client_ip, const char *role_name,
	const TssSdkUserExtData *ext_data )
{
    TssSdkAntiAddUserInfoV3 user_info;
    memset(&user_info, 0, sizeof(user_info));
    user_info.openid_.openid_ = open_id;
    user_info.openid_.openid_len_ = openid_len;
    user_info.openid_.openid_[user_info.openid_.openid_len_ ] = 0;
    // platid, 0: IOS, 1: Android
    user_info.plat_id_ = plat_id;
    user_info.world_id_ = world_id;
	user_info.role_id_ = rold_id;
	user_info.client_ver_ = client_ver;
	user_info.client_ip_ = client_ip;
	user_info.role_name_ = role_name;
	user_info.user_ext_data_ = ext_data;

	return anti_interf_->add_user_(&user_info);
}

int anti_del_user(unsigned char *open_id, unsigned char openid_len,
	unsigned char plat_id, unsigned int world_id,
	TSS_UINT64 rold_id, const TssSdkUserExtData *ext_data)
{
    TssSdkAntiDelUserInfoV3 user_info;
    memset(&user_info, 0, sizeof(user_info));
    user_info.openid_.openid_ = open_id;
    user_info.openid_.openid_len_ =openid_len;
    user_info.openid_.openid_[user_info.openid_.openid_len_] = 0;
    // platid, 0: IOS, 1: Android
    user_info.plat_id_ = plat_id;
    user_info.world_id_ = world_id;
    user_info.role_id_ = rold_id;
    user_info.user_ext_data_ = ext_data;

    return anti_interf_->del_user_(&user_info);
}

int on_recv_anti_data(unsigned char *open_id, unsigned char openid_len,
	unsigned char plat_id, unsigned int world_id,
	TSS_UINT64 rold_id, const unsigned char *anti_data, unsigned int anti_data_len,
	const TssSdkUserExtData *ext_data)
{
    // call the recv data interface of sdk anti-hacking service to recv package
    TssSdkAntiRecvDataInfoV3 pkg_info;
    memset(&pkg_info, 0, sizeof(pkg_info));
    pkg_info.openid_.openid_ = open_id;
    pkg_info.openid_.openid_len_ = openid_len;
    pkg_info.openid_.openid_[pkg_info.openid_.openid_len_] = 0;
    // platid, 0: IOS, 1: Android
    pkg_info.plat_id_ = plat_id;
    pkg_info.world_id_ = world_id;
    pkg_info.role_id_ = rold_id;
	pkg_info.user_ext_data_ = ext_data;
	pkg_info.anti_data_ = anti_data;
    pkg_info.anti_data_len_ = anti_data_len;

    return anti_interf_->recv_anti_data_(&pkg_info);
}

int judge_user_input_name(char *msg, int msg_len, int *result)
{
	int ret;
	TssSdkUicNameUserInputInfoV3 info;
	memset(&info, 0, sizeof(info));
	info.msg_len_ = msg_len;
	info.msg_ = msg;
	info.door_level_ = 1;
	info.if_replace_ = 0;
	info.local_language_ = 0;

	if (uic_interf_ == NULL)
	{
		return -1;
	}

	ret = uic_interf_->judge_user_input_name_(&info);
	if (ret != 0)
	{
		return ret;
	}
	*result = info.msg_result_flag_;
	return 0;
}

void proc()
{
	busi_interf_->proc_();
}

*/
import "C"
import (
	"encoding/binary"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

// ErrNotInited SDK未初始化
var ErrNotInited = errors.New("SDK not inited")

var stopSig chan bool

// Init 初始化
func Init(uid uint64) {
	ret := C.init_anti_interf(C.uint(uid))
	if ret != 0 {
		panic("TSSDK 初始化失败")
	}

	stopSig = make(chan bool, 1)

	go loop()
}

// Destroy 销毁
func Destroy() {
	stopSig <- true
	C.tss_sdk_unload()
}

func loop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for {
		select {
		case <-stopSig:
			return
		case <-ticker.C:
			C.proc()
		}
	}
}

// ErrAddPlayerFailed 增加用户失败
var ErrAddPlayerFailed = errors.New("TSSSDK add player failed")

// OnPlayerLogin 用户登录时
func OnPlayerLogin(entityid uint64, openid string, platid uint8, worldid uint32, roleid uint64, clientver uint32, clientip string, rolename string) error {
	ipStr := strings.Split(clientip, ":")[0]
	bits := strings.Split(ipStr, ".")
	var ip uint32
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		ip += uint32(b0) << 24
		ip += uint32(b1) << 16
		ip += uint32(b2) << 8
		ip += uint32(b3)
	} else {
		ip = 0
	}

	extData := C.TssSdkUserExtData{}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, entityid)
	extData.user_ext_data_[0] = C.char(data[0])
	extData.user_ext_data_[1] = C.char(data[1])
	extData.user_ext_data_[2] = C.char(data[2])
	extData.user_ext_data_[3] = C.char(data[3])
	extData.user_ext_data_[4] = C.char(data[4])
	extData.user_ext_data_[5] = C.char(data[5])
	extData.user_ext_data_[6] = C.char(data[6])
	extData.user_ext_data_[7] = C.char(data[7])
	extData.ext_data_len_ = 8

	ret := C.anti_add_user((*C.uchar)(unsafe.Pointer(C.CString(openid))), C.uchar(len(openid)),
		C.uchar(platid), C.uint(worldid), C.TSS_UINT64(roleid), C.uint(clientver), C.uint(ip),
		C.CString(rolename), &extData)

	if ret != 0 {
		return ErrAddPlayerFailed
	}

	return nil
}

// ErrDelPlayerFailed 删除用户失败
var ErrDelPlayerFailed = errors.New("TSSSDK del player failed")

// OnPlayerLogout 玩家登出时
func OnPlayerLogout(entityid uint64, openid string, platid uint8, worldid uint32, roleid uint64) error {
	extData := C.TssSdkUserExtData{}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, entityid)
	extData.user_ext_data_[0] = C.char(data[0])
	extData.user_ext_data_[1] = C.char(data[1])
	extData.user_ext_data_[2] = C.char(data[2])
	extData.user_ext_data_[3] = C.char(data[3])
	extData.user_ext_data_[4] = C.char(data[4])
	extData.user_ext_data_[5] = C.char(data[5])
	extData.user_ext_data_[6] = C.char(data[6])
	extData.user_ext_data_[7] = C.char(data[7])
	extData.ext_data_len_ = 8

	ret := C.anti_del_user((*C.uchar)(unsafe.Pointer(C.CString(openid))), C.uchar(len(openid)),
		C.uchar(platid), C.uint(worldid), C.TSS_UINT64(roleid), &extData)
	if ret != 0 {
		return ErrDelPlayerFailed
	}

	return nil
}

// ErrTransAntiDataFailed 删除用户失败
var ErrTransAntiDataFailed = errors.New("TSSSDK transport anti data failed")

// OnRecvAntiData 收到数据
func OnRecvAntiData(entityid uint64, openid string, platid uint8, worldid uint32, roleid uint64, antidata []byte) error {
	extData := C.TssSdkUserExtData{}
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, entityid)
	extData.user_ext_data_[0] = C.char(data[0])
	extData.user_ext_data_[1] = C.char(data[1])
	extData.user_ext_data_[2] = C.char(data[2])
	extData.user_ext_data_[3] = C.char(data[3])
	extData.user_ext_data_[4] = C.char(data[4])
	extData.user_ext_data_[5] = C.char(data[5])
	extData.user_ext_data_[6] = C.char(data[6])
	extData.user_ext_data_[7] = C.char(data[7])
	extData.ext_data_len_ = 8

	antiStr := string(antidata)
	ret := C.on_recv_anti_data((*C.uchar)(unsafe.Pointer(C.CString(openid))), C.uchar(len(openid)),
		C.uchar(platid), C.uint(worldid), C.TSS_UINT64(roleid),
		(*C.uchar)(unsafe.Pointer(C.CString(antiStr))), C.uint(len(antidata)), &extData)
	if ret != 0 {
		return ErrTransAntiDataFailed
	}

	return nil
}

// ErrJudgeNameFailed 检测失败
var ErrJudgeNameFailed = errors.New("TSSSDK UIC Judge Name Failed")

// JudgeUserInputName 检测名字是否合法
func JudgeUserInputName(name string) (bool, error) {
	msg := C.CString(name)
	msglen := C.int(len(name)) + 1
	result := C.int(0)

	ret := C.judge_user_input_name(msg, msglen, &result)
	if ret != 0 {
		return false, ErrJudgeNameFailed
	}

	if result != 0 {
		return false, nil
	}

	return true, nil
}
