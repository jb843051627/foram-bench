package handler

import (
	"net/http"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createSection(w http.ResponseWriter, r *http.Request) {
	var input service.SectionInput
	if !h.readBody(w, r, &input) {
		return
	}
	section, err := h.lab.CreateSection(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, section)
}

func (h *Handler) listSections(w http.ResponseWriter, r *http.Request) {
	sections, err := h.lab.ListSections(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, sections)
}

func (h *Handler) stainSection(w http.ResponseWriter, r *http.Request) {
	section, err := h.lab.StainSection(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, section)
}

func (h *Handler) reviewSectionState(w http.ResponseWriter, r *http.Request) {
	section, err := h.lab.MarkSectionReviewed(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, section)
}
