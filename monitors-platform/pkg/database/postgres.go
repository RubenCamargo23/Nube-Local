package database

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/org/monitors-platform/pkg/config"
)

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("parsear config DB: %w", err)
    }

    poolCfg.MaxConns = 25
    poolCfg.MinConns = 5
    poolCfg.MaxConnLifetime = 1 * time.Hour
    poolCfg.MaxConnIdleTime = 30 * time.Minute

    pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
    if err != nil {
        return nil, fmt.Errorf("crear pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("ping DB: %w", err)
    }

    return pool, nil
}
