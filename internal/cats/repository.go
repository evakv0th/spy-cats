package cats

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(cat Cat) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO cats (name, years_of_experience, breed, salary)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetAll() ([]Cat, error) {
	rows, err := r.db.Query(`SELECT id, name, years_of_experience, breed, salary FROM cats ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []Cat
	for rows.Next() {
		var c Cat
		if err := rows.Scan(&c.ID, &c.Name, &c.YearsOfExperience, &c.Breed, &c.Salary); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *Repository) GetByID(id int64) (*Cat, error) {
	var c Cat
	err := r.db.QueryRow(
		`SELECT id, name, years_of_experience, breed, salary FROM cats WHERE id=$1`, id,
	).Scan(&c.ID, &c.Name, &c.YearsOfExperience, &c.Breed, &c.Salary)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (r *Repository) UpdateSalary(id int64, salary float64) (int64, error) {
	res, err := r.db.Exec(`UPDATE cats SET salary=$1 WHERE id=$2`, salary, id)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return rows, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM cats WHERE id=$1`, id)
	return err
}
