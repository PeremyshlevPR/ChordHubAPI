package web

import (
	"chords_app/internal/web/handlers"

	"github.com/gorilla/mux"
)

const apiV1Prefix = "/api/v1"

func SetupRouter(userHandler *handlers.UserHandler) *mux.Router {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix(apiV1Prefix).Subrouter()

	apiRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	apiRouter.HandleFunc("/login", userHandler.Login).Methods("POST")
	apiRouter.HandleFunc("/refresh", userHandler.Refresh).Methods("POST")

	return r
}
