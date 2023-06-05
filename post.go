package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dv4mp1r3/sssg/config"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	GeneralHtmlData
	Path    string
	Folders []string
	Time    time.Time
	Tags    []string
}

type PageData struct {
	Post
	PublishDate    string
	PaginationData []PaginationElement
}

func getPosts(posts *[]Post, root string, dirs []string, maxLevel int, currentLevel int, c *config.Config) error {

	f, err := os.Open(root)
	if err != nil {
		return err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}

	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".md") {
			post := Post{
				Path:    file.Name(),
				Folders: dirs,
				//Tags:    GetPostTags(c, file.Name(), dirs),
			}
			pathInCategory := path.Join(JoinFolders(&post), file.Name())
			post.Time = getPostTime(pathInCategory, file, c)
			*posts = append(*posts, post)
		}

		if file.IsDir() && currentLevel <= maxLevel {
			copyDirs := append(dirs, file.Name())
			newRoot := filepath.Join(root, file.Name())
			getPosts(posts, newRoot, copyDirs, maxLevel, currentLevel+1, c)
		}
	}
	return nil
}

func getPostTime(key string, fInfo fs.FileInfo, c *config.Config) time.Time {
	if dateVal, ok := c.PostDates.CustomDates[key]; ok {
		t, err := time.Parse(c.PostDates.Layout, dateVal)
		if err == nil {
			return t
		}
	}
	return getCtime(fInfo)
}

func JoinFolders(p *Post) string {
	res := ""
	for _, fld := range p.Folders {
		res = path.Join(res, fld)
	}
	return res
}

func GenFullSourcePath(c *config.Config, post *Post) string {
	res := path.Join(c.SourcePath, "content", JoinFolders(post))
	return path.Join(res, post.Path)
}

func GenFullDestPath(c *config.Config, post *Post) string {
	return path.Join(c.ResultPath, JoinFolders(post))
}

func GenPostHtml(postPageMd *string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Attributes
	opts := html.RendererOptions{
		Flags:          html.UseXHTML,
		RenderNodeHook: renderHookDropCodeBlock,
	}
	parser := parser.NewWithExtensions(extensions)
	renderer := html.NewRenderer(opts)
	html := markdown.ToHTML([]byte(*postPageMd), parser, renderer)
	return string(html)
}

func renderHookDropCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	n, ok_n := node.(*ast.Link)
	if n != nil && ok_n {
		tmp := fixInternalUrls(&n.Destination)
		n.Destination = tmp
		return ast.GoToNext, false
	}

	i, ok_i := node.(*ast.Image)
	if i != nil && ok_i {
		i.Destination = fixInternalUrls(&i.Destination)
		return ast.GoToNext, false
	}

	return ast.GoToNext, false
}

func fixInternalUrls(destination *[]byte) []byte {
	c, err := config.GetInstance()
	if err != nil {
		return *destination
	}
	_url := string(*destination)
	if strings.Index(_url, "./") != 0 {
		return *destination
	}

	mdExtIndex := strings.LastIndex(_url, ".md")

	if mdExtIndex == len(_url)-3 {
		_url = _url[2:mdExtIndex]
		_url += ".html"
	}

	parentDirIndex := strings.LastIndex(_url, "../")
	if parentDirIndex != -1 {
		_url = _url[parentDirIndex+3:]
	}
	return []byte(fmt.Sprint(c.Url, "/", _url))
}

func GenPostUrl(c *config.Config, destPath *string) string {
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
