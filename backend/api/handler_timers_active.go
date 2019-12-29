package api

import (
	"expvar"
	"net/http"

	"go.uber.org/zap"
)

func (srv *Server) handleTimerActive() http.HandlerFunc {
	var (
		expTotal = expvar.NewInt("handle_timer_active_ok")
		expErr   = expvar.NewInt("handle_timer_active_err")
	)

	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			srv.logger.Error(
				"HANDLE_TIMER_ACTIVE_METHOD_FAIL",
				zap.String("method", req.Method),
			)
			srv.jsonErr(rw, "405 Method Not Allowed", http.StatusMethodNotAllowed)
			expErr.Add(1)
			return
		}

		timers, err := srv.timerService.ActiveTimers()
		if err != nil {
			srv.logger.Error(
				"HANDLE_TIMER_START_TIMERSERVICE_ACTIVE_TIMERS_FAIL",
				zap.Error(err),
			)
			srv.jsonErr(rw, "500 Internal Server Error", http.StatusInternalServerError)
			expErr.Add(1)
			return
		}

		srv.logger.Debug(
			"HANDLE_TIMER_ACTIVE",
			zap.Any("active", len(timers)),
		)

		srv.json(rw, timers, http.StatusOK)
		expTotal.Add(1)
	}
}
