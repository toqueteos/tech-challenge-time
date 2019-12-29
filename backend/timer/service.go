package timer

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Service defines the business logic operations that might be performed on a
// particular Timer.
type Service interface {
	Repo

	Start(int64) (*Timer, error)
	Stop(int64) (*Timer, error)

	ActiveTimers() ([]*Timer, error)
	TimersOfDay() ([]*Timer, error)
	TimersOfWeek() ([]*Timer, error)
	TimersOfMonth() ([]*Timer, error)
}

var _ Service = &service{}

type service struct {
	Repo
	db *sqlx.DB
}

func NewService(r Repo) Service {
	return &service{db: r.DB(), Repo: r}
}

func (s *service) Start(ID int64) (*Timer, error) {
	t, err := s.ByID(ID)
	if err != nil {
		return nil, err
	}

	t.DateStart.Time = time.Now()
	t.DateStart.Valid = true
	err = s.Update(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *service) Stop(ID int64) (*Timer, error) {
	t, err := s.ByID(ID)
	if err != nil {
		return nil, err
	}

	t.DateEnd.Time = time.Now()
	t.DateEnd.Valid = true
	err = s.Update(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *service) ActiveTimers() ([]*Timer, error) {
	const query = `
		SELECT *
		FROM timers
		WHERE date_start IS NOT NULL
		AND date_end IS NULL
		ORDER BY date_start ASC
		`

	var timers []*Timer
	err := s.db.Select(&timers, query)
	if err != nil {
		return nil, err
	}

	return timers, nil
}

// Pagination is tricky so I've skipped it entirely because the number of
// users will be just one and any pagination approach will add unnecessary
// complexity to our SQL queries too soon.
// - https://use-the-index-luke.com/no-offset
// - https://www.citusdata.com/blog/2016/03/30/five-ways-to-paginate/
func (s *service) timersOf(interval string) ([]*Timer, error) {
	const query = `
		SELECT *
		FROM timers
		WHERE date_start > current_timestamp - interval '%s'
		ORDER BY date_start ASC
		`

	var timers []*Timer
	err := s.db.Select(&timers, fmt.Sprintf(query, interval))
	if err != nil {
		return nil, err
	}

	return timers, nil
}

func (s *service) TimersOfDay() ([]*Timer, error) {
	return s.timersOf("1 day")
}

func (s *service) TimersOfWeek() ([]*Timer, error) {
	return s.timersOf("1 week")
}

func (s *service) TimersOfMonth() ([]*Timer, error) {
	return s.timersOf("1 month")
}
