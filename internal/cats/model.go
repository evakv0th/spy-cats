package cats

type Cat struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	YearsOfExperience int     `json:"years_of_experience"`
	Breed             string  `json:"breed"`
	Salary            float64 `json:"salary"`
}

type CreateCatRequest struct {
	Name              string  `json:"name" binding:"required,min=2,max=50"`
	YearsOfExperience int     `json:"years_of_experience" binding:"required,gte=0,lte=50"`
	Breed             string  `json:"breed" binding:"required"`
	Salary            float64 `json:"salary" binding:"required,gte=0"`
}

type UpdateSalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,gte=0"`
}
