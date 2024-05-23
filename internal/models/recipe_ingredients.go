package models

import (
	"database/sql"
	"errors"
	"github.com/vladComan0/tasty-byte/pkg/transactions"
)

type RecipeIngredientModelInterface interface {
	Associate(tx transactions.Transaction, recipeID, ingredientID int, quantity float64, unit string) error
	DissociateNotInList(tx transactions.Transaction, recipeID int, recipeIngredients []*FullIngredient) error
	getIngredientIDsForRecipe(tx transactions.Transaction, recipeID int) ([]int, error)
	deleteRecord(tx transactions.Transaction, recipeID, ingredientID int) error
	deleteRecordsByRecipe(tx transactions.Transaction, recipeID int) error
}

type RecipeIngredient struct {
	RecipeID     int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type RecipeIngredientModel struct {
	DB *sql.DB
}

func (m *RecipeIngredientModel) Associate(tx transactions.Transaction, recipeID, ingredientID int, quantity float64, unit string) error {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM recipe_ingredients WHERE recipe_id = ? AND ingredient_id = ?)", recipeID, ingredientID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = tx.Exec("UPDATE recipe_ingredients SET quantity = ?, unit = ? WHERE recipe_id = ? AND ingredient_id = ?", quantity, unit, recipeID, ingredientID)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit) VALUES (?, ?, ?, ?)", recipeID, ingredientID, quantity, unit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *RecipeIngredientModel) DissociateNotInList(tx transactions.Transaction, recipeID int, recipeIngredients []*FullIngredient) error {
	ingredientIDs, err := m.getIngredientIDsForRecipe(tx, recipeID)
	if err != nil {
		return err
	}

	ingredientMap := make(map[int]bool)
	for _, recipeIngredient := range recipeIngredients {
		ingredientMap[recipeIngredient.ID] = true
	}

	for _, ingredientID := range ingredientIDs {
		if !ingredientMap[ingredientID] {
			if err := m.deleteRecord(tx, recipeID, ingredientID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *RecipeIngredientModel) getIngredientIDsForRecipe(tx transactions.Transaction, recipeID int) ([]int, error) {
	var ingredientIDs []int

	rows, err := tx.Query("SELECT ingredient_id FROM recipe_ingredients WHERE recipe_id = ?", recipeID)
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
		var ingredientID int
		if err := rows.Scan(&ingredientID); err != nil {
			return nil, err
		}
		ingredientIDs = append(ingredientIDs, ingredientID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ingredientIDs, nil
}

func (m *RecipeIngredientModel) deleteRecord(tx transactions.Transaction, recipeID, ingredientID int) error {
	_, err := tx.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = ? AND ingredient_id = ?", recipeID, ingredientID)
	return err
}

func (m *RecipeIngredientModel) deleteRecordsByRecipe(tx transactions.Transaction, recipeID int) error {
	_, err := tx.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = ?", recipeID)
	return err
}
