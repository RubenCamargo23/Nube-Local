package spaces

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/org/monitors-platform/internal/periods"
)

var (
    ErrPeriodNotActive = errors.New("el período académico no está activo (RN-07)")
    ErrSpaceClosed     = errors.New("el espacio está cerrado (RN-09)")
)

type Service interface {
    Create(ctx context.Context, req CreateSpaceRequest, profesorID uuid.UUID) (*SpaceResponse, error)
    FindByProfesor(ctx context.Context, profesorID uuid.UUID) ([]SpaceResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*SpaceResponse, error)
    Close(ctx context.Context, id uuid.UUID) error
    IsActive(ctx context.Context, id uuid.UUID) (bool, error)
}

type service struct {
    repo       Repository
    periodsSvc periods.Service
}

func NewService(repo Repository, periodsSvc periods.Service) Service {
    return &service{repo: repo, periodsSvc: periodsSvc}
}

func (s *service) Create(ctx context.Context, req CreateSpaceRequest, profesorID uuid.UUID) (*SpaceResponse, error) {
    active, err := s.periodsSvc.IsActive(ctx, req.PeriodoID)
    if err != nil {
        return nil, err
    }
    if !active {
        return nil, ErrPeriodNotActive
    }
    return s.repo.Create(ctx, req, profesorID)
}

func (s *service) FindByProfesor(ctx context.Context, profesorID uuid.UUID) ([]SpaceResponse, error) {
    return s.repo.FindByProfesor(ctx, profesorID)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*SpaceResponse, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *service) Close(ctx context.Context, id uuid.UUID) error {
    return s.repo.Close(ctx, id)
}

func (s *service) IsActive(ctx context.Context, id uuid.UUID) (bool, error) {
    sp, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return false, err
    }
    return sp.Estado == "ACTIVE", nil
}
