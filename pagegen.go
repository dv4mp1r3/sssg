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

type GeneralHtmlData struct {
	Url     string
	Content string
	Config  config.Config
}

func GenPreviews(pageName string, posts *[]Post, c *config.Config, tpl template.Template) (string, string) {

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
			Post{
				GeneralHtmlData: GeneralHtmlData{
					Url:     post.Url,
					Content: previewText,
					Config:  *c,
				},
				Tags:    post.Tags,
				Folders: post.Folders,
				Time:    post.Time,
				Path:    post.Path,
			},
		)

		divs += b.String()
		divs += "\n"
	}

	pagePath := path.Join(c.ResultPath, pageName+".html")
	return divs, pagePath
}

func genPreviewByBreak(postContent string, c *config.Config) string {
	postParts := strings.Split(postContent, c.PreviewPageBreakString)
	if len(postParts) != 2 {
		return postContent
	}
	return postParts[0]
}

func genPreviewByLength(postContent string, c *config.Config) string {
	p := bluemonday.StrictPolicy()
	tmp := wordwrap.WrapString(p.Sanitize(postContent), uint(c.PreviewLength))
	if len(tmp) == 0 {
		return ""
	}
	return strings.Split(tmp, "\n")[0]
}

func GenPreviewText(postContent string, c *config.Config) string {
	if len(postContent) == 0 {
		return ""
	}

	if c.PreviewByPageBreak {
		return genPreviewByBreak(postContent, c)
	} else {
		return genPreviewByLength(postContent, c)
	}

}
