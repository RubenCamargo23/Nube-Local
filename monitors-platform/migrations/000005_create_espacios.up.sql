CREATE TABLE espacios (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tipo         VARCHAR(10) NOT NULL CHECK (tipo IN ('COURSE', 'PROJECT')),
    nombre       VARCHAR(200) NOT NULL,
    periodo_id   UUID NOT NULL REFERENCES periodos_academicos(id),
    fecha_inicio DATE NOT NULL,
    fecha_fin    DATE NOT NULL,
    profesor_id  UUID NOT NULL REFERENCES usuarios(id),
    observaciones TEXT,
    estado       VARCHAR(10) NOT NULL DEFAULT 'ACTIVE'
                 CHECK (estado IN ('ACTIVE', 'CLOSED')),
    creado_en    TIMESTAMP NOT NULL DEFAULT NOW()
);
