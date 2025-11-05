package missions

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateMission(m Mission) (int64, error) {
	query := `INSERT INTO missions (cat_id, name, is_complete) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRow(query, m.CatID, m.Name, m.IsComplete).Scan(&id)
	return id, err
}

func (r *Repository) CreateTarget(t Target) (int64, error) {
	query := `INSERT INTO targets (mission_id, name, country, notes, is_complete)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := r.db.QueryRow(query, t.MissionID, t.Name, t.Country, t.Notes, t.IsComplete).Scan(&id)
	return id, err
}

func (r *Repository) GetMissionByID(id int64) (*Mission, error) {
	m := Mission{}
	err := r.db.QueryRow(`SELECT id, cat_id, name, is_complete FROM missions WHERE id = $1`, id).
		Scan(&m.ID, &m.CatID, &m.Name, &m.IsComplete)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT id, mission_id, name, country, notes, is_complete FROM targets WHERE mission_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Target
		if err := rows.Scan(&t.ID, &t.MissionID, &t.Name, &t.Country, &t.Notes, &t.IsComplete); err != nil {
			return nil, err
		}
		m.Targets = append(m.Targets, t)
	}

	return &m, nil
}

func (r *Repository) DeleteMission(id int64) error {
	_, err := r.db.Exec(`DELETE FROM missions WHERE id = $1 AND cat_id IS NULL`, id)
	return err
}

func (r *Repository) MarkMissionComplete(id int64) error {
	res, err := r.db.Exec(`UPDATE missions SET is_complete = TRUE WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) UpdateTarget(t Target) error {
	query := `UPDATE targets 
	          SET notes = COALESCE($1, notes), is_complete = COALESCE($2, is_complete)
	          WHERE id = $3`
	res, err := r.db.Exec(query, t.Notes, t.IsComplete, t.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteTarget(id int64) error {
	res, err := r.db.Exec(`DELETE FROM targets WHERE id = $1 AND is_complete = FALSE`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("cannot delete completed target")
	}
	return nil
}

func (r *Repository) AddTarget(t Target) error {
	_, err := r.CreateTarget(t)
	return err
}

func (r *Repository) GetAllMissions() ([]Mission, error) {
	rows, err := r.db.Query(`SELECT id, cat_id, name, is_complete FROM missions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []Mission
	for rows.Next() {
		var m Mission
		err := rows.Scan(&m.ID, &m.CatID, &m.Name, &m.IsComplete)
		if err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func (r *Repository) AssignCat(missionID, catID int64) error {
	query := `UPDATE missions SET cat_id=$1 WHERE id=$2`
	res, err := r.db.Exec(query, catID, missionID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
