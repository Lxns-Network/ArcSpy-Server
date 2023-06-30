package main

import (
	"ArcSpy-Server/arcapi"
	"ArcSpy-Server/database"
	"github.com/gookit/config/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken != config.SubDataMap("api").Str("token") {
			respondWithJSON(w, 401, "unauthorized", nil)
			return
		}

		next(w, r)
	}
}

func setupRoutes(serverMux *http.ServeMux) {
	serverMux.HandleFunc("/player/sync", playerSyncHandler)
	serverMux.HandleFunc("/player/data", authMiddleware(playerDataHandler))
	serverMux.HandleFunc("/player/scores", authMiddleware(playerScoreHandler))
	serverMux.HandleFunc("/test", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")
		cookie, err := database.SelectPlayerCookie(userId)
		if err != nil {
			log.Errorf("Failed to get player data: %s", err.Error())
			return
		}
		data, err := arcapi.GetPlayerUserMe(cookie)
		if err != nil {
			log.Errorf("Failed to get player data: %s", err.Error())
			return
		}
		if !data["success"].(bool) {
			respondWithJSON(w, 400, "player session expired", nil)
			return
		}
		respondWithJSON(w, 200, "ok", data["value"])
	}))
	serverMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, 404, "endpoint not found", nil)
	})
}

func main() {
	initConfig()
	database.Create()

	serverPort := config.SubDataMap("server").Str("port")
	serverMux := http.NewServeMux()
	setupRoutes(serverMux)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: serverMux,
	}

	go func(server *http.Server) {
		log.Info("Running ArcSpy Server on http://localhost:" + serverPort)
		if err := server.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}(server)

	err := watchConfig()
	if err != nil {
		log.Panic("Failed to listen config file: ", err)
	}
}
