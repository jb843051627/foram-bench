package handler

import "net/http"

func (h *Handler) createReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.lab.GenerateReport(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, report)
}

func (h *Handler) getReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.lab.GetReport(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) exportReport(w http.ResponseWriter, r *http.Request) {
	body, err := h.lab.ExportReport(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(body))
}
