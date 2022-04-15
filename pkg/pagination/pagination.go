package pagination

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// DefaultPageSize specifies the default page size
	DefaultPageSize = 10
	// MaxPageSize specifies the maximum page size
	MaxPageSize = 100
	// PageVar specifies the query parameter name for page number
	PageVar = "page"
	// PageSizeVar specifies the query parameter name for page size
	PageSizeVar = "pageSize"
)

// Pages represents a paginated list of data items.
type Pages struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	PageCount  int         `json:"pageCount"`
	TotalCount int         `json:"totalCount"`
	Items      interface{} `json:"items"`
}

// New creates a new Pages instance.
// The page parameter is 1-based and refers to the current page index/number.
// The pageSize parameter refers to the number of items on each page.
// And the total parameter specifies the total number of data items.
// If total is less than 0, it means total is unknown.
func NewFromGinRequest(c *gin.Context, total int, items interface{}) *Pages {
	pageIndex, pageSize := GetPaginationParametersFromRequest(c)

	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	pageCount := -1
	if total >= 0 {
		pageCount = (total + pageSize - 1) / pageSize
	}
	paginatedResult := &Pages{
		Page:       pageIndex,
		PageSize:   pageSize,
		TotalCount: total,
		PageCount:  pageCount,
		Items:      items,
	}
	c.Header("Page Links", paginatedResult.BuildLinkHeader(c.Request.URL.Path, DefaultPageSize))
	return paginatedResult
}

func GetPaginationParametersFromRequest(c *gin.Context) (pageIndex, pageSize int) {
	pageIndex = parseInt(c.Query(PageVar), 1)
	pageSize = parseInt(c.Query(PageSizeVar), DefaultPageSize)
	return pageIndex, pageSize
}

// parseInt parses a string into an integer. If parsing is failed, defaultValue will be returned.
func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

// BuildLinkHeader returns an HTTP header containing the links about the pagination.
func (p *Pages) BuildLinkHeader(baseURL string, defaultPageSize int) string {
	links := p.BuildLinks(baseURL, defaultPageSize)
	header := ""
	if links[0] != "" {
		header += fmt.Sprintf("<%v>; rel=\"first\", ", links[0])
		header += fmt.Sprintf("<%v>; rel=\"prev\"", links[1])
	}
	if links[2] != "" {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf("<%v>; rel=\"next\"", links[2])
		if links[3] != "" {
			header += fmt.Sprintf(", <%v>; rel=\"last\"", links[3])
		}
	}
	return header
}

// BuildLinks returns the first, prev, next, and last links corresponding to the pagination.
// A link could be an empty string if it is not needed.
// For example, if the pagination is at the first page, then both first and prev links
// will be empty.
func (p *Pages) BuildLinks(baseURL string, defaultPageSize int) [4]string {
	var links [4]string
	pageCount := p.PageCount
	page := p.Page
	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}
	if page > 1 {
		links[0] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, 1)
		links[1] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page-1)
	}
	if pageCount >= 0 && page < pageCount {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
		links[3] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, pageCount)
	} else if pageCount < 0 {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
	}
	if pageSize := p.PageSize; pageSize != defaultPageSize {
		for i := 0; i < 4; i++ {
			if links[i] != "" {
				zap.L().Debug("ok")
				links[i] += fmt.Sprintf("&%v=%v", PageSizeVar, pageSize)
			}
		}
	}

	return links
}

func (p *Pages) SetItems(items interface{}) {
	p.Items = items
}
