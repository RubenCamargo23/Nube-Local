package users

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
    Create(ctx context.Context, nombre, email, passwordHash string) (*UserResponse, error)
    AssignRole(ctx context.Context, userID uuid.UUID, rol string) error
    FindAll(ctx context.Context) ([]UserResponse, error)
    FindByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
    FindRoles(ctx context.Context, id uuid.UUID) ([]string, error)
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, nombre, email, hash string) (*UserResponse, error) {
    var u UserResponse
    err := r.db.QueryRow(ctx,
        `INSERT INTO usuarios (nombre, email, password_hash)
         VALUES ($1, $2, $3)
         RETURNING id, nombre, email, creado_en`,
        nombre, email, hash,
    ).Scan(&u.ID, &u.Nombre, &u.Email, &u.CreadoEn)
    return &u, err
}

func (r *repository) AssignRole(ctx context.Context, userID uuid.UUID, rol string) error {
    _, err := r.db.Exec(ctx,
        `INSERT INTO usuario_roles (usuario_id, rol_id)
         SELECT $1, r.id FROM roles r WHERE r.nombre = $2
         ON CONFLICT DO NOTHING`,
        userID, rol,
    )
    return err
}

func (r *repository) FindAll(ctx context.Context) ([]UserResponse, error) {
    rows, err := r.db.Query(ctx,
        `SELECT id, nombre, email, creado_en FROM usuarios ORDER BY creado_en DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []UserResponse
    for rows.Next() {
        var u UserResponse
        if err := rows.Scan(&u.ID, &u.Nombre, &u.Email, &u.CreadoEn); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
    var u UserResponse
    err := r.db.QueryRow(ctx,
        `SELECT id, nombre, email, creado_en FROM usuarios WHERE id = $1`, id,
    ).Scan(&u.ID, &u.Nombre, &u.Email, &u.CreadoEn)
    return &u, err
}

func (r *repository) FindRoles(ctx context.Context, id uuid.UUID) ([]string, error) {
    rows, err := r.db.Query(ctx,
        `SELECT r.nombre FROM roles r
         JOIN usuario_roles ur ON ur.rol_id = r.id
         WHERE ur.usuario_id = $1`, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var roles []string
    for rows.Next() {
        var rol string
        rows.Scan(&rol)
        roles = append(roles, rol)
    }
    return roles, nil
}
