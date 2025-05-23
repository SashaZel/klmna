package models

import (
	"context"
	"database/sql"
	"errors"
	uuid "github.com/google/uuid"
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
	Solution   string    `json:"solution"`
	ProjectID  uuid.UUID `json:"project_id"`
	PoolID     uuid.UUID `json:"pool_id"`
}

type TaskInput struct {
	ID    uuid.UUID `json:"id"`
	Input string    `json:"input"`
}

func GetTaskInput(db *sql.DB, projectID string) (*TaskInput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	taskInput := &TaskInput{}
	selectSqlStatement := `
		SELECT id, input 
		FROM tasks 
		WHERE tasks.project_id = $1 AND 
		      tasks.solution is null AND 
			  ((tasks.assigned_at IS null) OR ($2 > (tasks.assigned_at + INTERVAL '30 minutes')))
	`
	err := db.QueryRowContext(ctx, selectSqlStatement, projectID, time.Now()).Scan(&taskInput.ID, &taskInput.Input)
	if err != nil {
		return nil, err
	}
	return taskInput, nil
}

func UpdateAssignDate(db *sql.DB, taskID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	updateSqlStatement := "UPDATE tasks SET assigned_at = $1 WHERE id = $2"
	_, err := db.ExecContext(ctx, updateSqlStatement, time.Now(), taskID)
	return err
}

func UpdateTaskSolution(db *sql.DB, taskID string, taskSolution string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	updateSqlStatement := "UPDATE tasks SET solution = $1 WHERE id = $2 AND solution IS null"
	result, err := db.ExecContext(ctx, updateSqlStatement, taskSolution, taskID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("no such row")
	}
	return err
}
