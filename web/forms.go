package web

import (
	"encoding/gob"
)

func init() {
	gob.Register(CreateUserForm{})
	gob.Register(LoginForm{})
	gob.Register(CreateUserForm{})
	gob.Register(CreatePostForm{})
	gob.Register(CreateCommentForm{})
}

type Errors map[string]string

type CreateUserForm struct {
	Username       string
	Password       string
	IsNotAvailable bool
	Err            Errors
}

func (f *CreateUserForm) Validate() bool {
	f.Err = Errors{}
	if f.Username == "" {
		f.Err["username"] = "please enter username"
	} else if f.IsNotAvailable {
		f.Err["username"] = "username is not available"
	}

	if f.Password == "" {
		f.Err["password"] = "please enter password"
	} else if len(f.Password) < 8 {
		f.Err["password"] = "password is too short (minimum is 8 characters)"
	}

	return len(f.Err) == 0
}

type LoginForm struct {
	Username             string
	Password             string
	IncorrectCredentials bool
	Err                  Errors
}

func (f *LoginForm) Validate() bool {
	f.Err = Errors{}

	if f.Username == "" {
		f.Err["username"] = "please enter username"
	} else if f.IncorrectCredentials {
		f.Err["warning"] = "incorrect username or password"
	}

	if f.Password == "" {
		f.Err["password"] = "please enter password"
	}

	return len(f.Err) == 0
}

type CreatePostForm struct {
	Title   string
	Content string
	Err     Errors
}

func (f *CreatePostForm) Validate() bool {
	f.Err = Errors{}
	if f.Title == `json:"title"` {
		f.Err["Title"] = "please enter a title"

	}
	if f.Content == `json:"content"` {
		f.Err["Content"] = "please enter a text"
	}

	return len(f.Err) == 0
}

type CreateThreadForm struct {
	Title       string
	Description string
	Err         Errors
}

func (f *CreateThreadForm) Validate() bool {
	f.Err = Errors{}
	if f.Title == `json:"title"` {
		f.Err["Title"] = "please enter a title"

	}
	if f.Description == `json:"description"` {
		f.Err["Description"] = "please enter a description"
	}

	return len(f.Err) == 0
}

type CreateCommentForm struct {
	Content string
	Err     Errors
}

func (f *CreateCommentForm) Validate() bool {
	f.Err = Errors{}
	if f.Content == `json:"content"` {
		f.Err["Content"] = "please enter a comment"
	}

	return len(f.Err) == 0
}
