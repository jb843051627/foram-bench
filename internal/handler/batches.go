package handler

import (
	"net/http"
	"strconv"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createBatch(w http.ResponseWriter, r *http.Request) {
	var input service.BatchInput
	if !h.readBody(w, r, &input) {
		return
	}
	batch, err := h.lab.CreateBatch(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, batch)
}

func (h *Handler) listBatches(w http.ResponseWriter, r *http.Request) {
	batches, err := h.lab.ListBatches(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, batches)
}

func (h *Handler) getBatch(w http.ResponseWriter, r *http.Request) {
	batch, err := h.lab.GetBatch(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}

func expectedRevision(r *http.Request) int {
	value, _ := strconv.Atoi(r.URL.Query().Get("revision"))
	return value
}

func (h *Handler) startBatch(w http.ResponseWriter, r *http.Request) {
	batch, err := h.lab.StartBatch(r.Context(), pathID(r), expectedRevision(r))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}

func (h *Handler) completeBatch(w http.ResponseWriter, r *http.Request) {
	batch, err := h.lab.CompleteBatch(r.Context(), pathID(r), expectedRevision(r))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}

func (h *Handler) blockBatch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Reason string `json:"reason"`
	}
	if !h.readBody(w, r, &body) {
		return
	}
	batch, err := h.lab.BlockBatch(r.Context(), pathID(r), expectedRevision(r), body.Reason)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, batch)
}

func (h *Handler) batchSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.lab.BatchSummary(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
