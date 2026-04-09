CREATE TABLE roles (
    id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nombre VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (nombre) VALUES
    ('ADMIN'),
    ('PROFESSOR'),
    ('MONITOR'),
    ('GRAD_ASSISTANT');
