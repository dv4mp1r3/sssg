package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dv4mp1r3/ovpngen/common"
)

func processInput() Config {
	configPath := flag.String("c", "config.json", "Config path")
	flag.Parse()

	configContent := common.ReadFile(*configPath)
	var c Config
	err := createConfig(configContent, &c)
	if err != nil {
		panic(err)
	}
	if !validateConfig(&c) {
		panic("Config is invalid. The programm will be stop")
	}
	return c
}

func needToAddCategory(url *string, label *string, categories *[]Category) bool {
	return len(*url) > 0 && len(*label) > 0 && IsUniqueCategory(*categories, url)
}

func writePaginationPages(posts *[]Post, pageTemplate *string, c *Config) {
	pages := len(*posts) / c.PostsPerPage
	currentPage := 0
	pageName := "index"
	for currentPage < pages {
		startPost := currentPage * c.PostsPerPage
		endPost := startPost + c.PostsPerPage
		pagePosts := (*posts)[startPost:endPost]
		if currentPage > 0 {
			pageName = fmt.Sprint(currentPage + 1)
		}
		pageUrls := GenPageUrls(pages, currentPage)
		currentPage++
		GenPaginationPages(
			pageName,
			*pageTemplate,
			&pagePosts,
			&pageUrls,
			c,
		)

	}
}

func writePost(post *Post, categories *[]Category, c *Config, pageTemplate *string) *Post {
	fmt.Println(post.Path)
	fp := GenFullSourcePath(c, post)
	cnt := common.ReadFile(fp)

	if len(cnt) == 0 {
		return post
	}

	//todo: fix
	_html := GenPostHtml(&cnt)
	postPage := strings.Replace(*pageTemplate, "{{.content}}", _html, 1)
	//todo: generate internal urls

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

func tryToUpdateCategories(categories *[]Category, post *Post) {
	catUrl := GetCategoryUrlByPost(post)
	catLabel := GetCategoryNameByPost(post)
	if needToAddCategory(&catUrl, &catLabel, categories) {
		*categories = append(*categories, Category{Url: catUrl, Label: catLabel})
	}
}

func main() {
	c := processInput()

	var categories []Category

	templateName := "page"
	pageTemplate := CreatePage(&c, templateName, true, nil)

	var posts []Post
	err := getPosts(&posts, filepath.Join(c.SourcePath, "content"), []string{}, 3, 1)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.Before(posts[j].Time)
	})
	if err != nil {
		fmt.Println(err)
	}

	for idx, post := range posts {
		posts[idx] = *writePost(&post, &categories, &c, &pageTemplate)
		tryToUpdateCategories(&categories, &post)
	}

	writePaginationPages(&posts, &pageTemplate, &c)

}
