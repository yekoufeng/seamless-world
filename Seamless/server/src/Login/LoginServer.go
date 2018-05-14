package main

import (
	"net"
	"net/http"
	"strings"
	"zeus/login"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"golang.org/x/net/netutil"
)

// Server 登陆服务器
type Server struct {
	*login.App

	addr          string
	port          string
	forceCreate   bool
	forceChannel  bool
	requireHost   bool
	msdkAddr      string
	allowVersion  string
	versionCheck  bool
	allowGrade    int
	initGrade     int
	whitelistfile string

	config *LoginConfig
}

// NewServer 创建一个新的登陆服务器
// 帧率为20帧
func NewServer() *Server {
	srv := new(Server)

	srv.port = viper.GetString("Login.Port")
	srv.addr = viper.GetString("Login.Addr")
	srv.App = login.NewApp()
	srv.loadConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if strings.Contains(e.Name, "server.json") {
			srv.loadConfig()
		}
	})

	return srv
}

// Run 启动
func (s *Server) Run() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	if viper.GetBool("Login.Log") {
		api.Use(&rest.AccessLogApacheMiddleware{
			Format: rest.CombinedLogFormat,
		})
	}
	if s.requireHost {
		api.Use(&login.RequireHostMiddleware{})
	}
	router, err := rest.MakeRouter(
		rest.Post("/login", s.loginHandler),
		rest.Post("/create", s.createHandler),
		rest.Post("/activate", s.AccountActive),
		rest.Post("/verify", s.cdkeyVerify),
	)
	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	go func() {
		l, err := net.Listen("tcp", s.addr+":"+s.port)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		maxConns := viper.GetInt("Login.MaxConns")
		if maxConns != 0 {
			l = netutil.LimitListener(l, maxConns)
		}

		err = http.Serve(l, api.MakeHandler())
		if err != nil {
			panic(err)
		}

		// err = http.ListenAndServe(s.addr+":"+s.port, api.MakeHandler())
		// if err != nil {
		// 	panic(err)
		// }
	}()

	return
}

func (s *Server) loadConfig() {
	s.forceCreate = viper.GetBool("Login.ForceCreate")
	s.forceChannel = viper.GetBool("Login.ForceChannel")
	s.requireHost = viper.GetBool("Login.RequireHost")
	s.allowGrade = viper.GetInt("Login.AllowGrade")
	s.versionCheck = viper.GetBool("Login.VersionCheck")
	s.allowVersion = viper.GetString("Login.AllowVersion")
	s.msdkAddr = viper.GetString("Config.MSDKAddr")
	s.initGrade = viper.GetInt("Login.InitGrade")
	s.whitelistfile = viper.GetString("Login.CdkeyFile")
	s.config = &LoginConfig{}
	s.config.HBEnable = viper.GetBool("Config.HeartBeat")

	log.Info("Load Config:")
	log.Info("ForceCreate:", s.forceCreate)
	log.Info("ForceChannel:", s.forceChannel)
	log.Info("RequireHost:", s.requireHost)
	log.Info("AllowGrade:", s.allowGrade)
	log.Info("InitGrade:", s.initGrade)
	log.Info("VersionCheck:", s.versionCheck)
	log.Info("AllowVersion:", s.allowVersion)
	log.Info("MSDKAddr:", s.msdkAddr)

	log.Info("Config:", s.config)
}
