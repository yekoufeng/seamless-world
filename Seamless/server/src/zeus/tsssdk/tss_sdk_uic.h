/**
 * @file:   tss_sdk_uic.h
 * @brief:  This header file provides the interface of uic plugin.
 * @copyright: 2011 Tencent. All Rights Reserved.
 */

#ifndef TSS_SDK_UIC_H_
#define TSS_SDK_UIC_H_
#include "tss_sdk.h"

#ifdef __cplusplus
extern "C"
{
#endif

static const unsigned int TSS_SDK_UIC_MAX_MESSAGE_LEN = 65536;          /* 最大消息长度 Maximum message length */
static const unsigned int TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN = 1024;     /* 最大回带字段长度 Maximum callback data length */
static const unsigned int TSS_SDK_UIC_MAX_GAME_DATA_LEN = 64000;        /* 最大扩展字段长度 Maximum extension data lengh */

#pragma pack(1)

/* 用户输入经过接口检查的判定结果 */
/* User input verdict by interface check */
typedef enum
{
    MSG_NORMAL_FLAG = 0,       /* 合法 legal */
    MSG_EVIL_FLAG = 1,         /* 不合法，不能显示 Illegal, can not be displayed */
    MSG_DIRTY_FLAG = 2,        /* 合法，但包含敏感词 legal, but contain sensitive words */
} UicMsgResultFlag;

/* 用户输入的类别 */
typedef enum
{
    CONTENT_CATEGORY_MAIL = 1,      /* 邮件 */
    CONTENT_CATEGORY_CHAT = 2,      /* 聊天 */
} UicMsgCategory;

/* 使用异步接口判断CHAT类型用户输入的参数结构体 */
/* Using an asynchronous interface to determine the CHAT type of user input parameters structure */
typedef struct
{
    TssSdkOpenid openid_;                   /* [in] 用户openid */
    unsigned char plat_id_;                 /* [in] plat_id, 0: IOS, 1: Android */
    unsigned int world_id_;                 /* [in] world id */
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
    int local_language_;                     /* [in] 指定要过滤的脏词的本地语言*/
} TssSdkUicChatUserInputInfoV3;

/* 使用异步接口判断CHAT类型用户输入后返回的回调函数参数结构体 */
/* Using an asynchronous interface to determine the callback function returned to the CHAT type of user input parameter structure */
typedef struct
{
    TssSdkOpenid openid_;                   /* [in] 用户openid */
    unsigned char plat_id_;                 /* [in] plat_id, 0: IOS, 1: Android */
    unsigned int world_id_;                 /* [in] world id */
    TSS_UINT64 role_id_;                    /* [in] role id */
    UicMsgResultFlag msg_result_flag_;      /* [in] 用户输入经过sdk的判定结果，开发方根据结果决定是否屏蔽 */
    /*      User input after the verdict of the sdk, developer to decide whether to shield in accordance with the results */
    int dirty_level_;                       /* [in] 给定字串的限制级别 */
    unsigned int msg_len_;                  /* [in] 消息长度 message length */
    unsigned char *msg_;                    /* [in] 消息内容，最长TSS_SDK_UIC_MAX_MESSAGE_LEN message content, max TSS_SDK_UIC_MAX_MESSAGE_LEN */
    unsigned short callback_data_len_;      /* [in] 回带字段长度 callback data length */
    unsigned char *callback_data_;          /* [in] 回带字段内容，最长TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN callback data content, maximum TSS_SDK_UIC_MAX_CALLBACK_DATA_LEN */
} TssSdkUicChatJudgeResultInfoV3;

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
    int local_language_;                     /* [in] 指定要过滤的脏词的本地语言*/
} TssSdkUicNameUserInputInfoV3;

/* 异步方式判断CHAT类型用户输入接口 */
/* Asynchronous way to judge the type CHAT user input interface */
typedef TssSdkProcResult(*TssSdkUicJudgeUserInputChatV3)(const TssSdkUicChatUserInputInfoV3 *input_info);

/* 异步方式下返回CHAT类型判定结果的回调函数 */
/* Return to the CHAT type to determine the results of a callback function in the asynchronous mode */
typedef TssSdkProcResult(*TssSdkUicChatOnJudgeResultV3)(const TssSdkUicChatJudgeResultInfoV3 *result_info);

/* 同步方式检查NAME类型用户输入接口 */
/* Synchronization check the NAME type of user input interface */
typedef TssSdkProcResult(*TssSdkUicJudgeUserInputNameV3)(TssSdkUicNameUserInputInfoV3 *input_info);

typedef struct
{

    TssSdkUicChatOnJudgeResultV3 on_chat_judge_result_;
} TssSdkUicInitInfoV3;

typedef struct
{
    TssSdkUicJudgeUserInputChatV3 judge_user_input_chat_;
    TssSdkUicJudgeUserInputNameV3 judge_user_input_name_;
} TssSdkUicInterfV3;

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
#define TSS_SDK_GET_UIC_INTERF_V3(init_data) \
    (const TssSdkUicInterfV3*)tss_sdk_get_busi_interf("tss_sdk_get_uic_interf_v3", (const TssSdkUicInitInfoV3 *)(init_data))

#pragma pack()

#ifdef __cplusplus
} /* end of extern "C" */
#endif

// 兼容旧版本的接口
#include "tss_sdk_uic_v2.h"

#endif   /*TSS_SDK_PLUGIN_UIC_H_*/

