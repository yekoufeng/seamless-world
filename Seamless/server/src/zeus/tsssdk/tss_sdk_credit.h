/**
* @file:   tss_sdk_credit.h
* @brief:  This header file provides the interface of credit plugin.
* @copyright: 2017 Tencent. All Rights Reserved.
*/

#ifndef TSS_SDK_CREDIT_H_
#define TSS_SDK_CREDIT_H_

#include "tss_sdk.h"

#ifdef __cplusplus
extern "C"
{
#endif

#pragma pack(1)

/* 信用度查询处理结果 */
typedef enum
{
    CREDIT_PROC_SUCC = 0,       /* 信用度查询成功 */
    CREDIT_PROC_FAIL = 1,       /* 信用度查询失败 */
    CREDIT_PROC_TIMEOUT = 2,    /* 信用度查询超时 */
} TssSdkCreditResultType;

/* 使用异步接口获取防刷结果用户输入的参数结构体 */
/* Using an asynchronous interface to get anti-brush result of user query parameters structure */
typedef struct
{
    unsigned int plat_id_;            /* [in] 平台 platid,0-IOS,1-Android,255-all */
    unsigned char openid[64];         /* [in] 用户openid */
    unsigned int id_;                 /* [in] 服务端回带id response id */
} TssSdkCreditQueryInfo;

/* 使用异步接口获取防刷结果用户输入后返回的回调函数参数结构体 */
/* Using an asynchronous interface to determine the callback function returned
   to the anti-brush result of user input parameter structure */
typedef struct
{
    unsigned int version_;            /* [out] 版本号 version */
    int result_;                      /* [out] 查询结果 result */
    unsigned int data_time_;          /* [out] 更新时间 update time(UNIX format) */
    unsigned char openid[64];         /* [out] 用户openid */
    unsigned int plat_id_;            /* [out] 平台 platid,0-IOS,1-Android,255-all */
    unsigned int payment_;            /* [out] 付费 payment */
    unsigned int score_;              /* [out] 分数 score(0-600) */
    unsigned int rank_;               /* [out] 排名 rank(0-100) */
    unsigned int stars_;              /* [out] 星级 stars(1-7) */
    unsigned int id_;                 /* [out] 服务端回带id response id */
} TssSdkCreditJudgeResultInfo;

typedef struct
{
    unsigned int plat_id_;            /* [in] platid,0-IOS,1-Android,255-all */
    unsigned char openid[64];         /* [in] openid */
    unsigned int seq_;                /* [in] seq */
} TssSdkPluginCreditQuery;


/* 异步方式查询信用度接口 */
/* Asynchronous way to query user credit status input interface */
typedef TssSdkProcResult(*TssSdkCreditJudgeQuery)(const TssSdkCreditQueryInfo *input_info);

/* 异步方式下返回信用度判定结果的回调函数 */
/* Return to the credit to determine the results of a callback function in the asynchronous mode */
typedef TssSdkProcResult(*TssSdkCreditOnJudgeResult)(const TssSdkCreditJudgeResultInfo *result_info);

typedef struct
{
    TssSdkCreditJudgeQuery credit_judge_query_;
} TssSdkCreditInterf;

typedef struct
{
    TssSdkCreditOnJudgeResult on_credit_judge_result_;
} TssSdkCreditInitInfo;

/* 获取信用度查询接口

参数说明
- init_data：信用度查询接口初始化信息

返回值：成功-信用度查询接口组指针，失败-NULL
*/
/* Get credit query interface

Parameter Description
- init_data：Credit query interface initialization information

return value：success-credit interface pointer, failure-NULL
*/
#define TSS_SDK_GET_CREDIT_INTERF(init_info) \
    (const TssSdkCreditInterf*)tss_sdk_get_busi_interf("tss_sdk_get_credit_interf", (const TssSdkCreditInitInfo *)(init_info))

#pragma pack()

#ifdef __cplusplus
} /* end of extern "C" */
#endif

#endif   /*TSS_SDK_CREDIT_H_*/

