package model

type Config struct {
	Content struct {
		Systems []string `yaml:"systems"`
	} `yaml:"content"`
}
