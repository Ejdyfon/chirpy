package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func customHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))

}

func (cfg *apiConfig) hitsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(200)
	str := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	rw.Write([]byte(str))

}

func (cfg *apiConfig) resetHitsHandler(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	str := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	rw.Write([]byte(str))

}

func main() {
	apiCfg := apiConfig{}
	srv := http.NewServeMux()
	srv.HandleFunc("GET /api/healthz", customHandler)
	srv.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	srv.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)
	srv.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	httpsrv := http.Server{}
	httpsrv.Handler = srv
	httpsrv.Addr = ":8080"
	httpsrv.ListenAndServe()

}
