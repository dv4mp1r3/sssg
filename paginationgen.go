package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type PaginationElement struct {
	Url     string
	Preview string
}

func GenPaginationPages(pageName string, pageTemplate string, posts *[]Post, pageUrls *string, c *Config) {
	cStr := "{{.content}}"
	divs := ""
	for _, post := range *posts {
		if len(post.Content) == 0 {
			continue
		}
		divContent := CreatePage(c, "pagination", false, nil)
		contentIndex := strings.Index(divContent, cStr)
		if contentIndex >= 0 {
			previewText := GenPreviewText(post.Content, c)
			divContent = ReplaceAtIndex(divContent, []rune(previewText), contentIndex, len(cStr))
			divs += divContent
			divs += "\n"
		}
	}

	if pageUrls != nil && len(*pageUrls) > 0 {
		divs += *pageUrls
	}

	if len(divs) > 0 {
		contentIndex := strings.Index(pageTemplate, cStr)
		if contentIndex >= 0 {
			pageContent := ReplaceAtIndex(pageTemplate, []rune(divs), contentIndex, len(cStr))
			pagePath := path.Join(c.ResultPath, pageName+".html")
			err := os.WriteFile(pagePath, []byte(pageContent), 0644)
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}
