package models

import (
	"context"
	"github.com/uptrace/bun"
)

type Pool struct {
	ID        int64    `bun:",pk" json:"id"`
	Name      string   `json:"name"`
	Input     string   `json:"input"`
	Output    string   `json:"output"`
	ProjectId int64    `json:"userId"`
	Project   *Project `bun:"rel:has-one" json:"project"`
}

func CreatePool(db *bun.DB, req *Pool) (*Pool, error) {
	ctx := context.Background()
	_, err := db.NewInsert().Model(req).Exec(ctx)
	if err != nil {
		return nil, err
	}

	pool := &Pool{}

	err = db.NewSelect().
		Model(pool).
		Where("pool.id = ?", req.ID).
		Scan(ctx)

	return pool, err
}
