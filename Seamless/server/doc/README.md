# 光荣使命服务器文档

## 编译

### 设置 GOPATH
```
	set GOPATH=E:\\svn\\mb\\zeus\\branch_1;E:\\svn\\timefire\trunk\server`
```

### 安装 MinGW-w64

不然报错：
```
go install -race Room
# zeus/unitypx
exec: "gcc": executable file not found in %PATH%
```

安装时须选择 x86_64. 安装完成后添加 mingw64/bin 到PATH.

LiteIDE可选择win64环境，然后查看->编辑当前环境，更改PATH:
```
GOROOT=d:\go
PATH=D:\Program Files\mingw-w64\x86_64-7.2.0-posix-seh-rt_v5-rev1\mingw64\bin;%GOROOT%\bin;%PATH%
```

### `install.bat` 生成测试版示例
```
PATH=D:\Program Files\mingw-w64\x86_64-7.2.0-posix-seh-rt_v5-rev1\mingw64\bin;%PATH%
set GOPATH=E:\\svn\\mb\\zeus\\branch_1;E:\\svn\\timefire_server

cd bin

go build -race Gateway

go install -race Login
go install -race Room
go install -race Center
go install -race Match
go install -race IDIPServer
go install -race DataCenter
go install -race Lobby

cd ..
```
注意，运营版本不应该加 -race, 会影响性能。

## 运行

测试启动5个服：

1. Gateway
1. Login
1. Room
1. Match
1. Lobby

Windows下关闭时，使用 Ctrl-C，如果直接关窗口并立即重启容易出现重复ID错误。

### 配置

更改配置见：[配置上的一些坑（必看）.md](3.配置文件/配置上的一些坑（必看）.md)

### 运行

Room 需要
* unitypx.dll
* PxFoundation_x64.dll
* PhysX3Cooking_x64.dll
* PhysX3_x64.dll
* PhysX3Common_x64.dll 缺少不会报错，但会异常退出

Lobby 需要
* qos_client.dll
* qostrans.dll
* ./tss_sdk.dll 动态加载
* 忽略 ./config 错误日志。如果复制 (zeus\branch_1\src\test\config -> bin\config) 则会出现名字检查错误无法进入游戏。
```
2017-11-30 18:56:24.288581[ERROR]conf_path_file config file is not exist, path=./config/tss_sdk_conf_path.xml
2017-11-30 18:56:24.294581[ERROR]conf_file config file is not exist, path=./config/tss_sdk_conf.xml
...
2017-11-30 18:56:24.323581[ERROR][zenlib] ZEN_OS::flock fail. operation =8 ret =-1
2017-11-30 18:56:24.327581[ERROR][FAIL RETRUN]Fail in file [.\zen_lock_file_lock.cpp|153],function:void
__cdecl ZEN_File_Lock::unlock(void),fail info:ZEN_OS::flock LOCK_UN,return -1,last error 2.
```
* 2个tss_sdk.dll中任选一个。第一次运行可能会有如下错误，再次运行会正常。TSSSDK会在c:\temp\下创建文件。
```
panic: TSSSDK 初始化失败

goroutine 1 [running]:
zeus/tsssdk.Init(0x283d)
        E:/svn/mb/zeus/branch_1/src/zeus/tsssdk/tsssdk.go:179 +0xfa
main.(*LobbySrv).Init(0xc04240a180, 0x0, 0x0)
        E:/svn/timefire_server/src/Lobby/LobbyServer.go:104 +0x43f
zeus/server.(*Server).init(0xc042410150, 0xe1c1e0, 0xc04240a180)
        E:/svn/mb/zeus/branch_1/src/zeus/server/server.go:99 +0x1f5
zeus/server.(*Server).Run(0xc042410150)
        E:/svn/mb/zeus/branch_1/src/zeus/server/server.go:122 +0x44
main.main()
        E:/svn/timefire_server/src/Lobby/Lobby.go:29 +0x172
```

### Unity客户端连接
Unity 客户端运行 Main 场景。
添加本服地址到 
client\timefire\Assets\RawResources\Lua\Tbx\ipconfig.lua

客户端设置角色名报: `onSetNameResult result : 2`, 因为
`tsssdk.JudgeUserInputName(name)`总是返回 "TSSSDK UIC Judge Name Failed"，
暂时禁用它。

