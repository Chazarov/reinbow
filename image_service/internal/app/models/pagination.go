package models

// Pagination описывает параметры постраничной выборки.
type Pagination struct {
	Page     int
	PageSize int
}