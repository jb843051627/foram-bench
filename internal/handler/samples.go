package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createSample(w http.ResponseWriter, r *http.Request) {
	var input service.SampleInput
	if !h.readBody(w, r, &input) {
		return
	}
	sample, err := h.lab.RegisterSample(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, sample)
}

func (h *Handler) listSamples(w http.ResponseWriter, r *http.Request) {
	samples, err := h.lab.ListSamples(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, samples)
}

func (h *Handler) getSample(w http.ResponseWriter, r *http.Request) {
	sample, err := h.lab.GetSample(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, sample)
}
