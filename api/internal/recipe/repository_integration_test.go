package recipe

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepositoryCreate(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)
	imageURL := "https://images.example.com/tom-yum.jpg"

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		ImageURL:     &imageURL,
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
		Ingredients:  []RecipeIngredient{{Description: "2 cups stock"}, {Description: "Prawns"}},
		Instructions: []RecipeInstruction{{Description: "Bring the stock to a simmer."}},
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Positive(t, created.ID)
	assert.Equal(t, imageURL, *created.ImageURL)
	assert.Equal(t, 0.0, created.AverageRating)
	assert.False(t, created.CreatedAt.IsZero())
	assert.False(t, created.UpdatedAt.IsZero())
	require.Len(t, created.Ingredients, 2)
	require.Len(t, created.Instructions, 1)

	for _, ingredient := range created.Ingredients {
		assert.Positive(t, ingredient.ID)
		assert.Equal(t, created.ID, ingredient.RecipeID)
		assert.False(t, ingredient.CreatedAt.IsZero())
	}
	for _, instruction := range created.Instructions {
		assert.Positive(t, instruction.ID)
		assert.Equal(t, created.ID, instruction.RecipeID)
		assert.False(t, instruction.CreatedAt.IsZero())
	}

	var persisted Recipe
	require.NoError(t, db.
		Preload("Creator").
		Preload("Ingredients.Recipe").
		Preload("Instructions.Recipe").
		First(&persisted, created.ID).Error)
	assert.Equal(t, creatorID, persisted.Creator.ID)
	for _, ingredient := range persisted.Ingredients {
		assert.Equal(t, persisted.ID, ingredient.Recipe.ID)
	}
	for _, instruction := range persisted.Instructions {
		assert.Equal(t, persisted.ID, instruction.Recipe.ID)
	}
	persisted.Creator = created.Creator
	for index := range persisted.Ingredients {
		persisted.Ingredients[index].Recipe = created.Ingredients[index].Recipe
	}
	for index := range persisted.Instructions {
		persisted.Instructions[index].Recipe = created.Instructions[index].Recipe
	}
	assert.Equal(t, *created, persisted)
}

func TestRepositoryCreateAllowsEmptyChildrenAndNilImageURL(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Plain rice",
		Description:  "Steamed rice.",
		DifficultyID: "easy",
		DurationID:   "10m",
		CreatorID:    creatorID,
	})

	require.NoError(t, err)
	assert.Nil(t, created.ImageURL)
	assert.Empty(t, created.Ingredients)
	assert.Empty(t, created.Instructions)

	var ingredientCount, instructionCount int64
	require.NoError(t, db.Model(&RecipeIngredient{}).Where("recipe_id = ?", created.ID).Count(&ingredientCount).Error)
	require.NoError(t, db.Model(&RecipeInstruction{}).Where("recipe_id = ?", created.ID).Count(&instructionCount).Error)
	assert.Zero(t, ingredientCount)
	assert.Zero(t, instructionCount)
}

func TestRepositoryCreateRollsBackWhenChildInsertFails(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)
	require.NoError(t, db.Exec("ALTER TABLE recipe_ingredients ADD CONSTRAINT recipe_ingredients_description_check CHECK (description <> 'invalid')").Error)

	_, err := repo.Create(context.Background(), Recipe{
		Name:         "Will fail",
		Description:  "This must not persist.",
		DifficultyID: "easy",
		DurationID:   "10m",
		CreatorID:    creatorID,
		Ingredients:  []RecipeIngredient{{Description: "invalid"}},
	})

	require.Error(t, err)
	var count int64
	require.NoError(t, db.Model(&Recipe{}).Where("name = ?", "Will fail").Count(&count).Error)
	assert.Zero(t, count)
}

