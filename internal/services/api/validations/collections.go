package validations

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
	"net/http"
	"strconv"
)

type colelctionQueryParams struct {
	Page   int
	Limit  int
	Recent bool
}

func ValidateCollectionsParams(c *gin.Context) (*colelctionQueryParams, *httperr.Error) {
	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = fmt.Sprintf("%d", defaultPage)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidPage)
	}

	if page == 0 {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidPage)
	}

	limitStr, ok := c.GetQuery("limit")
	if !ok {
		limitStr = fmt.Sprintf("%d", defaultLimit)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidLimit)
	}

	recent := false
	recentStr, ok := c.GetQuery("recent")
	if ok && recentStr == "true" {
		recent = true
	}

	return &colelctionQueryParams{
		Page:   page,
		Limit:  limit,
		Recent: recent,
	}, nil
}
