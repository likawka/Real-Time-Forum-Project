package services

import (
	"net/http"
	"strconv"
)

func GetSortingCriteria(sorting, param string) (sortBy, sqlSortExp string) {
	switch sorting {
	case "old":
		sortBy = "old"
		sqlSortExp = param + ".created_at ASC"
	case "popular":
		sortBy = "popular"
		sqlSortExp = param + ".rate DESC"
	case "name_ABC":
		sortBy = "name_ABC"
		sqlSortExp = "LOWER(" + param + ".nickname) ASC"
	default:
		sortBy = "new" 
		sqlSortExp = param + ".created_at DESC"
	}
	return sortBy, sqlSortExp
}

func ExtractPaginationParams(r *http.Request, baseSortType, sortWhat string) (sortBy, sqlSortExp string, page, pageSize int) {
	sortBy = r.URL.Query().Get("sort")
	pageParam := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeParam := r.URL.Query().Get("pageSize")
	pageSize, err = strconv.Atoi(pageSizeParam)
	if err != nil || pageSize <= 0 {
		pageSize = 20
	}

	// Validate sortBy to ensure it is one of "new", "old", or "popular"
	switch sortBy {
	case "new", "old", "popular", "name_ABC":
		// Retrieve the SQL sorting expression based on sortBy
		sortBy, sqlSortExp = GetSortingCriteria(sortBy, sortWhat)
	default:
		// If sortBy is unrecognized, fall back to baseSortType
		sortBy, sqlSortExp = GetSortingCriteria(baseSortType, sortWhat)
	}

	return sortBy, sqlSortExp, page, pageSize
}
