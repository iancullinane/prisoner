package utils

import (
	"sort"

	"github.com/iancullinane/prisoner/entity"
)

func SortByScore(entities []entity.Entity) {

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Score < entities[j].Score
	})
}
