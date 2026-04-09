CREATE TABLE reportes_pdf (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vinculacion_id    UUID NOT NULL REFERENCES vinculaciones(id),
    semana_inicio     DATE NOT NULL,
    ruta_pdf          VARCHAR(500),
    prompt_usado      TEXT,
    estado_generacion VARCHAR(15) NOT NULL DEFAULT 'PENDING'
                      CHECK (estado_generacion IN ('PENDING', 'PROCESSING', 'DONE', 'ERROR')),
    generado_en       TIMESTAMP,
    creado_en         TIMESTAMP NOT NULL DEFAULT NOW()
);
