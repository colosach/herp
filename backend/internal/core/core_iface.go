package core

import (
	"context"
	db "herp/db/sqlc"
)

type Querier interface {
	CreateBusiness(ctx context.Context, params db.CreateBusinessParams) (db.Business, error)
	GetBusiness(ctx context.Context, id int32) (db.Business, error)
	UpdateBusiness(ctx context.Context, params db.UpdateBusinessParams) (db.Business, error)
	DeleteBusiness(ctx context.Context, id int32) error
	ListBusinesses(ctx context.Context) ([]db.Business, error)
}

type CoreInterface interface {
	CreateBusiness(ctx context.Context, params db.CreateBusinessParams) (db.Business, error)
	GetBusiness(ctx context.Context, id int32) (db.Business, error)
	UpdateBusiness(ctx context.Context, params db.UpdateBusinessParams) (db.Business, error)
	DeleteBusiness(ctx context.Context, id int32) error
	ListBusinesses(ctx context.Context) ([]db.Business, error)
}