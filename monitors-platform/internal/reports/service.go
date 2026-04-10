package reports

import (
    "context"
    "encoding/json"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/hibiken/asynq"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
    Generate(ctx context.Context, req GenerateReportRequest) ([]*ReportResponse, error)
    FindByProfesor(ctx context.Context, profesorID uuid.UUID) ([]ReportResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*ReportResponse, error)
}

type service struct {
    repo         Repository
    ollamaClient *OllamaClient
    pdfPath      string
    db           *pgxpool.Pool
}

func NewService(repo Repository, ollamaClient *OllamaClient, pdfPath string, db *pgxpool.Pool) Service {
    return &service{repo: repo, ollamaClient: ollamaClient, pdfPath: pdfPath, db: db}
}

func (s *service) Generate(ctx context.Context, req GenerateReportRequest) ([]*ReportResponse, error) {
    vinculaciones, err := s.repo.FindVinculacionesByEspacio(ctx, req.EspacioID)
    if err != nil {
        return nil, err
    }

    parsedDate, err := time.Parse("2006-01-02", req.SemanaInicio)
    if err != nil {
        return nil, err
    }

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
    defer client.Close()

    var responses []*ReportResponse

    for _, vinculacionID := range vinculaciones {
        report, err := s.repo.Create(ctx, vinculacionID, parsedDate)
        if err != nil {
            continue
        }
        responses = append(responses, report)

        payload := ReportJobPayload{
            ReporteID:     report.ID.String(),
            VinculacionID: vinculacionID.String(),
            SemanaInicio:  req.SemanaInicio,
        }

        payloadBytes, _ := json.Marshal(payload)
        task := asynq.NewTask(TaskTypeGenerateReport, payloadBytes)
        
        // Enviar a la cola de redis
        _, err = client.EnqueueContext(ctx, task)
        if err != nil {
            s.repo.UpdateError(ctx, report.ID)
        }
    }

    return responses, nil
}

func (s *service) FindByProfesor(ctx context.Context, profesorID uuid.UUID) ([]ReportResponse, error) {
    return s.repo.FindByEspacioProfesor(ctx, profesorID)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*ReportResponse, error) {
    return s.repo.FindByID(ctx, id)
}
