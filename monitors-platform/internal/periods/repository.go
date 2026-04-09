package periods

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    Create(ctx context.Context, codigo, fechaInicio, fechaFin string) (*PeriodResponse, error)
    FindAll(ctx context.Context) ([]PeriodResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*PeriodResponse, error)
    Close(ctx context.Context, id uuid.UUID) error
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, codigo, fechaInicio, fechaFin string) (*PeriodResponse, error) {
    var p PeriodResponse
    err := r.db.QueryRow(ctx,
        `INSERT INTO periodos_academicos (codigo, fecha_inicio, fecha_fin)
         VALUES ($1, $2::date, $3::date)
         RETURNING id, codigo, fecha_inicio, fecha_fin, estado, creado_en`,
        codigo, fechaInicio, fechaFin,
    ).Scan(&p.ID, &p.Codigo, &p.FechaInicio, &p.FechaFin, &p.Estado, &p.CreadoEn)
    return &p, err
}

func (r *repository) FindAll(ctx context.Context) ([]PeriodResponse, error) {
    rows, err := r.db.Query(ctx,
        `SELECT id, codigo, fecha_inicio, fecha_fin, estado, creado_en
         FROM periodos_academicos ORDER BY creado_en DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var periods []PeriodResponse
    for rows.Next() {
        var p PeriodResponse
        rows.Scan(&p.ID, &p.Codigo, &p.FechaInicio, &p.FechaFin, &p.Estado, &p.CreadoEn)
        periods = append(periods, p)
    }
    return periods, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*PeriodResponse, error) {
    var p PeriodResponse
    err := r.db.QueryRow(ctx,
        `SELECT id, codigo, fecha_inicio, fecha_fin, estado, creado_en
         FROM periodos_academicos WHERE id = $1`, id,
    ).Scan(&p.ID, &p.Codigo, &p.FechaInicio, &p.FechaFin, &p.Estado, &p.CreadoEn)
    return &p, err
}

func (r *repository) Close(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.Exec(ctx,
        `UPDATE periodos_academicos SET estado = 'CLOSED' WHERE id = $1`, id)
    return err
}
