package recipe

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"wongnok/internal/httputil"
	"wongnok/internal/reqctx"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, creatorID uuid.UUID, recipe Recipe) (*Recipe, error)
	List(ctx context.Context, userID uuid.UUID, query GetRecipesQuery) ([]Recipe, int64, error)
	Get(ctx context.Context, id int, userID uuid.UUID) (*Recipe, error)
	Replace(ctx context.Context, id int, userID uuid.UUID, recipe Recipe) (*Recipe, error)
	Delete(ctx context.Context, id int, userID uuid.UUID) error
	Favorite(ctx context.Context, id int, userID uuid.UUID) error
	Unfavorite(ctx context.Context, id int, userID uuid.UUID) error
	Rate(ctx context.Context, id int, userID uuid.UUID, rating int) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// Create godoc
//
//	@Summary		สร้างสูตรอาหาร
//	@Description	สร้างสูตรอาหารจากข้อมูลที่ระบุ โดยผู้ใช้ที่ยืนยันตัวตนแล้วจะเป็นผู้สร้างสูตร
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateRecipeRequest	true	"Recipe data"
//	@Success		201		{object}	CreateRecipeResponse
//	@Failure		400		{object}	httputil.ErrorResponse
//	@Failure		401		{object}	httputil.ErrorResponse
//	@Failure		500		{object}	httputil.ErrorResponse
//	@Router			/recipes [post]
func (hdr *handler) Create(ctx *gin.Context) {
	creatorID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	var req CreateRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	recipe, err := hdr.service.Create(ctx.Request.Context(), creatorID, req.ToRecipe())
	if err != nil {
		if errors.Is(err, ErrInvalidReferenceData) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, NewCreateRecipeResponse(*recipe))
}

// GetRecipes godoc
//
//	@Summary		เรียกดูสูตรอาหารทั้งหมด
//	@Description	ค้นหาสูตรอาหารทั้งหมด และสามารถกรองด้วย ชื่อ, ความยาก และเรียงตามเวลาที่สร้างได้
//	@Tags			recipes
//	@Produce		json
//	@Param			name		query		string	false	"ชื่อของสูตรอาหาร"
//	@Param			difficulty	query		string	false	"Id ของความยากในการทำ"
//	@Param			favorite	query		bool	false	"กรองสูตรอาหารตามสถานะรายการโปรดของผู้ใช้ปัจจุบัน: true = เฉพาะที่ถูกใจไว้, false = เฉพาะที่ไม่ได้ถูกใจไว้ (ต้องยืนยันตัวตน)"
//	@Param			sort		query		string	false	"เรียงลำดับตามเวลาที่สร้าง"	Enums(ASC, DESC)	default(DESC)
//	@Param			page		query		int		false	"หน้าที่ต้องการแสดง"		minimum(1)			default(1)
//	@Param			limit		query		int		false	"จำนวนรายการต่อหน้า"		minimum(1)			maximum(100)	default(12)
//	@Success		200			{object}	ListRecipesResponse
//	@Failure		400			{object}	httputil.ErrorResponse
//	@Failure		404			{object}	httputil.ErrorResponse
//	@Failure		500			{object}	httputil.ErrorResponse
//	@Router			/recipes [get]
func (hdr *handler) GetRecipes(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	var query GetRecipesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: err.Error()})
		return
	}

	recipes, total, err := hdr.service.List(ctx.Request.Context(), userID, query)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidReferenceData):
			ctx.JSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})

		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})

		}
		return
	}

	ctx.JSON(http.StatusOK, NewListRecipesResponse(recipes, total))
}

// GetRecipe godoc
//
//	@Summary		เรียกดูสูตรอาหารแบบรายรายการ
//	@Description	ค้นหาสูตรอาหารจาก id แล้วคืนข้อมูลสูตรอาหารแบบสมบูรณ์
//	@Tags			recipes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Recipe ID"
//	@Success		200	{object}	RecipeResponse
//	@Failure		400	{object}	httputil.ErrorResponse
//	@Failure		401	{object}	httputil.ErrorResponse
//	@Failure		404	{object}	httputil.ErrorResponse
//	@Failure		500	{object}	httputil.ErrorResponse
//	@Router			/recipes/{id} [get]
func (hdr *handler) GetRecipe(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	recipe, err := hdr.service.Get(ctx.Request.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecipeNotFound):
			ctx.AbortWithStatusJSON(http.StatusNotFound, httputil.ErrorResponse{Message: "recipe not found"})

		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})

		}
		return
	}

	ctx.JSON(http.StatusOK, NewRecipeResponse(*recipe))
}

