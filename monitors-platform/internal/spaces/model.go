package spaces

import (
    "time"
    "github.com/google/uuid"
)

type CreateSpaceRequest struct {
    Tipo         string    `json:"tipo" binding:"required,oneof=COURSE PROJECT"`
    Nombre       string    `json:"nombre" binding:"required"`
    PeriodoID    uuid.UUID `json:"periodo_id" binding:"required"`
    FechaInicio  string    `json:"fecha_inicio" binding:"required"`
    FechaFin     string    `json:"fecha_fin" binding:"required"`
    Observaciones string   `json:"observaciones"`
}

type SpaceResponse struct {
    ID            uuid.UUID `json:"id"`
    Tipo          string    `json:"tipo"`
    Nombre        string    `json:"nombre"`
    PeriodoID     uuid.UUID `json:"periodo_id"`
    FechaInicio   time.Time `json:"fecha_inicio"`
    FechaFin      time.Time `json:"fecha_fin"`
    ProfesorID    uuid.UUID `json:"profesor_id"`
    Observaciones string    `json:"observaciones"`
    Estado        string    `json:"estado"`
    CreadoEn      time.Time `json:"creado_en"`
}
