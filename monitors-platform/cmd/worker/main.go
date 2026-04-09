package main

import (
    "log"

    "github.com/hibiken/asynq"
    "github.com/org/monitors-platform/internal/reports"
    "github.com/org/monitors-platform/pkg/config"
    "github.com/org/monitors-platform/pkg/database"
)

func main() {
    cfg := config.Load()

    db, err := database.NewPool(cfg)
    if err != nil {
        log.Fatalf("Error conectando a DB: %v", err)
    }
    defer db.Close()

    reportsRepo  := reports.NewRepository(db)
    ollamaClient := reports.NewOllamaClient(cfg.OllamaHost, cfg.OllamaModel, cfg.OllamaTimeoutSecs)
    worker       := reports.NewWorker(reportsRepo, ollamaClient, cfg.PDFPath, db)

    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: cfg.RedisAddr},
        asynq.Config{Concurrency: 5},
    )

    mux := asynq.NewServeMux()
    mux.HandleFunc(reports.TaskTypeGenerateReport, worker.ProcessTask)

    log.Println("Worker Asynq iniciado")
    if err := srv.Run(mux); err != nil {
        log.Fatalf("Error iniciando worker: %v", err)
    }
}
