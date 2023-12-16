package models

import (
	"database/sql"
	"errors"
	"time"
)

type Recipe struct {
	ID              int
	Name            string
	Description     string
	Instructions    string
	PreparationTime string
	CookingTime     string
	Portions        string
	CreatedAt       time.Time
}

type RecipeModel struct {
	DB *sql.DB
}

func (m *RecipeModel) Insert(name, description, instructions, preparationTime, cookingTime, portions string) (int, error) {
	stmt := `
    INSERT INTO recipes 
        (name, description, instructions, preparation_time, cooking_time, portions, created)
    VALUES 
        (?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())
`
	result, err := m.DB.Exec(stmt, name, description, instructions, preparationTime, cookingTime, portions)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *RecipeModel) Get(id int) (*Recipe, error) {
	recipe := &Recipe{}

	stmt := `
    SELECT 
        id, 
        name, 
        description, 
        instructions, 
        preparation_time, 
        cooking_time, 
        portions, 
        created
    FROM 
        recipes 
    WHERE 
        id = ?
`

	err := m.DB.QueryRow(stmt, id).Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.Description,
		&recipe.Instructions,
		&recipe.PreparationTime,
		&recipe.CookingTime,
		&recipe.Portions,
		&recipe.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return recipe, nil
}

func (m *RecipeModel) Update(recipe *Recipe) error {
	return nil
}

func (m *RecipeModel) Delete(id int) error { // possibly not needed
	return nil
}

func (m *RecipeModel) Latest() ([]*Recipe, error) {
	return nil, nil
}
