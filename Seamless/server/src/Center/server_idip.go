package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

// 启动idip服务
func (srv *Server) startIDIP() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	router, err := rest.MakeRouter(
		rest.Post("/idip", srv.idipQuery),
	)
	if err != nil {
		log.Error(err)
	}
	api.SetApp(router)
	err = http.ListenAndServe(srv.idipAddr+":"+srv.idipPort, api.MakeHandler())
	if err != nil {
		log.Error("listen error", err)
	}
}

// 处理idip指令
func (srv *Server) idipQuery(w rest.ResponseWriter, r *rest.Request) {

}
