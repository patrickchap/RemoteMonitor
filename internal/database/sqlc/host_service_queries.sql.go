// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: host_service_queries.sql

package database

import (
	"context"
	"database/sql"
)

const createHostService = `-- name: CreateHostService :one
INSERT INTO host_services(host_id, service_id)
VALUES(?, ?)
RETURNING id, host_id, service_id, active, schedule_number, schedule_unit, last_check, last_updated, status
`

type CreateHostServiceParams struct {
	HostID    sql.NullInt64 `json:"host_id"`
	ServiceID sql.NullInt64 `json:"service_id"`
}

func (q *Queries) CreateHostService(ctx context.Context, arg CreateHostServiceParams) (HostService, error) {
	row := q.db.QueryRowContext(ctx, createHostService, arg.HostID, arg.ServiceID)
	var i HostService
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.ServiceID,
		&i.Active,
		&i.ScheduleNumber,
		&i.ScheduleUnit,
		&i.LastCheck,
		&i.LastUpdated,
		&i.Status,
	)
	return i, err
}

const deleteHostService = `-- name: DeleteHostService :one
UPDATE host_services SET 
    active = 0
WHERE id = ?
RETURNING id, host_id, service_id, active, schedule_number, schedule_unit, last_check, last_updated, status
`

func (q *Queries) DeleteHostService(ctx context.Context, id int64) (HostService, error) {
	row := q.db.QueryRowContext(ctx, deleteHostService, id)
	var i HostService
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.ServiceID,
		&i.Active,
		&i.ScheduleNumber,
		&i.ScheduleUnit,
		&i.LastCheck,
		&i.LastUpdated,
		&i.Status,
	)
	return i, err
}

const getHostService = `-- name: GetHostService :one
SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.last_check, hs.last_updated, hs.status, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.id = ?
AND hs.active = 1
`

type GetHostServiceRow struct {
	ID             int64          `json:"id"`
	HostID         sql.NullInt64  `json:"host_id"`
	ServiceID      sql.NullInt64  `json:"service_id"`
	Active         sql.NullInt64  `json:"active"`
	ScheduleNumber sql.NullInt64  `json:"schedule_number"`
	ScheduleUnit   sql.NullString `json:"schedule_unit"`
	LastCheck      sql.NullTime   `json:"last_check"`
	LastUpdated    sql.NullTime   `json:"last_updated"`
	Status         sql.NullString `json:"status"`
	HostName       string         `json:"host_name"`
	ServiceName    sql.NullString `json:"service_name"`
}

func (q *Queries) GetHostService(ctx context.Context, id int64) (GetHostServiceRow, error) {
	row := q.db.QueryRowContext(ctx, getHostService, id)
	var i GetHostServiceRow
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.ServiceID,
		&i.Active,
		&i.ScheduleNumber,
		&i.ScheduleUnit,
		&i.LastCheck,
		&i.LastUpdated,
		&i.Status,
		&i.HostName,
		&i.ServiceName,
	)
	return i, err
}

const getHostServices = `-- name: GetHostServices :many
SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.last_check, hs.last_updated, hs.status, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.host_id = ?
AND hs.active = 1
`

type GetHostServicesRow struct {
	ID             int64          `json:"id"`
	HostID         sql.NullInt64  `json:"host_id"`
	ServiceID      sql.NullInt64  `json:"service_id"`
	Active         sql.NullInt64  `json:"active"`
	ScheduleNumber sql.NullInt64  `json:"schedule_number"`
	ScheduleUnit   sql.NullString `json:"schedule_unit"`
	LastCheck      sql.NullTime   `json:"last_check"`
	LastUpdated    sql.NullTime   `json:"last_updated"`
	Status         sql.NullString `json:"status"`
	HostName       string         `json:"host_name"`
	ServiceName    sql.NullString `json:"service_name"`
}

