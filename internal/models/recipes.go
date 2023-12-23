package models

import (
	"database/sql"
	"errors"
	"time"
)

type Recipe struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	Instructions    string    `json:"instructions,omitempty"`
	PreparationTime string    `json:"preparation_time,omitempty"`
	CookingTime     string    `json:"cooking_time,omitempty"`
	Portions        int       `json:"portions,omitempty"`
	CreatedAt       time.Time `json:"-"`
}

type RecipeModel struct {
	DB *sql.DB
}

func (m *RecipeModel) Insert(name, description, instructions, preparationTime, cookingTime string, portions int) (int, error) {
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
	stmt := `
	UPDATE recipes
	SET 
		name = ?, 
		description = ?, 
		instructions = ?, 
		preparation_time = ?, 
		cooking_time = ?, 
		portions = ?
	WHERE 
		id = ?
	`
	results, err := m.DB.Exec(
		stmt,
		recipe.Name,
		recipe.Description,
		recipe.Instructions,
		recipe.PreparationTime,
		recipe.CookingTime,
		recipe.Portions,
		recipe.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *RecipeModel) Delete(id int) error { // possibly not needed
	stmt := `
	DELETE FROM recipes
	WHERE id = ?
	`
	results, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *RecipeModel) Latest() ([]*Recipe, error) {
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
    ORDER BY
		id DESC
	LIMIT 
		10
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := []*Recipe{}
	for rows.Next() {
		recipe := &Recipe{}
		err := rows.Scan(&recipe.ID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Instructions,
			&recipe.PreparationTime,
			&recipe.CookingTime,
			&recipe.Portions,
			&recipe.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}
