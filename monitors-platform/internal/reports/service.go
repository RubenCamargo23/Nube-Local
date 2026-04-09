package reports

type Service struct{}

func NewService(repo Repository, ollamaClient *OllamaClient, pdfPath string, db interface{}) *Service {
    return &Service{}
}
