/**
 * @file:   tss_sdk_import.c
 * @brief:  load tss_sdk dynamic library.
 * @copyright: 2011 Tencent. All Rights Reserved.
 */

#include "tss_sdk.h"
#include <stdio.h>

typedef const TssSdkBusiInterf *(*TssSdkGetBusiInterf)(const TssSdkInitInfo* init_info);
#define TSS_SDK_GET_COMM_INTERF_SYM "tss_sdk_get_comm_interf"

typedef int (*TssSdkReleaseBusiInterf)();
#define TSS_SDK_RELEASE_BUSI_INTERF_SYM "tss_sdk_release_busi_interf"

static int g_is_sdk_loaded = 0;

#if defined(WIN32) || defined(_WIN64)

#include <windows.h>
#define LIB_NAME _T("tss_sdk")
HMODULE g_handle = NULL;
#define SUFFIX _T("dll")
#define snprintf _snprintf

int
tss_sdk_load_win(const TSS_TCHAR *shared_lib_file)
{
    if (g_handle != NULL)
    {
        // 已经加载，则不再加载
	    // Has been loaded, no longer load
        return 0;
    }

    g_handle = LoadLibraryEx(shared_lib_file, NULL, 0);
    if (g_handle == NULL)
    {
#if (defined(WIN32) || defined(_WIN64))
        _ftprintf(stderr,
                _T("Open DLL %s failed, error code is %d.\n"),
                shared_lib_file,
                GetLastError());
#else
        fprintf(stderr,
                "Open DLL %s failed, error code is %d.\n",
                shared_lib_file,
                GetLastError());
#endif
        return -1;
    }

    return 0;
}

void *
tss_sdk_get_syml_win(const char *syml_name)
{
    void *func = NULL;
    if (g_handle == NULL)
    {
        return NULL;
    }

    func = GetProcAddress(g_handle, syml_name);
    if (func == NULL)
    {
        fprintf(stderr,
                "GetProcAddress failed, error code is %d.\n",
                 GetLastError());
        return NULL;
    }

    return func;
}

#else /* #if defined(WIN32) || defined(_WIN64) */

#include <dlfcn.h>
void *g_handle = NULL;
#define LIB_NAME "libtss_sdk"
#define SUFFIX "so"
static int
tss_sdk_load_linux(const char *shared_lib_file)
{
    char *error = NULL;

    if (g_handle != NULL)
    {
        // 已经加载，则不再加载
        // Has been loaded, no longer load
        return 0;
    }

    g_handle = dlopen(shared_lib_file, RTLD_NOW);
    error = dlerror();
    if (error)
    {
        fprintf(stderr, "Open shared lib %s failed, %s.\n", shared_lib_file, error);
        return -1;
    }

    return 0;
}

void *
tss_sdk_get_syml_linux(const char *syml_name)
{
    void *func = NULL;
    char *error = NULL;

    if (g_handle == NULL)
    {
        return NULL;
    }

    func = dlsym(g_handle, syml_name);
    error = dlerror();
    if (error)
    {
        fprintf(stderr, "dlsym failed, %s.\n", error);
        return NULL;
    }

    return func;
}

#endif /* #if defined(WIN32) || defined(_WIN64) */

int
tss_sdk_load_impl(const TSS_TCHAR *shared_lib_file)
{
    int ret = 0;
    if (g_is_sdk_loaded)
    {
        return 0;
    }

    #if defined(WIN32) || defined(_WIN64)
    ret = tss_sdk_load_win(shared_lib_file);
    #else
    ret = tss_sdk_load_linux(shared_lib_file);
    #endif
    if (ret == 0)
    {
        g_is_sdk_loaded = 1;
    }

    return ret;
}

void *
tss_sdk_get_syml_impl(const char *syml_name)
{
    void *func = NULL;
    if (!g_is_sdk_loaded)
    {
        return NULL;
    }

#if defined(WIN32) || defined(_WIN64)
    func = tss_sdk_get_syml_win(syml_name);
#else
    func = tss_sdk_get_syml_linux(syml_name);
#endif

    return func;
}

// 获取sdk通用的接口
// Get the sdk common interface
static const TssSdkBusiInterf *
tss_sdk_get_comm_interf(const TssSdkInitInfo* init_info)
{
    TssSdkGetBusiInterf func = NULL;
    func = (TssSdkGetBusiInterf)tss_sdk_get_syml_impl(TSS_SDK_GET_COMM_INTERF_SYM);
    if (func != NULL)
    {
        return func(init_info);
    }

    return NULL;
}


const TssSdkBusiInterf *
tss_sdk_load(const TSS_TCHAR *shared_lib_dir, const TssSdkInitInfo *init_info)
{
    TSS_TCHAR shared_lib_file[1024] = {0};
    int ret = 0;

    if (shared_lib_dir == NULL || init_info == NULL)
    {
        return NULL;
    }
	
#if (defined(WIN32) || defined(_WIN64))
	_sntprintf(shared_lib_file,
             sizeof(shared_lib_file),
             _T("%s/%s.%s"),
             shared_lib_dir,
             LIB_NAME,
             SUFFIX);
#else
    snprintf(shared_lib_file,
             sizeof(shared_lib_file),
             "%s/%s.%s",
             shared_lib_dir,
             LIB_NAME,
             SUFFIX);
#endif

    ret = tss_sdk_load_impl(shared_lib_file);
    if (ret != 0)
    {
        // load library fail
        return NULL;
    }

    return tss_sdk_get_comm_interf(init_info);
}

int
tss_sdk_unload()
{
    // 如果还没load，就回去吧
    if (!g_is_sdk_loaded)
    {
        return 0;
    }

    int rc = 0;
    TssSdkReleaseBusiInterf func = (TssSdkReleaseBusiInterf)tss_sdk_get_syml_impl(TSS_SDK_RELEASE_BUSI_INTERF_SYM);
    if (func != NULL)
    {
        // 调用接口释放函数
        // Call Interface release function
        func();
    }

    #if defined(WIN32) || defined(_WIN64)
    rc = FreeLibrary(g_handle);
    if (!rc)
    {
        return -2;
    }
    #else
    rc = dlclose(g_handle);
    if (rc != 0)
    {
        return -2;
    }
    #endif
    g_is_sdk_loaded = 0;
    fprintf(stdout, "Unload TSS SDK OK.\n");
    return 0;
}

// 获取某一特定业务的接口
// get the interface of a particular business
typedef const void *(*TssSdkGetInterf)(const void *init_data);
const void *tss_sdk_get_busi_interf(const char *syml_name, const void *data)
{
    TssSdkGetInterf f = (TssSdkGetInterf)tss_sdk_get_syml_impl(syml_name);
    if (f != NULL)
    {
        return f(data);
    }

    // 没找到对应的函数
    // Did not find the corresponding function
    return NULL;
}



