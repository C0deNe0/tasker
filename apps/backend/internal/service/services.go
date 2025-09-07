package service

import (
	"github.com/C0deNe0/go-tasker/internal/lib/job"
	"github.com/C0deNe0/go-tasker/internal/repository"
	"github.com/C0deNe0/go-tasker/internal/server"
)

type Services struct {
	Auth     *AuthService
	Job      *job.JobService
	Todo     *TodoService
	Comment  *CommentService
	Category *CategoryService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:      s.Job,
		Auth:     authService,
		Todo:     NewTodoService(s, repos.Todo, repos.Category),
		Comment:  NewCommentService(s, repos.Comment, repos.Todo),
		Category: NewCategoryService(s, repos.Category),
	}, nil
}
