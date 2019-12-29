package api

import (
	"expvar"
	"net/http"
)

func (srv *Server) setRoutes() {
	fs := http.FileServer(http.Dir("static"))

	srv.router.Handle("/_/debug/vars", expvar.Handler())
	srv.router.Handle("/static/", http.StripPrefix("/static/", fs))
	srv.router.HandleFunc("/", srv.handleTimer())

	srv.router.HandleFunc("/timer/_fake", srv.handleTimerFake())
	srv.router.HandleFunc("/timer/start", srv.handleTimerStart())
	srv.router.HandleFunc("/timer/stop", srv.handleTimerStop())
	srv.router.HandleFunc("/timer/active", srv.handleTimerActive())
	srv.router.HandleFunc("/timer/analytics", srv.handleTimerAnalytics())
}
