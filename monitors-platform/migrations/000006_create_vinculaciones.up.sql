CREATE TABLE vinculaciones (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id      UUID NOT NULL REFERENCES usuarios(id),
    espacio_id      UUID NOT NULL REFERENCES espacios(id),
    rol             VARCHAR(20) NOT NULL CHECK (rol IN ('MONITOR', 'GRAD_ASSISTANT')),
    horas_semanales INT NOT NULL CHECK (horas_semanales > 0),
    profesor_id     UUID NOT NULL REFERENCES usuarios(id),
    activa          BOOLEAN NOT NULL DEFAULT TRUE,
    creado_en       TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vinculaciones_usuario ON vinculaciones(usuario_id);
CREATE INDEX idx_vinculaciones_espacio ON vinculaciones(espacio_id);
