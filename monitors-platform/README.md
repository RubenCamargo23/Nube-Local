# Monitors Platform Backend 🚀

Esta plataforma está construida en **Go 1.24**, y se compone de un `api` y un `worker` en segundo plano. Emplea tecnologías robustas como **PostgreSQL** para persistencia, **Redis** para mensajería en cola, y **Ollama** para generación con Modelos de Inteligencia Artificial locales (LLM).

---

## 📋 Requisitos para el Desarrollador

Para probar y contribuir a este repositorio debes contar con:
- **Docker y Docker Compose** (Muy recomendado, levanta todo el ecosistema con 1 comando).
- **Go 1.24+** (Únicamente si deseas compilar de manera manual y física tu entorno local).

## 🌍 Archivo de Entorno (`.env`)

Crea un archivo llamado `.env` en la raíz de esta carpeta (`monitors-platform/`) y pega la siguiente configuración oficial de desarrollo para integrar las bases de datos.

```ini
# ===== SERVIDOR =====
PORT=8080

# ===== POSTGRESQL =====
DB_HOST=localhost
DB_PORT=5432
DB_USER=monitors_user
DB_PASSWORD=monitors_pass
DB_NAME=monitors_db

# ===== REDIS =====
REDIS_ADDR=localhost:6379

# ===== SEGURIDAD JWT =====
JWT_SECRET=mi_clave_hiper_secreta_local
JWT_EXPIRATION_HOURS=24

# ===== INTELIGENCIA ARTIFICIAL =====
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=qwen2.5:3b
OLLAMA_TIMEOUT_SECS=300

# ===== ALMACENAMIENTO DE RECURSOS =====
STORAGE_PATH=./storage
PDF_PATH=./storage/pdfs
```

> 🔒 **Nota sobre Producción y el Pipeline:**  
> Exponer temporalmente estas claves en tu archivo `.env` para probar y correr **no afecta en lo absoluto a tu pipeline de Github Actions**. Durante la evaluación (CI), Github inyecta variables estériles en tiempo de ejecución. **NUNCA** debes incluir claves reales y confidenciales de tu servidor productivo aquí frente al repositorio público, usa este de arriba solo como una base de desarrollo para que todos en tu equipo compilen sin dolores de cabeza.

---

## 🐳 Guía de Despliegue con Docker (Super Fácil)

La manera más estructurada de iniciar tu app es dejar que Docker despierte a las Bases de Datos y compile los servicios de Go:

```bash
docker compose build api worker
docker compose up -d
```
Esto arrancará tu base de datos relacional y el servidor principal `api` quedará escuchando en `http://localhost:8080/api/v1/health`.

## 💻 Guía de Despliegue Manual (Para Programadores Go)

Si prefieres usar la extensión de Go en tu máquina y ejecutar el código directamente:

1. Levanta primero las bases de datos externas:
`docker compose up -d postgres redis`

2. Descarga y alinea todas las dependencias locales:
`go mod tidy`
`go mod download`

3. Inicia el CLI migrador oficial para construir las tablas de Postgres:
`migrate -path ./migrations -database "postgres://monitors_user:monitors_pass@localhost:5432/monitors_db?sslmode=disable" up`

4. ¡Prende el servidor!:
`go run cmd/api/main.go`
