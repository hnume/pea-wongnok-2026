package recipe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) HasActiveReferences(ctx context.Context, difficultyID, durationID string) (bool, error) {
	var difficultyCount int64
	if err := repo.db.WithContext(ctx).Model(&Difficulty{}).Where("id = ?", difficultyID).Count(&difficultyCount).Error; err != nil {
		return false, err
	}

	if difficultyCount == 0 {
		return false, nil
	}

	var durationCount int64
	if err := repo.db.WithContext(ctx).Model(&Duration{}).Where("id = ?", durationID).Count(&durationCount).Error; err != nil {
		return false, err
	}

	return durationCount > 0, nil
}

func (repo *repository) Create(ctx context.Context, recipe Recipe) (*Recipe, error) {
	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Create(&recipe).Error; err != nil {
			return err
		}

		for index := range recipe.Ingredients {
			recipe.Ingredients[index].RecipeID = recipe.ID
		}
		if len(recipe.Ingredients) > 0 {
			if err := tx.Omit(clause.Associations).Create(&recipe.Ingredients).Error; err != nil {
				return err
			}
		}

		for index := range recipe.Instructions {
			recipe.Instructions[index].RecipeID = recipe.ID
		}
		if len(recipe.Instructions) > 0 {
			if err := tx.Omit(clause.Associations).Create(&recipe.Instructions).Error; err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return nil, err

	}

	return &recipe, nil
}

func (repo *repository) Replace(ctx context.Context, recipe Recipe) (*Recipe, error) {
	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Recipe{}).Where("id = ?", recipe.ID).Updates(map[string]any{
			"name":          recipe.Name,
			"description":   recipe.Description,
			"image_url":     recipe.ImageURL,
			"difficulty_id": recipe.DifficultyID,
			"duration_id":   recipe.DurationID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Where("recipe_id = ?", recipe.ID).Delete(&RecipeIngredient{}).Error; err != nil {
			return err
		}
		if err := tx.Where("recipe_id = ?", recipe.ID).Delete(&RecipeInstruction{}).Error; err != nil {
			return err
		}

		for index := range recipe.Ingredients {
			recipe.Ingredients[index].ID = 0
			recipe.Ingredients[index].RecipeID = recipe.ID
		}
		if len(recipe.Ingredients) > 0 {
			if err := tx.Omit(clause.Associations).Create(&recipe.Ingredients).Error; err != nil {
				return err
			}
		}

		for index := range recipe.Instructions {
			recipe.Instructions[index].ID = 0
			recipe.Instructions[index].RecipeID = recipe.ID
		}
		if len(recipe.Instructions) > 0 {
			if err := tx.Omit(clause.Associations).Create(&recipe.Instructions).Error; err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return nil, fmt.Errorf("replace recipe %d: %w", recipe.ID, err)
	}

	return repo.FindByID(ctx, recipe.ID)
}

func (repo *repository) Delete(ctx context.Context, id int) error {
	result := repo.db.WithContext(ctx).Delete(&Recipe{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete recipe %d: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrRecipeNotFound
	}

	return nil
}

func (repo *repository) FindByID(ctx context.Context, id int) (*Recipe, error) {
	var recipe Recipe

	// Context
	db := repo.db.WithContext(ctx)

	// Preload
	db = db.Preload("Difficulty").Preload("Duration").Preload("Creator").Preload("Ingredients").Preload("Instructions")

	if err := db.First(&recipe, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecipeNotFound
		}

		return nil, fmt.Errorf("find recipe %d: %w", id, err)
	}

	var ratingTotal int64
	if err := repo.db.WithContext(ctx).Model(&RecipeRating{}).Where("recipe_id = ?", id).Count(&ratingTotal).Error; err != nil {
		return nil, fmt.Errorf("count ratings for recipe %d: %w", id, err)
	}

	recipe.RatingTotal = ratingTotal

	return &recipe, nil
}

func (repo *repository) Favorite(ctx context.Context, userID uuid.UUID, recipeID int) error {
	favorite := UserFavorite{UserID: userID, RecipeID: recipeID}

	if err := repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "recipe_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"deleted_at": nil,
			"updated_at": time.Now(),
		}),
	}).Create(&favorite).Error; err != nil {
		return fmt.Errorf("favorite recipe %d: %w", recipeID, err)
	}

	return nil
}

func (repo *repository) Unfavorite(ctx context.Context, userID uuid.UUID, recipeID int) error {
	if err := repo.db.WithContext(ctx).Unscoped().Where(
		"user_id = ? AND recipe_id = ?", userID, recipeID,
	).Delete(&UserFavorite{}).Error; err != nil {
		return fmt.Errorf("unfavorite recipe %d: %w", recipeID, err)
	}

	return nil
}

