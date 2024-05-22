package models

import "database/sql"

type RecipeIngredient struct {
	RecipeID     int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type RecipeIngredientModel struct {
	DB *sql.DB
}

func (m *RecipeIngredientModel) Associate(tx *sql.Tx, recipeID, ingredientID int, quantity float64, unit string) error {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM recipe_ingredients WHERE recipe_id = ? AND ingredient_id = ?)", recipeID, ingredientID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = tx.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit) VALUES (?, ?, ?, ?)", recipeID, ingredientID, quantity, unit)
	if err != nil {
		return err
	}

	return nil
}
