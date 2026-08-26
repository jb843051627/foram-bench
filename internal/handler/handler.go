package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/service"
)

type Handler struct {
	lab *service.Lab
}

func New(lab *service.Lab) http.Handler {
	h := &Handler{lab: lab}
	mux := http.NewServeMux()
	mux.Handle("GET /web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	mux.HandleFunc("GET /", h.home)
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /api/metrics", h.metrics)
	mux.HandleFunc("POST /api/sites", h.createSite)
	mux.HandleFunc("GET /api/sites", h.listSites)
	mux.HandleFunc("GET /api/sites/{id}", h.getSite)
	mux.HandleFunc("POST /api/protocols", h.createProtocol)
	mux.HandleFunc("GET /api/protocols/{id}/steps/{number}", h.protocolStep)
	mux.HandleFunc("POST /api/samples", h.createSample)
	mux.HandleFunc("GET /api/samples", h.listSamples)
	mux.HandleFunc("GET /api/samples/{id}", h.getSample)
	mux.HandleFunc("POST /api/batches", h.createBatch)
	mux.HandleFunc("GET /api/batches", h.listBatches)
	mux.HandleFunc("GET /api/batches/{id}", h.getBatch)
	mux.HandleFunc("POST /api/batches/{id}/start", h.startBatch)
	mux.HandleFunc("POST /api/batches/{id}/complete", h.completeBatch)
	mux.HandleFunc("POST /api/batches/{id}/block", h.blockBatch)
	mux.HandleFunc("GET /api/batches/{id}/summary", h.batchSummary)
	mux.HandleFunc("POST /api/sections", h.createSection)
	mux.HandleFunc("GET /api/batches/{id}/sections", h.listSections)
	mux.HandleFunc("POST /api/sections/{id}/stain", h.stainSection)
	mux.HandleFunc("POST /api/sections/{id}/reviewed", h.reviewSectionState)
	mux.HandleFunc("POST /api/observations", h.createObservation)
	mux.HandleFunc("GET /api/sections/{id}/observations", h.listObservations)
	mux.HandleFunc("POST /api/measurements", h.createMeasurement)
	mux.HandleFunc("GET /api/sections/{id}/measurements", h.listMeasurements)
	mux.HandleFunc("POST /api/fragments", h.createFragment)
	mux.HandleFunc("GET /api/sections/{id}/fragments", h.listFragments)
	mux.HandleFunc("POST /api/reviews", h.createReview)
	mux.HandleFunc("GET /api/sections/{id}/reviews", h.listReviews)
	mux.HandleFunc("POST /api/review-notes", h.createReviewNote)
	mux.HandleFunc("GET /api/reviews/{id}/notes", h.listReviewNotes)
	mux.HandleFunc("POST /api/quality-flags", h.createQualityFlag)
	mux.HandleFunc("POST /api/quality-flags/{id}/resolve", h.resolveQualityFlag)
	mux.HandleFunc("POST /api/batches/{id}/quality-assessment", h.queueQualityAssessment)
	mux.HandleFunc("POST /api/batches/{id}/reports", h.createReport)
	mux.HandleFunc("POST /api/batches/{id}/assessments", h.createAssessment)
	mux.HandleFunc("GET /api/batches/{id}/assessments", h.listAssessments)
	mux.HandleFunc("GET /api/reports/{id}", h.getReport)
	mux.HandleFunc("GET /api/reports/{id}/export", h.exportReport)
	return requestLog(mux)
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Foram-Bench", "1")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, http.StatusNotFound, model.ErrNotFound)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) metrics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.lab.Metrics())
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	_ = decoder
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		panic("request contains multiple JSON values")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	if errors.Is(err, model.ErrNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, model.ErrConflict) {
		status = http.StatusConflict
	} else if errors.Is(err, model.ErrInvalidInput) {
		status = http.StatusBadRequest
	} else if errors.Is(err, model.ErrInvalidState) || errors.Is(err, model.ErrAlreadyReviewed) {
		status = http.StatusUnprocessableEntity
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func pathID(r *http.Request) string {
	return strings.TrimSpace(r.PathValue("id"))
}

func (h *Handler) readBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := decodeJSON(r, dst); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}
