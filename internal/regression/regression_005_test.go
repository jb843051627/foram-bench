package regression

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/service"
	"github.com/jb843051627/foram-bench/internal/store"
)

func TestBug05_ConcurrentSectionAdvanceHasOneWinner(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	lab := service.NewLab(db)
	defer lab.Close()
	defer db.Close()
	now := time.Now()
	if err := db.SaveSection(model.ThinSection{ID: "SEC-05", BatchID: "B-05", Label: "A", ThicknessUM: 30, Stain: "rose", Status: model.SectionCut, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := lab.StainSection(context.Background(), "SEC-05"); results <- err }()
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("successful transitions=%d", success)
	}
}

func TestBug05_FinalSectionStateIsStained(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	lab := service.NewLab(db)
	defer lab.Close()
	now := time.Now()
	if err := db.SaveSection(model.ThinSection{ID: "SEC-05", BatchID: "B-05", Label: "A", ThicknessUM: 30, Stain: "rose", Status: model.SectionCut, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := lab.StainSection(context.Background(), "SEC-05"); err != nil {
		t.Fatal(err)
	}
	section, err := db.GetSection("SEC-05")
	if err != nil {
		t.Fatal(err)
	}
	if section.Status != model.SectionStained {
		t.Fatalf("status=%s", section.Status)
	}
}

func TestBug05_InvalidAdvanceReturnsStateError(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	lab := service.NewLab(db)
	defer lab.Close()
	now := time.Now()
	if err := db.SaveSection(model.ThinSection{ID: "SEC-05", BatchID: "B-05", Label: "A", ThicknessUM: 30, Stain: "rose", Status: model.SectionStained, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("invalid advance panicked: %v", recovered)
		}
	}()
	_, err = lab.AdvanceSection(context.Background(), "SEC-05", model.SectionCut)
	if !errors.Is(err, model.ErrInvalidState) {
		t.Fatalf("expected invalid state, got %v", err)
	}
}
