/**
 * @file:   tss_sdk_anti.h
 * @brief:  This header file provides the interface of TSS SDK.
 * @copyright: 2012 Tencent. All Rights Reserved.
 */

#ifndef TSS_SDK_ANTI_V1_H_
#define TSS_SDK_ANTI_V1_H_

#include "tss_sdk.h"
#include "tss_sdk_anti.h"

#ifdef __cplusplus
extern "C"
{
#endif /* end of __cplusplus */

#pragma pack(1)

/* 添加用户信息 */
/* Add user information */
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
    /* [in] 游戏客户端的版本 game client version */
    unsigned int client_ver_;
    /* [in] 游戏客户端ip game client ip */
    unsigned int client_ip_;
    /* [in] 用户当前的角色名 user's current role name */
    const char *role_name_;
} TssSdkAntiAddUserInfoEx;

/* 删除用户信息 */
/* delete user information */
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
} TssSdkAntiDelUserInfoEx;

/* 收到反外挂数据*/
/* recv anti data */
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
    /* [in] gamesvr收到的anti数据 */
    /* [in] anti data received by gamesvr */
    const unsigned char *anti_data_;
    /* [in] gamesvr收到的anti数据长度 */
    /* [in] length of anti data received by gamesvr */
    unsigned int anti_data_len_;
} TssSdkAntiRecvDataInfoEx;

/* 解密游戏数据包 */
/* decryption of packets */
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
    /* [in] gamesvr收到的游戏加密数据 */
    /* [in] game encrypt data received by gamesvr */
    const unsigned char *encrypt_data_;
    /* [in] gamesvr收到的游戏加密数据长度 */
    /* [in] length of encrypt data received by gamesvr */
    unsigned int encrypt_data_len_;
    /* [in/out] 用来存放解密后的游戏包的buf，空间由调用方分配 */
    /* [in/out] buf used to store the decrypted game package, space allocated by the caller */
    unsigned char *game_pkg_buf_;
    /* [in/out] 输入为game_pkg_buf_的size，输出为解密后的游戏包实际长度 */
    /* [in/out] input is size of game_pkg_buf_, output is the actual length of decrypted game package */
    unsigned int game_pkg_buf_len_;
} TssSdkAntiDecryptPkgInfoEx;

/* 游戏数据包信息 */
/* Game data package information */
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
    /* [in] gamesvr收到的游戏加密数据 */
    /* [in] 游戏包命令字 Game package command word */
    unsigned int cmd_id_;
    /* [in] 游戏数据包 Game data packets */
    const unsigned char *game_pkg_;
    /* [in] 游戏数据包长度 the length of game data packets */
    unsigned int game_pkg_len_;
} TssSdkAntiGamePkgInfoEx;

/* 加密数据包信息 */
/* Encrypted data packet information */
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
    /* [in] 游戏包命令字 Game package command word */
    unsigned int cmd_id_;
    /* [in] 游戏数据包 Game data packets */
    const unsigned char *game_pkg_;
    /* [in] 游戏数据包长度，最大长度要小于65000
    the length of game data packets, maximum length less than 65,000 */
    unsigned int game_pkg_len_;
    /* [in/out] 游戏数据包加密后的数据包，内存由调用方分配，最大64k
    encrypted game data package assembled into anti data,
    memory allocated by the caller, 64k at the maximum */
    unsigned char *encrypt_data_;
    /* [in/out] 输入时为encrypt_data_len_的长度，输出时为实际encrypt_data_len_使用的长度
    length of encrypt_data_len_ when input, actual length of encrypt_data_len_ when output */
    unsigned int encrypt_data_len_;
} TssSdkAntiEncryptPkgInfoEx;

/* 要发送的anti数据信息 */
/* anti data to be sent*/
typedef struct
{
    /* [in] openid*/
    TssSdkOpenid openid_;
    /* [in] plat_id, 0: IOS, 1: Android */
    unsigned char plat_id_;
    /* [in] anti数据长度 length of anti data */
    unsigned short anti_data_len_;
    /* [in] 需要发送到客户端的anti数据 anti data need to be sent to the client */
    const unsigned char *anti_data_;
} TssSdkAntiSendDataInfoEx;

