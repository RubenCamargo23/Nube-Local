CREATE TABLE periodos_academicos (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    codigo      VARCHAR(10) NOT NULL UNIQUE,
    fecha_inicio DATE NOT NULL,
    fecha_fin    DATE NOT NULL,
    estado       VARCHAR(10) NOT NULL DEFAULT 'ACTIVE'
                 CHECK (estado IN ('ACTIVE', 'CLOSED')),
    creado_en    TIMESTAMP NOT NULL DEFAULT NOW()
);
