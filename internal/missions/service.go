package missions

import (
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateMission(req CreateMissionRequest) (*Mission, error) {
	mission := Mission{
		CatID:      req.CatID,
		Name:       req.Name,
		IsComplete: req.IsComplete,
	}
	missionID, err := s.repo.CreateMission(mission)
	if err != nil {
		return nil, err
	}

	for _, t := range req.Targets {
		target := Target{
			MissionID:  missionID,
			Name:       t.Name,
			Country:    t.Country,
			Notes:      t.Notes,
			IsComplete: t.IsComplete,
		}
		if _, err := s.repo.CreateTarget(target); err != nil {
			return nil, err
		}
	}
	return s.repo.GetMissionByID(missionID)
}

func (s *Service) DeleteMission(id int64) error {
	return s.repo.DeleteMission(id)
}

func (s *Service) MarkMissionComplete(id int64) error {
	return s.repo.MarkMissionComplete(id)
}

func (s *Service) UpdateTarget(id int64, req UpdateTargetRequest) error {
	target := Target{
		ID:         id,
		IsComplete: req.IsComplete != nil && *req.IsComplete,
	}
	if req.Notes != nil {
		target.Notes = *req.Notes
	}
	return s.repo.UpdateTarget(target)
}

func (s *Service) DeleteTarget(id int64) error {
	return s.repo.DeleteTarget(id)
}

func (s *Service) AddTarget(missionID int64, req CreateTarget) error {
	mission, err := s.repo.GetMissionByID(missionID)
	if err != nil {
		return err
	}
	if mission.IsComplete {
		return errors.New("cannot add target to completed mission")
	}
	_, err = s.repo.CreateTarget(Target{
		MissionID: missionID,
		Name:      req.Name,
		Country:   req.Country,
		Notes:     req.Notes,
	})
	return err
}

func (s *Service) GetAllMissions() ([]Mission, error) {
	return s.repo.GetAllMissions()
}

func (s *Service) GetMissionByID(id int64) (*Mission, error) {
	return s.repo.GetMissionByID(id)
}

func (s *Service) AssignCat(missionID, catID int64) error {
	return s.repo.AssignCat(missionID, catID)
}
