package models

import (
	"database/sql"
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
	return 0, nil
}

func (m *RecipeModel) Get(id int) (*Recipe, error) {
	return nil, nil
}

func (m *RecipeModel) Delete(id int) error { // possibly not needed
	return nil
}

func (m *RecipeModel) Update(recipe *Recipe) error {
	return nil
}

func (m *RecipeModel) Latest() ([]*Recipe, error) {
	return nil, nil
}
