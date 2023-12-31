package validations

import (
	"encoding/json"
	"fmt"
	"github.com/unanoc/blockchain-indexer/internal/services/api/dtos"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
)

func ValidateCollectionsParams(c *gin.Context) (*dtos.GetCollectionQueryDtos, *httperr.Error) {
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

	return &dtos.GetCollectionQueryDtos{
		Page:   page,
		Limit:  limit,
		Recent: recent,
	}, nil
}

func ValidateCreateCollectionParams(c *gin.Context) (*dtos.CreateCollectionBodyDtos, *httperr.Error) {
	var collectionRequest dtos.CreateCollectionBodyDtos
	if err := c.ShouldBind(&collectionRequest); err != nil {
		errMessage, _ := json.Marshal(dtos.CreateBadRequestErrorDto(err))
		return nil, httperr.NewError(http.StatusBadRequest, string(errMessage))
	}

	mintedTimestamp := collectionRequest.MintedTimestamp
	if mintedTimestamp == 0 {
		mintedTimestamp = time.Now().Unix()
	}
	return &dtos.CreateCollectionBodyDtos{
		Name:            collectionRequest.Name,
		Slug:            collectionRequest.Slug,
		Contract:        collectionRequest.Contract,
		Metadata:        collectionRequest.Metadata,
		TokenCount:      collectionRequest.TokenCount,
		MintedTimestamp: collectionRequest.MintedTimestamp,
	}, nil
}
