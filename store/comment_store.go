package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) CreateComment(c *entity.ForumComment) error {
	if err := s.Get(c, `INSERT INTO forum_comments VALUE ($1, $2, $3) RETURNING *`,
		c.ID,
		c.PostID,
		c.Content,
	); err != nil {
		return fmt.Errorf("error while creating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) ReadComment(id uuid.UUID) (entity.ForumComment, error) {
	var comm entity.ForumComment
	if err := s.Get(&comm, `SELECT * FROM forum_comments WHERE id = $1`,
		comm.ID,
	); err != nil {
		return entity.ForumComment{}, fmt.Errorf("error while getting comment: %w", err)
	}
	return comm, nil
}

func (s *CommentStore) ReadCommentsByPost(postID uuid.UUID) ([]entity.ForumComment, error) {
	var comms []entity.ForumComment
	if err := s.Get(&comms, `SELECT * FROM forum_comments where post_id = $1`, postID); err != nil {
		return []entity.ForumComment{}, fmt.Errorf("error while gettning comments: %w", err)
	}
	return comms, nil
}

func (s *CommentStore) UpdateComment(c *entity.ForumComment) error {
	if err := s.Get(c, `UPDATE forum_comments SET post_id = $1, post_id = $2 WHERE id = $3 RETURNING *`,
		c.PostID,
		c.Content,
		c.ID); err != nil {
		return fmt.Errorf("error while updating comment: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	_, err := s.Exec(`DELETE FROM forum_comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error while deleting comment: %w", err)
	}
	return nil
}
