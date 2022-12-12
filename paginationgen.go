package main

import (
	"github.com/dv4mp1r3/sssg/config"
)

type PaginationElement struct {
	Type   string
	Custom string
	Value  int
}

func GenPaginationElements(pageCount int, activePage int, c *config.Config) []PaginationElement {
	currentPage := 0
	m := []PaginationElement{}
	if pageCount <= 1 {
		return m
	}

	for currentPage < pageCount {
		var t string
		if currentPage == activePage {
			t = "active"
		} else {
			t = "inactive"
		}
		m = append(m, PaginationElement{Type: t, Custom: "", Value: currentPage})
		currentPage++
	}
	return m
}
