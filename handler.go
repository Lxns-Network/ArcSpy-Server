package main

import (
	"ArcSpy-Server/database"
	mw "ArcSpy-Server/middleware"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func playerSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		mw.RespondWithJSON(w, 405, fmt.Sprintf("method %s is not allowed", r.Method), nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data database.SyncUserData
	err := decoder.Decode(&data)
	if err != nil {
		mw.RespondWithJSON(w, 400, "invalid post data", nil)
		return
	}

	if data.Player == nil {
		mw.RespondWithJSON(w, 400, "invalid post data", nil)
		return
	}

	err = database.InsertPlayerData(&data)
	if err != nil {
		mw.RespondWithJSON(w, 400, "invalid post data", nil)
		return
	}

	mw.RespondWithJSON(w, 200, "ok", nil)
}

func playerDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		mw.RespondWithJSON(w, 405, fmt.Sprintf("method %s is not allowed", r.Method), nil)
		return
	}

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		mw.RespondWithJSON(w, 400, "invalid query data", nil)
		return
	}

	playerData, err := database.SelectPlayerData(userId)
	if err != nil {
		mw.RespondWithJSON(w, 404, "player not found", nil)
		return
	}

	mw.RespondWithJSON(w, 200, "ok", playerData)
}

func playerScoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		mw.RespondWithJSON(w, 405, fmt.Sprintf("method %s is not allowed", r.Method), nil)
		return
	}

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		mw.RespondWithJSON(w, 400, "invalid query data", nil)
		return
	}

	playerScores, err := database.SelectPlayerScores(userId)
	if err != nil {
		log.Error(err)
		mw.RespondWithJSON(w, 404, "player not found", nil)
		return
	}

	mw.RespondWithJSON(w, 200, "ok", playerScores)
}
