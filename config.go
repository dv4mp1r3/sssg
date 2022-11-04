package main

import (
	"encoding/json"
	"fmt"
)

type menuElement struct {
	Label string
	Url   string
}

type Config struct {
	Label         string
	Menu          []menuElement
	PostsPerPage  int
	PreviewLength int
	SourcePath    string
	ResultPath    string
	StaticPath    string
	Url           string
}

func validateConfig(c *Config) bool {
	if len(c.ResultPath) == 0 || len(c.SourcePath) == 0 || len(c.Url) == 0 {
		return false
	}
	return true
}

func createConfig(s string, c *Config) error {
	var err = json.Unmarshal([]byte(s), &c)
	if err != nil {
		fmt.Println(fmt.Sprintln("Error on createConfig %s", err))
		return err
	}
	return nil
}
