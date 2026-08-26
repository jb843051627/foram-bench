package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createObservation(w http.ResponseWriter, r *http.Request) {
	var input service.ObservationInput
	if !h.readBody(w, r, &input) {
		return
	}
	observation, err := h.lab.RecordObservation(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, observation)
}

func (h *Handler) listObservations(w http.ResponseWriter, r *http.Request) {
	observations, err := h.lab.ListObservations(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, observations)
}
