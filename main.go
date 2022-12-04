package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dv4mp1r3/ovpngen/common"
	"github.com/dv4mp1r3/sssg/config"
)

func processInput() config.Config {
	configPath := flag.String("c", "config.json", "Config path")
	flag.Parse()

	configContent := common.ReadFile(*configPath)
	c, err := config.CreateConfig(configContent)
	if err != nil {
		panic(err)
	}
	if !config.ValidateConfig(&c) {
		panic("Config is invalid. The programm will be stop")
	}
	return c
}

func needToAddCategory(url *string, label *string, categories *[]Category) bool {
	return len(*url) > 0 && len(*label) > 0 && IsUniqueCategory(*categories, url)
}

func writePaginationPages(posts *[]Post, pageTemplate *string, c *config.Config, customPath string) {
	pLen := len(*posts)
	pages := pLen / c.PostsPerPage
	if pages == 0 && pLen > 0 {
		pages = 1
	}

	currentPage := 0
	pageName := "index"
	for currentPage < pages {
		startPost := currentPage * c.PostsPerPage
		endPost := startPost + c.PostsPerPage
		if endPost > pLen {
			endPost = pLen
		}
		pagePosts := (*posts)[startPost:endPost]
		if currentPage > 0 {
			pageName = fmt.Sprint(currentPage + 1)
		}
		paginationElements := GenPaginationElements(pages, currentPage, c)
		currentPage++
		if customPath != "" {
			pageName = path.Join(customPath, pageName)
		}
		GenPreviews(
			pageName,
			*pageTemplate,
			&pagePosts,
			&paginationElements,
			c,
		)

	}
}

func writePost(post *Post, categories *[]Category, c *config.Config, pageTemplate *string) *Post {
	const templateName = "post"

	fmt.Println(post.Path)
	fp := GenFullSourcePath(c, post)
	cnt := common.ReadFile(fp)

	if len(cnt) == 0 {
		return post
	}

	_html := GenPostHtml(&cnt)
	m := make(map[string]any)
	m[templateName] = PageData{DrawPagination: false, Content: _html, Menu: c.Menu, Time: "", Tags: post.Tags}
	//todo: убрать повторный рендеринг шаблона
	tmp := CreatePageFromFile(c, "page", true, m)
	postPage := CreatePage(c, templateName, tmp, false, m)

	destPath := GenFullDestPath(c, post)
	err := os.MkdirAll(destPath, 0755)
	if err != nil {
		fmt.Println(err)
		return post
	}

	destPath = path.Join(destPath, strings.Replace(post.Path, ".md", ".html", 1))
	err = os.WriteFile(destPath, []byte(postPage), 0644)
	if err != nil {
		fmt.Println(err)
	}

	post.Content = _html
	post.Url = GenPostUrl(c, &destPath)

	return post
}

func tryToUpdateCategories(categories *[]Category, post *Post, c *config.Config) {
	catUrl := GetCategoryUrlByPost(post, c)
	catLabel := GetCategoryNameByPost(post)
	if needToAddCategory(&catUrl, &catLabel, categories) {
		c := Category{Url: catUrl, Label: catLabel, Path: JoinFolders(post)}
		*categories = append(*categories, c)
	}
}

func copyStatic(c *config.Config) {
	source := path.Join(c.SourcePath, c.StaticPath)
	result := path.Join(c.ResultPath, c.StaticPath)
	exec.Command("rm", "-rf", result).Output()
	out, err := exec.Command("cp", "-R", source, result).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(out)
}

func main() {
	c := processInput()

	var categories []Category

	pageTemplate := CreatePageFromFile(&c, "page", true, nil)

	var posts []Post
	err := getPosts(&posts, filepath.Join(c.SourcePath, "content"), []string{}, 3, 1, &c)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.Before(posts[j].Time)
	})
	if err != nil {
		fmt.Println(err)
	}

	postsByCategory := make(map[string][]Post)
	for idx, post := range posts {
		posts[idx] = *writePost(&post, &categories, &c, &pageTemplate)
		tryToUpdateCategories(&categories, &post, &c)
		if len(post.Folders) > 0 {
			pbcKey := JoinFolders(&post)
			postsByCategory[pbcKey] = append(postsByCategory[pbcKey], post)
		}
	}

	writePaginationPages(&posts, &pageTemplate, &c, "")
	for idx := range categories {
		pbc := postsByCategory[categories[idx].Path]
		writePaginationPages(&pbc, &pageTemplate, &c, categories[idx].Path)
	}

	copyStatic(&c)

}
