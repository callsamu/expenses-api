package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiltersGetsSortColumn(t *testing.T) {
	t.Run("returns sort column", func(t *testing.T) {
		filters := Filters{
			Sort:         "field",
			SortSafelist: []string{"field"},
		}

		assert.Equal(t, "field", filters.sortColumn())
	})
	t.Run("panics when sort column is not in safelist", func(t *testing.T) {
		filters := Filters{
			Sort: "field",
		}

		assert.Panics(t, func() {
			filters.sortColumn()
		})
	})
}

func TestFiltersGetSortDirection(t *testing.T) {
	t.Run("gets ascendent sort direction when does not have prefix", func(t *testing.T) {
		filters := Filters{
			Sort: "field",
		}

		assert.Equal(t, "ASC", filters.sortDirection())
	})
	t.Run("gets ascendent sort direction when does not have prefix", func(t *testing.T) {
		filters := Filters{
			Sort: "-field",
		}

		assert.Equal(t, "DESC", filters.sortDirection())
	})
}

func TestFiltersGetOffset(t *testing.T) {
	filters := Filters{
		Page:     2,
		PageSize: 10,
	}

	assert.Equal(t, 10, filters.offset())
}
