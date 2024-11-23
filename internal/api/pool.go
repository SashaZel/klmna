package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"database/sql"
	"github.com/go-chi/chi/v5"
	"klmna/internal/db/models"
)

type CreatePoolRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Input       string `json:"input"`
	ProjectId   string `json:"project_id"`
}

type PoolResponse struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Data  *models.Pool `json:"data"`
}

func poolCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		poolID := chi.URLParam(r, "poolID")
		pgdb, ok := r.Context().Value("DB").(*sql.DB)
		if !ok {
			log.Printf("fail to connect DB")
			http.Error(w, http.StatusText(500), 500)
			return
		}

		pool, err := models.GetPool(pgdb, poolID)
		if err != nil {
			log.Printf("fail to get project %w \n", err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "pool", pool)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createPool(w http.ResponseWriter, r *http.Request) error {

	req := &CreatePoolRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	if req.Name == "" || req.ProjectId == "" {
		return errors.New("Bad request")
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	pool, err := models.CreatePool(pgdb, &models.NewPool{
		Name:        req.Name,
		Description: req.Description,
		ProjectId:   req.ProjectId,
	})
	if err != nil {
		return err
	}

	tasks := []string{}
	//   "input": "[\"{'foo':'bar'}\",\"bro1\"]",
	err = json.Unmarshal([]byte(req.Input), &tasks)
	if err != nil {
		return err
	}

	err = models.CreateTask(pgdb, tasks, pool.ID)
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

func getPool(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	pool, ok := ctx.Value("pool").(*models.Pool)
	if !ok {
		return errors.New("fail to get pool from context")
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to get DB")
	}
	poolWithTasks, err := models.GetPoolWithTasks(pgdb, pool.ID.String())
	if err != nil {
		return err
	}

	res := &PoolResponse{
		Ok:    true,
		Error: "",
		Data:  poolWithTasks,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}
