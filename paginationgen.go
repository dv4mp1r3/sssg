package main

import (
	"github.com/dv4mp1r3/sssg/config"
)

type PaginationElement struct {
	Type   string
	Custom string
	Value  int
}

func GenPaginationElements(pageCount int, activePage int, c *config.Config) string {
	elements := ""
	currentPage := 0
	m := make(map[string]any)
	if pageCount <= 1 {
		return elements
	}

	const tName = "pagination_element"
	for currentPage < pageCount {
		var t string
		if currentPage == activePage {
			t = "active"
		} else {
			t = "inactive"
		}
		currentPage++
		m[tName] = PaginationElement{Type: t, Custom: "", Value: currentPage}
		//elements += CreatePageFromFile(c, tName, false, m)
	}
	return elements
}
