package entity

import (
	"github.com/google/uuid"
)

type ForumThread struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}

type ForumPost struct {
	ID          uuid.UUID `db:"id"`
	ThreadID    uuid.UUID `db:"thread_id"`
	ThreadTitle string
	Title       string `db:"title"`
	Content     string `db:"content"`
	Count       int    `db:"comms_count"`
}

type ForumComment struct {
	ID      uuid.UUID `db:"id"`
	PostID  uuid.UUID `db:"post_id"`
	Content string    `db:"content"`
}

type ThreadStore interface {
	Threads() ([]ForumThread, error)
	CreateThread(t *ForumThread) error
	ReadThread(id uuid.UUID) (ForumThread, error)
	UpdateThread(t *ForumThread) error
	DeleteThread(id uuid.UUID) error
}

type PostStore interface {
	CreatePost(p *ForumPost) error
	ReadPost(id uuid.UUID) (ForumPost, error)
	ReadPostsByThread(threadID uuid.UUID) ([]ForumPost, error)
	UpdatePost(p *ForumPost) error
	DeletePost(id uuid.UUID) error
}

type CommentStore interface {
	CreateComment(c *ForumComment) error
	ReadComment(id uuid.UUID) (ForumComment, error)
	ReadCommentsByPost(postID uuid.UUID) ([]ForumComment, error)
	UpdateComment(c *ForumComment) error
	DeleteComment(id uuid.UUID) error
}

type Store interface {
	ThreadStore
	PostStore
	CommentStore
}
