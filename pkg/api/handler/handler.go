package handler

import (
	"fmt"
	"main/pkg/constants"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "main/docs"
	"main/pkg/repository"
	"net/http"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", constants.Host)),
	))

	r.HandleFunc("/requests", h.requests).Methods("GET")
	r.HandleFunc("/requests/{uid}", h.request).Methods("GET")
	r.HandleFunc("/repeat/{uid}", h.repeat).Methods("GET")
	r.HandleFunc("/scan/{uid}", h.scan).Methods("GET")

	r.Use(
		h.panicRecoveryMiddleware,
	)

	return r
}
