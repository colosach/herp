// Copyright 2025 The HERP Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package business

import (
	"context"
	"database/sql"
	db "herp/db/sqlc"
)

type Business struct {
	db      *sql.DB
	queries Querier
}

func NewBusiness(queries Querier, db *sql.DB) *Business {
	return &Business{
		queries: queries,
		db:      db,
	}
}

// CreateBusiness creates a new business with a default branch.
func (c *Business) CreateBusinessWithBranch(ctx context.Context, params db.CreateBusinessParams) (db.Business, db.Branch, error) {
	q, ok := c.queries.(*db.Queries)
	if !ok {
		return db.Business{}, db.Branch{}, nil
	}

	// Start a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return db.Business{}, db.Branch{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	txQueries := q.WithTx(tx)

	// Create the business
	business, err := txQueries.CreateBusiness(ctx, params)
	if err != nil {
		return db.Business{}, db.Branch{}, err
	}

	// Create a default branch for the business
	branchParams := db.CreateBranchParams{
		Name:       "Main Branch",
		BusinessID: business.ID,
	}
	branch, err := txQueries.CreateBranch(ctx, branchParams)
	if err != nil {
		return db.Business{}, db.Branch{}, err
	}

	return business, branch, nil
}

func (c *Business) CreateBusiness(ctx context.Context, params db.CreateBusinessParams) (db.Business, error) {
	return c.queries.CreateBusiness(ctx, params)
}

// GetBusiness retrieves a business by its ID.
func (c *Business) GetBusiness(ctx context.Context, id int32) (db.Business, error) {
	return c.queries.GetBusiness(ctx, id)
}

// UpdateBusiness updates an existing business.
func (c *Business) UpdateBusiness(ctx context.Context, params db.UpdateBusinessParams) (db.Business, error) {
	return c.queries.UpdateBusiness(ctx, params)
}

// DeleteBusiness deletes a business by its ID.
func (c *Business) DeleteBusiness(ctx context.Context, id int32) (db.Business, error) {
	return c.queries.DeleteBusiness(ctx, id)
}

// ListBusinesses lists all businesses.
func (c *Business) ListBusinesses(ctx context.Context) ([]db.Business, error) {
	return c.queries.ListBusinesses(ctx)
}

// --------Branch Methods-------- //

// CreateBranch creates a new branch.
func (c *Business) CreateBranch(ctx context.Context, params db.CreateBranchParams) (db.Branch, error) {
	return c.queries.CreateBranch(ctx, params)
}

// GetBranch retrieves a branch by its ID.
func (c *Business) GetBranch(ctx context.Context, id int32) (db.Branch, error) {
	return c.queries.GetBranch(ctx, id)
}

// UpdateBranch updates an existing branch.
func (c *Business) UpdateBranch(ctx context.Context, params db.UpdateBranchParams) (db.Branch, error) {
	return c.queries.UpdateBranch(ctx, params)
}

// DeleteBranch deletes a branch by its ID.
func (c *Business) DeleteBranch(ctx context.Context, id int32) (db.Branch, error) {
	return c.queries.DeleteBranch(ctx, id)
}

// ListBranch lists branches
func (c *Business) ListBranches(ctx context.Context) ([]db.Branch, error) {
	return c.queries.ListBranches(ctx)
}

// --------Sotore and StoreManager Methods-------- //

// CreateStore creates a new store.
func (c *Business) CreateStore(ctx context.Context, params db.CreateStoreParams) (db.Store, error) {
	return c.queries.CreateStore(ctx, params)
}

// CreateStoreManager creates a new store manager.
func (c *Business) CreateStoreManager(ctx context.Context, params db.CreateStoreManagerParams) (db.StoreManager, error) {
	return c.queries.CreateStoreManager(ctx, params)
}

// DeleteStore deletes a store by its ID.
func (c *Business) DeleteStore(ctx context.Context, id int32) error {
	return c.queries.DeleteStore(ctx, id)
}

// GetStoreByCode retrieves a store by its code.
func (c *Business) GetStoreByCode(ctx context.Context, storeCode string) (db.Store, error) {
	return c.queries.GetStoreByCode(ctx, storeCode)
}

// ListStores lists all stores.
func (c *Business) ListStores(ctx context.Context, params db.ListStoresParams) ([]db.Store, error) {
	return c.queries.ListStores(ctx, params)
}

// GetStoreByID retrieves a store by its ID.
func (c *Business) GetStoreByID(ctx context.Context, id int32) (db.Store, error) {
	return c.queries.GetStoreByID(ctx, id)
}

// UpdateStore updates an existing store.
func (c *Business) UpdateStore(ctx context.Context, params db.UpdateStoreParams) (db.Store, error) {
	return c.queries.UpdateStore(ctx, params)
}

// GetStoreManagerByID retrieves a store manager by its ID.
func (c *Business) GetStoreManagerByID(ctx context.Context, id int32) (db.StoreManager, error) {
	return c.queries.GetStoreManagerByID(ctx, id)
}

// DeleteStoreManager deletes a store manager by its ID.
func (c *Business) DeleteStoreManager(ctx context.Context, id int32) error {
	return c.queries.DeleteStoreManager(ctx, id)
}

func (c *Business) LogActivity(ctx context.Context, params db.LogActivityParams) (db.ActivityLog, error) {
	return c.queries.LogActivity(ctx, params)
}


// --------Activity Log Methods-------- //
func (c *Business) GetAcitivityLogs(ctx context.Context, limit int32) ([]db.ActivityLog, error) {
	return c.queries.GetActivityLogs(ctx, limit)
}