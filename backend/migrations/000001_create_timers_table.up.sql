CREATE TABLE IF NOT EXISTS timers (
    id bigserial PRIMARY KEY,
    name varchar(64) NOT NULL,
    date_start timestamp with time zone,
    date_end timestamp with time zone
);

CREATE INDEX idx_timers_date_start ON timers(date_start);
CREATE INDEX idx_timers_date_end ON timers(date_end);
