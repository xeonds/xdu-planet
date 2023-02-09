package util

type Config struct {
	Version int      `yaml:"version"`
	Feeds   []string `yaml:"feeds"`
}
