CREATE TABLE adjuntos (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tarea_id            UUID NOT NULL REFERENCES tareas(id) ON DELETE CASCADE,
    nombre_archivo      VARCHAR(255) NOT NULL,
    ruta_almacenamiento VARCHAR(500) NOT NULL,
    tipo_mime           VARCHAR(100),
    tamanio_bytes       BIGINT,
    creado_en           TIMESTAMP NOT NULL DEFAULT NOW()
);
