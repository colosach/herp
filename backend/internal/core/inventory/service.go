package inventory

import (
	"context"
	"database/sql"
	db "herp/db/sqlc"
)

type Inventory struct {
	db      *sql.DB
	queries Querier
}

func NewInventory(queries Querier, db *sql.DB) *Inventory {
	return &Inventory{
		queries: queries,
		db:      db,
	}
}

func (i *Inventory) CreateBrand(ctx context.Context, args db.CreateBrandParams) (db.Brand, error) {
	return i.queries.CreateBrand(ctx, args)
}

func (i *Inventory) LogActivity(ctx context.Context, args db.LogActivityParams) (db.ActivityLog, error) {
	return i.queries.LogActivity(ctx, args)
}

func (i *Inventory) CreateCategory(ctx context.Context, params db.CreateCategoryParams) (db.Category, error) {
	return i.queries.CreateCategory(ctx, params)
}

func (i *Inventory) GetCategory(ctx context.Context, id int32) (db.Category, error) {
	return i.queries.GetCategory(ctx, id)
}

func (i *Inventory) CreateItem(ctx context.Context, params db.CreateItemParams) (db.Item, error) {
	return i.queries.CreateItem(ctx, params)
}

func (i *Inventory) GetBrand(ctx context.Context, id int32) (db.Brand, error) {
	return i.queries.GetBrand(ctx, id)
}

func (i *Inventory) CreateVariation(ctx context.Context, params db.CreateVariationParams) (db.Variation, error) {
	return i.queries.CreateVariation(ctx, params)
}

func (i *Inventory) GetItem(ctx context.Context, id int32) (db.Item, error) {
	return i.queries.GetItem(ctx, id)
}
