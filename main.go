package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/m3ng9i/feedreader"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

//go:embed frontend/dist/*
var f embed.FS

// options
var updateDB bool

// database struct
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

// config file
type Config struct {
	Version int      `yaml:"version"`
	Feeds   []string `yaml:"feeds"`
}

func main() {
	flag.BoolVar(&updateDB, "fetch", false, "Fetch and generate feed database")
	flag.Parse()

	config := LoadConfig[Config]()
	feed := new(Feed)

	log.Println("Fetching feeds...")
	if updateDB {
		FetchFeed(feed, config)
		ExportDB(feed)
		return
	}
	go func() {
		FetchFeed(feed, config)
		ExportDB(feed)
	}()

	log.Println("Starting server...")
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/api/v1/feed", GetFeed(feed))
	r.Static("/db/", "/db/")
	r.StaticFile("/db.json", "./db.json")
	r.StaticFile("/index.json", "./index.json")

	subFS, _ := fs.Sub(f, "frontend/dist")
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(subFS))))

	crontab := cron.New(cron.WithSeconds())
	if _, err := crontab.AddFunc("0 15 * * * *", func() {
		FetchFeed(feed, config)
		ExportDB(feed)
	}); err != nil {
		log.Fatal("Failed to start feed update daemon")
	}
	crontab.Start()
	log.Println("Feed update daemon started.")
	log.Fatal(r.Run(":8192"))
}

func LoadConfig[Config any]() *Config {
	if _, err := os.Stat("config.yml"); err != nil {
		data, _ := yaml.Marshal(new(Config))
		os.WriteFile("config.yml", []byte(data), 0644)
		log.Fatal(errors.New("config file not found, a template file has been created"))
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("config file read failed")
	}
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("config file parse failed")
	}
	return config
}

func FetchFeed(feed *Feed, config *Config) {
	wg, mutex := new(sync.WaitGroup), new(sync.Mutex)

	for _, url := range config.Feeds {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			res, err := feedreader.Fetch(url)
			if err != nil {
				log.Println("Fetch RSS failed:", err)
				return
			}
			if res.Author == nil {
				res.Author = &feedreader.FeedPerson{Name: "Unknown"}
			}
			articles := make([]Article, len(res.Items))
			for i, item := range res.Items {
				articles[i] = Article{item.Title, item.PubDate, item.Content, item.Link}
			}
			mutex.Lock()
			feed.Author = append(feed.Author, Author{res.Title, res.Author.Email, res.Link, res.Description, articles})
			mutex.Unlock()
			log.Println("Fetched RSS:", url)
		}(url)
	}
	wg.Wait()
	feed.Update = time.Now()
	log.Println("Fetch RSS done.")
}

func ExportDB(feed *Feed) {
	log.Println("Exporting db...")
	data, err := json.Marshal(feed)
	if err != nil {
		log.Fatal("Failed to marshal db:", err)
	}
	if err := os.WriteFile("db.json", data, 0644); err != nil {
		log.Panic("Failed to write db:", err)
	}
	if err := os.MkdirAll("db", 0777); err != nil && !os.IsExist(err) {
		log.Fatal("Failed to create db directory:", err)
	}
	for i, author := range feed.Author {
		for j, article := range author.Article {
			fileName := fmt.Sprintf("db/%d_%d_%s.txt", i, j, url.PathEscape(article.Title))
			feed.Author[i].Article[j].Content = fileName
			if err := os.WriteFile(fileName, []byte(article.Content), 0644); err != nil {
				log.Println("Failed to write:", fileName)
			} else {
				log.Println("Wrote:", fileName)
			}
		}
	}
	data, err = json.Marshal(feed)
	if err != nil {
		log.Fatal("Failed to marshal index:", err)
	}
	if err := os.WriteFile("index.json", data, 0644); err != nil {
		log.Panic("Failed to write index:", err)
	}
	log.Println("Export db done.")
}

func GetFeed(feed *Feed) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, feed)
	}
}
