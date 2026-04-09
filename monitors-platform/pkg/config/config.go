package config

import (
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    DBHost              string
    DBPort              string
    DBUser              string
    DBPassword          string
    DBName              string
    DBSSLMode           string
    DatabaseURL         string
    RedisAddr           string
    JWTSecret           string
    JWTExpirationHours  int
    OllamaHost          string
    OllamaModel         string
    OllamaTimeoutSecs   int
    StoragePath         string
    PDFPath             string
    Port                string
}

func Load() *Config {
    godotenv.Load()

    jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
    ollamaTimeout, _ := strconv.Atoi(getEnv("OLLAMA_TIMEOUT_SECONDS", "120"))

    host := getEnv("DB_HOST", "localhost")
    port := getEnv("DB_PORT", "5432")
    user := getEnv("DB_USER", "monitors_user")
    pass := getEnv("DB_PASSWORD", "monitors_pass")
    name := getEnv("DB_NAME", "monitors_db")
    ssl  := getEnv("DB_SSLMODE", "disable")

    return &Config{
        DBHost:             host,
        DBPort:             port,
        DBUser:             user,
        DBPassword:         pass,
        DBName:             name,
        DBSSLMode:          ssl,
        DatabaseURL:        "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + name + "?sslmode=" + ssl,
        RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
        JWTSecret:          getEnv("JWT_SECRET", "secret"),
        JWTExpirationHours: jwtExp,
        OllamaHost:         getEnv("OLLAMA_HOST", "http://localhost:11434"),
        OllamaModel:        getEnv("OLLAMA_MODEL", "qwen2.5:3b"),
        OllamaTimeoutSecs:  ollamaTimeout,
        StoragePath:        getEnv("STORAGE_PATH", "./storage"),
        PDFPath:            getEnv("PDF_PATH", "./storage/pdfs"),
        Port:               getEnv("PORT", "8080"),
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}