/* 通知添加用户，游戏在用户登陆时调用 */
/* Notice to add a user, the game called when the user login */
typedef TssSdkProcResult(*TssSdkAddUserEx)(TssSdkAntiAddUserInfoEx *add_user_info);

/* 通知删除用户，游戏在用户退出时调用 */
/* Notice to delete a user, the game called when the user log out */
typedef TssSdkProcResult(*TssSdkDelUserEx)(const TssSdkAntiDelUserInfoEx *del_user_info);

/* 接收客户端上行anti数据，游戏收到反外挂数据时调用 */
/* recv anti data, the game called when recvd anti data */
typedef TssSdkProcResult(*TssSdkRecvDataFromClientEx)(TssSdkAntiRecvDataInfoEx *recv_pkg_info);

/* 判定上行游戏包是否作弊包, 0表示不是欺骗包，游戏继续后续处理 1表示是欺骗包，游戏需要将此游戏包丢弃 */
/* Determine whether the uplink game package is cheating package, */
/*0 means it is not deceive package which the game continue to processlater; */
/*1 means a deceive package which the game needs to discard */
typedef int (*TssSdkIsCheatPkgEx)(const TssSdkAntiGamePkgInfoEx *up_pkg_info);

/* 对上行的anti数据包进行解密 */
/* Decrypt for the uplink game data package */
typedef TssSdkAntiDecryptResult(*TssSdkDecryptPkgEx)(TssSdkAntiDecryptPkgInfoEx *decrypt_pkg_info);

/* 对需要加密的游戏包进行加密,如果不需要加密，则游戏自行处理，如果需要加密，需要游戏发送返回的anti数据  */
/* encrypt the game package which is needed, if not, the game process on its own; if needed, the game should send back the anti data */
typedef TssSdkAntiEncryptResult(*TssSdkEncryptPkgEx)(TssSdkAntiEncryptPkgInfoEx *down_pkg_info);

/* 发送Anti数据，需由游戏实现 */
/* send anti data, need to be implement by game */
typedef TssSdkProcResult(*TssSdkSendDataToClientEx)(const TssSdkAntiSendDataInfoEx *anti_data);

typedef struct
{
    /* 发送加密后的数据到客户端, 此函数需要游戏服务器实现*/
    /* Send encrypted data to the client, this function need to be implemented by game server */
    TssSdkSendDataToClientEx send_data_to_client_;
} TssSdkAntiInitInfoEx;

typedef struct
{
    /* 添加用户 Add user */
    TssSdkAddUserEx add_user_;

    /* 删除用户 delete user */
    TssSdkDelUserEx del_user_;

    /* 收到反外挂数据包 recv anti data package */
    TssSdkRecvDataFromClientEx recv_anti_data_;

    /* 判定是否欺骗包 Determine whether the deception package */
    TssSdkIsCheatPkgEx is_cheat_pkg_;

    /* 解密数据包 decrypt data package */
    TssSdkDecryptPkgEx decrypt_pkg_;

    /* 加密数据包 encrypt data package */
    TssSdkEncryptPkgEx encrypt_pkg_;
} TssSdkAntiInterfEx;

/* 获取Anti接口 */
/* 此宏必须在tss_sdk_load成功之后才能调用 */
/* Get Anti interface */
/* This macro must be called after the success of tss_sdk_load */
#define TSS_SDK_GET_ANTI_INTERF_EX(init_data) \
    (const TssSdkAntiInterfEx*)tss_sdk_get_busi_interf("tss_sdk_get_anti_interf_ex", (const TssSdkAntiInitInfoEx *)(init_data))

#pragma pack()

#ifdef __cplusplus
} /* end of extern "C" */
#endif /* end of __cplusplus */

#endif /* TSS_SDK_ANTI_H_ */

