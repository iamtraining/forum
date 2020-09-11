package store

import (
	"fmt"

	"github.com/iamtraining/forum/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

type Store struct {
	entity.ThreadStore
	entity.PostStore
	entity.CommentStore
}

func NewStore(dSN string) (*Store, error) {
	db, err := sqlx.Open("postgres", dSN)
	if err != nil {
		return nil, fmt.Errorf("failure while opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error while connecting to database: %w", err)
	}

	return &Store{
		ThreadStore:  &ThreadStore{DB: db},
		PostStore:    &PostStore{DB: db},
		CommentStore: &CommentStore{DB: db},
	}, nil
}
