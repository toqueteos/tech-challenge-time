package api

import (
	"expvar"
	"net/http"
)

func (srv *Server) handleTimer() http.HandlerFunc {
	var (
		expTotal = expvar.NewInt("handle_timer_index_ok")
	)

	return func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, "static/index.html")
		expTotal.Add(1)
	}
}
