package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/ejdyfon/chirpy/internal/auth"
	"github.com/ejdyfon/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	prm := database.CreatChirpParams{Body: cleaned, UserID: userID}

	user, err := cfg.db.CreatChirp(r.Context(), prm)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp(user))
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	s := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	if s != "" {
		authorId, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Cannot parse param", err)
			return
		}
		chirps, err = cfg.db.GetAllChirpsByAuthor(r.Context(), authorId)
	} else {
		var err error
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
			return
		}
	}

	if sortOrder != "" {
		if sortOrder == "asc" {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
		} else if sortOrder == "desc" {
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
		}
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}

	res := []Chirp{}
	for _, v := range chirps {
		res = append(res, Chirp(v))
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	pathParam := r.PathValue("chirpID")
	parmId, err := uuid.Parse(pathParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert param", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), parmId)

	if chirp.ID == uuid.Nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while fetching", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp(chirp))
}
