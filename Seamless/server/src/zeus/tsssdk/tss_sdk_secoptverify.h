/*
 @copyright  2004-2015  Apache License, Version 2.0
 @filename   tss_sdk_secoptverify.h
 @author     yunfeiyang
 @version
 @date       2016/03/16 11:17
 @brief
 @details    2016/03/16 yunfeiyang create
*/
#ifndef TSS_SDK_SENSITIVE_H_
#define TSS_SDK_SENSITIVE_H_

#include "tss_sdk.h"

#ifdef __cplusplus
extern "C"
{
#endif /* end of __cplusplus */

#define TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN  1024        /* 最大消息长度 Maximum message length */

#pragma pack(1)

/* 敏感操作保护态用户基本信息 */
typedef struct
{
    /* 用户openid */
    unsigned char openid_[64];
    /* 用户操作系统类型，0-IOS，1-Android */
    unsigned int plat_id_;
    /* 用户大区号，没有填0 */
    unsigned int world_id_;
} TssSdkSensitiveUserInfoBase;

/* 安全控制类型 */
typedef enum
{
    TYPE_SAFE_MODE_CONTROL = 1,         // 保护态控制
    TYPE_OTHER = 99,                    // 其他类型
} TssSdkSecurityType;


/* 保护态控制结果 */
typedef enum
{
    SECURITY_PASSED = 0,                   // 保护态直接放过
    SECURITY_NEED_SET_SAFE_MODE = 1,       // 设置保护态
    SECURITY_IN_LIMITED_TIME = 2,          // 还处于刚登录后的限制期，暂时不允许任何敏感操作
} TssSdkSecurityResultType;


/* 绑定手机号码标志 */
typedef enum
{
    BIND_PHONE_NO = 0,         // 未绑定手机号码
    BIND_PHONE_YES = 1,        // 已绑定手机号码
    BIND_PHONE_UNKNOWN = 2,    // 不确定是否绑定了手机号码
} TssSdkSensitiveBindPhoneType;


/* 保护态状态 */
typedef struct
{
    /* [in] 用户基本信息 */
    TssSdkSensitiveUserInfoBase user_info_;
    /* [out] 安全控制类型 */
    TssSdkSecurityType security_type_;
    /* [out] 安全控制结果 */
    TssSdkSecurityResultType security_result_;
    /* [out] 是否绑定已绑定手机号，0-未绑定，1-已绑定 */
    TssSdkSensitiveBindPhoneType bind_phone_num_flag_;
    /* [out] 弹窗类型， 类型值对应的弹窗ui由产品经理和游戏项目组协商 */
    int pop_up_type_;
    /* [in\out] 弹窗消息内容，最长TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN; message content, max TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN */
    char msg_[TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN];
} TssSdkSensitiveSafeModeInfo;

/* 同步方式查询用户保护态状态信息接口 */
typedef TssSdkProcResult(*TssSdkSensitiveQuerySensitiveSafeMode)(TssSdkSensitiveSafeModeInfo *input_info);


/* 发送短信验证保护态 */
typedef struct
{
    /* [in] 用户基本信息 */
    TssSdkSensitiveUserInfoBase user_info_;
    /* 登陆的客户端IP */
    char client_ip_[20];
} TssSdkSensitiveSMSInfo;

/* 敏感操作发送短信验证请求接口 */
typedef TssSdkProcResult(*TssSdkSensitiveSendSMSReq)(const TssSdkSensitiveSMSInfo *sms_request);


/* 短信验证请求信息 */
typedef struct
{
    /* [in] 用户基本信息 */
    TssSdkSensitiveUserInfoBase user_info_;
    /* [in] 回带sequence，此处建议填写操作流水号 */
    unsigned int seq_;
    /* [in] 短信验证码内容 */
    char sms_serial_[16];
} TssSdkSensitiveSMSVerifyInfo;

/* 发送短信验证请求接口 */
typedef TssSdkProcResult(*TssSdkSensitiveSMSVerifyReq)(const TssSdkSensitiveSMSVerifyInfo *verify_request);


/* 短信验证结果类型 */
typedef enum
{
    SMS_RESULT_VERIFY_PASSED = 0,           // 短信验证通过
    SMS_RESULT_NO_BING_PHONE = 1,           // 用户没有绑定手机
    SMS_RESULT_VERIFY_UNPASSED = 2,         // 验证不通过
    SMS_RESULT_TIMEOUT = 3,                 // 超时等内部错误
} TssSdkSensitiveSMSResultType;

/* 验证短信通过或者用户无手机的通知内容 */
typedef struct
{
    /* [in] 用户基本信息 */
    TssSdkSensitiveUserInfoBase user_info_;
    /* [in] 回带sequence，此处建议填写操作流水号 */
    unsigned int seq_;
    /* [in] 短信验证结果 */
    TssSdkSensitiveSMSResultType sms_result_;
    /* [in] 弹窗类型 */
    int pop_up_type_;
    /* [in] 弹窗消息内容，最长TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN; message content, max TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN */
    char msg_[TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN];
} TssSdkSensitiveSMSResult;

/* 游戏服务端需要实现的获取短信验证结果的回调函数 */
typedef TssSdkProcResult(*TssSdkSensitiveSetSMSVerifyResult)(const TssSdkSensitiveSMSResult *verify_result);


/* 绑定手机号码请求信息 */
typedef struct
{
    /* [in] 用户基本信息 */
    TssSdkSensitiveUserInfoBase user_info_;
    /* [in] 回带sequence，此处建议填写操作流水号 */
    unsigned int seq_;
    /* [in] 手机号码 */
    char mobile_phone_number_[16];
} TssSdkSensitiveBindPhoneNumInfo;

/* 绑定手机号码请求 */
typedef TssSdkProcResult(*TssSdkSensitiveBindPhoneNumReq)(const TssSdkSensitiveBindPhoneNumInfo *bind_request);


/* 绑定手机号码结果类型 */
typedef enum
{
    BIND_RESULT_SUCCESS = 0,           // 绑定手机号码成功
    BIND_RESULT_FAILED = 1,            // 绑定手机号码失败
    BIND_RESULT_TIMEOUT = 2,           // 超时等内部错误
} TssSdkSensitiveBindPhoneResultType;

/* 绑定手机号码结果 */
typedef struct
{
    /* [in] 用户基本信息 */
    TssSdkSensitiveUserInfoBase user_info_;
    /* [in] 回带sequence，此处建议填写操作流水号 */
    unsigned int seq_;
    /* [in] 绑定手机号码结果 */
    TssSdkSensitiveBindPhoneResultType bind_result_;
    /* [in] 弹窗类型 */
    int pop_up_type_;
    /* [in] 弹窗消息内容，最长TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN; message content, max TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN */
    char msg_[TSS_SDK_SAFEMODE_MAX_MESSAGE_LEN];
} TssSdkSensitiveBindPhoneNumResult;

/* 游戏服务端需要实现的获取短信验证结果的回调函数 */
typedef TssSdkProcResult(*TssSdkSensitiveSetBindPhoneNumResult)(const TssSdkSensitiveBindPhoneNumResult *bind_result);


/* 敏感操作接口 */
typedef struct
{
    /* 请求查询玩家的保护态状态信息 */
    TssSdkSensitiveQuerySensitiveSafeMode query_safe_mode_info_;

    /* 请求绑定手机号码 */
    TssSdkSensitiveBindPhoneNumReq req_bind_phone_num_;

    /* 请求发送敏感操作短信验证码 */
    TssSdkSensitiveSendSMSReq req_send_sms_;

    /* 请求对短信验证码进行校验 */
    TssSdkSensitiveSMSVerifyReq req_verify_sms_;
} TssSdkSensitiveInterf;

/* 敏感操作统一验证接口初始化数据 */
typedef struct
{
    /* 设置手机号码绑定结果 */
    TssSdkSensitiveSetBindPhoneNumResult set_bind_phone_num_result_;

    /* 设置短信验证码校验结果 */
    TssSdkSensitiveSetSMSVerifyResult set_sms_verify_result_;
} TssSdkSensitiveInitData;

#define TSS_SDK_GET_SENSITIVE_INTERF(init_data) \
    (const TssSdkSensitiveInterf*)tss_sdk_get_busi_interf("tss_sdk_get_sensitive_interf", (const TssSdkSensitiveInitData *)(init_data))


#pragma pack()

#ifdef __cplusplus
}
#endif /* end of __cplusplus */

#endif // TSS_SDK_SENSITIVE_H_

