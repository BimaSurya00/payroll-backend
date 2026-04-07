package helper

import (
	"fmt"
	"net/url"
	"strconv"

	"hris/internal/payroll/dto"
)

type PayrollPagination struct {
	CurrentPage  int                      `json:"currentPage"`
	Data         []dto.PayrollListResponse `json:"data"`
	FirstPageURL string                   `json:"firstPageUrl"`
	From         int                      `json:"from"`
	LastPage     int                      `json:"lastPage"`
	LastPageURL  string                   `json:"lastPageUrl"`
	Links        []Link                   `json:"links"`
	NextPageURL  *string                  `json:"nextPageUrl"`
	Path         string                   `json:"path"`
	PerPage      int                      `json:"perPage"`
	PrevPageURL  *string                  `json:"prevPageUrl"`
	To           int                      `json:"to"`
	Total        int64                    `json:"total"`
}

type Link struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Active bool    `json:"active"`
}

func BuildPayrollPagination(data []dto.PayrollListResponse, page, perPage int, total int64, path string) *PayrollPagination {
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	var firstPageURL, lastPageURL string
	if page > 1 {
		firstPageURL = fmt.Sprintf("%s?page=%d&per_page=%d", path, 1, perPage)
	} else {
		firstPageURL = fmt.Sprintf("%s?page=%d&per_page=%d", path, 1, perPage)
	}
	lastPageURL = fmt.Sprintf("%s?page=%d&per_page=%d", path, totalPages, perPage)

	from := (page-1)*perPage + 1
	to := from + len(data) - 1

	links := make([]Link, 0)
	links = append(links, Link{
		URL:    nil,
		Label:  "&laquo; Previous",
		Active: false,
	})

	for i := 1; i <= totalPages; i++ {
		pageURL := fmt.Sprintf("%s?page=%d&per_page=%d", path, i, perPage)
		links = append(links, Link{
			URL:    &pageURL,
			Label:  strconv.Itoa(i),
			Active: i == page,
		})
	}

	links = append(links, Link{
		URL:    nil,
		Label:  "Next &raquo;",
		Active: false,
	})

	var nextPageURL, prevPageURL *string
	if page < totalPages {
		nextPage := fmt.Sprintf("%s?page=%d&per_page=%d", path, page+1, perPage)
		nextPageURL = &nextPage
	}
	if page > 1 {
		prevPage := fmt.Sprintf("%s?page=%d&per_page=%d", path, page-1, perPage)
		prevPageURL = &prevPage
	}

	// Ensure path doesn't have query params
	cleanPath, _ := url.Parse(path)
	cleanPath.RawQuery = ""
	cleanPath.Fragment = ""

	return &PayrollPagination{
		CurrentPage:  page,
		Data:         data,
		FirstPageURL: firstPageURL,
		From:         from,
		LastPage:     totalPages,
		LastPageURL:  lastPageURL,
		Links:        links,
		NextPageURL:  nextPageURL,
		Path:         cleanPath.String(),
		PerPage:      perPage,
		PrevPageURL:  prevPageURL,
		To:           to,
		Total:        total,
	}
}
