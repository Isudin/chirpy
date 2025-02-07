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
	}

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerGetHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
