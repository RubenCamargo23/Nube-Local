package tasks

import (
    "time"
    "github.com/google/uuid"
)

type CreateTaskRequest struct {
    Titulo          string    `json:"titulo" binding:"required"`
    Descripcion     string    `json:"descripcion" binding:"required"`
    Estado          string    `json:"estado" binding:"required,oneof=abierto en_desarrollo finalizado"`
    SemanaInicio    string    `json:"semana_inicio" binding:"required"`
    HorasInvertidas int       `json:"horas_invertidas" binding:"required,min=1"`
    Observaciones   string    `json:"observaciones"`
}

type UpdateTaskRequest struct {
    Titulo          string `json:"titulo"`
    Descripcion     string `json:"descripcion"`
    Estado          string `json:"estado" binding:"omitempty,oneof=abierto en_desarrollo finalizado"`
    HorasInvertidas int    `json:"horas_invertidas" binding:"omitempty,min=1"`
    Observaciones   string `json:"observaciones"`
}

type TaskResponse struct {
    ID              uuid.UUID `json:"id"`
    VinculacionID   uuid.UUID `json:"vinculacion_id"`
    Titulo          string    `json:"titulo"`
    Descripcion     string    `json:"descripcion"`
    Estado          string    `json:"estado"`
    SemanaInicio    time.Time `json:"semana_inicio"`
    HorasInvertidas int       `json:"horas_invertidas"`
    Observaciones   string    `json:"observaciones"`
    ReporteTardio   bool      `json:"reporte_tardio"`
    CreadoEn        time.Time `json:"creado_en"`
    ActualizadoEn   time.Time `json:"actualizado_en"`
}