func TestRepositoryReplace(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
		Ingredients:  []RecipeIngredient{{Description: "2 cups stock"}},
		Instructions: []RecipeInstruction{{Description: "Bring the stock to a simmer."}},
	})
	require.NoError(t, err)
	newImageURL := "https://images.example.com/tom-yum-2.jpg"

	replaced, err := repo.Replace(context.Background(), Recipe{
		ID:           created.ID,
		Name:         "Tom yum soup, revisited",
		Description:  "An even brighter, spicier Thai soup.",
		ImageURL:     &newImageURL,
		DifficultyID: "hard",
		DurationID:   "60m",
		Ingredients:  []RecipeIngredient{{Description: "3 cups stock"}, {Description: "Lemongrass"}},
		Instructions: []RecipeInstruction{{Description: "Simmer harder."}},
	})

	require.NoError(t, err)
	require.NotNil(t, replaced)
	assert.Equal(t, created.ID, replaced.ID)
	assert.Equal(t, "Tom yum soup, revisited", replaced.Name)
	assert.Equal(t, "An even brighter, spicier Thai soup.", replaced.Description)
	require.NotNil(t, replaced.ImageURL)
	assert.Equal(t, newImageURL, *replaced.ImageURL)
	assert.Equal(t, "hard", replaced.DifficultyID)
	assert.Equal(t, "60m", replaced.DurationID)
	assert.Equal(t, creatorID, replaced.CreatorID)
	require.Len(t, replaced.Ingredients, 2)
	require.Len(t, replaced.Instructions, 1)
	for _, ingredient := range replaced.Ingredients {
		assert.Positive(t, ingredient.ID)
		assert.NotContains(t, []int{created.Ingredients[0].ID}, ingredient.ID)
	}

	var ingredientCount int64
	require.NoError(t, db.Model(&RecipeIngredient{}).Where("id = ?", created.Ingredients[0].ID).Count(&ingredientCount).Error)
	assert.Zero(t, ingredientCount)
}

func TestRepositoryReplaceReturnsRecipeNotFoundWhenMissing(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)

	replaced, err := repo.Replace(context.Background(), Recipe{
		ID:           404,
		Name:         "Ghost",
		Description:  "Does not exist.",
		DifficultyID: "easy",
		DurationID:   "10m",
	})

	assert.Nil(t, replaced)
	assert.ErrorIs(t, err, ErrRecipeNotFound)
}

func TestRepositoryFindByID(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
		Ingredients:  []RecipeIngredient{{Description: "2 cups stock"}},
		Instructions: []RecipeInstruction{{Description: "Bring the stock to a simmer."}},
	})
	require.NoError(t, err)

	t.Run("existing recipe", func(t *testing.T) {
		found, err := repo.FindByID(context.Background(), created.ID)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "medium", found.Difficulty.ID)
		assert.Equal(t, "30m", found.Duration.ID)
		assert.Equal(t, creatorID, found.Creator.ID)
		require.Len(t, found.Ingredients, 1)
		require.Len(t, found.Instructions, 1)
	})

	t.Run("missing recipe", func(t *testing.T) {
		_, err := repo.FindByID(context.Background(), created.ID+1000)

		assert.ErrorIs(t, err, ErrRecipeNotFound)
	})

	t.Run("soft-deleted recipe", func(t *testing.T) {
		require.NoError(t, db.Delete(&Recipe{}, created.ID).Error)

		_, err := repo.FindByID(context.Background(), created.ID)

		assert.ErrorIs(t, err, ErrRecipeNotFound)
	})
}

func TestRepositoryFindByIDReturnsRatingTotal(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	t.Run("no ratings yet", func(t *testing.T) {
		found, err := repo.FindByID(context.Background(), created.ID)

		require.NoError(t, err)
		assert.Zero(t, found.RatingTotal)
	})

	t.Run("after ratings are recorded", func(t *testing.T) {
		require.NoError(t, repo.Rate(context.Background(), createCreator(t, db), created.ID, 4))
		require.NoError(t, repo.Rate(context.Background(), createCreator(t, db), created.ID, 5))

		found, err := repo.FindByID(context.Background(), created.ID)

		require.NoError(t, err)
		assert.EqualValues(t, 2, found.RatingTotal)
	})
}

