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
}

type PoolResponse struct {
	Ok       bool         `json:"ok"`
	Error    string       `json:"error"`
	Pool     *models.Pool `json:"pool"`
	Template string       `json:"template"`
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

	ctx := r.Context()
	project, ok := ctx.Value("project").(*models.Project)
	if !ok {
		return errors.New("fail to get project from context")
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	tasks := []string{}
	err = json.Unmarshal([]byte(req.Input), &tasks)
	if err != nil {
		return err
	}

	pool, err := models.CreatePool(pgdb, &models.NewPool{
		Name:        req.Name,
		Description: req.Description,
		ProjectId:   project.ID.String(),
	}, tasks)
	if err != nil {
		return err
	}

	res := &PoolResponse{
		Ok:       true,
		Error:    "",
		Pool:     pool,
		Template: project.Template,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func getPool(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	project, ok := ctx.Value("project").(*models.Project)
	if !ok {
		return errors.New("fail to get project from context")
	}
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
		Ok:       true,
		Error:    "",
		Pool:     poolWithTasks,
		Template: project.Template,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}
