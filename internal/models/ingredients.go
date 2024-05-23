package models

import (
	"database/sql"
	"errors"
	"github.com/vladComan0/tasty-byte/pkg/transactions"
)

type IngredientModelInterface interface {
	GetByRecipeID(tx transactions.Transaction, recipeID int) ([]*FullIngredient, error)
	InsertIfNotExists(tx transactions.Transaction, name string) (int, error)
}

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

func (m *IngredientModel) GetByRecipeID(tx transactions.Transaction, recipeID int) ([]*FullIngredient, error) {
	var ingredients []*FullIngredient

	stmt := `
		SELECT i.id, i.name, ri.quantity, ri.unit
		FROM ingredients i INNER JOIN recipe_ingredients ri ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = ?
		`

	rows, err := tx.Query(stmt, recipeID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		ingredient := &FullIngredient{
			Ingredient: &Ingredient{},
		}

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Quantity,
			&ingredient.Unit,
		)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (m *IngredientModel) InsertIfNotExists(tx transactions.Transaction, name string) (int, error) {
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
