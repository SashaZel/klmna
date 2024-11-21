package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	// "github.com/uptrace/bun"

	"klmna/internal/db/models"
	// "github.com/jackc/pgx/v5/pgxpool"
	"database/sql"
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
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	projectId := queryParams["projectId"]
	if err != nil {
		return err
	}

	pgdb, ok := r.Context().Value("DB").(*sql.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	project, err := models.GetProject(pgdb, projectId[0])
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
