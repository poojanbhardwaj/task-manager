package service

import (
	"fmt"

	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/lib/aws"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/lib/job"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/repository"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/server"
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
	s.Job.SetAuthService(authService)

	var awsClient *aws.AWS
	if s.Config.Storage != nil {
		var err error
		awsClient, err = aws.NewAWS(s)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage: %w", err)
		}
	}

	return &Services{
		Job:      s.Job,
		Auth:     authService,
		Todo:     NewTodoService(s, repos.Todo, repos.Category, awsClient),
		Comment:  NewCommentService(s, repos.Comment, repos.Todo),
		Category: NewCategoryService(s, repos.Category),
	}, nil
}
