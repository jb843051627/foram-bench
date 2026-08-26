package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createQualityFlag(w http.ResponseWriter, r *http.Request) {
	var input service.QualityInput
	if !h.readBody(w, r, &input) {
		return
	}
	flag, err := h.lab.AddQualityFlag(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, flag)
}

func (h *Handler) resolveQualityFlag(w http.ResponseWriter, r *http.Request) {
	flag, err := h.lab.ResolveQualityFlag(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, flag)
}

func (h *Handler) queueQualityAssessment(w http.ResponseWriter, r *http.Request) {
	done := make(chan error, 1)
	if err := h.lab.QueueQualityAssessment(r.Context(), pathID(r), done); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "queued"})
}
