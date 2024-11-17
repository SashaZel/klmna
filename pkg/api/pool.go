package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/uptrace/bun"

	"klmna/pkg/db/models"
)

type CreatePoolRequest struct {
	Name      string `json:"name"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	ProjectId int64  `json:"project_id"`
}

type PoolResponse struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Data  *models.Pool `json:"data"`
}

func createPool(w http.ResponseWriter, r *http.Request) error {
	req := &CreatePoolRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	pgdb, ok := r.Context().Value("DB").(*bun.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}
	project, err := models.GetProject(pgdb, strconv.Itoa(int(req.ProjectId)))
	if err != nil {
		return err
	}

	pool, err := models.CreatePool(pgdb, &models.NewPool{
		Name:      req.Name,
		Input:     req.Input,
		Output:    req.Output,
		ProjectId: req.ProjectId,
		Project:   project,
	})
	if err != nil {
		return err
	}

	res := &PoolResponse{
		Ok:    true,
		Error: "",
		Data:  pool,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}
