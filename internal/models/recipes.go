package models

import (
	"database/sql"
	"errors"
	"sort"
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
	var results []*Recipe

	stmt := `
	SELECT 
		recipes.id,
		recipes.name,
		recipes.description,
		recipes.instructions,
		recipes.preparation_time,
		recipes.cooking_time,
		recipes.portions,
		recipes.created,
		tags.id,
		tags.name
	FROM
		recipes
	LEFT JOIN 
		recipe_tags ON recipes.id=recipe_tags.recipe_id
	LEFT JOIN
		tags ON recipe_tags.tag_id=tags.id`

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

	recipes := make(map[int]*Recipe)
	for rows.Next() {
		var (
			tag     = &Tag{}
			tagID   sql.NullInt64
			tagName sql.NullString
			recipe  = &Recipe{}
		)

		err := rows.Scan(
			&recipe.ID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Instructions,
			&recipe.PreparationTime,
			&recipe.CookingTime,
			&recipe.Portions,
			&recipe.CreatedAt,
			&tagID,
			&tagName,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := recipes[recipe.ID]; !ok {
			recipe.Tags = []*Tag{}
			recipes[recipe.ID] = recipe
		}

		if tagID.Valid && tagName.Valid {
			tag.ID = int(tagID.Int64)
			tag.Name = tagName.String
			recipes[recipe.ID].Tags = append(recipes[recipe.ID].Tags, tag)
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, recipe := range recipes {
		results = append(results, recipe)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results, nil
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
