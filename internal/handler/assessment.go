package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createAssessment(w http.ResponseWriter, r *http.Request) {
	assessment, err := h.lab.AssessBatch(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, assessment)
}

func (h *Handler) listAssessments(w http.ResponseWriter, r *http.Request) {
	items, err := h.lab.ListAssessments(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) createReviewNote(w http.ResponseWriter, r *http.Request) {
	var input service.ReviewNoteInput
	if !h.readBody(w, r, &input) {
		return
	}
	note, err := h.lab.AddReviewNote(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, note)
}

func (h *Handler) listReviewNotes(w http.ResponseWriter, r *http.Request) {
	items, err := h.lab.ListReviewNotes(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
