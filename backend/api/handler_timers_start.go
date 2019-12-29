package api

import (
	"expvar"
	"net/http"

	"go.uber.org/zap"
)

func (srv *Server) handleTimerStart() http.HandlerFunc {
	var (
		expTotal = expvar.NewInt("handle_timer_start_ok")
		expErr   = expvar.NewInt("handle_timer_start_err")
	)

	type request struct {
		Name string `json:"name"`
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			srv.logger.Error(
				"HANDLE_TIMER_START_METHOD_FAIL",
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
				"HANDLE_TIMER_START_DECODE_FAIL",
				zap.Error(err),
			)
			srv.jsonErr(rw, "400 Bad Request", http.StatusBadRequest)
			expErr.Add(1)
			return
		}

		srv.logger.Debug(
			"HANDLE_TIMER_START",
			zap.Any("input", input),
		)

		t, err := srv.timerService.Create(input.Name)
		if err != nil {
			srv.logger.Error(
				"HANDLE_TIMER_START_TIMERSERVICE_CREATE_FAIL",
				zap.String("name", input.Name),
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
