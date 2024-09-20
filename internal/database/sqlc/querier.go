// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package database

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateHost(ctx context.Context, arg CreateHostParams) (Host, error)
	CreateHostService(ctx context.Context, arg CreateHostServiceParams) (HostService, error)
	DeleteHostService(ctx context.Context, id int64) (HostService, error)
	GetHost(ctx context.Context, id int64) (Host, error)
	GetHostByHostname(ctx context.Context, hostName string) (Host, error)
	GetHostService(ctx context.Context, id int64) (GetHostServiceRow, error)
	GetHostServices(ctx context.Context, hostID sql.NullInt64) ([]GetHostServicesRow, error)
	GetHosts(ctx context.Context, arg GetHostsParams) ([]Host, error)
	GetHostsServicesToMonitor(ctx context.Context) ([]GetHostsServicesToMonitorRow, error)
	GetHostsWithServices(ctx context.Context, arg GetHostsWithServicesParams) ([]GetHostsWithServicesRow, error)
	GetInnactiveHostServiceByHostAndService(ctx context.Context, arg GetInnactiveHostServiceByHostAndServiceParams) (GetInnactiveHostServiceByHostAndServiceRow, error)
	GetServices(ctx context.Context) ([]Service, error)
	GetUserByEmail(ctx context.Context, email sql.NullString) (User, error)
	GetUserForAuth(ctx context.Context, email sql.NullString) (GetUserForAuthRow, error)
	ReactivateHostService(ctx context.Context, id int64) (HostService, error)
	UpdateHost(ctx context.Context, arg UpdateHostParams) (Host, error)
	UpdteHostServiceSchedule(ctx context.Context, arg UpdteHostServiceScheduleParams) (HostService, error)
}

var _ Querier = (*Queries)(nil)
