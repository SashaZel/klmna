package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

type NewPool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectId   string `json:"project_id"`
}

type Pool struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectId   string    `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
	Tasks       []*Task   `json:"tasks"`
	// TODO: add progress statistic by tasks
}

func CreatePool(db *sql.DB, req *NewPool, tasks []string) (*Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	id := uuid.New()
	createdAt := time.Now()
	insertSqlStatement := "INSERT INTO pools (id, name, description, created_at, project_id) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.ExecContext(ctx, insertSqlStatement, id, req.Name, req.Description, createdAt, req.ProjectId)
	if err != nil {
		return nil, err
	}

	pool := &Pool{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   createdAt,
		ProjectId:   req.ProjectId,
	}

	valueStrings := make([]string, 0, len(tasks))
	valueArgs := make([]interface{}, 0, len(tasks)*4)
	i := 0
	for _, task := range tasks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, task)
		valueArgs = append(valueArgs, pool.ProjectId)
		valueArgs = append(valueArgs, pool.ID)
		i++
	}
	sqlQuery := fmt.Sprintf("INSERT INTO tasks (created_at, input, project_id, pool_id) VALUES %s", strings.Join(valueStrings, ","))

	var params []interface{}

	for i := 0; i < len(valueArgs); i++ {
		var param sql.NamedArg
		param.Name = fmt.Sprintf("p%v", i+1)
		param.Value = valueArgs[i]
		params = append(params, param)
	}

	_, err = tx.Exec(sqlQuery, params...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func GetPool(db *sql.DB, poolID string) (*Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pool := &Pool{}
	selectSqlStatement := "SELECT id, name, description, project_id, created_at FROM pools WHERE id = $1"
	err := db.QueryRowContext(ctx, selectSqlStatement, poolID).Scan(&pool.ID, &pool.Name, &pool.Description, &pool.ProjectId, &pool.CreatedAt)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func GetPoolWithTasks(db *sql.DB, poolId string) (*Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pool := &Pool{}
	tasks := make([]*Task, 0)
	selectSqlStatement := `
		SELECT 
		pools.id, pools.name, pools.description, pools.created_at, pools.project_id,
		tasks.id, tasks.created_at, tasks.assigned_at, tasks.input, tasks.solution, tasks.project_id, tasks.pool_id
		FROM pools
		LEFT JOIN tasks ON pools.id = tasks.pool_id AND tasks.pool_id = $1
		ORDER BY tasks.created_at
    `
	rows, err := db.QueryContext(ctx, selectSqlStatement, poolId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		var taskID uuid.UUID
		var taskCreatedAt sql.NullTime
		var taskAssignedAt sql.NullTime
		var taskInput sql.NullString
		var taskSolution sql.NullString
		var taskProjectId uuid.UUID
		var taskPoolId uuid.UUID
		err := rows.Scan(
			&pool.ID,
			&pool.Name,
			&pool.Description,
			&pool.CreatedAt,
			&pool.ProjectId,
			&taskID,
			&taskCreatedAt,
			&taskAssignedAt,
			&taskInput,
			&taskSolution,
			&taskProjectId,
			&taskPoolId,
		)
		if err != nil {
			return nil, err
		}
		if taskID.String() != "" && taskInput.Valid {
			task.ID = taskID
			task.CreatedAt = taskCreatedAt.Time
			task.AssignedAt = taskAssignedAt.Time
			task.Input = taskInput.String
			task.Solution = taskSolution.String
			task.ProjectID = taskProjectId
			task.PoolID = taskPoolId
			tasks = append(tasks, task)
		}
	}
	pool.Tasks = tasks

	return pool, nil
}
