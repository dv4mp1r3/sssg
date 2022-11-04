package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dv4mp1r3/ovpngen/common"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mitchellh/go-wordwrap"
)

func main() {
	configPath := flag.String("c", "config.json", "Config path")
	flag.Parse()

	configContent := common.ReadFile(*configPath)
	var c Config
	err := createConfig(configContent, &c)
	if err != nil {
		panic(err)
	}
	if !validateConfig(&c) {
		fmt.Println("Config is invalid. The programm will be stop")
		return
	}

	var categories []Category

	templateName := "page"
	pageTemplate := createPage(&c, templateName, true)

	var posts []Post
	err = getPosts(&posts, filepath.Join(c.SourcePath, "content"), []string{}, 3, 1)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.Before(posts[j].Time)
	})
	if err != nil {
		fmt.Println(err)
	}
	for idx, post := range posts {
		fmt.Println(post.Path)
		fp := GenFullSourcePath(&c, &post)
		cnt := common.ReadFile(fp)

		if len(cnt) > 0 {
			//todo: fix
			_html := genPostHtml(&cnt)
			postPage := strings.Replace(pageTemplate, "{{.content}}", _html, 1)
			//todo: generate internal urls

			destPath := GenFullDestPath(&c, &post)
			err = os.MkdirAll(destPath, 0755)
			if err != nil {
				fmt.Println(err)
				continue
			}

			catUrl := getCategoryUrlByPost(&post)
			catLabel := getCategoryNameByPost(&post)
			if len(catUrl) > 0 && len(catLabel) > 0 {
				if isUniqueCategory(categories, &catUrl) {
					categories = append(categories, Category{Url: catUrl, Label: catLabel})
				}
			}

			destPath = path.Join(destPath, strings.Replace(post.Path, ".md", ".html", 1))
			err = os.WriteFile(destPath, []byte(postPage), 0644)
			if err != nil {
				fmt.Println(err)
			}

			posts[idx].Content = _html
			posts[idx].Url = genPostUrl(&c, &destPath)

		}

	}

	pages := len(posts) / c.PostsPerPage
	currentPage := 0
	pageName := "index"
	for currentPage < pages {
		startPost := currentPage * c.PostsPerPage
		endPost := startPost + c.PostsPerPage
		pagePosts := posts[startPost:endPost]
		fmt.Println(pagePosts)
		if currentPage > 0 {
			pageName = fmt.Sprint(currentPage + 1)
		}
		pageUrls := genPageUrls(pages, currentPage)
		currentPage++
		genPaginationPages(
			pageName,
			pageTemplate,
			&pagePosts,
			&pageUrls,
			&c,
		)

	}

	for _, category := range categories {
		fmt.Println(category.Label, category.Url)
	}

}

func genPageUrls(pageCount int, currentPage int) string {
	//todo: implement
	return ""
}

func genPreviewText(postContent string, c *Config) string {
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

func genPaginationPages(pageName string, pageTemplate string, posts *[]Post, pageUrls *string, c *Config) {
	cStr := "{{.content}}"
	divs := ""
	for _, post := range *posts {
		if len(post.Content) == 0 {
			continue
		}
		divContent := createPage(c, "pagination", false)
		contentIndex := strings.Index(divContent, cStr)
		if contentIndex >= 0 {
			previewText := genPreviewText(post.Content, c)
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

func renderHookDropCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// skip all nodes that are not CodeBlock nodes
	n, ok := node.(*ast.CodeBlock)
	if n != nil {
		fmt.Println(n.Content)
	}

	if !ok {
		return ast.GoToNext, false
	}

	// custom rendering logic for ast.CodeBlock. By doing nothing it won't be
	// present in the output

	return ast.GoToNext, true
}

func createPage(c *Config, templateName string, isIndexPage bool) string {
	pageContent := ""
	fp := filepath.Join(c.SourcePath, templateName+".html")
	templateContent := readTemplate(&fp, &templateName)
	templateObject := template.New(templateName)
	pageContent = templateContent
	if templateContent != "" {
		l := parseTemplate(templateObject, &templateContent)
		for idx, str := range l {
			str = strings.Replace(strings.TrimSpace(str), "{{.", "", 1)
			str = strings.Replace(str, "}}", "", 1)
			if len(str) == 0 {
				continue
			}
			if str != "content" && str != "url" {
				includedContent := createPage(c, str, false)
				pageContent = strings.ReplaceAll(pageContent, l[idx], includedContent)
			}
		}
	}

	return pageContent
}

func genPostHtml(postPageMd *string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.FencedCode
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: renderHookDropCodeBlock,
	}
	parser := parser.NewWithExtensions(extensions)
	renderer := html.NewRenderer(opts)
	html := markdown.ToHTML([]byte(*postPageMd), parser, renderer)
	return string(html)
}

func genPostUrl(c *Config, destPath *string) string {
	res := c.Url
	tmp := c.ResultPath
	if strings.HasPrefix(c.ResultPath, ".") {
		tmp = strings.TrimLeft(c.ResultPath, ".")
	}
	if strings.HasPrefix(tmp, "/") {
		tmp = strings.TrimLeft(tmp, "/")
	}
	resPath := strings.Replace(*destPath, tmp, "", 1)
	return res + strings.ReplaceAll(resPath, string(os.PathSeparator), "/")
}

func getCategoryNameByPost(post *Post) string {
	l := len(post.Folders)
	if l == 0 {
		return ""
	}
	name := strings.TrimSpace(post.Folders[l-1])
	return name
}

func isUniqueCategory(c []Category, url *string) bool {
	for _, category := range c {
		if category.Url == *url {
			return false
		}
	}
	return true
}

func getCategoryUrlByPost(post *Post) string {
	l := len(post.Folders)
	if l == 0 {
		return ""
	}
	url := ""
	for _, folder := range post.Folders {
		if len(url) > 0 {
			url += "/"
		}
		url += strings.TrimSpace(folder)
	}
	return url
}
