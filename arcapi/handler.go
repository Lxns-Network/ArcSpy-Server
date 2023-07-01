package arcapi

import (
	"ArcSpy-Server/database"
	mw "ArcSpy-Server/middleware"
	"net/http"
	"strconv"
)

func WebAPIHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		mw.RespondWithJSON(w, 400, "invalid query data", nil)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		mw.RespondWithJSON(w, 400, "invalid query data", nil)
		return
	}

	cookie, err := database.SelectPlayerCookie(userId)
	if err != nil {
		mw.RespondWithJSON(w, 404, "player not found", nil)
		return
	}

	var data *WebAPIResponse

	switch r.URL.Path {
	case "/webapi/user/me":
		data, err = GetPlayerUserMe(cookie)
		if err != nil {
			mw.RespondWithJSON(w, 500, "internal server error", nil)
			return
		}
	case "/webapi/score/song/me":
		songId := r.URL.Query().Get("song_id")
		difficultyStr := r.URL.Query().Get("difficulty")
		if songId == "" || difficultyStr == "" {
			mw.RespondWithJSON(w, 400, "invalid query data", nil)
			return
		}

		difficulty, err := strconv.Atoi(difficultyStr)
		if err != nil {
			mw.RespondWithJSON(w, 400, "invalid query data", nil)
			return
		}

		if difficulty < 0 || difficulty > 3 {
			mw.RespondWithJSON(w, 400, "invalid query data", nil)
			return
		}

		data, err = GetPlayerBestScore(cookie, userId, songId, difficulty)
		if err != nil {
			mw.RespondWithJSON(w, 500, "internal server error", nil)
			return
		}

		if !data.Success {
			mw.RespondWithJSON(w, 404, "score not found", nil)
			return
		}
	default:
		mw.RespondWithJSON(w, 404, "api not found", nil)
		return
	}

	if !data.Success {
		mw.RespondWithJSON(w, 400, "player session expired", nil)
		return
	}

	mw.RespondWithJSON(w, 200, "ok", data.Value)
}
