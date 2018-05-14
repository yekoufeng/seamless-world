package main

import (
	"net/http"
	"sync"
	"time"
	"zeus/dbservice"

	"zeus/iserver"

	jwt "github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// AdminServer 管理服务器
type AdminServer struct {
	serverMap   map[uint64]*iserver.ServerInfo
	serverMapRW sync.RWMutex
}

var srvInst *AdminServer

// GetAdminServer 获取管理服务器实例
func GetAdminServer() *AdminServer {
	if srvInst == nil {
		srvInst = &AdminServer{}
	}
	return srvInst
}

// Init 初始化服务器
func (srv *AdminServer) Init() {
	if err := srv.initBaseInfo(); err != nil {
		log.Error(err)
	}
}

// Start 启动服务器
func (srv *AdminServer) Start() {
	srv.startHTTP()
}

// 启动http服务
func (srv *AdminServer) startHTTP() {
	jwtMiddleware := &jwt.JWTMiddleware{
		Key:           []byte("emag suez"),
		Realm:         "zeus admin",
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: srv.verifyPassword,
	}

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
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login"
		},
		IfTrue: jwtMiddleware,
	})
	router, err := rest.MakeRouter(
		rest.Post("/login", jwtMiddleware.LoginHandler),
		rest.Post("/cmds", srv.ServerCmd),
		rest.Get("/servers", srv.ServerList),
	)
	if err != nil {
		log.Error(err)
	}
	api.SetApp(router)
	e := http.ListenAndServe(":3300", api.MakeHandler())
	if e != nil {
		log.Error("login listen error", e)
	}
}

func (srv *AdminServer) verifyPassword(username string, password string) bool {
	uid, err := dbservice.GetUID(username)
	if err != nil {
		log.Error(err)
		return false
	}
	if uid == 0 {
		log.Error("用户不存在", username)
		return false
	}
	if uid > 1000 {
		log.Error("无权限", username, uid)
		return false
	}

	return dbservice.Account(uid).VerifyPassword(password)
}
