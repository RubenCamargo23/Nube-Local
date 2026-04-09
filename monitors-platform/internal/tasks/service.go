package tasks

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/org/monitors-platform/internal/spaces"
    "github.com/org/monitors-platform/internal/assignments"
    "github.com/org/monitors-platform/pkg/week"
)

var (
    ErrTareaNoEditable  = errors.New("la tarea no es editable fuera de la semana activa")
    ErrReporteTardio    = errors.New("un reporte tardío no puede editarse (RN-33)")
    ErrEspacioCerrado   = errors.New("el espacio está cerrado (RN-09)")
)

type Service interface {
    Create(ctx context.Context, vinculacionID uuid.UUID, req CreateTaskRequest) (*TaskResponse, error)
    FindByVinculacion(ctx context.Context, vinculacionID uuid.UUID) ([]TaskResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*TaskResponse, error)
    Update(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
    repo           Repository
    assignmentsSvc assignments.Service
    spacesSvc      spaces.Service
}

func NewService(repo Repository, assignmentsSvc assignments.Service, spacesSvc spaces.Service) Service {
    return &service{repo: repo, assignmentsSvc: assignmentsSvc, spacesSvc: spacesSvc}
}

func (s *service) Create(ctx context.Context, vinculacionID uuid.UUID, req CreateTaskRequest) (*TaskResponse, error) {
    assignment, err := s.assignmentsSvc.FindByID(ctx, vinculacionID)
    if err != nil {
        return nil, err
    }

    active, err := s.spacesSvc.IsActive(ctx, assignment.EspacioID)
    if err != nil || !active {
        return nil, ErrEspacioCerrado
    }

    parsedDate, err := time.Parse("2006-01-02", req.SemanaInicio)
    if err != nil {
        return nil, errors.New("formato de semana_inicio inválido, usar YYYY-MM-DD")
    }

    semanaInicio := week.GetWeekStart(parsedDate)
    tardio := week.IsLateReport(semanaInicio)

    return s.repo.Create(ctx, vinculacionID, req, semanaInicio, tardio)
}

func (s *service) FindByVinculacion(ctx context.Context, vinculacionID uuid.UUID) ([]TaskResponse, error) {
    return s.repo.FindByVinculacion(ctx, vinculacionID)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*TaskResponse, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*TaskResponse, error) {
    task, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Bloquear si reporte tardío (RN-33)
    if task.ReporteTardio {
        return nil, ErrReporteTardio
    }

    // Bloquear si semana ya cerró (respuesta profesor)
    if week.IsLateReport(task.SemanaInicio) {
        return nil, ErrTareaNoEditable
    }

    return s.repo.Update(ctx, id, req)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
    task, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    // Bloquear si semana ya cerró (RN-34)
    if week.IsLateReport(task.SemanaInicio) {
        return ErrTareaNoEditable
    }

    return s.repo.Delete(ctx, id)
}
