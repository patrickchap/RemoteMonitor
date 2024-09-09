package viewmodels

import "database/sql"

type Host struct {
	ID            int64          `json:"id"`
	HostName      string         `json:"host_name"`
	CanonicalName sql.NullString `json:"canonical_name"`
	Url           sql.NullString `json:"url"`
	Ip            sql.NullString `json:"ip"`
	Ipv6          sql.NullString `json:"ipv6"`
	Location      sql.NullString `json:"location"`
	Os            sql.NullString `json:"os"`
	Active        sql.NullInt64  `json:"active"`
	LastUpdated   sql.NullTime   `json:"last_updated"`
}
