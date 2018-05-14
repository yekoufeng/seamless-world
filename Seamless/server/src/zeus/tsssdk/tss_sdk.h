/**
 * @file:   tss_sdk.h
 * @brief:  This header file provides the interface of TSS SDK.
 * @copyright: 2012 Tencent. All Rights Reserved.
 */

#ifndef TSS_SDK_H__
#define TSS_SDK_H__

#ifdef __cplusplus
extern "C"
{
#endif

#if (defined(WIN32) || defined(_WIN64))

#include <tchar.h>

#if _MSC_VER >= 1300
typedef unsigned long long  TSS_UINT64;
typedef long long   TSS_INT64;
#else /* _MSC_VER */
typedef unsigned __int64    TSS_UINT64;
typedef __int64 TSS_INT64;
#endif /* _MSC_VER */
typedef TCHAR TSS_TCHAR;

#else // (defined(WIN32) || defined(_WIN64))

#include <stdint.h>
typedef uint64_t TSS_UINT64;
typedef int64_t TSS_INT64;
typedef char TSS_TCHAR;

#endif

/* 接口通用处理结果 */
/* Interface general processing results */
typedef enum
{
    TSS_SDK_PROC_OK = 0,             /* 处理成功 Processing Successfully */
    TSS_SDK_PROC_INVALID_ARG = 1,    /* 无效参数 Invalid parameters */
    TSS_SDK_PROC_INTERNAL_ERR = 2,   /* 内部错误 Internal error */
    TSS_SDK_PROC_FAIL = 3,           /* 处理失败 Processing failed */
} TssSdkProcResult;

/* callback function of sdk proc */
typedef int (*TssSdkProc)();

/* SDK interface */
typedef struct
{
    /* tss_sdk 内部数据处理函数，需要游戏定时调用，
    频率至少需要达到10次/s，整个进程只能有一个处调用 */
    /* tss_sdk internal data processing function which gamesvr calls timely
    with frequency of at least 10 times/s, only be called by one place of the whole process */
    TssSdkProc proc_;
} TssSdkBusiInterf;

/*openid of the role*/
typedef struct
{
    unsigned char *openid_;
    unsigned char openid_len_;
} TssSdkOpenid;

typedef struct
{
    /*
    加载了tss_sdk的进程实体id,
    由于这个id与sdk到后端的通讯相关，所以
    如果同一台机器上有多个进程加载了tss_sdk，需要满足两个条件：
    1. 不同进程传入的id不一样。
    2. 同一进程在重启时这个实体id也需要保持与上次一致
    */
    /*
    Load process tss_sdk entity id,
    because this id is related with the sdk to the back-end communication
    On the same machine multiple process load tss_sdk need to meet two conditions:
    1. Different process passes different id
    2. The process also need to maintain consistency with the previous entity id when restart.
    */
    unsigned int unique_instance_id_;
    /*
    tss_sdk的配置路径
    */
    const char *tss_sdk_conf_;
} TssSdkInitInfo;

/*
*
* @fn     tss_sdk_load
* @brief  加载sdk load sdk
*
* @param  shared_lib_dir [in] tss_sdk 动态库目录 tss_sdk Dynamic library directory
* @param  init_info [in] sdk 初始化参数 sdk Initialization parameters
*
* @return NULL --加载失败，可能是目录不对，或者sdk初始化失败 load failed, maybe it is the wrong directory, or sdk failed to initialize
*         非NULL --加载成功，请保留返回的地址，后面调用proc接口需要 load successfully, Please keep the return address, the following needs it to call proc interface
*/
const TssSdkBusiInterf *tss_sdk_load(const TSS_TCHAR *shared_lib_dir,
                                     const TssSdkInitInfo *init_info);
/* SDK unload function */
int tss_sdk_unload();

/*
获取sdk的业务接口结构， 此函数不要直接使用，
业务的接口获取直接使用业务提供的宏，类似TSS_SDK_GET_XX_INTERF。
*/
/*
Get the sdk service interface structure, this function is not used directly,
use the business macri to obtain business interface, similar with TSS_SDK_GET_XX_INTERF
*/
const void *tss_sdk_get_busi_interf(const char *syml_name,
                                    const void *data);


/*
用户扩展数据

以前的版本sdk在调用回调函数的时候,只会传递openid给业务svr

但是不少业务svr内部是使用自己的uid标识用户的,为了接sdk必须
在加上一个openid到内部uid的映射表

为了降低业务svr的工作量,在sdk接口上额外添加一个用户扩展数据,
业务可以把内部uid放到里面,并在调用sdk接口时传递给sdk,
当sdk调用业务svr提供的回调函数时,会将该数据随openid一起传递
给业务svr,业务svr可以使用该数据做快速查找


注意: 如果uid存在在不同的openid间复用的可能,
那么在根据扩展数据查找到用户信息后,应该再核对一次openid,确保不会关联到错误的用户
*/

/*这里用来传输SDK关键数据*/
#define TSS_SDK_USER_EXT_DATA_MAX_LEN 1024

typedef struct
{
    /*
        用户扩展数据
    */
    char user_ext_data_[TSS_SDK_USER_EXT_DATA_MAX_LEN];
    /*
        扩展数据数据的长度
    */
    unsigned int ext_data_len_;
} TssSdkUserExtData;



#ifdef __cplusplus
} /* end of extern "C" */
#endif

#endif /* end of TSS_SDK_H__ */

