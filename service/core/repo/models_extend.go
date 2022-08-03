package core_repo

import "time"

type Column struct {
	Name      string `json:"name" validate:"required,alphanum,gte=2,lte=30"`
	Type      string `json:"type" validate:"required,oneof=integer smallint bigint decimal numeric real 'double precision' smallserial serial bigserial varchar char character text timestamp 'timestamp with time zone' 'timestamp without time zone' date 'time with time zone' time 'time without time zone' bool boolean bit 'bit varying' cidr inet macaddr macaddr8 json jsonb money uuid"`
	Length    int64  `json:"length"`
	Primary   bool   `json:"primary"`
	Unique    bool   `json:"unique"`
	Required  bool   `json:"required"`
	Precision int32  `json:"precision"`
	Scale     int32  `json:"scale"`
	Default   string `json:"default"`
}

type TableSchema struct {
	ID      int64     `json:"id" validate:"required,numeric"`
	Name    string    `json:"name" validate:"required,alphanum,gte=3,lte=60"`
	UserID  int64     `json:"user_id" validate:"required,numeric"`
	Columns []Column  `json:"columns"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
