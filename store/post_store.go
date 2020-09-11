package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/jmoiron/sqlx"
)

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) CreatePost(p *entity.ForumPost) error {
	if err := s.Get(p, `INSERT INTO forum_posts VALUES ($1, $2, $3, $4) RETURNING *`,
		p.ID,
		p.ThreadID,
		p.Title,
		p.Content,
	); err != nil {
		return fmt.Errorf("error while creating post: %w", err)
	}
	return nil
}

func (s *PostStore) ReadPost(id uuid.UUID) (entity.ForumPost, error) {
	var post entity.ForumPost
	if err := s.Get(&post, `SELECT * FROM forum_posts WHERE id = $1`, id); err != nil {
		return entity.ForumPost{}, fmt.Errorf("error while getting post: %w", err)
	}
	return post, nil
}

func (s *PostStore) ReadPostsByThread(threadID uuid.UUID) ([]entity.ForumPost, error) {
	var posts []entity.ForumPost
	var query = `
		SELECT
			forum_posts.*,
			COUNT(forum_comments.*) AS comms_count
		FROM forum_posts
		LEFT JOIN forum_comments ON forum_comments.post_id = forum_posts.id
		WHERE thread_id = $1
		GROUP BY forum_posts.id`
	if err := s.Select(&posts, query, threadID); err != nil {
		return []entity.ForumPost{}, fmt.Errorf("error while getting posts: %w", err)
	}
	return posts, nil
}

func (s *PostStore) UpdatePost(p *entity.ForumPost) error {
	if err := s.Get(p, `UPDATE forum_posts SET thread_id = $1, title = $2, content = $3 WHERE id = $4 RETURNING *`,
		p.ThreadID,
		p.Title,
		p.Content,
		p.ID); err != nil {
		return fmt.Errorf("error while updating post: %w", err)
	}
	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM forum_posts WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error while deleting post: %w", err)
	}
	return nil
}
