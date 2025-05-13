package pagination

import (
	"net/http"
	"strconv"
)

// DefaultLimit is the default number of items to return per page
const DefaultLimit = 10

// PaginationParams contains the pagination parameters
type PaginationParams struct {
	Offset int
	Limit  int
	Page   int
}

// ParsePaginationParams extracts pagination parameters from the request
func ParsePaginationParams(r *http.Request) PaginationParams {
	// Parse offset parameter
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	// Parse limit parameter with default value
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = DefaultLimit
	}

	// Parse page parameter (alternative to offset)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	// If page is provided but offset is not, calculate offset from page
	if r.URL.Query().Get("page") != "" && r.URL.Query().Get("offset") == "" {
		offset = (page - 1) * limit
	}

	return PaginationParams{
		Offset: offset,
		Limit:  limit,
		Page:   page,
	}
}

// SetPaginationHeaders sets the pagination headers on the response
func SetPaginationHeaders(w http.ResponseWriter, params PaginationParams, totalCount int) {
	w.Header().Set("X-Total-Count", strconv.Itoa(totalCount))

	// Calculate total pages
	totalPages := totalCount / params.Limit
	if totalCount%params.Limit > 0 {
		totalPages++
	}
	w.Header().Set("X-Total-Pages", strconv.Itoa(totalPages))

	// Current page (calculated from offset and limit)
	currentPage := (params.Offset / params.Limit) + 1
	w.Header().Set("X-Current-Page", strconv.Itoa(currentPage))

	// Page size
	w.Header().Set("X-Page-Size", strconv.Itoa(params.Limit))
}

// CreatePaginationResponse creates a standard pagination response
type PaginationResponse struct {
	TotalCount  int         `json:"total_count"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	Data        interface{} `json:"data"`
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(params PaginationParams, totalCount int, data interface{}) *PaginationResponse {
	// Calculate total pages
	totalPages := totalCount / params.Limit
	if totalCount%params.Limit > 0 {
		totalPages++
	}

	// Calculate current page from offset and limit
	currentPage := (params.Offset / params.Limit) + 1

	return &PaginationResponse{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    params.Limit,
		Data:        data,
	}
}
