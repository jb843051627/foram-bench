package regression

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/service"
	"github.com/jb843051627/foram-bench/internal/store"
)

func bug02Observation(id string) service.ObservationInput {
	return service.ObservationInput{ID: id, SectionID: "SEC-02", Observer: "microscope", Taxon: "Ammonia", Count: 2, Preservation: "good", Confidence: .9, ObservedAt: time.Now()}
}

func TestBug02_DuplicateBatchReturnsConflict(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	lab := service.NewLab(db)
	defer lab.Close()
	defer db.Close()
	if err := db.SaveSection(model.ThinSection{ID: "SEC-02", BatchID: "B-02", Label: "A", ThicknessUM: 30, Stain: "rose", Status: model.SectionStained, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("duplicate import panicked: %v", recovered)
		}
	}()
	_, err = lab.RecordObservationSet(context.Background(), "SEC-02", []service.ObservationInput{bug02Observation("O-02"), bug02Observation("O-02")})
	if !errors.Is(err, model.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestBug02_CancelledBatchReturnsCancellation(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	lab := service.NewLab(db)
	defer lab.Close()
	defer db.Close()
	if err := db.SaveSection(model.ThinSection{ID: "SEC-02-C", BatchID: "B-02", Label: "A", ThicknessUM: 30, Stain: "rose", Status: model.SectionStained, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("cancelled import panicked: %v", recovered)
		}
	}()
	_, err = lab.RecordObservationSet(ctx, "SEC-02-C", []service.ObservationInput{bug02Observation("O-02-C")})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestBug02_EmptyBatchReturnsInvalidInput(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	lab := service.NewLab(db)
	defer lab.Close()
	defer db.Close()
	if err := db.SaveSection(model.ThinSection{ID: "SEC-02-E", BatchID: "B-02", Label: "A", ThicknessUM: 30, Stain: "rose", Status: model.SectionStained, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("empty import panicked: %v", recovered)
		}
	}()
	_, err = lab.RecordObservationSet(context.Background(), "SEC-02-E", nil)
	if !errors.Is(err, model.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