// Replace godoc
//
//	@Summary		แก้ไขสูตรอาหารทั้งหมด
//	@Description	แทนที่ข้อมูลสูตรอาหารทั้งหมดด้วยข้อมูลที่ระบุ โดยผู้สร้างสูตรเท่านั้นที่แก้ไขได้
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Recipe ID"
//	@Param			request	body		ReplaceRecipeRequest	true	"Recipe data"
//	@Success		200		{object}	RecipeResponse
//	@Failure		400		{object}	httputil.ErrorResponse
//	@Failure		401		{object}	httputil.ErrorResponse
//	@Failure		403		{object}	httputil.ErrorResponse
//	@Failure		404		{object}	httputil.ErrorResponse
//	@Failure		500		{object}	httputil.ErrorResponse
//	@Router			/recipes/{id} [put]
func (hdr *handler) Replace(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	var req ReplaceRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	recipe, err := hdr.service.Replace(ctx.Request.Context(), id, userID, req.ToRecipe())
	if err != nil {
		switch {
		case errors.Is(err, ErrRecipeNotFound), errors.Is(err, ErrReferenceDataUnavailable):
			ctx.AbortWithStatusJSON(http.StatusNotFound, httputil.ErrorResponse{Message: "recipe not found"})

		case errors.Is(err, ErrForbidden):
			ctx.AbortWithStatusJSON(http.StatusForbidden, httputil.ErrorResponse{Message: "forbidden"})

		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})

		}
		return
	}

	ctx.JSON(http.StatusOK, NewRecipeResponse(*recipe))
}

// Favorite godoc
//
//	@Summary		เพิ่มสูตรอาหารในรายการโปรด
//	@Description	เพิ่มสูตรอาหารที่ระบุเข้ารายการโปรดของผู้ใช้ที่ยืนยันตัวตนแล้ว หากเคยเพิ่มไว้แล้วจะไม่มีผลซ้ำ
//	@Tags			recipes
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Recipe ID"
//	@Success		204
//	@Failure		400	{object}	httputil.ErrorResponse
//	@Failure		401	{object}	httputil.ErrorResponse
//	@Failure		500	{object}	httputil.ErrorResponse
//	@Router			/recipes/{id}/favorite [post]
func (hdr *handler) Favorite(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	if err := hdr.service.Favorite(ctx.Request.Context(), id, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Unfavorite godoc
//
//	@Summary		ลบสูตรอาหารออกจากรายการโปรด
//	@Description	ลบสูตรอาหารที่ระบุออกจากรายการโปรดของผู้ใช้ที่ยืนยันตัวตนแล้วแบบ hard delete หากไม่เคยเพิ่มไว้จะไม่มีผล
//	@Tags			recipes
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Recipe ID"
//	@Success		204
//	@Failure		400	{object}	httputil.ErrorResponse
//	@Failure		401	{object}	httputil.ErrorResponse
//	@Failure		500	{object}	httputil.ErrorResponse
//	@Router			/recipes/{id}/favorite [delete]
func (hdr *handler) Unfavorite(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	if err := hdr.service.Unfavorite(ctx.Request.Context(), id, userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Rate godoc
//
//	@Summary		ให้คะแนนสูตรอาหาร
//	@Description	ให้คะแนนสูตรอาหารที่ระบุ ผู้ใช้แต่ละคนให้คะแนนสูตรอาหารแต่ละสูตรได้เพียงครั้งเดียว คำขอซ้ำจะไม่มีผล และค่าเฉลี่ยของสูตรจะถูกคำนวณใหม่ทุกครั้งที่บันทึกคะแนน
//	@Tags			recipes
//	@Accept			json
//	@Security		BearerAuth
//	@Param			id		path	int					true	"Recipe ID"
//	@Param			request	body	RateRecipeRequest	true	"Rating data"
//	@Success		204
//	@Failure		400	{object}	httputil.ErrorResponse
//	@Failure		401	{object}	httputil.ErrorResponse
//	@Failure		500	{object}	httputil.ErrorResponse
//	@Router			/recipes/{id}/rating [post]
func (hdr *handler) Rate(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	var req RateRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	if err := hdr.service.Rate(ctx.Request.Context(), id, userID, req.Rating); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Delete godoc
//
//	@Summary		ลบสูตรอาหาร
//	@Description	ลบสูตรอาหารแบบ soft delete โดยผู้สร้างสูตรเท่านั้นที่ลบได้
//	@Tags			recipes
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Recipe ID"
//	@Success		204
//	@Failure		400	{object}	httputil.ErrorResponse
//	@Failure		401	{object}	httputil.ErrorResponse
//	@Failure		403	{object}	httputil.ErrorResponse
//	@Failure		404	{object}	httputil.ErrorResponse
//	@Failure		500	{object}	httputil.ErrorResponse
//	@Router			/recipes/{id} [delete]
func (hdr *handler) Delete(ctx *gin.Context) {
	userID, ok := reqctx.UserID(ctx.Request.Context())
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httputil.ErrorResponse{Message: "unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httputil.ErrorResponse{Message: "invalid request"})
		return
	}

	if err := hdr.service.Delete(ctx.Request.Context(), id, userID); err != nil {
		switch {
		case errors.Is(err, ErrRecipeNotFound):
			ctx.AbortWithStatusJSON(http.StatusNotFound, httputil.ErrorResponse{Message: "recipe not found"})

		case errors.Is(err, ErrForbidden):
			ctx.AbortWithStatusJSON(http.StatusForbidden, httputil.ErrorResponse{Message: "forbidden"})

		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "internal server error"})

		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
