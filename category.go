package main

import (
	"strings"

	"github.com/dv4mp1r3/sssg/config"
)

type Category struct {
	Url   string
	Label string
	Path  string
}

func IsUniqueCategory(c []Category, url *string) bool {
	for _, category := range c {
		if category.Url == *url {
			return false
		}
	}
	return true
}

func GetCategoryUrlByPost(post *Post, c *config.Config) string {
	l := len(post.Folders)
	if l == 0 {
		return ""
	}
	url := c.Url
	for _, folder := range post.Folders {
		url += "/"
		url += strings.TrimSpace(folder)
	}
	return url
}

func GetCategoryNameByPost(post *Post) string {
	l := len(post.Folders)
	if l == 0 {
		return ""
	}
	name := strings.TrimSpace(post.Folders[l-1])
	return name
}
