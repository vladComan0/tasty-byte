package main

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vladComan0/tasty-byte/internal/mocks"
	"github.com/vladComan0/tasty-byte/internal/models"
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
			name:           "Server Error",
			id:             3,
			mockReturn:     nil,
			mockReturnErr:  fmt.Errorf("server error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRecipes.EXPECT().Get(tc.id).Return(tc.mockReturn, tc.mockReturnErr)

			res, err := http.Get(fmt.Sprintf("%s/v1/recipes/%d", ts.URL, tc.id))
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

			res, err := http.Get(fmt.Sprintf("%s/v1/recipes", ts.URL))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}
