package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iamtraining/forum-golang/entity"
	"github.com/jmoiron/sqlx"
)

type ThreadStore struct {
	*sqlx.DB
}

func (s *ThreadStore) Threads() ([]entity.ForumThread, error) {
	var threads []entity.ForumThread
	if err := s.Select(&threads, `SELECT * FROM forum_threads`); err != nil {
		return []entity.ForumThread{}, fmt.Errorf("error while getting threads: %w", err)
	}
	return threads, nil
}

func (s *ThreadStore) CreateThread(t *entity.ForumThread) error {
	if err := s.Get(t, `INSERT INTO forum_threads VALUES ($1, $2, $3) RETURNING *`,
		t.ID,
		t.Title,
		t.Description); err != nil {
		return fmt.Errorf("error while creating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) ReadThread(id uuid.UUID) (entity.ForumThread, error) {
	var thread entity.ForumThread
	if err := s.Get(&thread, `SELECT * FROM forum_threads WHERE id = $1`, id); err != nil {
		return entity.ForumThread{}, fmt.Errorf("error while getting thread: %w", err)
	}
	return thread, nil
}

func (s *ThreadStore) UpdateThread(t *entity.ForumThread) error {
	if err := s.Get(t, `UPDATE forum_threads SET title = $1, description = $2 WHERE id = $3 RETURNING *`,
		t.Title,
		t.Description,
		t.ID); err != nil {
		return fmt.Errorf("error while updating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM forum_threads WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error while deleting thread: %w", err)
	}
	return nil
}
