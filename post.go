package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Path    string
	Folders []string
	Time    time.Time
	Url     string
	Content string
}

func getPosts(posts *[]Post, root string, dirs []string, maxLevel int, currentLevel int) error {

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
			post := Post{Path: file.Name(), Folders: dirs, Time: getCtime(file)}
			*posts = append(*posts, post)
		}

		if file.IsDir() && currentLevel <= maxLevel {
			copyDirs := append(dirs, file.Name())
			newRoot := filepath.Join(root, file.Name())
			currentLevel++
			getPosts(posts, newRoot, copyDirs, maxLevel, currentLevel)
		}
	}
	return nil
}

func getCtime(fInfo fs.FileInfo) time.Time {

	stat := fInfo.Sys().(*syscall.Stat_t)
	return time.Unix(int64(getCtimeSec(stat)), int64(getCtimeNSec(stat)))
}

func joinFolders(p *Post) string {
	res := ""
	for _, fld := range p.Folders {
		res = path.Join(res, fld)
	}
	return res
}

func GenFullSourcePath(c *Config, post *Post) string {
	res := path.Join(c.SourcePath, "content", joinFolders(post))
	return path.Join(res, post.Path)
}

func GenFullDestPath(c *Config, post *Post) string {
	return path.Join(c.ResultPath, joinFolders(post))
}

func GenPostHtml(postPageMd *string) string {
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

func GenPostUrl(c *Config, destPath *string) string {
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
