package main

import (
	"fmt"

	"github.com/dv4mp1r3/sssg/config"
)

type PaginationElement struct {
	Type   string
	Custom string
	Value  int
}

const prevType = "previous"
const nextType = "next"

func isFirstPage(currentPage int) bool {
	return currentPage == 0
}

func isLastPage(currentPage int, pageCount int) bool {
	return currentPage == pageCount-1
}

func hasNextPage(currentPage int, pageCount int) bool {
	return pageCount-currentPage >= 2
}

func genPageUrl(c *config.Config, val int, customPath string) string {
	if customPath != "" {
		customPath += "/"
	}
	if val == 1 {
		return fmt.Sprintf("%s/%s%s.html", c.Url, customPath, "index")
	}
	return fmt.Sprintf("%s/%s%d.html", c.Url, customPath, val)
}

func genMaxTwoPaginationButtons(pageCount int, activePage int, c *config.Config, m []PaginationElement, customPath string) []PaginationElement {
	hnp := hasNextPage(activePage, pageCount)
	ifp := isFirstPage(activePage)
	val := 0
	if ifp && hnp {
		val = activePage + 2
		m = append(m, PaginationElement{Type: nextType, Custom: genPageUrl(c, val, customPath), Value: val})
	} else if isLastPage(activePage, pageCount) {
		val = activePage
		m = append(m, PaginationElement{Type: prevType, Custom: genPageUrl(c, val, customPath), Value: val})
	} else {
		val = activePage
		m = append(m, PaginationElement{Type: prevType, Custom: genPageUrl(c, val, customPath), Value: val})
		if hnp {
			val = activePage + 2
			m = append(m, PaginationElement{Type: nextType, Custom: genPageUrl(c, val, customPath), Value: val})
		}
	}
	return m
}

func genAllPaginationButtons(pageCount int, activePage int, c *config.Config, m []PaginationElement, customPath string) []PaginationElement {
	currentPage := 0
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

func GenPaginationElements(pageCount int, activePage int, c *config.Config, customPath string) []PaginationElement {

	m := []PaginationElement{}
	if pageCount <= 1 {
		return m
	}

	if c.MaxTwoPaginationButtons {
		return genMaxTwoPaginationButtons(pageCount, activePage, c, m, customPath)
	}
	return genAllPaginationButtons(pageCount, activePage, c, m, customPath)

}
