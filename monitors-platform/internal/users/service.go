package users

import (
    "context"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type Service interface {
    Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
    AssignRole(ctx context.Context, userID uuid.UUID, req AssignRoleRequest) error
    FindAll(ctx context.Context) ([]UserResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    return s.repo.Create(ctx, req.Nombre, req.Email, string(hash))
}

func (s *service) AssignRole(ctx context.Context, userID uuid.UUID, req AssignRoleRequest) error {
    return s.repo.AssignRole(ctx, userID, req.Rol)
}

func (s *service) FindAll(ctx context.Context) ([]UserResponse, error) {
    users, err := s.repo.FindAll(ctx)
    if err != nil {
        return nil, err
    }
    for i := range users {
        roles, _ := s.repo.FindRoles(ctx, users[i].ID)
        users[i].Roles = roles
    }
    return users, nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
    u, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    u.Roles, _ = s.repo.FindRoles(ctx, id)
    return u, nil
}
