package periods

import (
    "context"
    "errors"
    "github.com/google/uuid"
)

var ErrPeriodClosed = errors.New("período cerrado")

type Service interface {
    Create(ctx context.Context, req CreatePeriodRequest) (*PeriodResponse, error)
    FindAll(ctx context.Context) ([]PeriodResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*PeriodResponse, error)
    Close(ctx context.Context, id uuid.UUID) error
    IsActive(ctx context.Context, id uuid.UUID) (bool, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, req CreatePeriodRequest) (*PeriodResponse, error) {
    return s.repo.Create(ctx, req.Codigo, req.FechaInicio, req.FechaFin)
}

func (s *service) FindAll(ctx context.Context) ([]PeriodResponse, error) {
    return s.repo.FindAll(ctx)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*PeriodResponse, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *service) Close(ctx context.Context, id uuid.UUID) error {
    return s.repo.Close(ctx, id)
}

func (s *service) IsActive(ctx context.Context, id uuid.UUID) (bool, error) {
    p, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return false, err
    }
    return p.Estado == "ACTIVE", nil
}
