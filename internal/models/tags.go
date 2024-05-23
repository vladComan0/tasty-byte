package models

import (
	"database/sql"
	"errors"
	"github.com/vladComan0/tasty-byte/pkg/transactions"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type TagModel struct {
	DB *sql.DB
}

func (m *TagModel) GetByRecipeID(tx transactions.Transaction, recipeID int) ([]*Tag, error) {
	var tags []*Tag

	stmt := `
		SELECT t.id, t.name
		FROM tags t INNER JOIN recipe_tags rt ON rt.tag_id = t.id
		WHERE rt.recipe_id = ?
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
		tag := &Tag{}

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) InsertIfNotExists(tx transactions.Transaction, name string) (int, error) {
	var id int
	if err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", name).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			result, err := tx.Exec("INSERT INTO tags(name) VALUES (?)", name)
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
