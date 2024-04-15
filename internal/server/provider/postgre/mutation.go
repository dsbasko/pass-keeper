package postgre

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dsbasko/pass-keeper/internal/model"
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
)

func (p *Postgre) CreateUser(ctx context.Context, email, passwordHash string) (resp model.User, err error) {
	defer errWrapper.PtrWithOP(&err, "postgre.Postgre.CreateUser")

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING user_id, email
	`

	row := p.db.QueryRow(ctx, query, email, passwordHash)
	if err = row.Scan(&resp.ID, &resp.Email); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return model.User{}, errWrapper.WithOP(ErrEmailExists, "row.Scan")
		}

		return model.User{}, errWrapper.WithOP(err, "row.Scan")
	}

	return resp, nil
}
