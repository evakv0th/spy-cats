package missions

type Mission struct {
	ID         int64    `json:"id"`
	CatID      *int64   `json:"cat_id,omitempty"`
	Name       string   `json:"name"`
	IsComplete bool     `json:"is_complete"`
	Targets    []Target `json:"targets,omitempty"`
}

type Target struct {
	ID         int64  `json:"id"`
	MissionID  int64  `json:"mission_id"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	Notes      string `json:"notes"`
	IsComplete bool   `json:"is_complete"`
}

type CreateMissionRequest struct {
	CatID      *int64         `json:"cat_id"`
	Name       string         `json:"name" binding:"required"`
	Targets    []CreateTarget `json:"targets" binding:"required"`
	IsComplete bool           `json:"is_complete"`
}

type CreateTarget struct {
	Name       string `json:"name" binding:"required"`
	Country    string `json:"country" binding:"required"`
	Notes      string `json:"notes"`
	IsComplete bool   `json:"is_complete"`
}

type UpdateMissionRequest struct {
	IsComplete *bool `json:"is_complete"`
}

type UpdateTargetRequest struct {
	IsComplete *bool   `json:"is_complete"`
	Notes      *string `json:"notes"`
}

type CreateTargetRequest struct {
	Name       string `json:"name" binding:"required"`
	Country    string `json:"country" binding:"required"`
	Notes      string `json:"notes"`
	IsComplete bool   `json:"is_complete"`
}

type AssignCatRequest struct {
	CatID int64 `json:"cat_id" binding:"required"`
}
