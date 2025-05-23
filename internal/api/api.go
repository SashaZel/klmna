package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func StartAPI(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, headersMiddleware, middleware.WithValue("DB", db))

	r.NotFound(http.HandlerFunc(notFound))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/project", func(r chi.Router) {

		r.Post("/", errorWrapper(createProject))

		r.Route("/{projectID}", func(r chi.Router) {
			r.Use(projectCtx)
			r.Get("/", errorWrapper(getProject))
			r.Put("/update", errorWrapper(updateProject))
			r.Delete("/delete", errorWrapper(deleteProject))
			r.Get("/random_task", errorWrapper(getRandomTask))

			r.Route("/pool", func(r chi.Router) {
				r.Post("/", errorWrapper(createPool))
				r.Route("/{poolID}", func(r chi.Router) {
					r.Use(poolCtx)
					r.Get("/", errorWrapper(getPool))
				})
			})
		})
	})

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", errorWrapper(getProjects))
	})

	r.Route("/task", func(r chi.Router) {
		r.Put("/solution", errorWrapper(saveTaskSolution))
	})

	return r
}

type ErrorResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: set appropriate headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
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
