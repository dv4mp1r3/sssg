package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template/parse"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mitchellh/go-wordwrap"
)

func listTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}

func readTemplate(path *string, name *string) string {
	tpl, err := os.ReadFile(*path)
	if err != nil {
		return ""
	}
	return string(tpl)

}

func parseTemplate(tObject *template.Template, tContent *string) []string {
	t := template.Must(tObject.Parse(*tContent))
	l := listTemplFields(t)
	return l
}

func GenPageUrls(pageCount int, currentPage int) string {
	//todo: implement
	return ""
}

func GenPreviewText(postContent string, c *Config) string {
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

func CreatePage(c *Config, templateName string, isIndexPage bool, data any) string {
	const fileDelemiter = "file_"
	pageContent := ""
	//
	if strings.Index(templateName, fileDelemiter) == 0 {
		_, templateName, _ = strings.Cut(templateName, fileDelemiter)
	}
	//
	fp := filepath.Join(c.SourcePath, templateName+".html")
	templateContent := readTemplate(&fp, &templateName)
	templateObject := template.New(templateName)
	pageContent = templateContent
	if templateContent != "" {
		l := parseTemplate(templateObject, &templateContent)
		for idx, str := range l {
			str = strings.Replace(strings.TrimSpace(str), "{{.", "", 1)
			str = strings.Replace(str, "}}", "", 1)
			if len(str) == 0 && strings.Index(str, fileDelemiter) != 0 {
				continue
			}
			//
			if str != "content" && str != "url" {
				includedContent := CreatePage(c, str, false, data)
				pageContent = strings.ReplaceAll(pageContent, l[idx], includedContent)
			}
		}
	}

	if data == nil || reflect.ValueOf(data).Kind() == reflect.Invalid {
		return pageContent
	} else {
		templateObject = template.Must(templateObject.Parse(pageContent))
		var b bytes.Buffer
		templateObject.Execute(&b, data)
		return b.String()
	}

}
