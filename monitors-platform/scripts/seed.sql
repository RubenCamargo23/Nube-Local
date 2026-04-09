-- Contraseña para todos: monitors123
-- Hash bcrypt de 'monitors123'
INSERT INTO usuarios (id, nombre, email, password_hash) VALUES
  ('11111111-1111-1111-1111-111111111111', 'Admin Sistema',    'admin@universidad.edu',    '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('22222222-2222-2222-2222-222222222222', 'Prof. García',      'garcia@universidad.edu',   '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('33333333-3333-3333-3333-333333333333', 'Prof. Martínez',    'martinez@universidad.edu', '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('44444444-4444-4444-4444-444444444444', 'Ana Rodríguez',     'ana@universidad.edu',      '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('55555555-5555-5555-5555-555555555555', 'Juan Pérez',        'juan@universidad.edu',     '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('66666666-6666-6666-6666-666666666666', 'María López',       'maria@universidad.edu',    '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('77777777-7777-7777-7777-777777777777', 'Carlos Gómez',      'carlos@universidad.edu',   '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO'),
  ('88888888-8888-8888-8888-888888888888', 'Estudiante GA y Monitor', 'gamonitor@universidad.edu', '$2a$10$T8Z/k6eY.o1QzZ1kQnL.heOONXF3x9T8v2YgH6Hxjz74t4DqF/0sO');

INSERT INTO usuario_roles (usuario_id, rol_id)
SELECT u.id, r.id FROM usuarios u, roles r
WHERE (u.email = 'admin@universidad.edu'    AND r.nombre = 'ADMIN')
   OR (u.email = 'garcia@universidad.edu'   AND r.nombre = 'PROFESSOR')
   OR (u.email = 'martinez@universidad.edu' AND r.nombre = 'PROFESSOR')
   OR (u.email = 'ana@universidad.edu'      AND r.nombre = 'MONITOR')
   OR (u.email = 'juan@universidad.edu'     AND r.nombre = 'MONITOR')
   OR (u.email = 'maria@universidad.edu'    AND r.nombre = 'GRAD_ASSISTANT')
   OR (u.email = 'carlos@universidad.edu'   AND r.nombre = 'GRAD_ASSISTANT')
   OR (u.email = 'gamonitor@universidad.edu' AND r.nombre = 'GRAD_ASSISTANT')
   OR (u.email = 'gamonitor@universidad.edu' AND r.nombre = 'MONITOR');

INSERT INTO periodos_academicos (id, codigo, fecha_inicio, fecha_fin, estado) VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2026-10', '2026-01-15', '2026-06-15', 'ACTIVE'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '2025-20', '2025-07-15', '2025-11-30', 'CLOSED');

INSERT INTO espacios (id, tipo, nombre, periodo_id, fecha_inicio, fecha_fin, profesor_id) VALUES
  ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'COURSE',  'Ingeniería de Software I',  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2026-01-20', '2026-06-10', '22222222-2222-2222-2222-222222222222'),
  ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'PROJECT', 'Proyecto de Investigación',  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2026-02-01', '2026-06-15', '22222222-2222-2222-2222-222222222222'),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'COURSE',  'Bases de Datos Avanzadas',   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '2026-01-20', '2026-06-10', '33333333-3333-3333-3333-333333333333');

-- Vinculaciones para caso de validación (GA=12h, Monitor=5h -> 40% de 12h = 4.8 ~ 5h OK)
INSERT INTO vinculaciones (id, usuario_id, espacio_id, rol, horas_semanales, profesor_id, activa) VALUES
  ('11111111-1111-1111-1111-11111111aaaa', '88888888-8888-8888-8888-888888888888', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'GRAD_ASSISTANT', 12, '22222222-2222-2222-2222-222222222222', TRUE),
  ('22222222-2222-2222-2222-22222222bbbb', '88888888-8888-8888-8888-888888888888', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'MONITOR', 5, '22222222-2222-2222-2222-222222222222', TRUE);
