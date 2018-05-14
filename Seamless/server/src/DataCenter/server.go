package main

import (
	"db"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// Server 服务器
type Server struct {
	innerAddr       string //内网地址
	innerPort       string //内网端口
	innerListen     string //内网实际监听地址
	innerListenPort string //内网实际监听端口
	outerAddr       string //外网地址
	outerPort       string //外网端口
	outerListen     string //外网实际监听地址
	outerListenPort string //外网实际监听端口
}

var srv *Server

// GetSrvInst 获取服务器实例
func GetSrvInst() *Server {
	if srv == nil {
		srv = &Server{}

		srv.innerAddr = viper.GetString("DataCenter.InnerAddr")
		srv.innerPort = viper.GetString("DataCenter.InnerPort")
		srv.innerListen = viper.GetString("DataCenter.InnerListen")
		srv.innerListenPort = viper.GetString("DataCenter.InnerListenPort")
		srv.outerAddr = viper.GetString("DataCenter.OuterAddr")
		srv.outerPort = viper.GetString("DataCenter.OuterPort")
		srv.outerListen = viper.GetString("DataCenter.OuterListen")
		srv.outerListenPort = viper.GetString("DataCenter.OuterListenPort")

		err := db.SetDataCenterAddr("DataCenterInnerAddr", srv.innerAddr+":"+srv.innerPort, "DataCenterOuterAddr", srv.outerAddr+":"+srv.outerPort)
		if err != nil {
			log.Error(err)
			log.Info("DataCenterAddr注册失败")
			return nil
		}

		log.Info("DataCenterAddr注册成功")
	}

	return srv
}

// Run 运行
func (srv *Server) Run() {
	srv.startInnerService()
	srv.startOuterService()

	log.Info("DataCenter Running!")
	return
}

// Stop 停止服务器
func (srv *Server) Stop() {
	log.Info("DataCenter Stoped")

	err := db.DelDataCenterAddr("DataCenterInnerAddr", "DataCenterOuterAddr")
	if err != nil {
		log.Error(err)
	}
	log.Info("DataCenterAddr Del!")
}

func (srv *Server) startInnerService() {
	log.Info("启动处理Room服 post过来每局的数据监听")
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	router, err := rest.MakeRouter(
	// rest.Post("/dataCenter", srv.postHandler), //有需要时在这加路由
	)
	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	go func() {
		err = http.ListenAndServe(srv.innerListen+":"+srv.innerListenPort, api.MakeHandler())
		if err != nil {
			panic(err)
		}
	}()
}

func (srv *Server) startOuterService() {
	log.Info("启动处理客户端请求Get的数据监听")
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	router, err := rest.MakeRouter(
	// rest.Get("/career/#season/#uid", srv.careerHandler), //有需要时在这加路由
	)
	if err != nil {
		panic(err)
	}
	api.SetApp(router)

	go func() {
		err = http.ListenAndServe(srv.outerListen+":"+srv.outerListenPort, api.MakeHandler())
		if err != nil {
			panic(err)
		}
	}()
}
