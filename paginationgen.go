package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dv4mp1r3/sssg/config"
)

type PaginationElement struct {
	Url     string
	Preview string
}

type pData struct {
	PaginationData string
	DrawPagination bool
	Menu           string
	Content        string
}

func GenPaginationPages(pageName string, pageTemplate string, posts *[]Post, pageUrls *string, c *config.Config) {

	divs := ""
	const templateName = "pagination"
	var m = make(map[string]any)
	for _, post := range *posts {
		if len(post.Content) == 0 {
			continue
		}
		m[templateName] = post
		divContent := CreatePageFromFile(c, templateName, false, m)
		contentIndex := strings.Index(divContent, post.Content)
		if contentIndex >= 0 {
			previewText := GenPreviewText(post.Content, c)
			divContent = ReplaceAtIndex(divContent, []rune(previewText), contentIndex, len(post.Content))
			divs += divContent
			divs += "\n"
		}
	}

	if pageUrls != nil && len(*pageUrls) > 0 {
		divs += *pageUrls
	}

	//todo: redo
	m[templateName] = pData{DrawPagination: true, PaginationData: "", Menu: "", Content: divs}
	pageContent := CreatePage(c, templateName, pageTemplate, false, m)
	pagePath := path.Join(c.ResultPath, pageName+".html")
	err := os.WriteFile(pagePath, []byte(pageContent), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
