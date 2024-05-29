package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vladComan0/tasty-byte/internal/mocks"
	"github.com/vladComan0/tasty-byte/internal/models"
	"io"
	"net/http"
	"testing"
)

var testRecipe = &models.Recipe{
	ID:              1,
	Name:            "Test Recipe",
	Description:     "Test Description",
	Instructions:    "Test Instructions",
	PreparationTime: "30m",
	CookingTime:     "1h",
	Portions:        4,
	Ingredients: []*models.FullIngredient{
		{
			Ingredient: &models.Ingredient{
				ID:   1,
				Name: "Test Ingredient",
			},
			Quantity: 1,
			Unit:     "cup",
		},
	},
	Tags: []*models.Tag{
		{
			ID:   1,
			Name: "Test Tag",
		},
	},
}

var testRecipes = []*models.Recipe{
	{
		ID:              2,
		Name:            "Test Recipe 2",
		Description:     "Test Description 2",
		Instructions:    "Test Instructions 2",
		PreparationTime: "40m",
		CookingTime:     "1h 10m",
		Portions:        5,
		Ingredients: []*models.FullIngredient{
			{
				Ingredient: &models.Ingredient{
					ID:   2,
					Name: "Test Ingredient 2",
				},
				Quantity: 2,
				Unit:     "cup",
			},
		},
		Tags: []*models.Tag{
			{
				ID:   2,
				Name: "Test Tag 2",
			},
		},
	},
	{
		ID:              3,
		Name:            "Test Recipe 3",
		Description:     "Test Description 3",
		Instructions:    "Test Instructions 3",
		PreparationTime: "50m",
		CookingTime:     "1h 20m",
		Portions:        6,
		Ingredients: []*models.FullIngredient{
			{
				Ingredient: &models.Ingredient{
					ID:   3,
					Name: "Test Ingredient 3",
				},
				Quantity: 3,
				Unit:     "cup",
			},
		},
		Tags: []*models.Tag{
			{
				ID:   3,
				Name: "Test Tag 3",
			},
		},
	},
}

type mockRecipeInput struct {
	Name            *string                  `json:"name"`
	Description     *string                  `json:"description"`
	Instructions    *string                  `json:"instructions"`
	PreparationTime *string                  `json:"preparation_time"`
	CookingTime     *string                  `json:"cooking_time"`
	Portions        *int                     `json:"portions"`
	Ingredients     []*models.FullIngredient `json:"ingredients"`
	Tags            []*models.Tag            `json:"tags"`
}

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := newTestApplication()

	mockRecipes := mocks.NewMockRecipeModelInterface(ctrl)
	app.recipes = mockRecipes

	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name           string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Ping Successful",
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "pong",
		},
		{
			name:           "Ping Failed Due to Database Connection Error",
			mockReturnErr:  fmt.Errorf("database connection error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRecipes.EXPECT().Ping().Return(tc.mockReturnErr)

			res, err := ts.Client().Get(fmt.Sprintf("%s/ping", ts.URL))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			if tc.expectedBody != "" {
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, string(body))
			}
		})
	}
}

func TestCreateRecipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := newTestApplication()

	mockRecipes := mocks.NewMockRecipeModelInterface(ctrl)

	app.recipes = mockRecipes

	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name           string
		recipe         *models.Recipe
		mockReturnID   int
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name: "Successful Recipe Creation",
			recipe: &models.Recipe{
				Name:            "Test Recipe",
				Description:     "Test Description",
				Instructions:    "Test Instructions",
				PreparationTime: "30m",
				CookingTime:     "1h",
				Portions:        4,
				Ingredients: []*models.FullIngredient{
					{
						Ingredient: &models.Ingredient{
							ID:   1,
							Name: "Test Ingredient",
						},
						Quantity: 1,
						Unit:     "cup",
					},
				},
				Tags: []*models.Tag{
					{
						ID:   1,
						Name: "Test Tag",
					},
				},
			},
			mockReturnID:   1,
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Failed Recipe Creation Due to Server Error",
			recipe: &models.Recipe{
				Name:            "Test Recipe",
				Description:     "Test Description",
				Instructions:    "Test Instructions",
				PreparationTime: "30m",
				CookingTime:     "1h",
				Portions:        4,
				Ingredients: []*models.FullIngredient{
					{
						Ingredient: &models.Ingredient{
							Name: "Test Ingredient",
						},
						Unit: "cup",
					},
				},
				Tags: []*models.Tag{
					{
						Name: "Test Tag",
					},
				},
			},
			mockReturnID:   0,
			mockReturnErr:  fmt.Errorf("server error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRecipes.EXPECT().Insert(tc.recipe).Return(tc.mockReturnID, tc.mockReturnErr)

			if tc.mockReturnErr == nil {
				mockRecipes.EXPECT().Get(tc.mockReturnID).Return(tc.recipe, nil)
			}

			input := mockRecipeInput{
				Name:            &tc.recipe.Name,
				Description:     &tc.recipe.Description,
				Instructions:    &tc.recipe.Instructions,
				PreparationTime: &tc.recipe.PreparationTime,
				CookingTime:     &tc.recipe.CookingTime,
				Portions:        &tc.recipe.Portions,
				Ingredients:     tc.recipe.Ingredients,
				Tags:            tc.recipe.Tags,
			}

			body, err := json.Marshal(input)
			assert.NoError(t, err)

			res, err := ts.Client().Post(fmt.Sprintf("%s/v1/recipes", ts.URL), "application/json", bytes.NewBuffer(body))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestGetRecipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := newTestApplication()

	// Create a new instance of your mock
	mockRecipes := mocks.NewMockRecipeModelInterface(ctrl)

	// Replace the real RecipeModelInterface with the mock
	app.recipes = mockRecipes

	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name           string
		id             int
		mockReturn     *models.Recipe
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Recipe Found",
			id:             1,
			mockReturn:     testRecipe,
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Recipe Not Found",
			id:             2,
			mockReturn:     nil,
			mockReturnErr:  models.ErrNoRecord,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Recipe Fetch Failed Due to Server Error",
			id:             3,
			mockReturn:     nil,
			mockReturnErr:  fmt.Errorf("server error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Recipe Fetch Failed Due to Bad Request",
			id:             -1,
			mockReturn:     nil,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.id > 0 {
				mockRecipes.EXPECT().Get(tc.id).Return(tc.mockReturn, tc.mockReturnErr)
			} else {
				mockRecipes.EXPECT().Get(tc.id).Times(0)
			}

			res, err := ts.Client().Get(fmt.Sprintf("%s/v1/recipes/%d", ts.URL, tc.id))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestListRecipes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := newTestApplication()

	mockRecipes := mocks.NewMockRecipeModelInterface(ctrl)

	app.recipes = mockRecipes

	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name           string
		mockReturn     []*models.Recipe
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Recipes Found",
			mockReturn:     testRecipes,
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Recipes Not Found",
			mockReturn:     nil,
			mockReturnErr:  models.ErrNoRecord,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Server Error",
			mockReturn:     nil,
			mockReturnErr:  fmt.Errorf("server error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRecipes.EXPECT().GetAll().Return(tc.mockReturn, tc.mockReturnErr)

			res, err := ts.Client().Get(fmt.Sprintf("%s/v1/recipes", ts.URL))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestUpdateRecipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := newTestApplication()

	mockRecipes := mocks.NewMockRecipeModelInterface(ctrl)
	app.recipes = mockRecipes

	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name           string
		recipe         *models.Recipe
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name: "Successful Recipe Update",
			recipe: &models.Recipe{
				ID:              1,
				Name:            "Updated Test Recipe",
				Description:     "Updated Test Description",
				Instructions:    "Updated Test Instructions",
				PreparationTime: "35m",
				CookingTime:     "1h 5m",
				Portions:        5,
				Ingredients: []*models.FullIngredient{
					{
						Ingredient: &models.Ingredient{
							ID:   1,
							Name: "Updated Test Ingredient",
						},
						Quantity: 2,
						Unit:     "cup",
					},
				},
				Tags: []*models.Tag{
					{
						ID:   1,
						Name: "Updated Test Tag",
					},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failed Recipe Update Due to Non-Existent Recipe",
			recipe: &models.Recipe{
				ID:              999, // Non-existent ID
				Name:            "Non-existent Test Recipe",
				Description:     "Non-existent Test Description",
				Instructions:    "Non-existent Test Instructions",
				PreparationTime: "30m",
				CookingTime:     "1h",
				Portions:        4,
				Ingredients: []*models.FullIngredient{
					{
						Ingredient: &models.Ingredient{
							ID:   1,
							Name: "Non-existent Test Ingredient",
						},
						Quantity: 1,
						Unit:     "cup",
					},
				},
				Tags: []*models.Tag{
					{
						ID:   1,
						Name: "Non-existent Test Tag",
					},
				},
			},
			mockReturnErr:  models.ErrNoRecord,
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Failed Recipe Update Due to Server Error",
			recipe: &models.Recipe{
				ID:              1,
				Name:            "Test Recipe",
				Description:     "Test Description",
				Instructions:    "Test Instructions",
				PreparationTime: "30m",
				CookingTime:     "1h",
				Portions:        4,
				Ingredients: []*models.FullIngredient{
					{
						Ingredient: &models.Ingredient{
							ID:   1,
							Name: "Test Ingredient",
						},
						Quantity: 1,
						Unit:     "cup",
					},
				},
				Tags: []*models.Tag{
					{
						ID:   1,
						Name: "Test Tag",
					},
				},
			},
			mockReturnErr:  fmt.Errorf("server error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Failed Recipe Update Due to Bad Request",
			recipe: &models.Recipe{
				ID:              -1,
				Name:            "Test Recipe",
				Description:     "Test Description",
				Instructions:    "Test Instructions",
				PreparationTime: "30m",
				CookingTime:     "1h",
				Portions:        4,
				Ingredients: []*models.FullIngredient{
					{
						Ingredient: &models.Ingredient{
							Name: "Test Ingredient",
						},
						Unit: "cup",
					},
				},
				Tags: []*models.Tag{
					{
						Name: "Test Tag",
					},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.recipe.ID > 0 {
				mockRecipes.EXPECT().Get(tc.recipe.ID).Return(tc.recipe, tc.mockReturnErr)

				if tc.mockReturnErr == nil {
					mockRecipes.EXPECT().Update(tc.recipe).Return(tc.mockReturnErr)
				}
			} else {
				mockRecipes.EXPECT().Get(tc.recipe.ID).Times(0)
				mockRecipes.EXPECT().Update(tc.recipe).Times(0)
			}

			input := mockRecipeInput{
				Name:            &tc.recipe.Name,
				Description:     &tc.recipe.Description,
				Instructions:    &tc.recipe.Instructions,
				PreparationTime: &tc.recipe.PreparationTime,
				CookingTime:     &tc.recipe.CookingTime,
				Portions:        &tc.recipe.Portions,
				Ingredients:     tc.recipe.Ingredients,
				Tags:            tc.recipe.Tags,
			}

			body, err := json.Marshal(input)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/recipes/%d", ts.URL, tc.recipe.ID), bytes.NewBuffer(body))
			assert.NoError(t, err)

			res, err := ts.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}

}

func TestDeleteRecipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := newTestApplication()

	mockRecipes := mocks.NewMockRecipeModelInterface(ctrl)
	app.recipes = mockRecipes

	ts := newTestServer(app.routes())
	defer ts.Close()

	testCases := []struct {
		name           string
		id             int
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Successful Recipe Deletion",
			id:             1,
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Failed Recipe Deletion Due to Non-Existent Recipe",
			id:             999, // Non-existent ID
			mockReturnErr:  models.ErrNoRecord,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Failed Recipe Deletion Due to Server Error",
			id:             1,
			mockReturnErr:  fmt.Errorf("server error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Failed Recipe Deletion Due to Bad Request",
			id:             0, // Invalid ID
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.id > 0 {
				mockRecipes.EXPECT().Delete(tc.id).Return(tc.mockReturnErr)
			} else {
				// Expect Delete not to be called if id is 0 or less
				mockRecipes.EXPECT().Delete(tc.id).Times(0)
			}
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/recipes/%d", ts.URL, tc.id), nil)
			assert.NoError(t, err)

			res, err := ts.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}
