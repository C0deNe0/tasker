package service

import (
	"github.com/C0deNe0/go-tasker/internal/lib/job"
	"github.com/C0deNe0/go-tasker/internal/repository"
	"github.com/C0deNe0/go-tasker/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}
