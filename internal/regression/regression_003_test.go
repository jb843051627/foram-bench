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

func setupBug03(t *testing.T) (*store.Store, *service.Lab) {
	t.Helper()
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	lab := service.NewLab(db)
	now := time.Now()
	if err := db.SaveSample(model.Sample{ID: "S-03", SiteCode: "SITE-03", DepthMM: 1, Material: "mudstone", CollectionDate: now, Location: "north", TimeZone: "UTC", Status: model.SampleRegistered, CreatedAt: now, UpdatedAt: now}); err != nil {
		lab.Close()
		db.Close()
		t.Fatal(err)
	}
	return db, lab
}

func TestBug03_ConcurrentBatchCreationHasOneWinner(t *testing.T) {
	db, lab := setupBug03(t)
	defer lab.Close()
	defer db.Close()
	var wg sync.WaitGroup
	results := make(chan error, 2)
	start := make(chan struct{})
	for _, operator := range []string{"alice", "bob"} {
		wg.Add(1)
		go func(operator string) {
			defer wg.Done()
			<-start
			_, err := lab.CreateBatch(context.Background(), service.BatchInput{ID: "B-03", SampleID: "S-03", Protocol: "acid", Operator: operator})
			results <- err
		}(operator)
	}
	close(start)
	wg.Wait()
	close(results)
	success, conflicts := 0, 0
	for err := range results {
		if err == nil {
			success++
		} else if errors.Is(err, model.ErrConflict) {
			conflicts++
		}
	}
	if success != 1 || conflicts != 1 {
		t.Fatalf("success=%d conflicts=%d", success, conflicts)
	}
}

func TestBug03_DuplicateDoesNotReplaceBatch(t *testing.T) {
	db, lab := setupBug03(t)
	defer lab.Close()
	defer db.Close()
	if _, err := lab.CreateBatch(context.Background(), service.BatchInput{ID: "B-03", SampleID: "S-03", Protocol: "acid", Operator: "first"}); err != nil {
		t.Fatal(err)
	}
	if _, err := lab.CreateBatch(context.Background(), service.BatchInput{ID: "B-03", SampleID: "S-03", Protocol: "acid", Operator: "second"}); !errors.Is(err, model.ErrConflict) {
		t.Fatalf("err=%v", err)
	}
	batch, err := db.GetBatch("B-03")
	if err != nil {
		t.Fatal(err)
	}
	if batch.Operator != "first" {
		t.Fatalf("operator=%s", batch.Operator)
	}
}

func TestBug03_BatchCountRemainsOne(t *testing.T) {
	db, lab := setupBug03(t)
	defer lab.Close()
	defer db.Close()
	if _, err := lab.CreateBatch(context.Background(), service.BatchInput{ID: "B-03", SampleID: "S-03", Protocol: "acid", Operator: "first"}); err != nil {
		t.Fatal(err)
	}
	items, err := db.ListBatchesBySample("S-03")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("batches=%d", len(items))
	}
}
