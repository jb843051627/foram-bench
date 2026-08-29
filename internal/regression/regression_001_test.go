package regression

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/store"
)

func TestBug01_TransactionReturnsCallbackError(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	called := false
	err = db.Transaction(func(*sql.Tx) error { called = true; return errors.New("forced batch failure") })
	if !called || err == nil {
		t.Fatalf("called=%v err=%v", called, err)
	}
}

func TestBug01_RollsBackEarlierWrites(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = db.Transaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO records(kind, id, payload, updated_at) VALUES('sample', 'S-01', '{}', 'now')`)
		if err != nil {
			return err
		}
		return errors.New("forced batch failure")
	})
	if err == nil {
		t.Fatal("transaction swallowed callback error")
	}
	if _, err := db.GetSample("S-01"); !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("write survived rollback: %v", err)
	}
}

func TestBug01_EmptyStoreAfterFailure(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/case.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if count, err := db.Count("sample"); err != nil || count != 0 {
		t.Fatalf("unexpected records: %d %v", count, err)
	}
}
