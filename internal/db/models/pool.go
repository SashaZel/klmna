package models

import (
	"context"
	uuid "github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type NewPool struct {
	Name      string   `json:"name"`
	Input     string   `json:"input"`
	Output    string   `json:"output"`
	ProjectId string   `json:"project_id"`
	Project   *Project `bun:"rel:belongs-to" json:"project"`
}

type Pool struct {
	*NewPool
	ID        uuid.UUID `bun:",pk" json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatePool(db *bun.DB, req *NewPool) (*Pool, error) {
	ctx := context.Background()
	id := uuid.New()
	createdPool := &Pool{
		NewPool:   req,
		ID:        id,
		CreatedAt: time.Now(),
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
