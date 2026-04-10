package tasks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, vinculacionID uuid.UUID, req CreateTaskRequest, semanaInicio time.Time, tardio bool) (*TaskResponse, error)
	FindByVinculacion(ctx context.Context, vinculacionID uuid.UUID) ([]TaskResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*TaskResponse, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Nuevas consultas faltantes
	FindByVinculacionSemana(ctx context.Context, vinculacionID uuid.UUID, semana time.Time) ([]TaskResponse, error)
	ResumenByUsuario(ctx context.Context, usuarioID uuid.UUID, profesorID uuid.UUID) ([]TaskResponse, error)
	EstadoSemanal(ctx context.Context, espacioID uuid.UUID, semana time.Time) ([]map[string]interface{}, error)
	FindByUsuario(ctx context.Context, usuarioID uuid.UUID) ([]TaskResponse, error)
	AddAttachment(ctx context.Context, tareaID uuid.UUID, nombreArchivo string, ruta string, mime string, size int64) error
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, vinculacionID uuid.UUID, req CreateTaskRequest, semanaInicio time.Time, tardio bool) (*TaskResponse, error) {
	var t TaskResponse
	err := r.db.QueryRow(ctx,
		`INSERT INTO tareas (vinculacion_id, titulo, descripcion, estado, semana_inicio, horas_invertidas, observaciones, reporte_tardio)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING id, vinculacion_id, titulo, descripcion, estado, semana_inicio,
                   horas_invertidas, COALESCE(observaciones, ''), reporte_tardio, creado_en, actualizado_en`,
		vinculacionID, req.Titulo, req.Descripcion, req.Estado,
		semanaInicio, req.HorasInvertidas, req.Observaciones, tardio,
	).Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
		&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
		&t.CreadoEn, &t.ActualizadoEn)
	return &t, err
}

func (r *repository) FindByVinculacion(ctx context.Context, vinculacionID uuid.UUID) ([]TaskResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, vinculacion_id, titulo, descripcion, estado, semana_inicio,
                horas_invertidas, COALESCE(observaciones, ''), reporte_tardio, creado_en, actualizado_en
         FROM tareas WHERE vinculacion_id = $1 ORDER BY semana_inicio DESC`, vinculacionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []TaskResponse
	for rows.Next() {
		var t TaskResponse
		rows.Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
			&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
			&t.CreadoEn, &t.ActualizadoEn)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*TaskResponse, error) {
	var t TaskResponse
	err := r.db.QueryRow(ctx,
		`SELECT id, vinculacion_id, titulo, descripcion, estado, semana_inicio,
                horas_invertidas, COALESCE(observaciones, ''), reporte_tardio, creado_en, actualizado_en
         FROM tareas WHERE id = $1`, id,
	).Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
		&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
		&t.CreadoEn, &t.ActualizadoEn)
	return &t, err
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error) {
	var t TaskResponse
	err := r.db.QueryRow(ctx,
		`UPDATE tareas SET
            titulo = COALESCE(NULLIF($1,''), titulo),
            descripcion = COALESCE(NULLIF($2,''), descripcion),
            estado = COALESCE(NULLIF($3,''), estado),
            horas_invertidas = CASE WHEN $4 > 0 THEN $4 ELSE horas_invertidas END,
            observaciones = COALESCE(NULLIF($5,''), observaciones),
            actualizado_en = NOW()
         WHERE id = $6
         RETURNING id, vinculacion_id, titulo, descripcion, estado, semana_inicio,
                   horas_invertidas, COALESCE(observaciones, ''), reporte_tardio, creado_en, actualizado_en`,
		req.Titulo, req.Descripcion, req.Estado, req.HorasInvertidas, req.Observaciones, id,
	).Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
		&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
		&t.CreadoEn, &t.ActualizadoEn)
	return &t, err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tareas WHERE id = $1`, id)
	return err
}

func (r *repository) FindByVinculacionSemana(ctx context.Context, vinculacionID uuid.UUID, semana time.Time) ([]TaskResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, vinculacion_id, titulo, descripcion, estado, semana_inicio,
                horas_invertidas, COALESCE(observaciones, ''), reporte_tardio, creado_en, actualizado_en
         FROM tareas WHERE vinculacion_id = $1 AND semana_inicio = $2 ORDER BY creado_en DESC`, vinculacionID, semana)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []TaskResponse
	for rows.Next() {
		var t TaskResponse
		rows.Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
			&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
			&t.CreadoEn, &t.ActualizadoEn)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *repository) ResumenByUsuario(ctx context.Context, usuarioID uuid.UUID, profesorID uuid.UUID) ([]TaskResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT t.id, t.vinculacion_id, t.titulo, t.descripcion, t.estado, t.semana_inicio,
                t.horas_invertidas, COALESCE(t.observaciones, ''), t.reporte_tardio, t.creado_en, t.actualizado_en
         FROM tareas t
         JOIN vinculaciones v ON t.vinculacion_id = v.id
         WHERE v.usuario_id = $1 AND v.profesor_id = $2 ORDER BY t.semana_inicio DESC, t.creado_en DESC`, usuarioID, profesorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []TaskResponse
	for rows.Next() {
		var t TaskResponse
		rows.Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
			&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
			&t.CreadoEn, &t.ActualizadoEn)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *repository) EstadoSemanal(ctx context.Context, espacioID uuid.UUID, semana time.Time) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT v.usuario_id, u.nombre, COALESCE(SUM(t.horas_invertidas), 0) as horas, 
                COUNT(t.id) as cantidad_tareas
         FROM vinculaciones v
         JOIN usuarios u ON v.usuario_id = u.id
         LEFT JOIN tareas t ON t.vinculacion_id = v.id AND t.semana_inicio = $2
         WHERE v.espacio_id = $1
         GROUP BY v.usuario_id, u.nombre`, espacioID, semana)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var uid uuid.UUID
		var nombre string
		var horas, tareas int
		if err := rows.Scan(&uid, &nombre, &horas, &tareas); err != nil {
			return nil, err
		}
		list = append(list, map[string]interface{}{
			"usuario_id":      uid,
			"nombre":          nombre,
			"horas_reportadas": horas,
			"tareas_reportadas": tareas,
		})
	}
	return list, nil
}

func (r *repository) FindByUsuario(ctx context.Context, usuarioID uuid.UUID) ([]TaskResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT t.id, t.vinculacion_id, t.titulo, t.descripcion, t.estado, t.semana_inicio,
                t.horas_invertidas, COALESCE(t.observaciones, ''), t.reporte_tardio, t.creado_en, t.actualizado_en
         FROM tareas t
         JOIN vinculaciones v ON t.vinculacion_id = v.id
         WHERE v.usuario_id = $1 ORDER BY t.semana_inicio DESC, t.creado_en DESC`, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []TaskResponse
	for rows.Next() {
		var t TaskResponse
		rows.Scan(&t.ID, &t.VinculacionID, &t.Titulo, &t.Descripcion, &t.Estado,
			&t.SemanaInicio, &t.HorasInvertidas, &t.Observaciones, &t.ReporteTardio,
			&t.CreadoEn, &t.ActualizadoEn)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *repository) AddAttachment(ctx context.Context, tareaID uuid.UUID, nombreArchivo string, ruta string, mime string, size int64) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO adjuntos (tarea_id, nombre_archivo, ruta_almacenamiento, tipo_mime, tamanio_bytes)
         VALUES ($1, $2, $3, $4, $5)`,
		tareaID, nombreArchivo, ruta, mime, size)
	return err
}
