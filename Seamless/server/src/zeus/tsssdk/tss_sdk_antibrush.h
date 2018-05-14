/**
* @file:   tss_sdk_antibrush.h
* @brief:  This header file provides the interface of anti-brush plugin.
* @copyright: 2011 Tencent. All Rights Reserved.
*/

#ifndef TSS_SDK_ANTIBRUSH_H_
#define TSS_SDK_ANTIBRUSH_H_

#include "tss_sdk.h"

#ifdef __cplusplus
extern "C"
{
#endif

static const unsigned int TSS_SDK_ANTIBRUSH_MAX_ACCOUNT_ID_LEN = 64; /* 最大account_id长度 Maximum account_id length */
static const unsigned int TSS_SDK_ANTIBRUSH_MAX_RUID_LEN = 64;       /* 最大ruid长度 Maximum ruid length */
static const unsigned int TSS_SDK_ANTIBRUSH_MAX_DEVICE_ID_LEN = 64;  /* 最大device_id长度 Maximum device_id length */

#pragma pack(1)

/* 用户帐号类型 */
/* User Account Type */
typedef enum
{
    ANTIRUSG_ACCOUNT_OTHER   = 0, /* 其他 other account */
    ANTIRUSG_ACCOUNT_QQ_OPEN = 1, /* QQ开放帐号 qq openplat account(openid) */
    ANTIRUSG_ACCOUNT_WECHAT  = 2, /* 微信开放帐号 wechat openplat account(openid) */
    ANTIRUSG_ACCOUNT_QQ      = 3, /* QQ号 qq account */
    ANTIRUSG_ACCOUNT_PHONE   = 4, /* 手机号 phone number */
    ANTIRUSG_ACCOUNT_SELF    = 5, /* 自有帐号 self account */
} AntiBrushAccountType;

/* 防刷处理结果 */
typedef enum
{
    ANTIBRUSH_PROC_SUCC = 0, /*防刷结果查询成功 */
    ANTIBRUSH_PROC_FAIL = 1, /* 防刷结果查询失败 */
    ANTIBRUSH_PROC_TIMEOUT = 2, /* 防刷查询超时 */
} TssSdkAntiBrushResultType;

/* 玩家基本信息 User Base Info */
typedef struct
{
    AntiBrushAccountType account_type_;     /* [in] 帐号类型 Account Type */
    unsigned char account_id_[TSS_SDK_ANTIBRUSH_MAX_ACCOUNT_ID_LEN]; /* [in] 可填玩家openid、微信填开放帐号，QQ号，手机号 User Openid or QQ Account */
    unsigned int gameid_;                   /* [in] 产品id（安全侧分配）Game Product ID */
    unsigned int appid_;                    /* [in] 安全分配id（安全侧分配）Application ID */
    unsigned int event_id_;                 /* [in] 活动id Event ID */
    unsigned int plat_;                     /* [in] 设备平台：0-IOS，1-Android */
    unsigned int clinet_ip_;                /* [in] 用户IP User Device IP */
    unsigned int area_;                     /* [in] 用户平台：1-微信，2-手Q wechat=1, mobile-qq=2*/
    unsigned int oidb_gameid_;              /* [in] oidb游戏id（使用手Q平台下openid时使用，可以联系安全侧提供）Global Gameid */
} TssSdkAntiBrushUserInfo;

/* 玩家业务信息 User Business Info */
typedef struct
{
    unsigned int cost_;                     /* [in] 30天内付费记录 */
    unsigned int fight_;                    /* [in] 战斗力 Fight */
    unsigned char ruid_[TSS_SDK_ANTIBRUSH_MAX_RUID_LEN]; /* [in] 角色号 Role ID */
    unsigned int level_;                    /* [in] 等级 Level */
    unsigned int register_ts_;              /* [in] 注册时间戳 Register Timestamp */
    unsigned int online_time_;              /* [in] 在线时长 Online Duration */
    unsigned int friend_count_;             /* [in] 好友数 Friends Count */
    unsigned char device_id_[TSS_SDK_ANTIBRUSH_MAX_DEVICE_ID_LEN]; /* [in] 设备ID Device ID */
} TssSdkAntiBrushBusiInfo;

/* 使用异步接口获取防刷结果用户输入的参数结构体 */
/* Using an asynchronous interface to get anti-brush result of user query parameters structure */
typedef struct
{
    TssSdkAntiBrushUserInfo user_info_;     /* [in] 玩家基本信息 User Base Info */
    TssSdkAntiBrushBusiInfo busi_info_;     /* [in] 玩家业务信息 User Business Info */
    unsigned int id_;                       /* [in] 服务端回带id response id */
} TssSdkAntiBrushQueryInfo;

/* 使用异步接口获取防刷结果用户输入后返回的回调函数参数结构体 */
/* Using an asynchronous interface to determine the callback function returned
   to the anti-brush result of user input parameter structure */
typedef struct
{
    unsigned char account_id_[TSS_SDK_ANTIBRUSH_MAX_ACCOUNT_ID_LEN];
    /* [in] 可填玩家openid、微信填开放帐号，QQ号，手机号 User Openid or QQ Account */
    unsigned int result_;                   /* [in] 服务端处理结果 anti-brush result */
    unsigned int level_;                   /* [in] 防刷结果 anti-brush result */
    unsigned int id_;                       /* [in] 服务端回带id response id */
} TssSdkAntiBrushJudgeResultInfo;

/* 异步方式查询防刷接口 */
/* Asynchronous way to query user anti-brush status input interface */
typedef TssSdkProcResult(*TssSdkAntiBrushJudgeQuery)(const TssSdkAntiBrushQueryInfo *input_info);

/* 异步方式下返回防刷判定结果的回调函数 */
/* Return to the anti-brush to determine the results of a callback function in the asynchronous mode */
typedef TssSdkProcResult(*TssSdkAntiBrushOnJudgeResult)(const TssSdkAntiBrushJudgeResultInfo *result_info);

typedef struct
{
    TssSdkAntiBrushJudgeQuery antibrush_judge_query_;
} TssSdkAntiBrushInterf;

typedef struct
{
    TssSdkAntiBrushOnJudgeResult on_antibrush_judge_result_;
} TssSdkAntiBrushInitInfo;

/* 获取用户输入控制接口

参数说明
- init_data：用户输入控制接口初始化信息

返回值：成功-用户输入控制接口组指针，失败-NULL
*/
/* Get user input control interface

Parameter Description
- init_data：User input control interface initialization information

return value：success-user input control interface pointer, failure-NULL
*/
#define TSS_SDK_GET_ANTIBRUSH_INTERF(init_info) \
    (const TssSdkAntiBrushInterf*)tss_sdk_get_busi_interf("tss_sdk_get_antibrush_interf", (const TssSdkAntiBrushInitInfo *)(init_info))

#pragma pack()

#ifdef __cplusplus
} /* end of extern "C" */
#endif

#endif   /*TSS_SDK_ANTIBRUSH_H_*/

