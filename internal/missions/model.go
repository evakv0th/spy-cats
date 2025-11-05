package missions

// Mission represents a spy mission
type Mission struct {
	ID         int64    `json:"id" example:"1"`
	CatID      *int64   `json:"cat_id,omitempty" example:"5"`
	Name       string   `json:"name" example:"Operation Stealth"`
	IsComplete bool     `json:"is_complete" example:"false"`
	Targets    []Target `json:"targets,omitempty"`
}

// Target represents a mission target
type Target struct {
	ID         int64  `json:"id" example:"1"`
	MissionID  int64  `json:"mission_id" example:"1"`
	Name       string `json:"name" example:"Agent Smith"`
	Country    string `json:"country" example:"Russia"`
	Notes      string `json:"notes" example:"High priority target"`
	IsComplete bool   `json:"is_complete" example:"false"`
}

// CreateMissionRequest represents the request to create a new mission
type CreateMissionRequest struct {
	CatID      *int64         `json:"cat_id" example:"5"`
	Name       string         `json:"name" binding:"required" example:"Operation Stealth"`
	Targets    []CreateTarget `json:"targets" binding:"required"`
	IsComplete bool           `json:"is_complete" example:"false"`
}

// CreateTarget represents the target data when creating a mission
type CreateTarget struct {
	Name       string `json:"name" binding:"required" example:"Agent Smith"`
	Country    string `json:"country" binding:"required" example:"Russia"`
	Notes      string `json:"notes" example:"High priority target"`
	IsComplete bool   `json:"is_complete" example:"false"`
}

// UpdateMissionRequest represents the request to update a mission
type UpdateMissionRequest struct {
	IsComplete *bool `json:"is_complete" example:"true"`
}

// UpdateTargetRequest represents the request to update a target
type UpdateTargetRequest struct {
	IsComplete *bool   `json:"is_complete" example:"true"`
	Notes      *string `json:"notes" example:"Mission accomplished"`
}

// CreateTargetRequest represents the request to create a target
type CreateTargetRequest struct {
	Name       string `json:"name" binding:"required" example:"Agent Smith"`
	Country    string `json:"country" binding:"required" example:"Russia"`
	Notes      string `json:"notes" example:"High priority target"`
	IsComplete bool   `json:"is_complete" example:"false"`
}

// AssignCatRequest represents the request to assign a cat to a mission
type AssignCatRequest struct {
	CatID int64 `json:"cat_id" binding:"required" example:"5"`
}
