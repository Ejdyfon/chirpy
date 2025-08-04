package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type validateChirpRequest struct {
	Body string `json:"body"`
}

type validateChirpResponseError struct {
	Error string `json:"error"`
}

type validateChirpResponse struct {
	Valid bool `json:"valid"`
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

func validateChirp(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/json; charset=utf-8")
	decoder := json.NewDecoder(req.Body)
	request := validateChirpRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		rw.WriteHeader(500)
		respBody := validateChirpResponseError{Error: "Something went wrong"}
		dat, _ := json.Marshal(respBody)
		rw.Write(dat)
		return
	}

	if len(request.Body) > 140 {
		rw.WriteHeader(400)
		respBody := validateChirpResponseError{Error: "Chirp is too long"}
		dat, _ := json.Marshal(respBody)
		rw.Write(dat)
		return
	}

	rw.WriteHeader(200)
	respBody := validateChirpResponse{Valid: true}
	dat, _ := json.Marshal(respBody)
	rw.Write(dat)
	return

}

func main() {
	apiCfg := apiConfig{}
	srv := http.NewServeMux()
	srv.HandleFunc("GET /api/healthz", customHandler)
	srv.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	srv.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)
	srv.HandleFunc("POST /api/validate_chirp", validateChirp)
	srv.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	httpsrv := http.Server{}
	httpsrv.Handler = srv
	httpsrv.Addr = ":8080"
	httpsrv.ListenAndServe()

}
