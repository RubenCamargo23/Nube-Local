package auth

import (
    "context"
    "errors"

    "golang.org/x/crypto/bcrypt"
    jwtpkg "github.com/org/monitors-platform/pkg/jwt"
)

var ErrInvalidCredentials = errors.New("credenciales inválidas")

type Service interface {
    Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

type service struct {
    repo       Repository
    jwtManager *jwtpkg.Manager
}

func NewService(repo Repository, jwtManager *jwtpkg.Manager) Service {
    return &service{repo: repo, jwtManager: jwtManager}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    user, err := s.repo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, ErrInvalidCredentials
    }

    roles, err := s.repo.FindRolesByUserID(ctx, user.ID)
    if err != nil {
        return nil, err
    }

    token, err := s.jwtManager.Generate(user.ID, user.Email, roles)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token: token,
        User:  UserDTO{ID: user.ID.String(), Email: user.Email, Roles: roles},
    }, nil
}
