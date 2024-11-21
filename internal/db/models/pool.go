package models

import (
	"context"
	"database/sql"
	"time"
	
	uuid "github.com/google/uuid"
)

type NewPool struct {
	Name      string `json:"name"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	ProjectId string `json:"project_id"`
}

type Pool struct {
	ID uuid.UUID `json:"id"`
	Name      string `json:"name"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	ProjectId string         `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatePool(db *sql.DB, req *NewPool) (*Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	project := &Project{}
	checkIfProjectExistSqlStatement := "SELECT id FROM projects WHERE id = $1"
	err := db.QueryRowContext(ctx, checkIfProjectExistSqlStatement, req.ProjectId).Scan(&project.ID)
	if err != nil {
		return nil, err
	}

	id := uuid.New()
	insertSqlStatement := "INSERT INTO pools (id, name, created_at, input, output, project_id) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = db.ExecContext(ctx, insertSqlStatement, id, req.Name, time.Now(), req.Input, req.Output, req.ProjectId)
	if err != nil {
		return nil, err
	}

	pool := &Pool{}
	selectSqlStatement := "SELECT id, name, created_at, input, output, project_id FROM pools WHERE id = $1"
	err = db.QueryRowContext(ctx, selectSqlStatement, id).Scan(&pool.ID, &pool.Name, &pool.CreatedAt, &pool.Input, &pool.Output, &pool.ProjectId)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
