# Monitors Platform Backend - Guía de Pruebas E2E (Reportes con IA)

Este documento describe el flujo completo para probar la funcionalidad de generación de reportes asistidos por IA, desde la autenticación hasta la descarga del PDF.

## 1. Requisitos Previos
*   Docker y Docker Compose instalados.
*   Contenedores en ejecución: `docker-compose up -d`.
*   Modelo Ollama descargado: `docker exec monitors-platform-ollama-1 ollama pull qwen2.5:3b`.

## 2. Usuarios de Prueba
Se han configurado los siguientes usuarios con contraseña `password123`:
*   **Profesor:** `garcia@universidad.edu`
*   **Monitor:** `gamonitor@universidad.edu`

---

## 3. Flujo de Ejecución (Paso a Paso)

### Paso 0: Preparación de Base de Datos

**1. Ejecutar Migraciones (Opcional si la API no lo hace automáticamente):**
**Instalación del comando `migrate`:**
*   **macOS:** `brew install golang-migrate`
*   **Windows:** `scoop install migrate` o `choco install golang-migrate`

**Ejecución manual:**
```bash
# macOS / Linux / Windows (Git Bash)
migrate -path ./migrations -database "postgres://monitors_user:monitors_pass@localhost:5432/monitors_db?sslmode=disable" up
```

**2. Poblar Datos de Prueba (Seed):**
Para que el flujo funcione con los IDs de prueba, debemos poblar la base de datos:
```bash
docker exec -i monitors-platform-postgres-1 psql -U monitors_user -d monitors_db < migrations/seed.sql
```
*   **Validación DB:** `SELECT count(*) FROM usuarios;` (Debe devolver al menos 2).

### Paso A: Autenticación (Login)
Primero, obtenemos los tokens JWT para ambos roles.

**Login Profesor:**
```bash
curl -X POST http://localhost:80/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email": "garcia@universidad.edu", "password": "password123"}'
```

*   **Validación Logs:** Ejecuta `docker logs monitors-platform-api-1` y busca `Login exitoso para: garcia@universidad.edu`.
*   **Validación DB:** Verifica los roles del usuario con esta consulta:
    ```sql
    SELECT u.email, r.nombre as rol 
    FROM usuarios u 
    JOIN usuario_roles ur ON u.id = ur.usuario_id 
    JOIN roles r ON ur.rol_id = r.id 
    WHERE u.email = 'garcia@universidad.edu';
    ```

### Paso B: Creación de Tareas (Monitor)
Para generar un reporte, el monitor debe haber registrado tareas en un espacio activo.

```bash
curl -X POST http://localhost:80/api/v1/vinculaciones/11111111-1111-1111-1111-11111111aaaa/tareas \
     -H "Authorization: Bearer [TOKEN_MONITOR]" \
     -H "Content-Type: application/json" \
     -d '{
       "titulo": "Desarrollo de Endpoints API",
       "descripcion": "Implementación de lógica para reportes",
       "estado": "finalizado",
       "semana_inicio": "2026-03-30",
       "horas_invertidas": 5
     }'
```
*   **Validación Logs:** Ejecuta `docker logs monitors-platform-api-1` y busca `Tarea creada exitosamente`.
*   **Validación DB:** `SELECT * FROM tareas WHERE vinculacion_id = '11111111-1111-1111-1111-11111111aaaa';`.

### Paso C: Disparo de Generación de Reporte (Profesor)
El profesor solicita la generación del reporte para una semana específica.

```bash
curl -X POST http://localhost:80/api/v1/profesor/reportes/generar \
     -H "Authorization: Bearer [TOKEN_PROFESOR]" \
     -H "Content-Type: application/json" \
     -d '{
       "espacio_id": "cccccccc-cccc-cccc-cccc-cccccccccccc",
       "semana_inicio": "2026-03-30"
     }'
```
*   **Validación API:** Ejecuta `docker logs monitors-platform-api-1` y busca `Reporte encolado correctamente`.
*   **Validación DB:** `SELECT id, estado_generacion FROM reportes_pdf WHERE estado_generacion = 'PENDING' LIMIT 1;`.

### Paso D: Procesamiento (Background Worker)
El worker de Asynq tomará la tarea y llamará a Ollama. 

