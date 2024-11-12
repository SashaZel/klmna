package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"

	"klmna/pkg/db/models"
)

func StartAPI(pgdb *pg.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	r.Route("/project", func(r chi.Router) {
		r.Get("/", getProject)
		r.Post("/", createProject)
	})

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", getProjects)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hit the main"))
	})

	return r
}

type CreateProjectRequest struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type ProjectResponse struct {
	Ok    bool            `json:"ok"`
	Error string          `json:"error"`
	Data  *models.Project `json:"data"`
}

func createProject(w http.ResponseWriter, r *http.Request) {
	req := &CreateProjectRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &ProjectResponse{
			Ok:    false,
			Error: "bad request",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error in encode res to malformed req %v \n")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &ProjectsResponse{
			Ok:    false,
			Error: "fail to get DB context",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending DB connection error response %v \n", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	project, err := models.CreateProject(pgdb, &models.Project{
		Name:     req.Name,
		Template: req.Template,
	})
	if err != nil {
		res := &ProjectResponse{
			Ok:    false,
			Error: "fail to create project at DB request",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error in encode fail resp of creating project %v \n", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := &ProjectResponse{
		Ok:    true,
		Error: "",
		Data:  project,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error in encode positive response %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	// projectId, ok := r.URL.Query()["projectId"]
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	projectId := queryParams["projectId"]
	if err != nil {
		res := &ProjectResponse{
			Ok:    false,
			Error: "malformed request",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending error resp at malformed req %v \n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &ProjectsResponse{
			Ok:    false,
			Error: "fail to get DB context",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending DB connection error response %v \n", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	project, err := models.GetProject(pgdb, projectId[0])
	if err != nil {
		res := &ProjectResponse{
			Ok:    false,
			Error: "error in getting project from db",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error in sending error form req to db %v \n", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := &ProjectResponse{
		Ok:    true,
		Error: "",
		Data:  project,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error in encode positive response %v \n", err)
	}
	w.WriteHeader(http.StatusOK)
	return
}

type ProjectsResponse struct {
	Ok    bool              `json:"ok"`
	Error string            `json:"error"`
	Data  []*models.Project `json:"data"`
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &ProjectsResponse{
			Ok:    false,
			Error: "fail to get DB context",
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending DB connection error response %v \n", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projects, err := models.GetProjects(pgdb)
	if err != nil {
		res := &ProjectsResponse{
			Ok:    false,
			Error: err.Error(),
			Data:  nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending db model error response %v \n", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := &ProjectsResponse{
		Ok:    true,
		Error: "",
		Data:  projects,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encode response %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
