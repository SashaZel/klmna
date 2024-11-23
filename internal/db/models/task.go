package models

import (
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
	ProjectID  uuid.UUID `json:"project_id"`
	PoolID     uuid.UUID `json:"pool_id"`
}

func CreateTask(db *sql.DB, tasks []string, poolId uuid.UUID, projectID string) error {
	valueStrings := make([]string, 0, len(tasks))
	valueArgs := make([]interface{}, 0, len(tasks)*4)
	i := 0
	for _, task := range tasks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, task)
		valueArgs = append(valueArgs, projectID)
		valueArgs = append(valueArgs, poolId)
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

	_, err := db.Exec(sqlQuery, params...)
	return err
}
