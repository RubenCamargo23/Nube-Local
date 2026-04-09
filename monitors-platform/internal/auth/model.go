package auth

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  UserDTO `json:"user"`
}

type UserDTO struct {
    ID    string   `json:"id"`
    Email string   `json:"email"`
    Roles []string `json:"roles"`
}
