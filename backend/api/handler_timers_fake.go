package api

import (
	"math/rand"
	"net/http"
	"time"
)

// handleTimerFake fills the database with fake data
func (srv *Server) handleTimerFake() http.HandlerFunc {
	var names = []string{
		"Dinner with Jonas",
		"Pizzas with Emil",
		"Powernap",
		"Gone fishin'",
		"Review AWS Bills",
		"Bughunting",
		"Coffee break",
		"PTO",
		"Beach 🏖",
		"Reading a book",
		"R&D",
		"Doctor appointment",
	}
	var months = []time.Month{
		time.January,
		time.February,
		time.March,
		time.April,
		time.May,
		time.June,
		time.July,
		time.August,
		time.September,
		time.October,
		time.November,
		time.December,
	}
	type response struct {
		Ok     int `json:"ok"`
		Errors int `json:"errors"`
		NoEnd  int `json:"no_end"`
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		var (
			ok     int
			errors int
			noEnd  int
		)

		rand.Seed(time.Now().UnixNano())

		const year = 2019
		for m := 9; m <= 12; m++ {
			for d := 1; d <= 25; d++ {
				name := names[rand.Intn(len(names))]

				t, err := srv.timerService.Create(name)
				if err != nil {
					errors++
					continue
				}

				hh := rand.Intn(24)
				mm := rand.Intn(60)
				ss := rand.Intn(60)

				start := time.Date(year, months[m-1], d, hh, mm, ss, 0, time.UTC)
				t.DateStart.Time = start
				t.DateStart.Valid = true

				if rand.Float64() < 0.9 {
					end := start.Add(time.Duration(rand.Intn(3*24)) * time.Hour)
					t.DateEnd.Time = end
					t.DateEnd.Valid = true
				} else {
					t.DateEnd.Valid = false
					noEnd++
				}

				err = srv.timerService.Update(t)
				if err != nil {
					errors++
					continue
				}

				ok++
			}
		}

		resp := response{
			Ok:     ok,
			Errors: errors,
			NoEnd:  noEnd,
		}
		srv.json(rw, resp, http.StatusOK)
	}
}
