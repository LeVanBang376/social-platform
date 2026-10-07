package response

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Pagination
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.TotalPages = int(
		(total + int64(p.PerPage) - 1) / int64(p.PerPage),
	)
}

func NewPagination(r *http.Request) *Pagination {
	page, err := strconv.Atoi(
		r.URL.Query().Get("page"),
	)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(
		r.URL.Query().Get("per_page"),
	)
	if err != nil || perPage < 1 {
		perPage = 10
	}

	if perPage > 200 {
		perPage = 200
	}

	return &Pagination{
		Page:    page,
		PerPage: perPage,
	}
}

// Response
type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type PaginatedResponse[T any] struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       T           `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type NonDataResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSON[T any](
	w http.ResponseWriter,
	status int,
	message string,
	data T,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Response[T]{
		Message: message,
		Data:    data,
	})
}

func PaginatedJSON[T any](
	w http.ResponseWriter,
	status int,
	message string,
	data T,
	pagination *Pagination,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(PaginatedResponse[T]{
		Code:       status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func NonDataJSON(
	w http.ResponseWriter,
	status int,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(NonDataResponse{
		Code:    status,
		Message: message,
	})
}
