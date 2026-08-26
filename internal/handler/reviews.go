package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	var input service.ReviewInput
	if !h.readBody(w, r, &input) {
		return
	}
	review, err := h.lab.ReviewSection(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, review)
}

func (h *Handler) listReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.lab.ListReviews(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, reviews)
}