*   **Validación Logs IA:** Ejecuta `docker logs monitors-platform-worker-1`:
    *   Busca: `Enviando prompt a Ollama...` (Confirmación de entrada).
    *   Busca: `Ollama respondió...` (Confirmación de salida).
*   **Validación DB Física:** `SELECT ruta_pdf, prompt_usado FROM reportes_pdf WHERE estado_generacion = 'DONE';`.

### Paso E: Descarga del Reporte (Profesor)
Una vez que el estado sea `DONE`, descarga el archivo final:

```bash
curl -X GET http://localhost:80/api/v1/profesor/reportes/[ID_REPORTE]/descargar \
     -H "Authorization: Bearer [TOKEN_PROFESOR]" \
     --output reporte_semanal.pdf
```
*   **Validación Física:** Verifica que el archivo exista en el volumen: `docker exec monitors-platform-worker-1 ls -lh /app/storage/pdfs`.

---

## 4. Troubleshooting y Notas Técnicas

### Manejo de Caracteres Especiales (Encoding)
Se implementó el uso de `UnicodeTranslatorFromDescriptor("")` en `internal/reports/pdf.go` para asegurar que el texto generado por la IA en UTF-8 se traduzca correctamente al charset ISO-8859-1 del PDF, evitando caracteres corruptos en palabras con acentos o eñes.

### Rendimiento de IA (Ollama en CPU)
Dado que Ollama se ejecuta en CPU dentro de Docker:
*   La concurrencia del worker se ha fijado en `1` para maximizar el uso de recursos por tarea.
*   El timeout del cliente HTTP de Ollama se ha extendido a `600s` (10 min) en el archivo `.env`.

### Monitoreo de Tareas
Puedes usar `redis-cli` para ver el estado de las colas de Asynq:
```bash
docker exec -it monitors-platform-redis-1 redis-cli
```

---

## 5. Arquitectura e Interacción de Componentes

El flujo de generación de reportes sigue un modelo asíncrono para garantizar la estabilidad del sistema:

### Diagrama de Flujo
```mermaid
sequenceDiagram
    participant P as Profesor
    participant A as API (Gin)
    participant R as Redis (Asynq)
    participant W as Worker (Go)
    participant DB as PostgreSQL
    participant O as Ollama (IA)
    participant V as Storage (Volume)

    P->>A: Solicita Reporte
    A->>DB: Crea Registro (PENDING)
    A->>R: Encola Tarea
    W->>R: Recoge Tarea
    W->>DB: Obtiene datos de la semana
    W->>O: Genera Resumen con IA
    W->>V: Crea PDF corregido (UTF-8)
    W->>DB: Actualiza a DONE + Ruta PDF
    P->>A: Descarga PDF
```

### Cómo verificar cada componente durante el flujo:

1.  **API (Validación de Proceso):**
    *   Verifica que los logs de la API muestren la recepción del POST: `docker logs monitors-platform-api-1`.
    *   Consulta la DB para ver el registro inicial: `SELECT * FROM reportes_pdf WHERE estado_generacion = 'PENDING';`.

2.  **Redis (Cola de espera):**
    *   Verifica que la tarea entró a la cola: `docker exec -it monitors-platform-redis-1 redis-cli LLEN asynq:queues:default`.

3.  **Worker (Lógica de Negocio):**
    *   Sigue el procesamiento en tiempo real: `docker logs -f monitors-platform-worker-1`.
    *   Debes ver: `Procesando reporte...` -> `Llamando a Ollama...` -> `Generando PDF...`.

4.  **Ollama (Inferencia de IA):**
    *   Verifica que el modelo esté cargado y procesando: `docker exec monitors-platform-ollama-1 ollama ps`.
    *   Observa el consumo de recursos: `docker stats monitors-platform-ollama-1` (el CPU debería subir significativamente).

5.  **Storage (Persistencia física):**
    *   Verifica que el archivo PDF se creó en el volumen: `docker exec monitors-platform-worker-1 ls -lh /app/storage/pdfs`.

---

### Secuencia de Prueba E2E Realizada (Ejemplo Real)

Esta es la secuencia exacta utilizada durante la estabilización del sistema para verificar el flujo:

1.  **Limpieza inicial de tareas:**
    *   Se aseguró que el usuario `gamonitor@universidad.edu` tuviera su contraseña corregida a `password123`.

