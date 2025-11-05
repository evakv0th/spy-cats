-- +goose Up
CREATE TABLE cats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    years_of_experience INT NOT NULL CHECK (years_of_experience >= 0),
    breed TEXT NOT NULL,
    salary NUMERIC(10,2) NOT NULL CHECK (salary >= 0)
);

CREATE TABLE missions (
    id SERIAL PRIMARY KEY,
    cat_id INT REFERENCES cats(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    is_complete BOOLEAN DEFAULT FALSE
);

CREATE TABLE targets (
    id SERIAL PRIMARY KEY,
    mission_id INT NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    notes TEXT,
    is_complete BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_missions_cat_id ON missions(cat_id);

CREATE INDEX idx_targets_mission_id ON targets(mission_id);

CREATE INDEX idx_missions_is_complete ON missions(is_complete);

CREATE INDEX idx_targets_is_complete ON targets(is_complete);

CREATE INDEX idx_targets_country ON targets(country);

-- +goose Down
DROP INDEX IF EXISTS idx_targets_country;
DROP INDEX IF EXISTS idx_targets_is_complete;
DROP INDEX IF EXISTS idx_missions_is_complete;
DROP INDEX IF EXISTS idx_targets_mission_id;
DROP INDEX IF EXISTS idx_missions_cat_id;

DROP TABLE IF EXISTS targets;
DROP TABLE IF EXISTS missions;
DROP TABLE IF EXISTS cats;
