package routes

import (
	"net/http"
	"procon_web_service/src/common/middleware"
	"procon_web_service/src/judge/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/judge", handlers.JudgeHandler).Methods(http.MethodPost)

	return router
}
