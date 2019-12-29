package null

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Time implements encoding/json {Marshaler,Unmarshaler} interfaces for sql.NullTime
// We could just use gopkg.in/guregu/null.v3 and not implement this ourselves.
type Time struct {
	sql.NullTime
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.NullTime.Time)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	t.Valid = x != nil
	if t.Valid {
		t.Time = *x
	}

	return nil
}
