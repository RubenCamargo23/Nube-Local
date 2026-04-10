package reports

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "path/filepath"
    "time"

    "github.com/google/uuid"
    "github.com/hibiken/asynq"
    "github.com/jackc/pgx/v5/pgxpool"
)

const TaskTypeGenerateReport = "report:generate"

type Worker struct {
    repo         Repository
    ollamaClient *OllamaClient
    pdfPath      string
    db           *pgxpool.Pool
}

func NewWorker(repo Repository, ollama *OllamaClient, pdfPath string, db *pgxpool.Pool) *Worker {
    return &Worker{repo: repo, ollamaClient: ollama, pdfPath: pdfPath, db: db}
}

func (w *Worker) ProcessTask(ctx context.Context, t *asynq.Task) error {
    var payload ReportJobPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }

    reporteID, _ := uuid.Parse(payload.ReporteID)
    vinculacionID, _ := uuid.Parse(payload.VinculacionID)
    semana, _ := time.Parse("2006-01-02", payload.SemanaInicio)

    // Actualizar estado a PROCESSING
    log.Printf("Procesando reporte: %s", reporteID)
    w.db.Exec(ctx, `UPDATE reportes_pdf SET estado_generacion='PROCESSING' WHERE id=$1`, reporteID)

    // Obtener nombre del monitor
    var nombreMonitor string
    err := w.db.QueryRow(ctx, 
        `SELECT u.nombre FROM vinculaciones v 
         JOIN usuarios u ON u.id = v.usuario_id 
         WHERE v.id = $1`, vinculacionID).Scan(&nombreMonitor)
    if err != nil {
        nombreMonitor = vinculacionID.String() // Fallback al ID si algo falla
    }

    // Consultar tareas de la semana
    log.Printf("Consultando tareas para vinculacion %s, semana %s", vinculacionID, payload.SemanaInicio)
    rows, err := w.db.Query(ctx,
        `SELECT titulo, descripcion, observaciones, horas_invertidas
         FROM tareas
         WHERE vinculacion_id = $1 AND semana_inicio = $2`,
        vinculacionID, semana)
    if err != nil {
        w.repo.UpdateError(ctx, reporteID)
        return err
    }
    defer rows.Close()

    type tarea struct {
        Titulo, Descripcion, Observaciones string
        Horas                              int
    }
    var tareas []tarea
    for rows.Next() {
        var t tarea
        rows.Scan(&t.Titulo, &t.Descripcion, &t.Observaciones, &t.Horas)
        tareas = append(tareas, t)
    }

    if len(tareas) == 0 {
        w.repo.UpdateError(ctx, reporteID)
        return fmt.Errorf("no hay tareas para la semana %s", payload.SemanaInicio)
    }

    // Construir prompt (RN-43: solo tareas, descripciones, observaciones, horas)
    prompt := fmt.Sprintf(
        "Genera un resumen profesional en español del trabajo semanal realizado durante la semana del %s. "+
            "El resumen debe ser claro, objetivo y destacar los logros. No evalúes el desempeño. "+
            "Actividades realizadas:\n\n", payload.SemanaInicio)

    for i, t := range tareas {
        prompt += fmt.Sprintf("%d. %s\n   Descripción: %s\n   Horas: %d\n   Observaciones: %s\n\n",
            i+1, t.Titulo, t.Descripcion, t.Horas, t.Observaciones)
    }

    // Llamar a Ollama
    log.Printf("Enviando prompt a Ollama para reporte %s:\n%s", reporteID, prompt)
    resumen, err := w.ollamaClient.Generate(ctx, prompt)
    if err != nil {
        log.Printf("Error Ollama para reporte %s: %v", reporteID, err)
        w.repo.UpdateError(ctx, reporteID)
        return err
    }
    log.Printf("Ollama respondió para reporte %s:\n%s", reporteID, resumen)

    // Generar PDF
    log.Printf("Generando PDF para reporte %s en %s", reporteID, w.pdfPath)
    // Nombre único usando reporteID para evitar sobreescritura
    outputPath := filepath.Join(w.pdfPath, fmt.Sprintf("reporte_%s.pdf", reporteID))
    if err := GeneratePDF(outputPath, nombreMonitor, payload.SemanaInicio, resumen); err != nil {
        w.repo.UpdateError(ctx, reporteID)
        return err
    }

    // Persistir resultado
    log.Printf("Persitiendo resultado para reporte %s", reporteID)
    return w.repo.UpdateDone(ctx, reporteID, outputPath, prompt, time.Now())
}
