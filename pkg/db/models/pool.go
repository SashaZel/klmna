package models

import (
	"context"
	"github.com/uptrace/bun"
	uuid "github.com/google/uuid"
)

type NewPool struct {
	Name      string   `json:"name"`
	Input     string   `json:"input"`
	Output    string   `json:"output"`
	ProjectId int64    `json:"userId"`
	Project   *Project `bun:"rel:has-one" json:"project"`
}

type Pool struct {
	*NewPool
	ID        uuid.UUID    `bun:",pk" json:"id"`
}

func CreatePool(db *bun.DB, req *NewPool) (*Pool, error) {
	ctx := context.Background()
	id := uuid.New()
	createdPool := &Pool{
		NewPool: req,
        ID: id,
	}
	_, err := db.NewInsert().Model(createdPool).Exec(ctx)
	if err != nil {
		return nil, err
	}

	pool := &Pool{}

	err = db.NewSelect().
		Model(pool).
		Where("pool.id = ?", id).
		Scan(ctx)

	return pool, err
}
