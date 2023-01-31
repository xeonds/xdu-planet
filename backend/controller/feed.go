package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/m3ng9i/feedreader"
	"gopkg.in/yaml.v3"
)

type FeedConfig struct {
	Version int      `yaml:"version"`
	Feeds   []string `yaml:"feeds"`
}

type Article struct {
	Title   string
	Time    time.Time
	Content string
	Url     string
}

type Feed struct {
	Version int       `json:"version"`
	Data    []Article `json:"data"`
	Update  time.Time `json:"update"`
}

var config *FeedConfig
var feed *Feed

func init() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Println("No config file found, creating...")
		config = &FeedConfig{Version: 1, Feeds: []string{""}}
		saveConfig()
	}
	if err2 := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err2)
	}
	db, err := os.ReadFile("db.json")
	if err != nil {
		log.Println("No db found, creating...")
		feed = &Feed{Version: 1, Data: nil, Update: time.Now()}
		saveDB()
	}
	if err2 := json.Unmarshal(db, &feed); err != nil {
		log.Println(err2)
	}
}

func FetchRawXml(c *gin.Context) {
	c.JSON(http.StatusOK, feed)
}

func FetchFeed() {
	feed.Data = nil
	for _, data := range config.Feeds {
		res, err := feedreader.Fetch(data)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, item := range res.Items {
				feed.Data = append(feed.Data, Article{Title: item.Title, Time: item.PubDate, Content: getContent(item.Content), Url: item.Link})
			}
		}
	}
	saveDB()
	feed.Update = time.Now()
}

func GenPage(c *gin.Context) {
	if len(feed.Data) == 0 {
		FetchFeed()
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"list": feed.Data,
	})
}

func AddFeed(url string) {
	config.Feeds = append(config.Feeds, url)
	saveConfig()
}

func saveConfig() {
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	if err := os.WriteFile("config.yaml", data, 0644); err != nil {
		log.Panic("Failed to write config")
	}
}
func saveDB() {
	data, err := json.Marshal(feed)
	if err != nil {
		log.Fatal("Failed to marshal db: ", err)
	}
	if err := os.WriteFile("db.json", data, 0644); err != nil {
		log.Panic("Failed to write db")
	}
}

func getContent(content string) string {
	doc, err := goquery.NewDocument(content)
	if err != nil {
		log.Print("Compile article content error.")
		return content
	}
	res := doc.Find("p")
	return res.Text()
}
