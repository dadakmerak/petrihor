package services

import (
	"github.com/dadakmerak/petrihor/pkg/maps"
	"github.com/dadakmerak/petrihor/pkg/sqlx"
)

type Service struct {
	repo *sqlx.Repository
}

func NewService(repos *sqlx.Repository) *Service {
	return &Service{
		repo: repos,
	}
}

func (s *Service) List() (maps.ListResponses, error) {
	return s.repo.List()
}

func (s *Service) Detail() (maps.Response, error) {
	return s.repo.Detail()
}

func (s *Service) Create() error {
	return s.repo.Create()
}

func (s *Service) Update() error {
	return s.repo.Update()
}

func (s *Service) Delete() error {
	return s.repo.Delete()
}
