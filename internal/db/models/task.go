package models

import (
	"context"
	"database/sql"
	"fmt"
	uuid "github.com/google/uuid"
	"strings"
	"time"
)

type NewTask struct {
	Input string `json:"input"`
}

type Task struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	AssignedAt time.Time `json:"assigned_at"`
	Input      string    `json:"input"`
	Output     string    `json:"output"`
	PoolID     uuid.UUID `json:"pool_id"`
}

func CreateTask(db *sql.DB, tasks []string, poolId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pool := &Pool{}
	checkIfPoolExistSqlStatement := "SELECT id FROM pools WHERE id = $1"
	err := db.QueryRowContext(ctx, checkIfPoolExistSqlStatement, poolId).Scan(&pool.ID)
	if err != nil {
		return err
	}

	valueStrings := make([]string, 0, len(tasks))
	valueArgs := make([]interface{}, 0, len(tasks)*3)
	i := 0
	for _, task := range tasks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, task)
		valueArgs = append(valueArgs, poolId)
		i++
	}
	sqlQuery := fmt.Sprintf("INSERT INTO tasks (created_at, input, pool_id) VALUES %s", strings.Join(valueStrings, ","))

	var params []interface{}

	for i := 0; i < len(valueArgs); i++ {
		var param sql.NamedArg
		param.Name = fmt.Sprintf("p%v", i+1)
		param.Value = valueArgs[i]
		params = append(params, param)
	}

	_, err = db.Exec(sqlQuery, params...)
	return err
}
