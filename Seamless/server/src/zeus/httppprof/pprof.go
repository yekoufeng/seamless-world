package httppprof

import (
	"net/http"
	"runtime/pprof"
)

func handlerGoroutine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}

func handlerHeap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("heap")
	p.WriteTo(w, 1)
}

func handlerThreadCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("threadcreate")
	p.WriteTo(w, 1)
}

func handlerBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("block")
	p.WriteTo(w, 1)
}

func StartPProf(addr string) {
	http.HandleFunc("/goroutine", handlerGoroutine)
	http.HandleFunc("/heap", handlerHeap)
	http.HandleFunc("/threadcreate", handlerThreadCreate)
	http.HandleFunc("/block", handlerBlock)
	http.ListenAndServe(addr, nil)
}
