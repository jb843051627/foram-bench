package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createMeasurement(w http.ResponseWriter, r *http.Request) {
	var input service.MeasurementInput
	if !h.readBody(w, r, &input) {
		return
	}
	measurement, err := h.lab.RecordMeasurement(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, measurement)
}

func (h *Handler) listMeasurements(w http.ResponseWriter, r *http.Request) {
	items, err := h.lab.ListMeasurements(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) createFragment(w http.ResponseWriter, r *http.Request) {
	var input service.FragmentInput
	if !h.readBody(w, r, &input) {
		return
	}
	fragment, err := h.lab.RecordFragment(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, fragment)
}

func (h *Handler) listFragments(w http.ResponseWriter, r *http.Request) {
	items, err := h.lab.ListFragments(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