func (repo *repository) Rate(ctx context.Context, userID uuid.UUID, recipeID int, score float64) error {
	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&RecipeRating{}).Where("user_id = ? AND recipe_id = ?", userID, recipeID).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return nil
		}

		if err := tx.Create(&RecipeRating{UserID: userID, RecipeID: recipeID, Score: score}).Error; err != nil {
			return err
		}

		var average float64
		if err := tx.Model(&RecipeRating{}).Where("recipe_id = ?", recipeID).Select("COALESCE(AVG(score), 0)").Scan(&average).Error; err != nil {
			return err
		}

		return tx.Model(&Recipe{}).Where("id = ?", recipeID).Update("average_rating", average).Error

	}); err != nil {
		return fmt.Errorf("rate recipe %d: %w", recipeID, err)
	}

	return nil
}

func (repo *repository) DifficultyExists(ctx context.Context, id string) (bool, error) {
	var count int64

	if err := repo.db.WithContext(ctx).Model(&Difficulty{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check difficulty %q: %w", id, err)
	}

	return count > 0, nil
}

func (repo *repository) List(ctx context.Context, userID uuid.UUID, query GetRecipesQuery) ([]Recipe, int64, error) {
	db := repo.db.WithContext(ctx).Model(&Recipe{})

	if query.Name != "" {
		db = db.Where("name ILIKE ?", ("%" + query.Name + "%"))
	}

	if query.Difficulty != "" {
		db = db.Where("difficulty_id = ?", query.Difficulty)
	}

	if query.Favorite != nil {
		favoriteIDs := repo.db.Model(&UserFavorite{}).Select("recipe_id").Where("user_id = ?", userID)

		if *query.Favorite {
			db = db.Where("id IN (?)", favoriteIDs)
		} else {
			db = db.Where("id NOT IN (?)", favoriteIDs)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count recipes: %w", err)
	}

	order := "created_at DESC"
	if query.Sort != "" {
		order = fmt.Sprintf("created_at %s", query.Sort)
	}

	// Preload
	db = db.Preload("Difficulty").Preload("Duration").Preload("Creator").Preload("Ingredients").Preload("Instructions")

	// Order
	db = db.Order(order)

	// Pagination
	db = db.Offset((query.Page - 1) * query.Limit).Limit(query.Limit)

	// Find
	var recipes []Recipe
	if err := db.Find(&recipes).Error; err != nil {
		return nil, 0, fmt.Errorf("list recipes: %w", err)
	}

	favorited, err := repo.favoritedRecipeIDs(ctx, userID, recipes)
	if err != nil {
		return nil, 0, fmt.Errorf("list recipes: %w", err)
	}

	ratingTotals, err := repo.ratingTotals(ctx, recipes)
	if err != nil {
		return nil, 0, fmt.Errorf("list recipes: %w", err)
	}

	for index := range recipes {
		recipes[index].IsFavorite = favorited[recipes[index].ID]
		recipes[index].RatingTotal = ratingTotals[recipes[index].ID]
	}

	return recipes, total, nil
}

func (repo *repository) ratingTotals(ctx context.Context, recipes []Recipe) (map[int]int64, error) {
	if len(recipes) == 0 {
		return map[int]int64{}, nil
	}

	ids := make([]int, len(recipes))
	for index, recipe := range recipes {
		ids[index] = recipe.ID
	}

	var rows []struct {
		RecipeID int
		Total    int64
	}

	// Model
	db := repo.db.WithContext(ctx).Model(&RecipeRating{})

	// Where
	db = db.Select("recipe_id, COUNT(*) AS total").Where("recipe_id IN ?", ids)

	// Scan
	if err := db.Group("recipe_id").Scan(&rows).Error; err != nil {
		return nil, err
	}

	totals := make(map[int]int64, len(rows))
	for _, row := range rows {
		totals[row.RecipeID] = row.Total
	}

	return totals, nil
}

func (repo *repository) favoritedRecipeIDs(ctx context.Context, userID uuid.UUID, recipes []Recipe) (map[int]bool, error) {
	if len(recipes) == 0 {
		return map[int]bool{}, nil
	}

	ids := make([]int, len(recipes))
	for index, recipe := range recipes {
		ids[index] = recipe.ID
	}

	var favoritedIDs []int
	if err := repo.db.WithContext(ctx).Model(&UserFavorite{}).Where(
		"user_id = ? AND recipe_id IN ?", userID, ids,
	).Pluck("recipe_id", &favoritedIDs).Error; err != nil {
		return nil, err
	}

	favorited := make(map[int]bool, len(favoritedIDs))
	for _, id := range favoritedIDs {
		favorited[id] = true
	}

	return favorited, nil
}

func (repo *repository) IsFavorite(ctx context.Context, userID uuid.UUID, recipeID int) (bool, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&UserFavorite{}).
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check favorite recipe %d: %w", recipeID, err)
	}

	return count > 0, nil
}
