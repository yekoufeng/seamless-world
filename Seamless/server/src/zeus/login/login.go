package login

import (
	"math"
	"zeus/dbservice"
	"zeus/events"
	"zeus/iserver"
	"zeus/server"
	"zeus/serverMgr"

	log "github.com/cihub/seelog"
)

// App 登录服务
type App struct {
	*events.GlobalEvents
}

// NewApp 创建新的登录服务器
func NewApp() *App {
	app := &App{}
	app.GlobalEvents = events.NewGlobalEventsInst()
	return app
}

// Start 启动登录服务
func (app *App) Start() {
	// app.GlobalEvents.AddListener("servernotify", app, "RefreshGatewayInfo")
}

// Stop 停止登录服务
func (app *App) Stop() {
	// app.GlobalEvents.RemoveListenerByObjInst(app)
}

// DoCreateNewUser 创建新用户
// 根据用户名和密码创建新用户
// 创建成功后返回uid
func (app *App) DoCreateNewUser(user, pwd string, grade uint32) (uint64, error) {
	uid, err := server.CreateNewUID()
	if err != nil {
		log.Error(err)
		return math.MaxUint64, err
	}

	err = dbservice.Account(uid).SetPassword(pwd)
	if err != nil {
		log.Error(err)
		return math.MaxUint64, err
	}

	err = dbservice.Account(uid).SetUsername(user)
	if err != nil {
		log.Error(err)
		return math.MaxUint64, err
	}

	err = dbservice.Account(uid).SetGrade(grade)
	if err != nil {
		log.Error(err)
		return math.MaxUint64, err
	}

	return uint64(uid), nil
}

// RefreshGatewayInfo 刷新Gateway信息
// func (app *App) RefreshGatewayInfo() {
// 	serverMgr.GetServerMgr().GetServerList()
// }

// GetGatewayAddr 获取网关信息
func (app *App) GetGatewayAddr() (string, error) {
	//TODO 网关信息的更新机制
	// app.GlobalEvents.HandleEvent()

	gateway, err := serverMgr.GetServerMgr().GetServerByType(iserver.ServerTypeGateway)
	if err != nil {
		return "", err
	}
	return gateway.OuterAddress, nil
}
