package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"

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

func writePaginationPages(posts *[]Post, tplParam template.Template, c *config.Config, customPath string) {
	pLen := len(*posts)
	pages := pLen / c.PostsPerPage
	if pages == 0 && pLen > 0 {
		pages = 1
	} else if pages > 1 {
		pages += 1
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
		previewDivs, pagePath := GenPreviews(
			pageName,
			&pagePosts,
			c,
			tplParam,
		)
		writePage(
			PageData{
				Post: Post{
					GeneralHtmlData: GeneralHtmlData{
						Content: previewDivs,
						Config:  *c,
					},
				},
				PaginationData: paginationElements,
			},
			pagePath,
			&tplParam,
		)
	}
}

func writePage(pgd PageData, destPath string, tplParam *template.Template) string {
	if destPath == "" {
		return destPath
	}

	var b bytes.Buffer
	tplParam.ExecuteTemplate(
		&b,
		genTemplateName("page"),
		pgd,
	)

	err := os.WriteFile(destPath, b.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return destPath
}

func genPostFile(post *Post, categories *[]Category, c *config.Config, tplParam *template.Template) *Post {

	fmt.Println(post.Path)
	fp := GenFullSourcePath(c, post)
	cnt := common.ReadFile(fp)

	if len(cnt) == 0 {
		return post
	}

	_html := GenPostHtml(&cnt)
	destPath := makePagePath(
		GenFullDestPath(c, post),
		strings.Replace(post.Path, ".md", ".html", 1),
	)
	post.GeneralHtmlData.Content = _html
	post.GeneralHtmlData.Config = *c
	writePage(
		PageData{
			PaginationData: []PaginationElement{},
			Post:           *post,
			PublishDate:    "",
		},
		destPath,
		tplParam,
	)

	post.Content = _html
	post.Url = GenPostUrl(c, &destPath)

	return post
}

func makePagePath(destPath string, postFilename string) string {
	err := os.MkdirAll(destPath, 0755)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return path.Join(destPath, postFilename)
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
	if runtime.GOOS == "windows" {
		exec.Command("rmdir", "/s", result).Output()
		source = winPathFormat(source)
		result = winPathFormat(result)
		out, err := exec.Command("xcopy", source, result, "/E", "/H").Output()
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Println(out)
	} else {
		exec.Command("rm", "-rf", result).Output()
		out, err := exec.Command("cp", "-R", source, result).Output()
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Println(out)
	}
}

func winPathFormat(path string) string {
	return fmt.Sprintf(".\\%s\\", strings.Replace(path, "/", "\\", 1))
}

func genTemplateName(name string) string {
	return name + ".html"
}

func main() {
	c := processInput()

	var categories []Category

	var a, _ = template.ParseGlob(path.Join(c.SourcePath, "*.html"))
	var posts []Post
	err := getPosts(&posts, filepath.Join(c.SourcePath, "content"), []string{}, 3, 1, &c)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})
	if err != nil {
		fmt.Println(err)
	}

	postsByCategory := make(map[string][]Post)
	for idx, post := range posts {
		posts[idx] = *genPostFile(&post, &categories, &c, a)
		tryToUpdateCategories(&categories, &post, &c)
		pbcKey := ""
		if len(post.Folders) > 0 {
			pbcKey = JoinFolders(&post)
			postsByCategory[pbcKey] = append(postsByCategory[pbcKey], post)
		}
		posts[idx].Tags = GetPostTags(&c, path.Join(pbcKey, posts[idx].Path))
		TryToUpdateTagInfo(posts[idx].Tags, &post)
	}

	writePaginationPages(&posts, *a, &c, "")
	for idx := range categories {
		pbc := postsByCategory[categories[idx].Path]
		sort.Slice(pbc, func(i, j int) bool {
			return pbc[i].Time.After(pbc[j].Time)
		})
		writePaginationPages(&pbc, *a, &c, categories[idx].Path)
	}

	for tag := range GetUniqueTags() {
		pbt := postsByTag[tag]
		sort.Slice(pbt, func(i, j int) bool {
			return pbt[i].Time.After(pbt[j].Time)
		})
		makePagePath(path.Join(c.ResultPath, tag), "")
		writePaginationPages(&pbt, *a, &c, tag)
	}

	copyStatic(&c)

}
