package business

import (
	"context"
	db "herp/db/sqlc"
)

type Querier interface {
	CreateBusiness(ctx context.Context, params db.CreateBusinessParams) (db.Business, error)
	GetBusiness(ctx context.Context, id int32) (db.Business, error)
	UpdateBusiness(ctx context.Context, params db.UpdateBusinessParams) (db.Business, error)
	DeleteBusiness(ctx context.Context, id int32) (db.Business, error)
	ListBusinesses(ctx context.Context) ([]db.Business, error)
	CreateBranch(ctx context.Context, params db.CreateBranchParams) (db.Branch, error)
	GetBranch(ctx context.Context, id int32) (db.Branch, error)
	UpdateBranch(ctx context.Context, params db.UpdateBranchParams) (db.Branch, error)
	DeleteBranch(ctx context.Context, id int32) (db.Branch, error)
	ListBranches(ctx context.Context) ([]db.Branch, error)
	CreateStore(ctx context.Context, params db.CreateStoreParams) (db.Store, error)
	CreateStoreManager(ctx context.Context, params db.CreateStoreManagerParams) (db.StoreManager, error)
	DeleteStore(ctx context.Context, id int32) error
	GetStoreByCode(ctx context.Context, storeCode string) (db.Store, error)
	ListStores(ctx context.Context, params db.ListStoresParams) ([]db.Store, error)
	GetStoreByID(ctx context.Context, id int32) (db.Store, error)
	UpdateStore(ctx context.Context, params db.UpdateStoreParams) (db.Store, error)
	GetStoreManagerByID(ctx context.Context, id int32) (db.StoreManager, error)
	DeleteStoreManager(ctx context.Context, id int32) error
	LogActivity(ctx context.Context, params db.LogActivityParams) (db.ActivityLog, error)
	GetActivityLogs(ctx context.Context, limit int32) ([]db.ActivityLog, error)
}

type BusinessInterface interface {
	CreateBusinessWithBranch(ctx context.Context, params db.CreateBusinessParams) (db.Business, db.Branch, error)
	CreateBusiness(ctx context.Context, params db.CreateBusinessParams) (db.Business, error)
	GetBusiness(ctx context.Context, id int32) (db.Business, error)
	UpdateBusiness(ctx context.Context, params db.UpdateBusinessParams) (db.Business, error)
	DeleteBusiness(ctx context.Context, id int32) (db.Business, error)
	ListBusinesses(ctx context.Context) ([]db.Business, error)
	CreateBranch(ctx context.Context, params db.CreateBranchParams) (db.Branch, error)
	GetBranch(ctx context.Context, id int32) (db.Branch, error)
	UpdateBranch(ctx context.Context, params db.UpdateBranchParams) (db.Branch, error)
	DeleteBranch(ctx context.Context, id int32) (db.Branch, error)
	ListBranches(ctx context.Context) ([]db.Branch, error)
	CreateStore(ctx context.Context, params db.CreateStoreParams) (db.Store, error)
	CreateStoreManager(ctx context.Context, params db.CreateStoreManagerParams) (db.StoreManager, error)
	DeleteStore(ctx context.Context, id int32) error
	GetStoreByCode(ctx context.Context, storeCode string) (db.Store, error)
	ListStores(ctx context.Context, params db.ListStoresParams) ([]db.Store, error)
	GetStoreByID(ctx context.Context, id int32) (db.Store, error)
	UpdateStore(ctx context.Context, params db.UpdateStoreParams) (db.Store, error)
	GetStoreManagerByID(ctx context.Context, id int32) (db.StoreManager, error)
	DeleteStoreManager(ctx context.Context, id int32) error
	LogActivity(ctx context.Context, params db.LogActivityParams) (db.ActivityLog, error)
}


