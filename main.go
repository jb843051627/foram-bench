package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jb843051627/foram-bench/internal/handler"
	"github.com/jb843051627/foram-bench/internal/service"
	"github.com/jb843051627/foram-bench/internal/store"
)

func main() {
	path := os.Getenv("FORAM_BENCH_DB")
	if path == "" {
		path = "data/foram-bench.db"
	}
	repository, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	app := service.NewLab(repository)
	defer app.Close()
	addr := os.Getenv("FORAM_BENCH_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("foram-bench listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.New(app)))
}
