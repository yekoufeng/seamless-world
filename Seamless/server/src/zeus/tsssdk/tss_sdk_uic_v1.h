/**
 * @file:   tss_sdk_uic.h
 * @brief:  This header file provides the interface of uic plugin.
 * @copyright: 2011 Tencent. All Rights Reserved.
 */

#ifndef TSS_SDK_UIC_V1_H_
#define TSS_SDK_UIC_V1_H_

#include "tss_sdk.h"
#include "tss_sdk_uic.h"

#ifdef __cplusplus
extern "C"
{
#endif

#pragma pack(1)

/* 使用异步接口判断CHAT类型用户输入的参数结构体 */
/* Using an asynchronous interface to determine the CHAT type of user input parameters structure */
typedef struct
{
    TssSdkOpenid openid_;                   /* [in] 用户openid */
    unsigned char plat_id_;                 /* [in] plat_id, 0: IOS, 1: Android */
    UicMsgCategory msg_category_;           /* [in] 消息内容的类别：邮件，聊天等 */
    unsigned int channel_id_;               /* [in] 发言具体频道号 */
    unsigned int client_ip_;                /* [in] 客户端ip */
    TSS_UINT64 role_id_;                    /* [in] 角色编号 */
    unsigned int role_level_;               /* [in] 角色等级 */
    unsigned short role_name_len_;          /* [in] 角色名长度 */
    unsigned char *role_name_;              /* [in] 角色名 */
    unsigned int msg_len_;                  /* [in] 消息长度，最长TSS_SDK_UIC_MAX_MESSAGE_LEN message length, maximum TSS_SDK_UIC_MAX_MESSAGE_LEN */
    unsigned char *msg_;                    /* [in] 消息内容 message content */
    int door_level_;                        /* [in] 限制门槛级别 */
    unsigned short callback_data_len_;      /* [in] 回带字段长度 callback data length */
    unsigned char *callback_data_;          /* [in] 回带字段内容，最长TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN callback data content, maximum TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN */
    unsigned short game_data_len_;          /* [in] 扩展字段长度 extension data length */
    unsigned char *game_data_;              /* [in] 扩展字段内容，最长TSS_SDK_UIC_MAX_GAME_DATA_LEN extension data content, max TSS_SDK_UIC_MAX_GAME_DATA_LEN */
} TssSdkUicChatUserInputInfo;

/* 使用异步接口判断CHAT类型用户输入后返回的回调函数参数结构体 */
/* Using an asynchronous interface to determine the callback function returned to the CHAT type of user input parameter structure */
typedef struct
{
    TssSdkOpenid openid_;                   /* [in] 用户openid */
    unsigned char plat_id_;                 /* [in] plat_id, 0: IOS, 1: Android */
    UicMsgResultFlag msg_result_flag_;      /* [in] 用户输入经过sdk的判定结果，开发方根据结果决定是否屏蔽 */
    /*      User input after the verdict of the sdk, developer to decide whether to shield in accordance with the results */
    int dirty_level_;                       /* [in] 给定字串的限制级别 */
    unsigned int msg_len_;                  /* [in] 消息长度 message length */
    unsigned char *msg_;                    /* [in] 消息内容，最长TSS_SDK_UIC_MAX_MESSAGE_LEN message content, max TSS_SDK_UIC_MAX_MESSAGE_LEN */
    unsigned short callback_data_len_;      /* [in] 回带字段长度 callback data length */
    unsigned char *callback_data_;          /* [in] 回带字段内容，最长TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN callback data content, maximum TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN */
} TssSdkUicChatJudgeResultInfo;

/* 使用同步接口判断NAME类型用户输入的参数结构体 */
/* Synchronous interface to determine the NAME type of user input parameters structure */
typedef struct
{
    unsigned int msg_len_;                  /* [in/out] 消息长度 message length*/
    unsigned char *msg_;                    /* [in/out] 消息内容，最长TSS_SDK_UIC_MAX_MESSAGE_LEN message content, max TSS_SDK_UIC_MAX_MESSAGE_LEN */
    int door_level_;                        /* [in] 限制门槛级别，如果不关注，填1 */
    char if_replace_;                       /* [in] 如果包含敏感词，是否替换成“*”。如果是，填1；如果否，填0 */
    /*      If it contains sensitive words, whether to replace with the "*". If yes, fill in 1; if not, fill in 0*/
    UicMsgResultFlag msg_result_flag_;      /* [out] 用户输入经过敏感词检查的判定结果 User input after the verdict of the sensitive words check */
    int dirty_level_;                       /* [out] 给定字串的限制级别 */
} TssSdkUicNameUserInputInfo;

/* 异步方式判断CHAT类型用户输入接口 */
/* Asynchronous way to judge the type CHAT user input interface */
typedef TssSdkProcResult(*TssSdkUicJudgeUserInputChat)(const TssSdkUicChatUserInputInfo *input_info);

/* 异步方式下返回CHAT类型判定结果的回调函数 */
/* Return to the CHAT type to determine the results of a callback function in the asynchronous mode */
typedef TssSdkProcResult(*TssSdkUicChatOnJudgeResult)(const TssSdkUicChatJudgeResultInfo *result_info);

/* 同步方式检查NAME类型用户输入接口 */
/* Synchronization check the NAME type of user input interface */
typedef TssSdkProcResult(*TssSdkUicJudgeUserInputName)(TssSdkUicNameUserInputInfo *input_info);

typedef struct
{
    TssSdkUicChatOnJudgeResult on_chat_judge_result_;
} TssSdkUicInitInfo;

typedef struct
{
    TssSdkUicJudgeUserInputChat judge_user_input_chat_;
    TssSdkUicJudgeUserInputName judge_user_input_name_;
} TssSdkUicInterf;

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
#define TSS_SDK_GET_UIC_INTERF(init_data) \
    (const TssSdkUicInterf*)tss_sdk_get_busi_interf("tss_sdk_get_uic_interf", (const TssSdkUicInitInfo *)(init_data))

#pragma pack()

#ifdef __cplusplus
} /* end of extern "C" */
#endif

#endif   /*TSS_SDK_PLUGIN_UIC_H_*/

