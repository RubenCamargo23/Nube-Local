package reports

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    Create(ctx context.Context, vinculacionID uuid.UUID, semana time.Time) (*ReportResponse, error)
    FindByEspacioProfesor(ctx context.Context, profesorID uuid.UUID) ([]ReportResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*ReportResponse, error)
    UpdateDone(ctx context.Context, id uuid.UUID, rutaPDF, prompt string, generadoEn time.Time) error
    UpdateError(ctx context.Context, id uuid.UUID) error
    FindVinculacionesByEspacio(ctx context.Context, espacioID uuid.UUID) ([]uuid.UUID, error)
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, vinculacionID uuid.UUID, semana time.Time) (*ReportResponse, error) {
    var rp ReportResponse
    err := r.db.QueryRow(ctx,
        `INSERT INTO reportes_pdf (vinculacion_id, semana_inicio, estado_generacion)
         VALUES ($1, $2, 'PENDING')
         RETURNING id, vinculacion_id, semana_inicio, estado_generacion, creado_en`,
        vinculacionID, semana,
    ).Scan(&rp.ID, &rp.VinculacionID, &rp.SemanaInicio, &rp.EstadoGeneracion, &rp.CreadoEn)
    return &rp, err
}

func (r *repository) FindByEspacioProfesor(ctx context.Context, profesorID uuid.UUID) ([]ReportResponse, error) {
    rows, err := r.db.Query(ctx,
        `SELECT rp.id, rp.vinculacion_id, rp.semana_inicio, rp.ruta_pdf,
                rp.estado_generacion, rp.generado_en, rp.creado_en
         FROM reportes_pdf rp
         JOIN vinculaciones v ON v.id = rp.vinculacion_id
         WHERE v.profesor_id = $1
         ORDER BY rp.creado_en DESC`, profesorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []ReportResponse
    for rows.Next() {
        var rp ReportResponse
        rows.Scan(&rp.ID, &rp.VinculacionID, &rp.SemanaInicio, &rp.RutaPDF,
            &rp.EstadoGeneracion, &rp.GeneradoEn, &rp.CreadoEn)
        list = append(list, rp)
    }
    return list, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*ReportResponse, error) {
    var rp ReportResponse
    err := r.db.QueryRow(ctx,
        `SELECT id, vinculacion_id, semana_inicio, ruta_pdf,
                estado_generacion, generado_en, creado_en
         FROM reportes_pdf WHERE id = $1`, id,
    ).Scan(&rp.ID, &rp.VinculacionID, &rp.SemanaInicio, &rp.RutaPDF,
        &rp.EstadoGeneracion, &rp.GeneradoEn, &rp.CreadoEn)
    return &rp, err
}

func (r *repository) UpdateDone(ctx context.Context, id uuid.UUID, ruta, prompt string, gen time.Time) error {
    _, err := r.db.Exec(ctx,
        `UPDATE reportes_pdf SET ruta_pdf=$1, prompt_usado=$2, estado_generacion='DONE', generado_en=$3 WHERE id=$4`,
        ruta, prompt, gen, id)
    return err
}

func (r *repository) UpdateError(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.Exec(ctx,
        `UPDATE reportes_pdf SET estado_generacion='ERROR' WHERE id=$1`, id)
    return err
}

func (r *repository) FindVinculacionesByEspacio(ctx context.Context, espacioID uuid.UUID) ([]uuid.UUID, error) {
    rows, err := r.db.Query(ctx,
        `SELECT id FROM vinculaciones WHERE espacio_id = $1 AND activa = true`, espacioID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var ids []uuid.UUID
    for rows.Next() {
        var id uuid.UUID
        rows.Scan(&id)
        ids = append(ids, id)
    }
    return ids, nil
}
