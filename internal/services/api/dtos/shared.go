package dtos

import (
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseDto struct {
	Success      bool     `json:"success"`
	FullMessages []string `json:"full_messages"`
}

type ErrorDto struct {
	BaseDto
	Errors map[string]interface{} `json:"errors"`
}

func CreatePageMeta(loadedItemsCount, page, page_size, totalItemsCount int) map[string]interface{} {
	page_meta := map[string]interface{}{}
	page_meta["offset"] = (page - 1) * page_size
	page_meta["requested_page_size"] = page_size
	page_meta["current_page_number"] = page
	page_meta["current_items_count"] = loadedItemsCount

	page_meta["prev_page_number"] = 1
	total_pages_count := int(math.Ceil(float64(totalItemsCount) / float64(page_size)))
	page_meta["total_pages_count"] = total_pages_count

	if page < total_pages_count {
		page_meta["has_next_page"] = true
		page_meta["next_page_number"] = page + 1
	} else {
		page_meta["has_next_page"] = false
		page_meta["next_page_number"] = 1
	}
	if page > 1 {
		page_meta["prev_page_number"] = page - 1
	} else {
		page_meta["has_prev_page"] = false
		page_meta["prev_page_number"] = 1
	}

	page_meta["next_page_url"] = fmt.Sprintf("/?page=%d&page_size=%d", page_meta["next_page_number"], page_meta["requested_page_size"])
	page_meta["prev_page_url"] = fmt.Sprintf("/?page=%d&page_size=%d", page_meta["prev_page_number"], page_meta["requested_page_size"])

	response := gin.H{
		"success":   true,
		"page_meta": page_meta,
	}

	return response
}

func CreatePagedResponse(resources []interface{}, resource_name string, page, page_size, totalItemsCount int) map[string]interface{} {

	response := CreatePageMeta(len(resources), page, page_size, totalItemsCount)
	response[resource_name] = resources
	return response
}

// This should only be called when we have an Error that is returned from a ShouldBind which contains a lot of information
// other kind of errors should use other functions such as CreateDetailedErrorDto
func CreateBadRequestErrorDto(err error) ErrorDto {
	res := ErrorDto{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	res.FullMessages = make([]string, len(errs))
	count := 0
	for _, v := range errs {
		if v.ActualTag() == "required" {
			var message = fmt.Sprintf("%v is required", v.Field)
			res.Errors[v.Field()] = message
			res.FullMessages[count] = message
		} else {
			var message = fmt.Sprintf("%v has to be %v", v.Field, v.ActualTag)
			res.Errors[v.Field()] = message
			res.FullMessages = append(res.FullMessages, message)
		}
		count++
	}
	return res
}
