package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/uptrace/bun"

	"klmna/pkg/db/models"
)

func StartAPI(pgdb *bun.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	r.NotFound(http.HandlerFunc(notFound))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/pool", func(r chi.Router) {
		r.Post("/", errorWrapper(createPool))
	})

	r.Route("/project", func(r chi.Router) {
		r.Get("/", errorWrapper(getProject))
		r.Post("/", errorWrapper(createProject))
	})

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", errorWrapper(getProjects))
	})

	return r
}

func errorWrapper[F ~func(w http.ResponseWriter, r *http.Request) error](wrapped F) func(w http.ResponseWriter, r *http.Request) {
	internal := func(w http.ResponseWriter, r *http.Request) {
		err := wrapped(w, r)
		if err != nil {
			log.Printf("internal error %w \n", err)
			res := &ProjectsResponse{
				Ok:    false,
				Error: "an error occurred",
				Data:  nil,
			}
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(res)
			if err != nil {
				log.Printf("error sending error response %v \n", err)
			}
		}
	}
	return internal
}

func notFound(w http.ResponseWriter, r *http.Request) {
	res := &ProjectsResponse{
		Ok:    false,
		Error: "not found",
		Data:  nil,
	}
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error sending 404 response %v \n", err)
	}
}

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

func createProject(w http.ResponseWriter, r *http.Request) error {
	req := &CreateProjectRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	pgdb, ok := r.Context().Value("DB").(*bun.DB)
	if !ok {
		return errors.New("fail to connect DB")
	}

	project, err := models.CreateProject(pgdb, &models.Project{
		ID:       req.ID,
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

	pgdb, ok := r.Context().Value("DB").(*bun.DB)
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

type ProjectsResponse struct {
	Ok    bool              `json:"ok"`
	Error string            `json:"error"`
	Data  []*models.Project `json:"data"`
}

func getProjects(w http.ResponseWriter, r *http.Request) error {
	pgdb, ok := r.Context().Value("DB").(*bun.DB)
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
