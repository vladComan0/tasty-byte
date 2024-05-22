package models

import (
	"database/sql"
	"errors"
)

// FullIngredient abstracts away the two models for storing ingredients and their quantities/units
type FullIngredient struct {
	*Ingredient
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type Ingredient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type IngredientModel struct {
	DB *sql.DB
}

func (m *IngredientModel) InsertIfNotExists(tx *sql.Tx, name string) (int, error) {
	var id int
	if err := tx.QueryRow("SELECT id FROM ingredients WHERE name = ?", name).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			result, err := tx.Exec("INSERT INTO ingredients(name) VALUES (?)", name)
			if err != nil {
				return 0, err
			}
			id64, err := result.LastInsertId()
			if err != nil {
				return 0, err
			}
			id = int(id64)
		} else {
			return 0, err
		}
	}

	return id, nil
}
