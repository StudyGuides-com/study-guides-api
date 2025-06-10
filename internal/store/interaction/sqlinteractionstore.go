package interaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

)

type SqlInteractionStore struct {
	db *pgxpool.Pool
}

func (s *SqlInteractionStore) Interact(ctx context.Context) error {
	return nil
}






