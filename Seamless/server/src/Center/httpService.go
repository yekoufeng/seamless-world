package main

import (
	"db"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/cihub/seelog"
)

type HttpService struct {
	srv *Server
}

func NewHttpService(srv *Server) *HttpService {
	ret := &HttpService{
		srv: srv,
	}

	//go ret.Start()

	return ret
}

func (h *HttpService) Start() {
	h.startHttpService()
}

func (h *HttpService) startHttpService() error {
	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)
	router, err := rest.MakeRouter(
		rest.Post("/add", h.Add), rest.Post("/del", h.Del), rest.Post("/print", h.Print),
	)

	if err != nil {
		log.Error(err)
		return err
	}

	api.SetApp(router)

	err = http.ListenAndServe(":"+"8099", api.MakeHandler())
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("http 监听成功!")
	return nil
}

// AnnuonceData 公告数据
type AnnuonceDataTest struct {
	ID           uint64
	ServerID     uint64
	PlatID       uint64
	StartTime    string
	EndTime      string
	InternalTime int64
	Content      string
}

// AccountActiveRet 账号激活返回格式
type AnnuonceOperateRet struct {
	Ret string
}

func GetTimeStamp(timeStr string) (int64, error) {
	loc, _ := time.LoadLocation("Local")
	timeStamp, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return timeStamp.Unix(), err
}

func (h *HttpService) Add(w rest.ResponseWriter, r *rest.Request) {

	msg := AnnuonceDataTest{}
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", 400)
		return
	}

	ret := &AnnuonceOperateRet{Ret: "Add Success"}

	paras := &db.AnnuonceData{
		ID:           msg.ID,
		ServerID:     uint32(msg.ServerID),
		PlatID:       uint8(msg.PlatID),
		InternalTime: uint32(msg.InternalTime),
		Content:      msg.Content,
	}

	if tmpStartTime, err := GetTimeStamp(msg.StartTime); err != nil {
		paras.StartTime = uint32(tmpStartTime)
		ret.Ret = err.Error()
		w.WriteJson(ret)
		return
	}

	if tmpEndTime, err := GetTimeStamp(msg.EndTime); err != nil {
		paras.EndTime = uint32(tmpEndTime)
		ret.Ret = err.Error()
		w.WriteJson(ret)
		return
	}

	//if h.srv.annuonceMgr.AddAnnuouce(paras) == false {
	//	ret.Ret = "id repeat!"
	//}

	w.WriteJson(ret)
}

type DelAnnuoncingID struct {
	ID uint64
}

func (h *HttpService) Del(w rest.ResponseWriter, r *rest.Request) {

	msg := DelAnnuoncingID{}
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		log.Error(err)
		rest.Error(w, "参数错误", 400)
		return
	}

	ret := &AnnuonceOperateRet{Ret: "Del Success"}
	if h.srv.annuonceMgr.DelAnnuoucing(msg.ID) == false {
		ret.Ret = "id isn't in the annuoncing queue!"
	}

	w.WriteJson(ret)
}

func (h *HttpService) Print(w rest.ResponseWriter, r *rest.Request) {

	db.PrintAnnuoncingData()

	w.WriteJson(&AnnuonceOperateRet{Ret: "Print Success"})
}
