package reports

import (
    "time"
    "github.com/google/uuid"
)

type GenerateReportRequest struct {
    EspacioID    uuid.UUID `json:"espacio_id" binding:"required"`
    SemanaInicio string    `json:"semana_inicio" binding:"required"`
}

type ReportResponse struct {
    ID               uuid.UUID  `json:"id"`
    VinculacionID    uuid.UUID  `json:"vinculacion_id"`
    SemanaInicio     time.Time  `json:"semana_inicio"`
    RutaPDF          string     `json:"ruta_pdf,omitempty"`
    EstadoGeneracion string     `json:"estado_generacion"`
    GeneradoEn       *time.Time `json:"generado_en,omitempty"`
    CreadoEn         time.Time  `json:"creado_en"`
}

// Payload del job Asynq
type ReportJobPayload struct {
    ReporteID     string `json:"reporte_id"`
    VinculacionID string `json:"vinculacion_id"`
    SemanaInicio  string `json:"semana_inicio"`
}
