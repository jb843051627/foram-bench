package handler

import (
	"net/http"
	"strconv"

	"github.com/jb843051627/foram-bench/internal/service"
)

func (h *Handler) createSite(w http.ResponseWriter, r *http.Request) {
	var input service.SiteInput
	if !h.readBody(w, r, &input) {
		return
	}
	site, err := h.lab.RegisterSite(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, site)
}

func (h *Handler) listSites(w http.ResponseWriter, r *http.Request) {
	activeOnly, _ := strconv.ParseBool(r.URL.Query().Get("active"))
	sites, err := h.lab.ListSites(r.Context(), activeOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, sites)
}

func (h *Handler) getSite(w http.ResponseWriter, r *http.Request) {
	site, err := h.lab.GetSite(r.Context(), pathID(r))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, site)
}

func (h *Handler) createProtocol(w http.ResponseWriter, r *http.Request) {
	var input service.ProtocolInput
	if !h.readBody(w, r, &input) {
		return
	}
	protocol, err := h.lab.CreateProtocol(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, protocol)
}

func (h *Handler) protocolStep(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil {
		panic(err)
	}
	step, err := h.lab.ProtocolStep(r.Context(), r.PathValue("id"), number)
	if err != nil {
		panic(err)
	}
	writeJSON(w, http.StatusOK, step)
}
