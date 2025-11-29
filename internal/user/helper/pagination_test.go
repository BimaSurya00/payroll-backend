package helper

import (
	"testing"
)

type TestItem struct {
	ID   int
	Name string
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name    string
		data    []TestItem
		page    int
		perPage int
		total   int64
		path    string
	}{
		{
			name: "First page with data",
			data: []TestItem{
				{ID: 1, Name: "Item 1"},
				{ID: 2, Name: "Item 2"},
			},
			page:    1,
			perPage: 2,
			total:   10,
			path:    "/api/v1/items",
		},
		{
			name:    "Empty result set",
			data:    []TestItem{},
			page:    1,
			perPage: 10,
			total:   0,
			path:    "/api/v1/items",
		},
		{
			name: "Last page",
			data: []TestItem{
				{ID: 9, Name: "Item 9"},
				{ID: 10, Name: "Item 10"},
			},
			page:    5,
			perPage: 2,
			total:   10,
			path:    "/api/v1/items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPagination(tt.data, tt.page, tt.perPage, tt.total, tt.path)

			// Verify basic fields
			if result.CurrentPage != tt.page {
				t.Errorf("Expected current page %d, got %d", tt.page, result.CurrentPage)
			}

			if result.PerPage != tt.perPage {
				t.Errorf("Expected per page %d, got %d", tt.perPage, result.PerPage)
			}

			if result.Total != tt.total {
				t.Errorf("Expected total %d, got %d", tt.total, result.Total)
			}

			if len(result.Data) != len(tt.data) {
				t.Errorf("Expected data length %d, got %d", len(tt.data), len(result.Data))
			}

			// Verify last page calculation
			expectedLastPage := 0
			if tt.total > 0 {
				expectedLastPage = int((tt.total + int64(tt.perPage) - 1) / int64(tt.perPage))
			}
			if result.LastPage != expectedLastPage {
				t.Errorf("Expected last page %d, got %d", expectedLastPage, result.LastPage)
			}

			// Verify from and to
			if tt.total == 0 {
				if result.From != 0 {
					t.Errorf("Expected from 0, got %d", result.From)
				}
				if result.To != 0 {
					t.Errorf("Expected to 0, got %d", result.To)
				}
			} else {
				expectedFrom := (tt.page-1)*tt.perPage + 1
				expectedTo := expectedFrom + len(tt.data) - 1

				if result.From != expectedFrom {
					t.Errorf("Expected from %d, got %d", expectedFrom, result.From)
				}
				if result.To != expectedTo {
					t.Errorf("Expected to %d, got %d", expectedTo, result.To)
				}
			}

			// Verify next page URL
			if tt.page < result.LastPage && result.NextPageURL == nil {
				t.Error("Expected next page URL but got nil")
			}
			if tt.page >= result.LastPage && result.NextPageURL != nil {
				t.Error("Expected nil next page URL but got value")
			}

			// Verify previous page URL
			if tt.page > 1 && result.PrevPageURL == nil {
				t.Error("Expected previous page URL but got nil")
			}
			if tt.page == 1 && result.PrevPageURL != nil {
				t.Error("Expected nil previous page URL but got value")
			}

			// Verify links
			if len(result.Links) < 2 {
				t.Error("Expected at least 2 links (previous and next)")
			}
		})
	}
}

func TestPaginationLinks(t *testing.T) {
	data := []TestItem{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
	}

	result := NewPagination(data, 3, 2, 10, "/api/v1/items")

	// Check that links include page numbers around current page
	hasCurrentPage := false
	for _, link := range result.Links {
		if link.Active && link.Label == "3" {
			hasCurrentPage = true
			break
		}
	}

	if !hasCurrentPage {
		t.Error("Expected links to include current page as active")
	}
}