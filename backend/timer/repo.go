package timer

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	DB() *sqlx.DB
	ByID(id int64) (*Timer, error)

	Create(name string) (*Timer, error)
	Update(timer *Timer) error
	Delete(id int64) error
}

var _ Repo = &repo{}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) DB() *sqlx.DB {
	return r.db
}

func (r *repo) ByID(id int64) (*Timer, error) {
	var u Timer
	err := r.db.Get(&u, `SELECT * FROM timers WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repo) Create(name string) (*Timer, error) {
	var id int64
	query := `INSERT INTO timers (id, name, date_start) VALUES (DEFAULT, $1, $2) RETURNING id`
	err := r.db.Get(&id, query, name, time.Now())
	if err != nil {
		return nil, err
	}

	return r.ByID(id)
}

func (r *repo) Update(timer *Timer) error {
	if timer.DateStart.Valid {
		timer.DateStart.Time = timer.DateStart.Time.UTC()
	}
	if timer.DateEnd.Valid {
		timer.DateEnd.Time = timer.DateEnd.Time.UTC()
	}

	query := `UPDATE timers
		SET
			name=$2,
			date_start=$3,
			date_end=$4
		WHERE id=$1`
	res, err := r.db.Exec(query, timer.ID, timer.Name, timer.DateStart, timer.DateEnd)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows != 1 {
		return ErrTimerInvalidID
	}

	return err
}

func (r *repo) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM timers WHERE id=$1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return ErrTimerInvalidID
	}

	return err
}
