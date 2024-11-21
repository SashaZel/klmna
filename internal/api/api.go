package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// "github.com/uptrace/bun"
	// "github.com/jackc/pgx/v5/pgxpool"
	"database/sql"
)

func StartAPI(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.WithValue("DB", db))

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

type ErrorResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func errorWrapper[F ~func(w http.ResponseWriter, r *http.Request) error](wrapped F) func(w http.ResponseWriter, r *http.Request) {
	internal := func(w http.ResponseWriter, r *http.Request) {
		err := wrapped(w, r)
		if err != nil {
			log.Printf("internal error %w \n", err)
			res := &ErrorResponse{
				Ok:    false,
				Error: "an error occurred",
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
	res := &ErrorResponse{
		Ok:    false,
		Error: "not found",
	}
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error sending 404 response %v \n", err)
	}
}