2.  **Registro de actividad (Monitor):**
    ```bash
curl -X POST http://localhost:80/api/v1/vinculaciones/11111111-1111-1111-1111-11111111aaaa/tareas \
     -H "Authorization: Bearer [TOKEN_MONITOR]" \
     -H "Content-Type: application/json" \
     -d '{
       "titulo": "Desarrollo de Endpoints API",
       "descripcion": "Implementación de lógica para reportes",
       "estado": "finalizado",
       "semana_inicio": "2026-03-30",
       "horas_invertidas": 5
     }'
```

3.  **Solicitud de reporte (Profesor):**
    ```bash
    # Disparar reporte para la semana 2026-03-30
    curl -X POST http://localhost:80/api/v1/profesor/reportes/generar \
         -H "Authorization: Bearer [TOKEN_PROFESOR]" \
         -d '{"espacio_id": "cccccccc-cccc-cccc-cccc-cccccccccccc", "semana_inicio": "2026-03-30"}'
    ```

4.  **Monitoreo del Worker:**
    Apareció el log: `Ollama respondió para reporte ee260024...` seguido del texto interpretado.

5.  **Confirmación y Descarga:**
    El reporte cambió a estado `DONE` y se pudo descargar exitosamente validando que los acentos en "realizó" y "implementación" se vieran correctamente.

---

## 7. Validaciones en Base de Datos

Durante el flujo, puedes realizar estas consultas para asegurar la integridad de los datos:

| Paso | Tabla | Qué Validar | Consulta SQL |
| :--- | :--- | :--- | :--- |
| **Inicio** | `tareas` | Que existan tareas con `horas_invertidas > 0` para la semana. | `SELECT * FROM tareas WHERE vinculacion_id = '...' AND semana_inicio = '...';` |
| **Disparo** | `reportes_pdf` | Que el registro se cree en estado `PENDING`. | `SELECT id, estado_generacion FROM reportes_pdf ORDER BY creado_en DESC LIMIT 1;` |
| **IA** | `reportes_pdf` | Que el prompt usado se guarde (una vez terminado). | `SELECT prompt_usado FROM reportes_pdf WHERE id = '...';` |
| **Final** | `reportes_pdf` | Que la `ruta_pdf` no sea nula y el estado sea `DONE`. | `SELECT ruta_pdf, generado_en FROM reportes_pdf WHERE id = '...';` |

---

## 9. Guía de Colaboración y Pipeline (Equipo de 4 Integrantes)

Para garantizar la estabilidad del proyecto y el éxito del pipeline de CI/CD, se define el siguiente esquema de trabajo:

### Estructura de Trabajo por Roles

| Integrante | Rol | Actividades Principales | Rama Base |
| :--- | :--- | :--- | :--- |
| **Ruben Camargo** | **DevOps & Security** | Infraestructura Docker, GitHub Actions, Seguridad JWT y Middleware de acceso. | `feature/infra-security` |
| **Diego Rojas** | **Core Backend Dev** | Modelos DB, Repositorios, Servicios de Espacios/Tareas y reglas de negocio. | `feature/core-domain` |
| **David Rojas** | **AI & Worker Dev** | Lógica de Asynq (Redis), Integración con Ollama y generación de PDFs. | `feature/ai-worker` |
| **Brian Martinez** | **QA & Doc Lead** | Pruebas unitarias, E2E con Postman, actualización de README y Swagger. | `feature/qa-docs` |

### Estrategia de Ramas (Gitflow)
*   **`main`**: Código de producción. Solo se actualiza tras un release exitoso.
*   **`develop`**: Rama principal de integración.
*   **`feature/*`**: Ramas de desarrollo individual. Se fusionan a `develop` mediante **Pull Request**.

### Reglas para NO romper el Pipeline
1.  **Pruebas Locales:** Nunca hacer `push` sin ejecutar `go test ./...` localmente.
2.  **Validación Docker:** Asegurar que el proyecto compila con `docker-compose build` antes de integrar.
3.  **Migraciones Atómicas:** Cualquier cambio en la DB debe ir acompañado de su archivo correspondiente en `/migrations`.
4.  **Mocks para IA:** En el pipeline de CI, las pruebas de IA deben usar Mocks para evitar la dependencia de Ollama y el consumo excesivo de CPU.
5.  **Code Review:** Al menos un compañero debe aprobar el Pull Request antes de fusionar a `develop`.

