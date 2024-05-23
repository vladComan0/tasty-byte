package models

import (
	"database/sql"
	"errors"
	"github.com/vladComan0/tasty-byte/pkg/transactions"
	"log"
	"sort"
	"time"
)

type RecipeModelInterface interface {
	Ping() error
	Insert(name, description, instructions, preparationTime, cookingTime string, portions int, ingredients []*FullIngredient, tags []*Tag) (int, error)
	GetAll() ([]*Recipe, error)
	GetWithTx(tx transactions.Transaction, id int) (*Recipe, error)
	Get(id int) (*Recipe, error)
	Update(recipe *Recipe) error
	Delete(id int) error
}

type Recipe struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	Instructions    string            `json:"instructions,omitempty"`
	PreparationTime string            `json:"preparation_time,omitempty"`
	CookingTime     string            `json:"cooking_time,omitempty"`
	Portions        int               `json:"portions,omitempty"`
	CreatedAt       time.Time         `json:"-"`
	Ingredients     []*FullIngredient `json:"ingredients,omitempty"`
	Tags            []*Tag            `json:"tags,omitempty"`
}

type RecipeModel struct {
	DB                    *sql.DB
	IngredientModel       IngredientModelInterface
	RecipeIngredientModel RecipeIngredientModelInterface
	TagModel              TagModelInterface
	RecipeTagModel        RecipeTagModelInterface
}

func (m *RecipeModel) Ping() error {
	return m.DB.Ping()
}

func (m *RecipeModel) Insert(name, description, instructions, preparationTime, cookingTime string, portions int, ingredients []*FullIngredient, tags []*Tag) (int, error) {
	var recipeID int
	err := transactions.WithTransaction(m.DB, func(tx transactions.Transaction) error {
		stmt := `
		INSERT INTO recipes 
			(name, description, instructions, preparation_time, cooking_time, portions, created)
		VALUES 
			(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())
		`
		result, err := tx.Exec(stmt, name, description, instructions, preparationTime, cookingTime, portions)
		if err != nil {
			return err
		}

		recipeID64, err := result.LastInsertId()
		if err != nil {
			return err
		}
		recipeID = int(recipeID64)

		// Must be updated to use batch inserts or reduce the number of SQL inserts through another method
		for _, ingredient := range ingredients {
			ingredientID, err := m.IngredientModel.InsertIfNotExists(tx, ingredient.Name)
			if err != nil {
				return err
			}
			if err := m.RecipeIngredientModel.Associate(tx, recipeID, ingredientID, ingredient.Quantity, ingredient.Unit); err != nil {
				return err
			}
		}
		
		// Must be updated to use batch inserts or reduce the number of SQL inserts through another method
		for _, tag := range tags {
			tagID, err := m.TagModel.InsertIfNotExists(tx, tag.Name)
			if err != nil {
				return err
			}
			if err := m.RecipeTagModel.Associate(tx, recipeID, tagID); err != nil {
				return err
			}
		}

		return nil
	})

	return recipeID, err
}

