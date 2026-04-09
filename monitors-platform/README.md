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

## 📊 Modelo de Datos y Relaciones (Estructura)

La plataforma utiliza un modelo relacional estructurado en PostgreSQL. Estas son sus tablas principales y cómo se conectan:

*   **`usuarios`**: Almacena a todas las personas del sistema (Estudiantes, Profesores, Administrativos). Se autentica vía JWT.
*   **`espacios` y `periodos_academicos`**: Definen las asignaturas, departamentos académicos y los semestres en los que ocurren las tutorías.
*   **`vinculaciones` (Tabla Pivote Central)**: Une a un estudiante (`usuario_id`) con su asignatura (`espacio_id`) y su profesor a cargo (`profesor_id`). Aquí se estipulan las reglas de negocio como el tipo (`MONITOR` / `GRAD_ASSISTANT`) y la carga máxima de `horas_semanales`.
*   **`tareas`**: El registro semanal (timesheet) del monitor asociado a una `vinculación`. Contiene el registro de las `horas_invertidas`, el `estado`, y si el reporte fue enviado tardíamente.
*   **`roles`, `adjuntos` y `reportes_pdf`**: Tablas de soporte para guardar control de permisos y rutas de la generación asíncrona de archivos PDF procesados por Redis.

---

## 🚀 Paso a Paso: Cómo arrancar y probar el proyecto en local

Para poner a marchar la aplicación en tu computadora de manera exitosa, sigue este orden estrictamente:

### Paso 1: Inicializar el `.env`
Primero, busca en esta misma carpeta un archivo llamado `.env.production`. Hazle una copia exacta a ese archivo y cámbiale el nombre para que quede llamándose únicamente `.env`. (El sistema detectará que ese es tu archivo de llaves para poder funcionar).

### Paso 2: Levantar el Servidor y las Bases de Datos
Abre la terminal aquí dentro de `monitors-platform/` e invoca a **Docker**. Él hará la tarea sucia de construir nuestra base de datos, el gestor de colas y el backend al mismo tiempo:

```bash
# 1. Empaqueta el código y sus dependencias
docker compose build api worker

# 2. Enciende el sistema completo
docker compose up -d
```
*(Tip: Si algún día solo quieres apagarlo todo al terminar tu jornada, escribe `docker compose down -v`)*.

### Paso 3: Comprobar la Vida del Servidor (Ping)
Espera unos 15 segundos a que la base de datos termine de despertar, abre tu navegador web o herramienta como POSTMAN, y dirígete a nuestra ruta médica:

👉 **[http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)**

Si el navegador te responde con un mensaje exitoso de "Status OK", ¡Felicidades! Tu ecosistema está perfectamente compilado y listo para que interactúes con toda la API local.

### Alternativa: Desarrollo Manual sin Docker
Si no quieres empaquetar la app de Go en Docker sino compilarla tú mismo en caliente, haz esto:
1. `docker compose up -d postgres redis` *(Prende solas las BD)*.
2. `go mod tidy` *(Descarga librerías manuales)*.
3. `migrate -path ./migrations -database "postgres://monitors_user:monitors_pass@localhost:5432/monitors_db?sslmode=disable" up`
4. `go run cmd/api/main.go` *(Enciende el servidor nativo).*

---

## 📡 Endpoints de la API y Guía de Postman

La API sigue el estándar REST bajo el prefijo `http://localhost:8080/api/v1`. Para consumirla, usa **Postman** siguiendo este método:

### 1. Autenticación Global en Postman (JWT)
Salvo el Ping y el Login, TODAS las rutas están protegidas. 
1. Realiza una petición `POST` a `http://localhost:8080/api/v1/auth/login` pasándole en el Body (formato raw JSON) el email y password del usuario. Te devolverá en verde un `token`.
2. Copia todo ese texto del token.
3. En Postman, ve a tu nueva petición (Ej: ver tareas), entra a la pestaña intermedia llamada **"Authorization"**.
4. En el campo *Type* elige **"Bearer Token"**. Pega tu token en la casilla derecha. *(Ahora estás autenticado y el servidor sabe quién eres).*

### 2. Catálogo Oficial de Rutas

**🌐 Públicas**
| Método | Endpoint | Acción |
| :--- | :--- | :--- |
| **GET** | `/api/v1/health` | Ping. Devuelve información de que el server está vivo y su versión. |
| **POST**| `/api/v1/auth/login` | Recibe `{"email": "", "password": ""}` y devuelve el token JWT del actor. |

**🧑‍💻 Administrador (`ADMIN`)**
| Método | Endpoint | Acción |
| :--- | :--- | :--- |
| **POST** | `/api/v1/admin/usuarios` | Crea una persona nueva en el sistema. |
| **GET** | `/api/v1/admin/usuarios` | Lista todo el padrón de usuarios. |
| **POST** | `/api/v1/admin/usuarios/:id/roles` | Le asigna o transfiere un rol directivo a un usuario. |
| **POST** | `/api/v1/admin/periodos` | Inaugura un periodo académico nuevo (ej. "2024-1"). |
| **GET** | `/api/v1/admin/periodos` | Lista el historial de periodos académicos. |
| **PATCH**| `/api/v1/admin/periodos/:id/cerrar` | Da por finalizado un periodo y prohíbe nuevas asignaciones ahí. |
| **PATCH**| `/api/v1/admin/vinculaciones/:id` | Bloquea o actualiza la vinculación de un estudiante. |

**🎓 Profesor (`PROFESSOR`)**
| Método | Endpoint | Acción |
| :--- | :--- | :--- |
| **POST** | `/api/v1/espacios` | Crea un espacio académico nuevo o asignatura (Ej. "Física I"). |
| **GET** | `/api/v1/espacios` | Devuelve todas las asignaturas a dictar de ese profesor. |
| **POST** | `/api/v1/espacios/:id/vinculaciones`| Contrata a un Estudiante/Asistente y le adjudica horas. |
| **GET** | `/api/v1/profesor/vinculaciones` | Lista los estudiantes subalternos al profesor actual. |
| **GET** | `/api/v1/profesor/vinculaciones/:id/tareas`| Le permite al Profesor auditar las tareas que ha subido un monitor. |
| **POST** | `/api/v1/profesor/reportes/generar`| Dispara el Worker de Inteligencia Artificial (Ollama) para crear PDFs. |
| **GET** | `/api/v1/profesor/reportes/:id/descargar`| Sirve el archivo binario PDF final con las firmas y resúmenes de AI. |

**📝 Estudiantes (`MONITOR` / `GRAD_ASSISTANT`)**
| Método | Endpoint | Acción |
| :--- | :--- | :--- |
| **GET** | `/api/v1/me/vinculaciones` | Devuelve los "Trabajos" asignados del estudiante que tiene la sesión iniciada. |
| **POST** | `/api/v1/vinculaciones/:id/tareas`| Sube el registro (timesheet) de qué hizo esta semana, y cuántas horas invirtió. |
| **GET** | `/api/v1/vinculaciones/:id/tareas`| Lista todas las tareas documentadas del monitor para esa vinculación. |
| **GET** | `/api/v1/me/tareas/historial` | Balance general de todo el tiempo y tareas invertidas por él. |
| **POST** | `/api/v1/tareas/:id/adjuntos` | Sube evidencias físicas (PDF, imágenes) para respaldar una labor. |
