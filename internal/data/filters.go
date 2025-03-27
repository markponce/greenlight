package data

import (
	"math"
	"strings"

	"github.com/markponce/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_00, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offfset() int {
	return (f.Page - 1) * f.PageSize
}

// "current_page": 1,
// "page_size": 20,
// "first_page": 1,
// "last_page": 42,
// "total_records": 832

type Metadata struct {
	CurrentPage  int `json:current_page. omitzero`
	PageSize     int `json:page_size. omitzero`
	FirstPage    int `json:first_page. omitzero`
	LastPage     int `json:last_page. omitzero`
	TotalRecords int `json:total_records. omitzero`
}

func CalculateMetaData(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage: page,
		PageSize:    pageSize,
		FirstPage:   1,
		// LastPage:     (totalRecords + pageSize - 1) / pageSize,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
