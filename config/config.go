package config

import (
	"encoding/json"
	"fmt"
)

type MenuElement struct {
	Label string
	Url   string
}

type Config struct {
	Label         string
	Menu          []MenuElement
	PostsPerPage  int
	PreviewLength int
	SourcePath    string
	ResultPath    string
	StaticPath    string
	Url           string
}

var c Config

func ValidateConfig(c *Config) bool {
	if len(c.ResultPath) == 0 || len(c.SourcePath) == 0 || len(c.Url) == 0 {
		return false
	}
	return true
}

func CreateConfig(s string) (Config, error) {
	var err = json.Unmarshal([]byte(s), &c)

	if err != nil {
		fmt.Printf("Error on createConfig %s\n", err)
		return c, err
	}
	return c, nil
}

func GetInstance() (*Config, error) {
	return &c, nil
}
