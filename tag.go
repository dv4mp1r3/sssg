package main

import (
	"github.com/dv4mp1r3/sssg/config"
)

type Tag struct {
	Key string
	Url string
}

var postsByTag = make(map[string][]Post)

func GetPostTags(c *config.Config, filename string) []string {
	if c.Tags[filename] != nil {
		return c.Tags[filename]
	}
	return []string{}
}

func TryToUpdateTagInfo(tags []string, post *Post) {
	for _, tag := range tags {
		if postsByTag[tag] == nil {
			postsByTag[tag] = []Post{}
		}
		postsByTag[tag] = append(postsByTag[tag], *post)
	}
}

func GetUniqueTags() map[string][]Post {
	return postsByTag
}
