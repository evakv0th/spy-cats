package cats

// Cat represents a spy cat
type Cat struct {
	ID                int64   `json:"id" example:"1"`
	Name              string  `json:"name" example:"Whiskers"`
	YearsOfExperience int     `json:"years_of_experience" example:"5"`
	Breed             string  `json:"breed" example:"Siamese"`
	Salary            float64 `json:"salary" example:"50000.0"`
}

// CreateCatRequest represents the request to create a new cat
type CreateCatRequest struct {
	Name              string  `json:"name" binding:"required,min=2,max=50" example:"Whiskers"`
	YearsOfExperience int     `json:"years_of_experience" binding:"required,gte=0,lte=50" example:"5"`
	Breed             string  `json:"breed" binding:"required" example:"Siamese"`
	Salary            float64 `json:"salary" binding:"required,gte=0" example:"50000.0"`
}

// UpdateSalaryRequest represents the request to update a cat's salary
type UpdateSalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,gte=0" example:"60000.0"`
}
