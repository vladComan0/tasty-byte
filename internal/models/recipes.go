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
	Tags            []*Tag    `json:"tags,omitempty"`
}

type RecipeModel struct {
	DB             *sql.DB
	TagModel       *TagModel
	RecipeTagModel *RecipeTagModel
}

func (m *RecipeModel) Insert(name, description, instructions, preparationTime, cookingTime string, portions int, tags []*Tag) (int, error) {
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

	recipeID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	recipeID := int(recipeID64)

	for _, tag := range tags {
		tagID, err := m.TagModel.InsertIfNotExists(tag.Name)
		if err != nil {
			return 0, err
		}
		if err := m.RecipeTagModel.Associate(recipeID, tagID); err != nil {
			return 0, err
		}
	}
	return recipeID, nil
}

func (m *RecipeModel) GetAll() ([]*Recipe, error) {
	var recipes []*Recipe
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
		recipes`

	rows, err := m.DB.Query(stmt)
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
		recipe := &Recipe{}
		err := rows.Scan(
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
			return nil, err
		}

		tags, err := m.TagModel.GetByRecipeID(recipe.ID)
		if err != nil {
			return nil, err
		}

		recipe.Tags = tags

		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
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

	tags, err := m.TagModel.GetByRecipeID(recipe.ID)
	if err != nil {
		return nil, err
	}

	recipe.Tags = tags

	return recipe, nil
}

func (m *RecipeModel) Update(recipe *Recipe) error {
	existingRecipe, err := m.Get(recipe.ID)
	if err != nil {
		return err
	}

	if existingRecipe == nil {
		return ErrNoRecord
	}

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
	_, err = m.DB.Exec(
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

	for _, tag := range recipe.Tags {
		tagID, err := m.TagModel.InsertIfNotExists(tag.Name)
		if err != nil {
			return err
		}
		tag.ID = tagID

		if err := m.RecipeTagModel.Associate(recipe.ID, tag.ID); err != nil {
			return err
		}
	}

	// Delete any associations in the recipe_tags table that are not in the updated Recipe struct
	if err := m.RecipeTagModel.DissociateNotInList(recipe.ID, recipe.Tags); err != nil {
		return err
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

	if err := m.RecipeTagModel.deleteRecordsByRecipe(id); err != nil {
		return err
	}

	return nil
}