func (q *Queries) GetHostServices(ctx context.Context, hostID sql.NullInt64) ([]GetHostServicesRow, error) {
	rows, err := q.db.QueryContext(ctx, getHostServices, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetHostServicesRow{}
	for rows.Next() {
		var i GetHostServicesRow
		if err := rows.Scan(
			&i.ID,
			&i.HostID,
			&i.ServiceID,
			&i.Active,
			&i.ScheduleNumber,
			&i.ScheduleUnit,
			&i.LastCheck,
			&i.LastUpdated,
			&i.Status,
			&i.HostName,
			&i.ServiceName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInnactiveHostServiceByHostAndService = `-- name: GetInnactiveHostServiceByHostAndService :one
SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.last_check, hs.last_updated, hs.status, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.host_id = ?
AND hs.service_id = ?
AND hs.active = 0
`

type GetInnactiveHostServiceByHostAndServiceParams struct {
	HostID    sql.NullInt64 `json:"host_id"`
	ServiceID sql.NullInt64 `json:"service_id"`
}

type GetInnactiveHostServiceByHostAndServiceRow struct {
	ID             int64          `json:"id"`
	HostID         sql.NullInt64  `json:"host_id"`
	ServiceID      sql.NullInt64  `json:"service_id"`
	Active         sql.NullInt64  `json:"active"`
	ScheduleNumber sql.NullInt64  `json:"schedule_number"`
	ScheduleUnit   sql.NullString `json:"schedule_unit"`
	LastCheck      sql.NullTime   `json:"last_check"`
	LastUpdated    sql.NullTime   `json:"last_updated"`
	Status         sql.NullString `json:"status"`
	HostName       string         `json:"host_name"`
	ServiceName    sql.NullString `json:"service_name"`
}

func (q *Queries) GetInnactiveHostServiceByHostAndService(ctx context.Context, arg GetInnactiveHostServiceByHostAndServiceParams) (GetInnactiveHostServiceByHostAndServiceRow, error) {
	row := q.db.QueryRowContext(ctx, getInnactiveHostServiceByHostAndService, arg.HostID, arg.ServiceID)
	var i GetInnactiveHostServiceByHostAndServiceRow
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.ServiceID,
		&i.Active,
		&i.ScheduleNumber,
		&i.ScheduleUnit,
		&i.LastCheck,
		&i.LastUpdated,
		&i.Status,
		&i.HostName,
		&i.ServiceName,
	)
	return i, err
}

const reactivateHostService = `-- name: ReactivateHostService :one
UPDATE host_services SET 
    active = 1
WHERE id = ?
RETURNING id, host_id, service_id, active, schedule_number, schedule_unit, last_check, last_updated, status
`

func (q *Queries) ReactivateHostService(ctx context.Context, id int64) (HostService, error) {
	row := q.db.QueryRowContext(ctx, reactivateHostService, id)
	var i HostService
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.ServiceID,
		&i.Active,
		&i.ScheduleNumber,
		&i.ScheduleUnit,
		&i.LastCheck,
		&i.LastUpdated,
		&i.Status,
	)
	return i, err
}

const updteHostServiceSchedule = `-- name: UpdteHostServiceSchedule :one
UPDATE host_services SET 
    schedule_number = ?,
    schedule_unit = ?
WHERE id = ?
RETURNING id, host_id, service_id, active, schedule_number, schedule_unit, last_check, last_updated, status
`

type UpdteHostServiceScheduleParams struct {
	ScheduleNumber sql.NullInt64  `json:"schedule_number"`
	ScheduleUnit   sql.NullString `json:"schedule_unit"`
	ID             int64          `json:"id"`
}

func (q *Queries) UpdteHostServiceSchedule(ctx context.Context, arg UpdteHostServiceScheduleParams) (HostService, error) {
	row := q.db.QueryRowContext(ctx, updteHostServiceSchedule, arg.ScheduleNumber, arg.ScheduleUnit, arg.ID)
	var i HostService
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.ServiceID,
		&i.Active,
		&i.ScheduleNumber,
		&i.ScheduleUnit,
		&i.LastCheck,
		&i.LastUpdated,
		&i.Status,
	)
	return i, err
}
