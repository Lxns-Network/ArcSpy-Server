package main

import (
	"ArcSpy-Server/arcapi"
	"ArcSpy-Server/database"
	mw "ArcSpy-Server/middleware"
	"github.com/gookit/config/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func setupRoutes(serverMux *http.ServeMux) {
	serverMux.HandleFunc("/player/sync", mw.ApplyMiddleware(playerSyncHandler,
		mw.RateLimitMiddleware(mw.RateLimitOptions{
			RateLimit:  1,
			BurstLimit: 1,
		})))
	serverMux.HandleFunc("/player/data", mw.ApplyMiddleware(playerDataHandler,
		mw.AuthMiddleware()))
	serverMux.HandleFunc("/player/scores", mw.ApplyMiddleware(playerScoreHandler,
		mw.AuthMiddleware()))
	serverMux.HandleFunc("/webapi/", mw.ApplyMiddleware(arcapi.WebAPIHandler,
		mw.AuthMiddleware()))
	serverMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mw.RespondWithJSON(w, 404, "api not found", nil)
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
