-- Seed Script para arranque inmediato del flujo de reportes
-- Contraseña para todos los usuarios: password123 (HAsh SHA1)

-- 1. Usuarios Base
INSERT INTO usuarios (id, nombre, email, password_hash) VALUES
('22222222-2222-2222-2222-222222222222', 'Profesor Garcia', 'garcia@universidad.edu', '$2a$10$3l943zTUHOPIIJhqYwVIU.Sp/o6iGDShH2iNneyVNPr3uLatvw4ya'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Estudiante GA y Monitor', 'gamonitor@universidad.edu', '$2a$10$3l943zTUHOPIIJhqYwVIU.Sp/o6iGDShH2iNneyVNPr3uLatvw4ya')
ON CONFLICT (id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- 2. Asignar Roles
INSERT INTO usuario_roles (usuario_id, rol_id) VALUES
('22222222-2222-2222-2222-222222222222', (SELECT id FROM roles WHERE nombre = 'PROFESSOR')),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', (SELECT id FROM roles WHERE nombre = 'MONITOR'))
ON CONFLICT DO NOTHING;

-- 3. Periodo Académico
INSERT INTO periodos_academicos (id, codigo, fecha_inicio, fecha_fin, estado) VALUES
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '2026-1', '2026-01-01', '2026-06-30', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 4. Espacio Académico
INSERT INTO espacios (id, tipo, nombre, periodo_id, fecha_inicio, fecha_fin, profesor_id, estado) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'COURSE', 'Diseño de Sistemas API', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '2026-01-01', '2026-06-30', '22222222-2222-2222-2222-222222222222', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 5. Vinculación
INSERT INTO vinculaciones (id, usuario_id, espacio_id, rol, horas_semanales, profesor_id, activa) VALUES
('11111111-1111-1111-1111-11111111aaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'MONITOR', 12, '22222222-2222-2222-2222-222222222222', true)
ON CONFLICT (id) DO NOTHING;
