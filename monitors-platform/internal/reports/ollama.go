package reports

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type OllamaRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
}

type OllamaResponse struct {
    Response string `json:"response"`
}

type OllamaClient struct {
    host    string
    model   string
    timeout time.Duration
}

func NewOllamaClient(host, model string, timeoutSecs int) *OllamaClient {
    return &OllamaClient{
        host:    host,
        model:   model,
        timeout: time.Duration(timeoutSecs) * time.Second,
    }
}

func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
    body, _ := json.Marshal(OllamaRequest{
        Model:  c.model,
        Prompt: prompt,
        Stream: false,
    })

    httpClient := &http.Client{Timeout: c.timeout}
    req, err := http.NewRequestWithContext(ctx, "POST", c.host+"/api/generate", bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("error llamando Ollama: %w", err)
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var ollamaResp OllamaResponse
    if err := json.Unmarshal(data, &ollamaResp); err != nil {
        return "", err
    }

    return ollamaResp.Response, nil
}
