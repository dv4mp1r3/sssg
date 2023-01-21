package main

import (
	"path"

	"github.com/dv4mp1r3/sssg/config"
)

type Tag struct {
	Key string
	Url string
}

var uniqueTags = make(map[string][]string)

func GetPostTags(c *config.Config, filename string, dirs []string) []string {
	p := ""
	for _, dir := range dirs {
		p = path.Join(p, dir)
	}

	p = path.Join(p, filename)
	for url, _ := range c.Tags {
		if url == p {
			tryToUpdateTagInfo(&p, c.Tags[url])
			return c.Tags[url]
		}
	}
	return []string{}

}

func tryToUpdateTagInfo(postPath *string, tags []string) {
	for _, tag := range tags {
		if uniqueTags[tag] == nil {
			uniqueTags[tag] = []string{}
			uniqueTags[tag] = append(uniqueTags[tag], *postPath)
			continue
		}
		var updateTagInfo = true
		for _, existedPost := range uniqueTags[tag] {
			if postPath == &existedPost {
				updateTagInfo = false
				break
			}
		}
		if updateTagInfo {
			uniqueTags[tag] = append(uniqueTags[tag], *postPath)
		}
	}
}

func GetUniqueTags() map[string][]string {
	return uniqueTags
}
