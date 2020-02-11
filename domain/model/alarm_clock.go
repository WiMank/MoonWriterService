package model

import "github.com/jackc/pgtype"

type AlarmClock struct {
	UserId      int          `db:"user_id"`
	Enabled     bool         `db:"enabled"`
	AlarmTime   string       `db:"alarm_time"`
	AlarmDays   pgtype.JSONB `db:"alarm_days"`
	Description string       `db:"description"`
}
