package models

import (
	"database/sql"
	"errors"
	"github.com/vladComan0/tasty-byte/pkg/transactions"
)

type RecipeTag struct {
	RecipeID int `json:"recipe_id"`
	TagID    int `json:"tag_id"`
}

type RecipeTagModel struct {
	DB *sql.DB
}

func (m *RecipeTagModel) Associate(tx transactions.Transaction, recipeID, tagID int) error {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM recipe_tags WHERE recipe_id = ? AND tag_id = ?)", recipeID, tagID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = tx.Exec("INSERT INTO recipe_tags (recipe_id, tag_id) VALUES (?, ?)", recipeID, tagID)
	if err != nil {
		return err
	}

	return nil
}

func (m *RecipeTagModel) DissociateNotInList(tx transactions.Transaction, recipeID int, recipeTags []*Tag) error {
	tagIDs, err := m.getTagIDsForRecipe(tx, recipeID)
	if err != nil {
		return err
	}

	tagMap := make(map[int]bool)
	for _, recipeTag := range recipeTags {
		tagMap[recipeTag.ID] = true
	}

	for _, tagID := range tagIDs {
		if !tagMap[tagID] {
			if err := m.deleteRecord(tx, recipeID, tagID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *RecipeTagModel) getTagIDsForRecipe(tx transactions.Transaction, recipeID int) ([]int, error) {
	var tagIDs []int

	rows, err := tx.Query("SELECT tag_id FROM recipe_tags WHERE recipe_id = ?", recipeID)
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
		var tagID int
		if err := rows.Scan(&tagID); err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tagIDs, nil
}

func (m *RecipeTagModel) deleteRecord(tx transactions.Transaction, recipeID, tagID int) error {
	_, err := tx.Exec("DELETE FROM recipe_tags WHERE recipe_id = ? AND tag_id = ?", recipeID, tagID)
	return err
}

func (m *RecipeTagModel) deleteRecordsByRecipe(tx transactions.Transaction, recipeID int) error {
	_, err := tx.Exec("DELETE FROM recipe_tags WHERE recipe_id = ?", recipeID)
	return err
}
