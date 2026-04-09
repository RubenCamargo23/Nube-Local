package assignments

import (
    "context"
    "errors"
    "math"

    "github.com/google/uuid"
    "github.com/org/monitors-platform/internal/spaces"
)

var (
    ErrExcedeHorasGA        = errors.New("excede el límite de 22h como asistente graduado (RN-14)")
    ErrMaxMonitorias        = errors.New("excede el límite de 3 monitorías simultáneas (RN-15)")
    ErrExcedeHorasMonitor   = errors.New("excede el límite de 12h como monitor (RN-16)")
    ErrReglaCombinada       = errors.New("horas de monitoría exceden el 40% de horas GA (RN-17)")
    ErrEspacioCerrado       = errors.New("el espacio está cerrado")
)

type Service interface {
    Create(ctx context.Context, espacioID, profesorID uuid.UUID, req CreateAssignmentRequest) (*AssignmentResponse, error)
    FindByEspacio(ctx context.Context, espacioID uuid.UUID) ([]AssignmentResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*AssignmentResponse, error)
    Update(ctx context.Context, id uuid.UUID, horas int) error
}

type service struct {
    repo      Repository
    spacesSvc spaces.Service
}

func NewService(repo Repository, spacesSvc spaces.Service) Service {
    return &service{repo: repo, spacesSvc: spacesSvc}
}

func (s *service) Create(ctx context.Context, espacioID, profesorID uuid.UUID, req CreateAssignmentRequest) (*AssignmentResponse, error) {
    // Verificar espacio activo
    active, err := s.spacesSvc.IsActive(ctx, espacioID)
    if err != nil || !active {
        return nil, ErrEspacioCerrado
    }

    // Iniciar transacción con SELECT FOR UPDATE
    tx, err := s.repo.BeginTx(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

    if req.Rol == "GRAD_ASSISTANT" {
        // RN-14: límite 22h GA
        horasGA, err := s.repo.SumarHorasPorRol(ctx, tx, req.UsuarioID, "GRAD_ASSISTANT")
        if err != nil {
            return nil, err
        }
        if horasGA+req.HorasSemanales > 22 {
            return nil, ErrExcedeHorasGA
        }
    }

    if req.Rol == "MONITOR" {
        // RN-15: máx 3 monitorías
        countMonitorias, err := s.repo.ContarMonitorias(ctx, tx, req.UsuarioID)
        if err != nil {
            return nil, err
        }
        if countMonitorias >= 3 {
            return nil, ErrMaxMonitorias
        }

        // RN-16: límite 12h monitor
        horasMonitor, err := s.repo.SumarHorasPorRol(ctx, tx, req.UsuarioID, "MONITOR")
        if err != nil {
            return nil, err
        }
        if horasMonitor+req.HorasSemanales > 12 {
            return nil, ErrExcedeHorasMonitor
        }

        // RN-17 + RN-18: regla combinada con GA
        horasGA, err := s.repo.SumarHorasPorRol(ctx, tx, req.UsuarioID, "GRAD_ASSISTANT")
        if err != nil {
            return nil, err
        }
        if horasGA > 0 {
            maxMonitor := int(math.Ceil(float64(horasGA) * 0.4))
            if horasMonitor+req.HorasSemanales > maxMonitor {
                return nil, ErrReglaCombinada
            }
        }
    }

    assignment, err := s.repo.Create(ctx, tx, espacioID, req.UsuarioID, profesorID, req.Rol, req.HorasSemanales)
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }

    return assignment, nil
}

func (s *service) FindByEspacio(ctx context.Context, espacioID uuid.UUID) ([]AssignmentResponse, error) {
    return s.repo.FindByEspacio(ctx, espacioID)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*AssignmentResponse, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, horas int) error {
    return s.repo.Update(ctx, id, horas)
}
