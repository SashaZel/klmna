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

type ProjectsResponse struct {
	Ok    bool              `json:"ok"`
	Error string            `json:"error"`
	Data  []*models.Project `json:"data"`
}

func getProjects(w http.ResponseWriter, r *http.Request) error {
	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	projects, err := models.GetProjects(pgdb)
	if err != nil {
		return err
	}

	res := &ProjectsResponse{
		Ok:    true,
		Error: "",
		Data:  projects,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

type CreateProjectRequest struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Template string `json:"template"`
}

type ProjectResponse struct {
	Ok    bool            `json:"ok"`
	Error string          `json:"error"`
	Data  *models.Project `json:"data"`
}

func projectCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "projectID")
		pgdb, ok := r.Context().Value("DB").(*sql.DB)
		if !ok {
			log.Printf("fail to connect DB")
			http.Error(w, http.StatusText(500), 500)
			return
		}

		project, err := models.GetProject(pgdb, projectID)
		if err != nil {
			log.Printf("fail to get project %w \n", err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "project", project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createProject(w http.ResponseWriter, r *http.Request) error {
	req := &CreateProjectRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	project, err := models.CreateProject(pgdb, &models.NewProject{
		Name:     req.Name,
		Template: req.Template,
	})
	if err != nil {
		return err
	}

	res := &ProjectResponse{
		Ok:    true,
		Error: "",
		Data:  project,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func getProject(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	project, ok := ctx.Value("project").(*models.Project)
	if !ok {
		return errors.New("fail to get project from co")
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to get project from co")
	}
	projectWithPools, err := models.GetProjectWithPools(pgdb, project.ID.String())
	if err != nil {
		return err
	}

	res := &ProjectResponse{
		Ok:    true,
		Error: "",
		Data:  projectWithPools,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}
