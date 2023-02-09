package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m3ng9i/feedreader"
	"gopkg.in/yaml.v3"
	"xyz.xeonds/xdu-planet/util"
)

// Data model for feed
type Article struct {
	Title   string
	Time    time.Time
	Content string
	Url     string
}

type Author struct {
	Name        string
	Email       string
	Uri         string
	Description string
}

type Feed struct {
	Version int       `json:"version"`
	Article []Article `json:"article"`
	Author  []Author  `json:"author"`
	Update  time.Time `json:"update"`
}

var config *util.Config
var feed *Feed

func init() {
	// Load/Initialize the config
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Println("No config file found, creating...")
		config = &util.Config{Version: 1, Feeds: []string{""}}
		saveConfig()
	}
	if err2 := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err2)
	}
	// Load articles from db
	db, err := os.ReadFile("db.json")
	if err != nil {
		log.Println("No db found, creating...")
		feed = &Feed{1, nil, nil, time.Now()}
		saveDB()
		FetchFeed()
	}
	if err2 := json.Unmarshal(db, &feed); err != nil {
		log.Println(err2)
	}
	if time.Until(feed.Update) < -30*time.Minute {
		log.Println("Database exceeds store time. Updating...")
		FetchFeed()
	}
}

// Controllers
func FetchRawFeed(c *gin.Context) {
	c.JSON(http.StatusOK, feed)
}

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"list": feed.Article,
	})
}

// Utils
func FetchFeed() {
	log.Println("Fetching feeds...")
	// Clear current database
	feed.Article = nil
	feed.Author = nil
	// Fetch articles, authors from feeds
	for _, data := range config.Feeds {
		res, err := feedreader.Fetch(data)
		if err != nil {
			log.Println("Fetch RSS failed: ", err)
			// TODO: Load cache and return
		} else {
			// !!! res.Author might be NULL !!!
			if res.Author != nil {
				feed.Author = append(feed.Author, Author{res.Author.Name, res.Author.Email, res.Author.Uri, res.Description})
			}
			for _, item := range res.Items {
				feed.Article = append(feed.Article, Article{Title: item.Title, Time: item.PubDate, Content: getContent(item.Content), Url: item.Link})
			}
		}
	}
	// Sort articles by time desc
	sort.SliceStable(feed.Article, func(i, j int) bool {
		return feed.Article[i].Time.Unix() > feed.Article[j].Time.Unix()
	})
	feed.Update = time.Now()
	saveDB()
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
	reg := regexp.MustCompile("<p>(.*?)</p>")
	arr := reg.FindStringSubmatch(content)
	if len(arr) == 0 {
		log.Println("Parse content error: ", arr)
		return ""
	}
	return arr[len(arr)-1]
}
