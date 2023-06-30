package main

import (
	"ArcSpy-Server/database"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func playerSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithJSON(w, 405, fmt.Sprintf("method %s is not allowed", r.Method), nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data database.SyncUserData
	err := decoder.Decode(&data)
	if err != nil {
		respondWithJSON(w, 400, "invalid post data", nil)
		return
	}

	if data.Player == nil {
		respondWithJSON(w, 400, "invalid post data", nil)
		return
	}

	err = database.InsertPlayerData(&data)
	if err != nil {
		respondWithJSON(w, 400, "invalid post data", nil)
		return
	}

	respondWithJSON(w, 200, "ok", nil)
}

func playerDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithJSON(w, 405, fmt.Sprintf("method %s is not allowed", r.Method), nil)
		return
	}

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		respondWithJSON(w, 400, "invalid query data", nil)
		return
	}

	playerData, err := database.SelectPlayerData(userId)
	if err != nil {
		respondWithJSON(w, 404, "player not found", nil)
		return
	}

	respondWithJSON(w, 200, "ok", playerData)
}

func playerScoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithJSON(w, 405, fmt.Sprintf("method %s is not allowed", r.Method), nil)
		return
	}

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		respondWithJSON(w, 400, "invalid query data", nil)
		return
	}

	playerScores, err := database.SelectPlayerScores(userId)
	if err != nil {
		log.Error(err)
		respondWithJSON(w, 404, "player not found", nil)
		return
	}

	respondWithJSON(w, 200, "ok", playerScores)
}

func respondWithJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	if code != 200 {
		w.WriteHeader(code)
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
	_, _ = w.Write(response)
}
