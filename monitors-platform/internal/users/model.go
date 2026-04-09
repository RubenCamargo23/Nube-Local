package users

import (
    "time"
    "github.com/google/uuid"
)

type CreateUserRequest struct {
    Nombre   string `json:"nombre" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type AssignRoleRequest struct {
    Rol string `json:"rol" binding:"required"`
}

type UserResponse struct {
    ID       uuid.UUID `json:"id"`
    Nombre   string    `json:"nombre"`
    Email    string    `json:"email"`
    Roles    []string  `json:"roles"`
    CreadoEn time.Time `json:"creado_en"`
}
