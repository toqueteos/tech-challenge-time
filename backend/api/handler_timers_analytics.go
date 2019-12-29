package api

import (
	"expvar"
	"net/http"

	"pento-tech-challenge/backend/timer"

	"go.uber.org/zap"
)

func (srv *Server) handleTimerAnalytics() http.HandlerFunc {
	var (
		expTotal = expvar.NewInt("handle_timer_analytics_day_ok")
		expErr   = expvar.NewInt("handle_timer_analytics_day_err")
	)

	return func(rw http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		period := q.Get("period")

		var (
			timers []*timer.Timer
			err    error
		)

		switch period {
		case "day":
			timers, err = srv.timerService.TimersOfDay()
		case "week":
			timers, err = srv.timerService.TimersOfWeek()
		case "month":
			timers, err = srv.timerService.TimersOfMonth()
		}

		if err != nil {
			srv.logger.Error(
				"HANDLE_TIMER_ANALYTICS_FETCH_FAIL",
				zap.String("period", period),
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
