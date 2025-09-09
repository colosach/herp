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
package core

import (
	"context"
	db "herp/db/sqlc"
)

type Core struct {
	queries Querier
}

func NewCore(queries Querier) *Core {
	return &Core{
		queries: queries,
	}
}

// CreateBusiness creates a new business.
func (c *Core) CreateBusiness(ctx context.Context, params db.CreateBusinessParams) (db.Business, error) {
	return c.queries.CreateBusiness(ctx, params)
}

// GetBusiness retrieves a business by its ID.
func (c *Core) GetBusiness(ctx context.Context, id int32) (db.Business, error) {
	return c.queries.GetBusiness(ctx, id)
}

// UpdateBusiness updates an existing business.
func (c *Core) UpdateBusiness(ctx context.Context, params db.UpdateBusinessParams) (db.Business, error) {
	return c.queries.UpdateBusiness(ctx, params)
}

// DeleteBusiness deletes a business by its ID.
func (c *Core) DeleteBusiness(ctx context.Context, id int32) error {
	return c.queries.DeleteBusiness(ctx, id)
}

// ListBusinesses lists all businesses.
func (c *Core) ListBusinesses(ctx context.Context) ([]db.Business, error) {
	return c.queries.ListBusinesses(ctx)
}
