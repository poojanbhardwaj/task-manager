package repository

import "github.com/poojanbhardwaj/task-manager/apps/backend/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