func TestRepositoryIsFavorite(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	userID := createCreator(t, db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	t.Run("not favorited", func(t *testing.T) {
		isFavorite, err := repo.IsFavorite(context.Background(), userID, created.ID)

		require.NoError(t, err)
		assert.False(t, isFavorite)
	})

	t.Run("favorited", func(t *testing.T) {
		require.NoError(t, repo.Favorite(context.Background(), userID, created.ID))

		isFavorite, err := repo.IsFavorite(context.Background(), userID, created.ID)

		require.NoError(t, err)
		assert.True(t, isFavorite)
	})
}

func TestRepositoryList(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	userID := createCreator(t, db)
	creatorID := createCreator(t, db)

	favorited, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	notFavorited, err := repo.Create(context.Background(), Recipe{
		Name:         "Plain rice",
		Description:  "Steamed rice.",
		DifficultyID: "easy",
		DurationID:   "10m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	require.NoError(t, repo.Favorite(context.Background(), userID, favorited.ID))

	t.Run("marks isFavorite per recipe for the caller", func(t *testing.T) {
		recipes, total, err := repo.List(context.Background(), userID, GetRecipesQuery{Pagination: Pagination{Page: 1, Limit: 100}})

		require.NoError(t, err)
		assert.EqualValues(t, 2, total)
		require.Len(t, recipes, 2)

		isFavoriteByID := make(map[int]bool, len(recipes))
		for _, recipe := range recipes {
			isFavoriteByID[recipe.ID] = recipe.IsFavorite
		}
		assert.True(t, isFavoriteByID[favorited.ID])
		assert.False(t, isFavoriteByID[notFavorited.ID])
	})

	t.Run("favorite filter returns only the caller's favorited recipes", func(t *testing.T) {
		recipes, total, err := repo.List(context.Background(), userID, GetRecipesQuery{
			Pagination: Pagination{Page: 1, Limit: 100},
			Favorite:   boolPtr(true),
		})

		require.NoError(t, err)
		assert.EqualValues(t, 1, total)
		require.Len(t, recipes, 1)
		assert.Equal(t, favorited.ID, recipes[0].ID)
		assert.True(t, recipes[0].IsFavorite)
	})

	t.Run("favorite=false filter returns only the caller's non-favorited recipes", func(t *testing.T) {
		recipes, total, err := repo.List(context.Background(), userID, GetRecipesQuery{
			Pagination: Pagination{Page: 1, Limit: 100},
			Favorite:   boolPtr(false),
		})

		require.NoError(t, err)
		assert.EqualValues(t, 1, total)
		require.Len(t, recipes, 1)
		assert.Equal(t, notFavorited.ID, recipes[0].ID)
		assert.False(t, recipes[0].IsFavorite)
	})

	t.Run("favorite filter is scoped to the requesting user", func(t *testing.T) {
		otherUserID := createCreator(t, db)

		recipes, total, err := repo.List(context.Background(), otherUserID, GetRecipesQuery{
			Pagination: Pagination{Page: 1, Limit: 100},
			Favorite:   boolPtr(true),
		})

		require.NoError(t, err)
		assert.Zero(t, total)
		assert.Empty(t, recipes)
	})

	t.Run("favorite=false filter is scoped to the requesting user", func(t *testing.T) {
		otherUserID := createCreator(t, db)

		recipes, total, err := repo.List(context.Background(), otherUserID, GetRecipesQuery{
			Pagination: Pagination{Page: 1, Limit: 100},
			Favorite:   boolPtr(false),
		})

		require.NoError(t, err)
		assert.EqualValues(t, 2, total)
		require.Len(t, recipes, 2)
	})

	t.Run("includes rating total per recipe", func(t *testing.T) {
		require.NoError(t, repo.Rate(context.Background(), createCreator(t, db), favorited.ID, 5))

		recipes, _, err := repo.List(context.Background(), userID, GetRecipesQuery{Pagination: Pagination{Page: 1, Limit: 100}})

		require.NoError(t, err)
		ratingTotalByID := make(map[int]int64, len(recipes))
		for _, recipe := range recipes {
			ratingTotalByID[recipe.ID] = recipe.RatingTotal
		}
		assert.EqualValues(t, 1, ratingTotalByID[favorited.ID])
		assert.Zero(t, ratingTotalByID[notFavorited.ID])
	})

	t.Run("paginates with page and limit", func(t *testing.T) {
		firstPage, total, err := repo.List(context.Background(), userID, GetRecipesQuery{Pagination: Pagination{Page: 1, Limit: 1}})
		require.NoError(t, err)
		assert.EqualValues(t, 2, total)
		require.Len(t, firstPage, 1)

		secondPage, total, err := repo.List(context.Background(), userID, GetRecipesQuery{Pagination: Pagination{Page: 2, Limit: 1}})
		require.NoError(t, err)
		assert.EqualValues(t, 2, total)
		require.Len(t, secondPage, 1)
		assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID)

		thirdPage, total, err := repo.List(context.Background(), userID, GetRecipesQuery{Pagination: Pagination{Page: 3, Limit: 1}})
		require.NoError(t, err)
		assert.EqualValues(t, 2, total)
		assert.Empty(t, thirdPage)
	})
}

func TestRepositoryDelete(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	t.Run("existing recipe", func(t *testing.T) {
		require.NoError(t, repo.Delete(context.Background(), created.ID))

		var deletedAt gorm.DeletedAt
		require.NoError(t, db.Unscoped().Model(&Recipe{}).Where("id = ?", created.ID).Select("deleted_at").Scan(&deletedAt).Error)
		assert.True(t, deletedAt.Valid)

		_, err := repo.FindByID(context.Background(), created.ID)
		assert.ErrorIs(t, err, ErrRecipeNotFound)
	})

	t.Run("missing recipe", func(t *testing.T) {
		err := repo.Delete(context.Background(), created.ID+1000)

		assert.ErrorIs(t, err, ErrRecipeNotFound)
	})

	t.Run("already soft-deleted recipe", func(t *testing.T) {
		err := repo.Delete(context.Background(), created.ID)

		assert.ErrorIs(t, err, ErrRecipeNotFound)
	})
}

func TestRepositoryHasActiveReferences(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)

	t.Run("active difficulty and duration", func(t *testing.T) {
		active, err := repo.HasActiveReferences(context.Background(), "easy", "10m")

		require.NoError(t, err)
		assert.True(t, active)
	})

	t.Run("missing difficulty", func(t *testing.T) {
		active, err := repo.HasActiveReferences(context.Background(), "missing", "10m")

		require.NoError(t, err)
		assert.False(t, active)
	})

	t.Run("missing duration", func(t *testing.T) {
		active, err := repo.HasActiveReferences(context.Background(), "easy", "missing")

		require.NoError(t, err)
		assert.False(t, active)
	})

	t.Run("soft-deleted difficulty", func(t *testing.T) {
		require.NoError(t, db.Exec("UPDATE difficulties SET deleted_at = now() WHERE id = ?", "hard").Error)

		active, err := repo.HasActiveReferences(context.Background(), "hard", "10m")

		require.NoError(t, err)
		assert.False(t, active)
	})

	t.Run("soft-deleted duration", func(t *testing.T) {
		require.NoError(t, db.Exec("UPDATE durations SET deleted_at = now() WHERE id = ?", "60m").Error)

		active, err := repo.HasActiveReferences(context.Background(), "easy", "60m")

		require.NoError(t, err)
		assert.False(t, active)
	})
}

func TestRepositoryFavorite(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	userID := createCreator(t, db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	t.Run("new favorite", func(t *testing.T) {
		require.NoError(t, repo.Favorite(context.Background(), userID, created.ID))

		var count int64
		require.NoError(t, db.Model(&UserFavorite{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).Count(&count).Error)
		assert.EqualValues(t, 1, count)
	})

	t.Run("already favorited recipe", func(t *testing.T) {
		require.NoError(t, repo.Favorite(context.Background(), userID, created.ID))

		var count int64
		require.NoError(t, db.Model(&UserFavorite{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).Count(&count).Error)
		assert.EqualValues(t, 1, count)
	})

	t.Run("re-favoriting an unfavorited recipe undoes the soft delete", func(t *testing.T) {
		require.NoError(t, db.Model(&UserFavorite{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).Update("deleted_at", "now()").Error)

		require.NoError(t, repo.Favorite(context.Background(), userID, created.ID))

		var deletedAt gorm.DeletedAt
		require.NoError(t, db.Unscoped().Model(&UserFavorite{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).Select("deleted_at").Scan(&deletedAt).Error)
		assert.False(t, deletedAt.Valid)
	})

	t.Run("missing recipe", func(t *testing.T) {
		err := repo.Favorite(context.Background(), userID, created.ID+1000)

		assert.Error(t, err)
	})
}

func TestRepositoryUnfavorite(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	userID := createCreator(t, db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	t.Run("existing favorite", func(t *testing.T) {
		require.NoError(t, repo.Favorite(context.Background(), userID, created.ID))

		require.NoError(t, repo.Unfavorite(context.Background(), userID, created.ID))

		var count int64
		require.NoError(t, db.Unscoped().Model(&UserFavorite{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).Count(&count).Error)
		assert.Zero(t, count)
	})

	t.Run("missing favorite", func(t *testing.T) {
		err := repo.Unfavorite(context.Background(), userID, created.ID)

		assert.NoError(t, err)
	})
}

func TestRepositoryRate(t *testing.T) {
	db := newIntegrationDB(t)
	repo := NewRepository(db)
	creatorID := createCreator(t, db)

	created, err := repo.Create(context.Background(), Recipe{
		Name:         "Tom yum soup",
		Description:  "A bright, spicy Thai soup.",
		DifficultyID: "medium",
		DurationID:   "30m",
		CreatorID:    creatorID,
	})
	require.NoError(t, err)

	t.Run("new rating", func(t *testing.T) {
		userID := createCreator(t, db)

		require.NoError(t, repo.Rate(context.Background(), userID, created.ID, 4))

		var count int64
		require.NoError(t, db.Model(&RecipeRating{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).Count(&count).Error)
		assert.EqualValues(t, 1, count)

		var recipe Recipe
		require.NoError(t, db.First(&recipe, created.ID).Error)
		assert.Equal(t, 4.0, recipe.AverageRating)
	})

	t.Run("already rated recipe does not record a duplicate or change the score", func(t *testing.T) {
		userID := createCreator(t, db)
		require.NoError(t, repo.Rate(context.Background(), userID, created.ID, 2))

		require.NoError(t, repo.Rate(context.Background(), userID, created.ID, 5))

		var score float64
		require.NoError(t, db.Model(&RecipeRating{}).Where("user_id = ? AND recipe_id = ?", userID, created.ID).
			Select("score").Scan(&score).Error)
		assert.Equal(t, 2.0, score)
	})

	t.Run("recomputes the recipe average across all raters", func(t *testing.T) {
		recipe, err := repo.Create(context.Background(), Recipe{
			Name:         "Green curry",
			Description:  "A rich, coconutty curry.",
			DifficultyID: "medium",
			DurationID:   "30m",
			CreatorID:    creatorID,
		})
		require.NoError(t, err)

		require.NoError(t, repo.Rate(context.Background(), createCreator(t, db), recipe.ID, 3))
		require.NoError(t, repo.Rate(context.Background(), createCreator(t, db), recipe.ID, 5))

		var updated Recipe
		require.NoError(t, db.First(&updated, recipe.ID).Error)
		assert.Equal(t, 4.0, updated.AverageRating)
	})

	t.Run("missing recipe", func(t *testing.T) {
		err := repo.Rate(context.Background(), createCreator(t, db), created.ID+1000, 5)

		assert.Error(t, err)
	})
}

func newIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()
	migrationPath := writeMigrationScript(t)
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "wongnok_test",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
			},
			Files: []testcontainers.ContainerFile{{
				HostFilePath:      migrationPath,
				ContainerFilePath: "/docker-entrypoint-initdb.d/001_migrations.sql",
				FileMode:          0o644,
			}},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(ctx)) })

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)
	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=wongnok_test sslmode=disable", host, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return db
}

func writeMigrationScript(t *testing.T) string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, paths)

	var statements []string
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		require.NoError(t, err)
		up, _, _ := strings.Cut(string(contents), "-- +goose Down")
		up = strings.TrimPrefix(up, "-- +goose Up")
		statements = append(statements, strings.TrimSpace(up))
	}

	scriptPath := filepath.Join(t.TempDir(), "migrations.sql")
	require.NoError(t, os.WriteFile(scriptPath, []byte(strings.Join(statements, "\n\n")), 0o600))
	return scriptPath
}

func createCreator(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()
	creatorID := uuid.New()
	require.NoError(t, db.Exec(
		"INSERT INTO users (id, email, uid) VALUES (?, ?, ?)",
		creatorID,
		fmt.Sprintf("%s@example.com", creatorID),
		fmt.Sprintf("uid-%s", creatorID),
	).Error)
	return creatorID
}
