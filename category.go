package main

import "strings"

type Category struct {
	Url   string
	Label string
}

func IsUniqueCategory(c []Category, url *string) bool {
	for _, category := range c {
		if category.Url == *url {
			return false
		}
	}
	return true
}

func GetCategoryUrlByPost(post *Post) string {
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

func GetCategoryNameByPost(post *Post) string {
	l := len(post.Folders)
	if l == 0 {
		return ""
	}
	name := strings.TrimSpace(post.Folders[l-1])
	return name
}
