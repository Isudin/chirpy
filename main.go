package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Isudin/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		println("error loading environmental variables: " + err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		println("error connecting to database: " + err.Error())
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        database.New(db),
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerGetHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("GET /api/chirps", cfg.handlerListChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.handlerGetChirpById)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
