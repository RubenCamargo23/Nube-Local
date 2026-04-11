package auth

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/google/uuid"
)

type Repository interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindRolesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type User struct {
    ID           uuid.UUID
    Email        string
    PasswordHash string
}

type repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
    return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var u User
    err := r.db.QueryRow(ctx,
        `SELECT id, email, password_hash FROM usuarios WHERE email = $1`, email,
    ).Scan(&u.ID, &u.Email, &u.PasswordHash)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *repository) FindRolesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
    rows, err := r.db.Query(ctx,
        `SELECT r.nombre FROM roles r
         JOIN usuario_roles ur ON ur.rol_id = r.id
         WHERE ur.usuario_id = $1`, userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    roles := []string{}
    for rows.Next() {
        var role string
        if err := rows.Scan(&role); err != nil {
            return nil, err
        }
        roles = append(roles, role)
    }
    return roles, nil
}