func (m *RecipeModel) GetAll() ([]*Recipe, error) {
	var results []*Recipe
	recipes := make(map[int]*Recipe)

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
		ingredients.id,
		ingredients.name,
		recipe_ingredients.quantity,
		recipe_ingredients.unit,
		tags.id,
		tags.name
	FROM
		recipes
	LEFT JOIN
		recipe_ingredients ON recipes.id = recipe_ingredients.recipe_id
	LEFT JOIN
		ingredients ON recipe_ingredients.ingredient_id = ingredients.id
	LEFT JOIN 
		recipe_tags ON recipes.id = recipe_tags.recipe_id
	LEFT JOIN
		tags ON recipe_tags.tag_id = tags.id`

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
		var (
			ingredientID       sql.NullInt64
			ingredientName     sql.NullString
			ingredientQuantity sql.NullFloat64
			ingredientUnit     sql.NullString
			tagID              sql.NullInt64
			tagName            sql.NullString
			recipe             = &Recipe{}
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
			&ingredientID,
			&ingredientName,
			&ingredientQuantity,
			&ingredientUnit,
			&tagID,
			&tagName,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := recipes[recipe.ID]; !exists {
			recipe.Ingredients = []*FullIngredient{}
			recipe.Tags = []*Tag{}
			recipes[recipe.ID] = recipe
		}

		if ingredientID.Valid && ingredientName.Valid && ingredientQuantity.Valid && ingredientUnit.Valid {
			ingredient := &FullIngredient{
				Ingredient: &Ingredient{
					ID:   int(ingredientID.Int64),
					Name: ingredientName.String,
				},
				Quantity: ingredientQuantity.Float64,
				Unit:     ingredientUnit.String,
			}
			recipes[recipe.ID].Ingredients = append(recipes[recipe.ID].Ingredients, ingredient)
		}

		if tagID.Valid && tagName.Valid {
			tag := &Tag{
				ID:   int(tagID.Int64),
				Name: tagName.String,
			}
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

func (m *RecipeModel) GetWithTx(tx transactions.Transaction, id int) (*Recipe, error) {
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

	err := tx.QueryRow(stmt, id).Scan(
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

	ingredients, err := m.IngredientModel.GetByRecipeID(tx, recipe.ID)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients = ingredients

	tags, err := m.TagModel.GetByRecipeID(tx, recipe.ID)
	if err != nil {
		return nil, err
	}

	recipe.Tags = tags

	return recipe, nil
}

func (m *RecipeModel) Get(id int) (*Recipe, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("could not rollback %v", err)
		}
	}()

	return m.GetWithTx(tx, id)
}

func (m *RecipeModel) Update(recipe *Recipe) error {
	return transactions.WithTransaction(m.DB, func(tx transactions.Transaction) error {
		existingRecipe, err := m.GetWithTx(tx, recipe.ID)
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
		_, err = tx.Exec(
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

		// Must be updated to use batch inserts or reduce the number of SQL inserts through another method
		for _, ingredient := range recipe.Ingredients {
			ingredientID, err := m.IngredientModel.InsertIfNotExists(tx, ingredient.Name)
			if err != nil {
				return err
			}
			ingredient.ID = ingredientID

			if err := m.RecipeIngredientModel.Associate(tx, recipe.ID, ingredient.ID, ingredient.Quantity, ingredient.Unit); err != nil {
				return err
			}
		}

		// Must be updated to use batch inserts or reduce the number of SQL inserts through another method
		for _, tag := range recipe.Tags {
			tagID, err := m.TagModel.InsertIfNotExists(tx, tag.Name)
			if err != nil {
				return err
			}
			tag.ID = tagID

			if err := m.RecipeTagModel.Associate(tx, recipe.ID, tag.ID); err != nil {
				return err
			}
		}

		// Delete any associations in the recipe_ingredients table that are not in the updated Recipe struct
		if err := m.RecipeIngredientModel.DissociateNotInList(tx, recipe.ID, recipe.Ingredients); err != nil {
			return err
		}

		// Delete any associations in the recipe_tags table that are not in the updated Recipe struct
		if err := m.RecipeTagModel.DissociateNotInList(tx, recipe.ID, recipe.Tags); err != nil {
			return err
		}

		return nil
	})
}

func (m *RecipeModel) Delete(id int) error {
	return transactions.WithTransaction(m.DB, func(tx transactions.Transaction) error {
		stmt := `
		DELETE FROM recipes
		WHERE id = ?
		`
		results, err := tx.Exec(stmt, id)
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

		if err := m.RecipeIngredientModel.deleteRecordsByRecipe(tx, id); err != nil {
			return err
		}

		if err := m.RecipeTagModel.deleteRecordsByRecipe(tx, id); err != nil {
			return err
		}

		return nil
	})
}
