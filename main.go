package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/m3ng9i/feedreader"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

//go:embed frontend/dist/*
var f embed.FS
var config *Config
var feed *Feed

type Feed struct {
	Version int       `json:"version"`
	Author  []Author  `json:"author"`
	Update  time.Time `json:"update"`
}
type Author struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Uri         string    `json:"uri"`
	Description string    `json:"description"`
	Article     []Article `json:"article"`
}
type Article struct {
	Title   string    `json:"title"`
	Time    time.Time `json:"time"`
	Content string    `json:"content"`
	Url     string    `json:"url"`
}

type Config struct {
	Version int      `yaml:"version"`
	Feeds   []string `yaml:"feeds"`
}

func init() {
	// Load/Initialize the config
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Println("No config file found, creating...")
		config = &Config{Version: 1, Feeds: []string{""}}
		SaveConfig()
	}
	if err2 := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err2)
	}
	// Load articles from db
	db, err := os.ReadFile("db.json")
	if err != nil {
		log.Println("No db found, creating...")
		feed = &Feed{1, nil, time.Now()}
		SaveDB()
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

func main() {
	var updateDB bool

	flag.BoolVar(&updateDB, "fetch", false, "Fetch and update feed database")
	flag.Parse()

	switch {
	case updateDB:
		FetchFeed()
	default:
		r := gin.Default()
		r.Use(cors.Default())
		initRouter(r)
		crontab := cron.New(cron.WithSeconds())
		if _, err := crontab.AddFunc("0 15 * * * *", FetchFeed); err != nil {
			panic("Failed to start feed update daemon")
		}
		crontab.Start()
		panic(r.Run(":8192"))
	}
}

// utils
func SaveConfig() {
	data, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal("Failed to marshal config: ", err)
	}
	if err := os.WriteFile("config.yml", data, 0644); err != nil {
		log.Panic("Failed to write config")
	}
}
func SaveDB() {
	// write to db.json
	data, err := json.Marshal(feed)
	if err != nil {
		log.Fatal("Failed to marshal db: ", err)
	}
	if err := os.WriteFile("db.json", data, 0644); err != nil {
		log.Panic("Failed to write db")
	}
	// create index db
	var feedCopy Feed
	json.Unmarshal(data, &feedCopy) // deep copy the array
	if err := os.MkdirAll("db", 0644); err != nil && !os.IsExist(err) {
		log.Fatalln("Failed to create db directory...")
	}
	for i, author := range feed.Author {
		for j, article := range author.Article {
			fileName := fmt.Sprintf("db/%d_%d_%s.txt", i, j, article.Title)
			if err := os.WriteFile(fileName, []byte(article.Content), 0644); err != nil {
				log.Println("Failed to write to db directory")
			}
			feedCopy.Author[i].Article[j].Content = fileName
		}
	}
	data, err = json.Marshal(feedCopy)
	if err != nil {
		log.Fatal("Failed to marshal index: ", err)
	}
	if err := os.WriteFile("index.json", data, 0644); err != nil {
		log.Panic("Failed to write db")
	}
}
func FetchFeed() {
	log.Println("Fetching feeds...")
	feed.Author = nil
	for _, data := range config.Feeds {
		res, err := feedreader.Fetch(data)
		if err != nil {
			log.Println("Fetch RSS failed: ", err)
		} else {
			// !!! res.Author might be NULL !!!
			if res.Author == nil {
				res.Author = &feedreader.FeedPerson{Name: "Unknown"}
			}
			var articles []Article
			for _, item := range res.Items {
				articles = append(articles, Article{item.Title, item.PubDate, item.Content, item.Link})
			}
			feed.Author = append(feed.Author, Author{res.Title, res.Author.Email, res.Link, res.Description, articles})
		}
	}
	feed.Update = time.Now()
	SaveDB()
	log.Println("Fetch RSS done.")
}
func initRouter(r *gin.Engine) {
	// APIs
	apiRouter := r.Group("/api/v1")
	{
		apiRouter.GET("/feed", GetFeed)             // Get feed url list
		apiRouter.PUT("/comment", GetFeed)          // Send comment by article ID
		apiRouter.GET("/comment/:comment", GetFeed) // Get comment by article ID
		apiRouter.PUT("/feed", AddFeed)
	}
	r.Static("/db/", "/db/")
	r.StaticFile("/db.json", "./db.json")
	r.StaticFile("/index.json", "./index.json")

	subFS, err := fs.Sub(f, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(subFS))))
}

// Controllers
func GetFeed(c *gin.Context) {
	c.JSON(http.StatusOK, feed)
}
func AddFeed(c *gin.Context) {
	var url string
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.Feeds = append(config.Feeds, url)
	SaveConfig()
	FetchFeed()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
