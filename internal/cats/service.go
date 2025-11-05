package cats

import (
	"fmt"
	"spy-cats/internal/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCat(req CreateCatRequest) (int64, error) {
	ok, err := utils.CatBreedExists(req.Breed)
	if err != nil {
		return 0, fmt.Errorf("failed to validate breed: %w", err)
	}
	if !ok {
		return 0, fmt.Errorf("invalid cat breed: %s", req.Breed)
	}

	cat := Cat{
		Name:              req.Name,
		YearsOfExperience: req.YearsOfExperience,
		Breed:             req.Breed,
		Salary:            req.Salary,
	}
	return s.repo.Create(cat)
}

func (s *Service) GetAllCats() ([]Cat, error) {
	return s.repo.GetAll()
}

func (s *Service) GetCat(id int64) (*Cat, error) {
	return s.repo.GetByID(id)
}

func (s *Service) UpdateSalary(id int64, salary float64) error {
	rowsAffected, err := s.repo.UpdateSalary(id, salary)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("cat not found with id %d", id)
	}
	return nil
}

func (s *Service) DeleteCat(id int64) error {
	return s.repo.Delete(id)
}
