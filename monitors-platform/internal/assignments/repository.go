package assignments

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    SumarHorasPorRol(ctx context.Context, tx pgx.Tx, usuarioID uuid.UUID, rol string) (int, error)
    ContarMonitorias(ctx context.Context, tx pgx.Tx, usuarioID uuid.UUID) (int, error)
    Create(ctx context.Context, tx pgx.Tx, espacioID, usuarioID, profesorID uuid.UUID, rol string, horas int) (*AssignmentResponse, error)
    FindByEspacio(ctx context.Context, espacioID uuid.UUID) ([]AssignmentResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*AssignmentResponse, error)
    Update(ctx context.Context, id uuid.UUID, horas int) error
    BeginTx(ctx context.Context) (pgx.Tx, error)
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return &repository{db: db} }

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
    return r.db.Begin(ctx)
}

func (r *repository) SumarHorasPorRol(ctx context.Context, tx pgx.Tx, usuarioID uuid.UUID, rol string) (int, error) {
    var total int
    err := tx.QueryRow(ctx,
        `SELECT COALESCE(SUM(horas_semanales), 0)
         FROM vinculaciones
         WHERE usuario_id = $1 AND rol = $2 AND activa = true`,
        usuarioID, rol,
    ).Scan(&total)
    return total, err
}

func (r *repository) ContarMonitorias(ctx context.Context, tx pgx.Tx, usuarioID uuid.UUID) (int, error) {
    var total int
    err := tx.QueryRow(ctx,
        `SELECT COUNT(*) FROM vinculaciones
         WHERE usuario_id = $1 AND rol = 'MONITOR' AND activa = true`,
        usuarioID,
    ).Scan(&total)
    return total, err
}

func (r *repository) Create(ctx context.Context, tx pgx.Tx, espacioID, usuarioID, profesorID uuid.UUID, rol string, horas int) (*AssignmentResponse, error) {
    var a AssignmentResponse
    err := tx.QueryRow(ctx,
        `INSERT INTO vinculaciones (usuario_id, espacio_id, rol, horas_semanales, profesor_id)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, usuario_id, espacio_id, rol, horas_semanales, profesor_id, activa, creado_en`,
        usuarioID, espacioID, rol, horas, profesorID,
    ).Scan(&a.ID, &a.UsuarioID, &a.EspacioID, &a.Rol, &a.HorasSemanales,
        &a.ProfesorID, &a.Activa, &a.CreadoEn)
    return &a, err
}

func (r *repository) FindByEspacio(ctx context.Context, espacioID uuid.UUID) ([]AssignmentResponse, error) {
    rows, err := r.db.Query(ctx,
        `SELECT id, usuario_id, espacio_id, rol, horas_semanales, profesor_id, activa, creado_en
         FROM vinculaciones WHERE espacio_id = $1 ORDER BY creado_en DESC`, espacioID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []AssignmentResponse
    for rows.Next() {
        var a AssignmentResponse
        rows.Scan(&a.ID, &a.UsuarioID, &a.EspacioID, &a.Rol, &a.HorasSemanales,
            &a.ProfesorID, &a.Activa, &a.CreadoEn)
        list = append(list, a)
    }
    return list, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*AssignmentResponse, error) {
    var a AssignmentResponse
    err := r.db.QueryRow(ctx,
        `SELECT id, usuario_id, espacio_id, rol, horas_semanales, profesor_id, activa, creado_en
         FROM vinculaciones WHERE id = $1`, id,
    ).Scan(&a.ID, &a.UsuarioID, &a.EspacioID, &a.Rol, &a.HorasSemanales,
        &a.ProfesorID, &a.Activa, &a.CreadoEn)
    return &a, err
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, horas int) error {
    _, err := r.db.Exec(ctx,
        `UPDATE vinculaciones SET horas_semanales = $1 WHERE id = $2`, horas, id)
    return err
}
