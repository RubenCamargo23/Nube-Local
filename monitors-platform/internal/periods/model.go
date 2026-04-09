package periods

import (
    "time"
    "github.com/google/uuid"
)

type CreatePeriodRequest struct {
    Codigo      string `json:"codigo" binding:"required"`
    FechaInicio string `json:"fecha_inicio" binding:"required"`
    FechaFin    string `json:"fecha_fin" binding:"required"`
}

type PeriodResponse struct {
    ID          uuid.UUID `json:"id"`
    Codigo      string    `json:"codigo"`
    FechaInicio time.Time `json:"fecha_inicio"`
    FechaFin    time.Time `json:"fecha_fin"`
    Estado      string    `json:"estado"`
    CreadoEn   time.Time `json:"creado_en"`
}
