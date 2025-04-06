package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"database/sql"
	"klmna/internal/db/models"
)

type TaskResponse struct {
	Ok       bool              `json:"ok"`
	Error    string            `json:"error"`
	Input    *models.TaskInput `json:"input"`
	Template string            `json:"template"`
}

type TaskSolutionRequest struct {
	TaskID   string `json:"task_id"`
	Solution string `json:"solution"`
}

func saveTaskSolution(w http.ResponseWriter, r *http.Request) error {
	req := &TaskSolutionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	if req.TaskID == "" || req.Solution == "" {
		return errors.New("Bad request")
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	err = models.UpdateTaskSolution(pgdb, req.TaskID, req.Solution)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
