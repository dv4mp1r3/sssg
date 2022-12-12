package main

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/dv4mp1r3/sssg/config"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mitchellh/go-wordwrap"
)

type previewData struct {
	Url     string
	Content string
}

func GenPreviews(pageName string, posts *[]Post, paginationElements *string, c *config.Config, tpl template.Template) (string, string) {

	divs := ""
	const templateName = "post_preview"
	var m = make(map[string]any)
	for _, post := range *posts {
		if len(post.Content) == 0 {
			continue
		}
		m[templateName] = post

		previewText := GenPreviewText(post.Content, c)

		var b bytes.Buffer
		tpl.ExecuteTemplate(
			&b,
			fmt.Sprint(templateName, ".html"),
			previewData{
				Url:     post.Url,
				Content: previewText,
			},
		)

		divs += b.String()
		divs += "\n"
	}

	pagePath := path.Join(c.ResultPath, pageName+".html")
	return divs, pagePath
}

func GenPreviewText(postContent string, c *config.Config) string {
	if len(postContent) == 0 {
		return ""
	}
	p := bluemonday.StrictPolicy()
	tmp := wordwrap.WrapString(p.Sanitize(postContent), uint(c.PreviewLength))
	if len(tmp) == 0 {
		return ""
	}
	return strings.Split(tmp, "\n")[0]
}
