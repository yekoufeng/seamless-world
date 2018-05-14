package admin

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// IAdmin 控制台服务接口
type IAdmin interface {
	ConfigHTTPAdmin(addr string, port uint64, admin string)
}

// IServerCommand 后代服务器类处理命令行的接口
type IServerCommand interface {
	HandleCommand([]string) *CmdResp
}

// iServerAdminCtrl 服务器需要提供的控制接口
type iServerAdminCtrl interface {
	HandleCommand([]string) *CmdResp
	GetSrvID() uint64
	GetToken() string
}

// CmdResp 命令行运行结果
type CmdResp struct {
	Result    int
	ResultStr string
	Data      []byte
}

// CmdReq 命令行
type CmdReq struct {
	cmd    string
	isHTTP bool
}

// Admin 控制台实现类
type Admin struct {
	srvAdminCtrl iServerAdminCtrl

	addr         string
	port         uint64
	adminDomain  []string
	consoleCmds  chan *CmdReq
	consoleRespC chan *CmdResp
}

func (r *CmdResp) String() string {
	return fmt.Sprintf("%+v", *r)
}

// NewAdminApp 创建控制台app
func NewAdminApp(srv iServerAdminCtrl) *Admin {
	app := new(Admin)
	app.srvAdminCtrl = srv
	app.consoleCmds = make(chan *CmdReq)
	app.consoleRespC = make(chan *CmdResp)
	return app
}

// ConfigHTTPAdmin 配置控制台相关
// port: http命令行端口
// admin: 允许执行http命令的域名或IP地址
func (app *Admin) ConfigHTTPAdmin(addr string, port uint64, admin string) {
	app.addr = addr
	app.port = port
	app.adminDomain = strings.Split(admin, ",")
}

// ProcCmds 处理命令行
func (app *Admin) ProcCmds() {
	select {
	case cmd := <-app.consoleCmds:
		app.handleCmd(cmd)
	default:
		return
	}
}

// StartHTTP 启动HTTP处理命令行服务
func (app *Admin) StartHTTP() {
	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods:                []string{"GET", "POST", "PUT"},
		AllowedHeaders:                []string{"Accept", "Content-Type", "X-Custom-Header", "Origin", "authorization"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	router, err := rest.MakeRouter(
		rest.Post("/exec", app.execHandler),
	)
	if err != nil {
		log.Error(err)
	}
	api.SetApp(router)
	addr := fmt.Sprintf("%s:%d", app.addr, app.port)
	log.Info("Start Admin Service", addr)
	e := http.ListenAndServe(addr, api.MakeHandler())
	if e != nil {
		log.Error("Start admin http service failed", e)
	}
}

// StartConsole 启动命令行服务
func (app *Admin) StartConsole() {
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, err := reader.ReadLine()
		if err != nil || data == nil {
			continue
		}
		cmd := string(data)
		if strings.HasPrefix(cmd, "cmd") && len(cmd) >= 5 {
			cmd = string(data[4:])
			app.consoleCmds <- &CmdReq{
				cmd:    cmd,
				isHTTP: false,
			}
		}
	}
}

// GetConsolePort 获取控制台端口
func (app *Admin) GetConsolePort() uint64 {
	return app.port
}

func (app *Admin) execHandler(w rest.ResponseWriter, r *rest.Request) {
	var cmd struct {
		ServerID uint64
		Command  string
		Token    string
	}
	err := r.DecodeJsonPayload(&cmd)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", 400)
		return
	}

	if cmd.ServerID != app.srvAdminCtrl.GetSrvID() {
		log.Errorf("Server id not match. Target:%d My:%d\n", cmd.ServerID, app.srvAdminCtrl.GetSrvID())
		rest.Error(w, "参数错误, 不是目标服务器", 400)
		return
	}

	if cmd.Token != app.srvAdminCtrl.GetToken() {
		log.Error("Token error", cmd.Token)
		rest.Error(w, "Token错误", 400)
		return
	}

	if !app.verifyDomain(r.Header.Get("Origin")) {
		log.Error("invalid Domain", r.Header.Get("Origin"))
		rest.Error(w, "无效Domain", 400)
		return
	}

	app.consoleCmds <- &CmdReq{
		cmd:    cmd.Command,
		isHTTP: true,
	}
	resp := <-app.consoleRespC
	w.WriteJson(resp)
}

func (app *Admin) handleCmd(cmd *CmdReq) {
	log.Info("Admin: HandlerCmd", cmd.cmd)

	c := strings.Split(cmd.cmd, " ")
	resp := app.srvAdminCtrl.HandleCommand(c)

	// 未知命令行
	if resp == nil {
		resp = &CmdResp{}
		resp.Result = -1
		resp.ResultStr = "未知命令行"
	}

	if cmd.isHTTP {
		app.consoleRespC <- resp
	} else {
		fmt.Println(resp)
	}
}

func (app *Admin) verifyDomain(origin string) bool {
	for _, v := range app.adminDomain {
		if strings.Contains(origin, v) {
			return true
		}
	}
	return false
}
