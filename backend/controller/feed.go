package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/m3ng9i/feedreader"
	"gopkg.in/yaml.v3"
)

type FeedConfig struct {
	Version int      `yaml:"version"`
	Feeds   []string `yaml:"feeds"`
}

var config *FeedConfig

func init() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err, "No config file found, creating...")
		config = &FeedConfig{Version: 1, Feeds: []string{""}}
		data, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err, "Failed to marshal config")
		}
		if err := os.WriteFile("config.yaml", data, 0644); err != nil {
			log.Panic("Failed to write config")
		}
	}
	if err2 := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err2)
	}
}

func FetchFeed(c *gin.Context) {
	result := make([]feedreader.Feed, 0)
	for index, feed := range config.Feeds {
		res, err := feedreader.Fetch(feed)
		if err != nil {
			log.Fatal(err)
		} else {
			result[index] = *res
		}
	}
	c.JSON(http.StatusOK, result)
}

func AddFeed(feed string) {
	config.Feeds = append(config.Feeds, feed)
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	if err := os.WriteFile("config.yaml", data, 0644); err != nil {
		log.Panic("Failed to write config")
	}
}
