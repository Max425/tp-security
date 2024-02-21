package handler

import (
	"github.com/gorilla/mux"
	"main/pkg/model/dto"
	"net/http"
)

// @Summary get 1 request
// @Tags request
// @Accept  json
// @Produce  json
// @Param uid path string true "request UID"
// @Success 200 {object} core.Request
// @Failure 500 {object} string
// @Router /requests/{uid} [get]
func (h *Handler) request(w http.ResponseWriter, r *http.Request) {
	UID, has := mux.Vars(r)["uid"]
	if !has {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	orders, err := h.repo.GetRequestByID(r.Context(), UID)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, orders)
}

// @Summary get all requests
// @Tags request
// @Accept  json
// @Produce  json
// @Success 200 {object} []core.Request
// @Failure 500 {object} string
// @Router /requests [get]
func (h *Handler) requests(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.GetAllRequests(r.Context())
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, orders)
}

// @Summary repeat request
// @Tags request
// @Accept  json
// @Produce  json
// @Param uid path string true "request UID"
// @Success 200 {object} core.Request
// @Failure 500 {object} string
// @Router /repeat/{uid} [get]
func (h *Handler) repeat(w http.ResponseWriter, r *http.Request) {
	UID, has := mux.Vars(r)["uid"]
	if !has {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	orders, err := h.repo.GetRequestByID(r.Context(), UID)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, orders)
}

// @Summary scan request
// @Tags request
// @Accept  json
// @Produce  json
// @Param uid path string true "request UID"
// @Success 200 {object} bool
// @Failure 500 {object} string
// @Router /scan/{uid} [get]
func (h *Handler) scan(w http.ResponseWriter, r *http.Request) {
}
