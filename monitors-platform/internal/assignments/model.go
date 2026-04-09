package assignments

import (
    "time"
    "github.com/google/uuid"
)

type CreateAssignmentRequest struct {
    UsuarioID      uuid.UUID `json:"usuario_id" binding:"required"`
    Rol            string    `json:"rol" binding:"required,oneof=MONITOR GRAD_ASSISTANT"`
    HorasSemanales int       `json:"horas_semanales" binding:"required,min=1"`
}

type AssignmentResponse struct {
    ID             uuid.UUID `json:"id"`
    UsuarioID      uuid.UUID `json:"usuario_id"`
    EspacioID      uuid.UUID `json:"espacio_id"`
    Rol            string    `json:"rol"`
    HorasSemanales int       `json:"horas_semanales"`
    ProfesorID     uuid.UUID `json:"profesor_id"`
    Activa         bool      `json:"activa"`
    CreadoEn       time.Time `json:"creado_en"`
}
