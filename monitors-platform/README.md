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

---

## 🐘 Conectar a la Base de Datos (DBeaver)

Al ejecutar Docker, la base de datos se expone en tu máquina local. Usa las siguientes credenciales para explorarla visualmente usando [DBeaver](https://dbeaver.io/):

- **Host (Servidor):** `localhost`
- **Port (Puerto):** `5432`
- **Database (Base de Datos):** `monitors_db`
- **Username (Usuario):** `monitors_user`
- **Password (Contraseña):** `monitors_pass`

### Instrucciones para Mac y Windows:
1. Descarga, instala y abre el programa **DBeaver**.
2. Haz clic en el ícono de enchufe en la esquina superior izquierda llamado **"Nueva Conexión"** (New Database Connection).
3. Selecciona el ícono del elefante de **PostgreSQL** y dale en Siguiente.
4. Llena el formulario copiando y pegando los datos de la lista de arriba exactamente como están.
5. Haz clic en el botón de abajo **"Probar conexión"** (Test Connection). 
   *(Ojo: Si es la primera vez que lo usas, te puede salir un aviso pidiendo descargar los Drivers de Postgres, solo dale clic en el botón azul "Download").*
6. Da clic en **"Finalizar"** (Finish). Podrás abrir la conexión creada en el menú izquierdo, ingresar a los esquemas, y visualizar o alterar las tablas de la aplicación en tiempo real.

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
