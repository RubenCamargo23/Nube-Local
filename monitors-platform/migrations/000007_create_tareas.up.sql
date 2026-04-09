CREATE TABLE tareas (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vinculacion_id  UUID NOT NULL REFERENCES vinculaciones(id),
    titulo          VARCHAR(200) NOT NULL,
    descripcion     TEXT NOT NULL,
    estado          VARCHAR(15) NOT NULL DEFAULT 'abierto'
                    CHECK (estado IN ('abierto', 'en_desarrollo', 'finalizado')),
    semana_inicio   DATE NOT NULL,
    horas_invertidas INT NOT NULL CHECK (horas_invertidas >= 1),
    observaciones   TEXT,
    reporte_tardio  BOOLEAN NOT NULL DEFAULT FALSE,
    creado_en       TIMESTAMP NOT NULL DEFAULT NOW(),
    actualizado_en  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tareas_vinculacion ON tareas(vinculacion_id);
CREATE INDEX idx_tareas_semana ON tareas(semana_inicio);
