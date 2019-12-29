package api

import (
	"expvar"
	"net/http"

	"go.uber.org/zap"
)

func (srv *Server) handleTimerStop() http.HandlerFunc {
	var (
		expTotal = expvar.NewInt("handle_timer_stop_ok")
		expErr   = expvar.NewInt("handle_timer_stop_err")
	)

	type request struct {
		ID int64 `json:"id"`
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			srv.logger.Error(
				"HANDLE_TIMER_STOP_METHOD_FAIL",
				zap.String("method", req.Method),
			)
			srv.jsonErr(rw, "405 Method Not Allowed", http.StatusMethodNotAllowed)
			expErr.Add(1)
			return
		}

		var input request
		err := srv.decode(req, &input)
		if err != nil {
			srv.logger.Error(
				"HANDLE_TIMER_STOP_DECODE_FAIL",
				zap.Error(err),
			)
			srv.jsonErr(rw, "400 Bad Request", http.StatusBadRequest)
			expErr.Add(1)
			return
		}

		srv.logger.Debug(
			"HANDLE_TIMER_STOP",
			zap.Any("input", input),
		)

		t, err := srv.timerService.Stop(input.ID)
		if err != nil {
			srv.logger.Error(
				"HANDLE_TIMER_STOP_TIMERSERVICE_STOP_FAIL",
				zap.Int64("id", input.ID),
				zap.Error(err),
			)
			srv.jsonErr(rw, "500 Internal Server Error", http.StatusInternalServerError)
			expErr.Add(1)
			return
		}

		srv.json(rw, t, http.StatusOK)
		expTotal.Add(1)
	}
}
