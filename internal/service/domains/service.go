package domains

import domains2 "github.com/sunwild/domain-checker_api/pkg/domains"

type Service struct {
	repo Repo
}

type Repo interface {
	GetAllDomains() ([]*domains2.Domain, error)
	GetDomainById(domain int) (*domains2.Domain, error)
	AddDomain(domain *domains2.Domain) error
	DeleteDomainById(domain int) error
	UpdateDomainById(domain *domains2.Domain) error
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllDomains() ([]*domains2.Domain, error) {
	return s.repo.GetAllDomains()
}

func (s *Service) GetDomainById(domain int) (*domains2.Domain, error) {
	return s.repo.GetDomainById(domain)
}

func (s *Service) AddDomain(name *domains2.Domain) error {
	return s.repo.AddDomain(name)
}

func (s *Service) DeleteDomainById(domain int) error {
	return s.repo.DeleteDomainById(domain)
}

func (s *Service) UpdateDomainById(domain *domains2.Domain) error {
	return s.repo.UpdateDomainById(domain)
}
