package model

import (
	"fmt"
	"strconv"
	"strings"
)

type SearchQuery struct {
	Title     string
	Completed bool
}

func NewSearchQuery(q string) (*SearchQuery, error) {
	var searchQuery SearchQuery
	s := strings.Fields(q)
	if len(s) == 0 {
		return nil, fmt.Errorf("%v", "Failed to parse query.")
	}

	for _, kv := range s {
		ss := strings.SplitN(kv, ":", 2)
		if len(ss) != 2 {
			continue
		}

		k, v := ss[0], ss[1]
		if k == "title" {
			searchQuery.Title = v
		}
		if k == "completed" {
			b, e := strconv.ParseBool(v)
			if e != nil {
				return nil, fmt.Errorf("%v", "Failed to parse parameter Completed to boolean.")
			}
			searchQuery.Completed = b
		}
	}

	return &searchQuery, nil
}
