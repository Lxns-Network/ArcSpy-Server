package arcapi

import (
	"ArcSpy-Server/database"
	mw "ArcSpy-Server/middleware"
	"net/http"
)

func WebAPIHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	userId := r.URL.Query().Get("user_id")
	cookie, err := database.SelectPlayerCookie(userId)
	if err != nil {
		mw.RespondWithJSON(w, 404, "player not found", nil)
		return
	}

	switch r.URL.Path {
	case "/webapi/user/me":
		data, err = GetPlayerUserMe(cookie)
	default:
		mw.RespondWithJSON(w, 404, "api not found", nil)
		return
	}

	if err != nil {
		mw.RespondWithJSON(w, 500, "internal server error", nil)
		return
	}

	if !data["success"].(bool) {
		mw.RespondWithJSON(w, 400, "player session expired", nil)
		return
	}

	mw.RespondWithJSON(w, 200, "ok", data["value"])
}
