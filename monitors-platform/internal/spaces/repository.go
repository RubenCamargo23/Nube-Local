package spaces

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    Create(ctx context.Context, req CreateSpaceRequest, profesorID uuid.UUID) (*SpaceResponse, error)
    FindByProfesor(ctx context.Context, profesorID uuid.UUID) ([]SpaceResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*SpaceResponse, error)
    Close(ctx context.Context, id uuid.UUID) error
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, req CreateSpaceRequest, profesorID uuid.UUID) (*SpaceResponse, error) {
    var s SpaceResponse
    err := r.db.QueryRow(ctx,
        `INSERT INTO espacios (tipo, nombre, periodo_id, fecha_inicio, fecha_fin, profesor_id, observaciones)
         VALUES ($1, $2, $3, $4::date, $5::date, $6, $7)
         RETURNING id, tipo, nombre, periodo_id, fecha_inicio, fecha_fin, profesor_id, observaciones, estado, creado_en`,
        req.Tipo, req.Nombre, req.PeriodoID, req.FechaInicio, req.FechaFin, profesorID, req.Observaciones,
    ).Scan(&s.ID, &s.Tipo, &s.Nombre, &s.PeriodoID, &s.FechaInicio, &s.FechaFin,
        &s.ProfesorID, &s.Observaciones, &s.Estado, &s.CreadoEn)
    return &s, err
}

func (r *repository) FindByProfesor(ctx context.Context, profesorID uuid.UUID) ([]SpaceResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, tipo, nombre, periodo_id, fecha_inicio, fecha_fin,
                profesor_id, COALESCE(observaciones, ''), estado, creado_en
         FROM espacios WHERE profesor_id = $1 ORDER BY creado_en DESC`, profesorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var spaces []SpaceResponse
	for rows.Next() {
		var s SpaceResponse
		err := rows.Scan(&s.ID, &s.Tipo, &s.Nombre, &s.PeriodoID, &s.FechaInicio, &s.FechaFin,
			&s.ProfesorID, &s.Observaciones, &s.Estado, &s.CreadoEn)
		if err != nil {
			return nil, err
		}
		spaces = append(spaces, s)
	}
	return spaces, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*SpaceResponse, error) {
	var s SpaceResponse
	err := r.db.QueryRow(ctx,
		`SELECT id, tipo, nombre, periodo_id, fecha_inicio, fecha_fin,
                profesor_id, COALESCE(observaciones, ''), estado, creado_en
         FROM espacios WHERE id = $1`, id,
	).Scan(&s.ID, &s.Tipo, &s.Nombre, &s.PeriodoID, &s.FechaInicio, &s.FechaFin,
		&s.ProfesorID, &s.Observaciones, &s.Estado, &s.CreadoEn)
	return &s, err
}

func (r *repository) Close(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.Exec(ctx,
        `UPDATE espacios SET estado = 'CLOSED' WHERE id = $1`, id)
    return err
}
