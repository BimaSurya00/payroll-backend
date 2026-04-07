package helper

import (
	"fmt"
	"math"

	"hris/internal/schedule/dto"
)

type Pagination[T any] struct {
	CurrentPage  int     `json:"currentPage"`
	Data         []T     `json:"data"`
	FirstPageURL string  `json:"firstPageUrl"`
	From         int     `json:"from"`
	LastPage     int     `json:"lastPage"`
	LastPageURL  string  `json:"lastPageUrl"`
	Links        []Link  `json:"links"`
	NextPageURL  *string `json:"nextPageUrl"`
	Path         string  `json:"path"`
	PerPage      int     `json:"perPage"`
	PrevPageURL  *string `json:"prevPageUrl"`
	To           int     `json:"to"`
	Total        int64   `json:"total"`
}

type Link struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Active bool    `json:"active"`
}

func NewPagination(data []*dto.ScheduleResponse, page, perPage int, total int64, path string) *Pagination[*dto.ScheduleResponse] {
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	from := (page-1)*perPage + 1
	to := from + len(data) - 1

	if total == 0 {
		from = 0
		to = 0
	}

	firstPageURL := fmt.Sprintf("%s?page=1&per_page=%d", path, perPage)
	lastPageURL := fmt.Sprintf("%s?page=%d&per_page=%d", path, lastPage, perPage)

	var nextPageURL *string
	if page < lastPage {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", path, page+1, perPage)
		nextPageURL = &url
	}

	var prevPageURL *string
	if page > 1 {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", path, page-1, perPage)
		prevPageURL = &url
	}

	links := generateLinks(page, lastPage, path, perPage)

	// Convert slice
	dataSlice := make([]*dto.ScheduleResponse, len(data))
	for i, v := range data {
		dataSlice[i] = v
	}

	return &Pagination[*dto.ScheduleResponse]{
		CurrentPage:  page,
		Data:         dataSlice,
		FirstPageURL: firstPageURL,
		From:         from,
		LastPage:     lastPage,
		LastPageURL:  lastPageURL,
		Links:        links,
		NextPageURL:  nextPageURL,
		Path:         path,
		PerPage:      perPage,
		PrevPageURL:  prevPageURL,
		To:           to,
		Total:        total,
	}
}

func generateLinks(currentPage, lastPage int, path string, perPage int) []Link {
	links := []Link{
		{
			URL:    nil,
			Label:  "&laquo; Previous",
			Active: false,
		},
	}

	if currentPage > 1 {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", path, currentPage-1, perPage)
		links[0].URL = &url
	}

	start := max(1, currentPage-2)
	end := min(lastPage, currentPage+2)

	for i := start; i <= end; i++ {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", path, i, perPage)
		links = append(links, Link{
			URL:    &url,
			Label:  fmt.Sprintf("%d", i),
			Active: i == currentPage,
		})
	}

	nextLink := Link{
		URL:    nil,
		Label:  "Next &raquo;",
		Active: false,
	}

	if currentPage < lastPage {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", path, currentPage+1, perPage)
		nextLink.URL = &url
	}

	links = append(links, nextLink)

	return links
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
