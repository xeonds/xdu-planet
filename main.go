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
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/m3ng9i/feedreader"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"xyz.xeonds/xdu-planet/avalon"
	"xyz.xeonds/xdu-planet/lib"
	"xyz.xeonds/xdu-planet/model"
)

//go:embed frontend/dist/*
var f embed.FS

// options
var updateDB bool

func main() {
	// 解析命令行参数
	flag.BoolVar(&updateDB, "fetch", false, "Fetch and generate feed database")
	flag.Parse()

	config := lib.LoadConfig[model.Config]()
	db := lib.NewDB(&config.DatabaseConfig, func(db *gorm.DB) error {
		return db.AutoMigrate(&model.Comment{})
	})
	feed := new(model.Feed)

	log.Println("Fetching feeds...")
	if updateDB {
		feed = FetchFeed(config)
		ExportDB(feed)
		return
	}
	// 启用超时自动屏蔽评论
	if config.AvalonGuard.EnableGraveTimer {
		go avalon.GraveTimer(db)
	}
	// 启用根据关键词过滤被举报评论
	if config.AvalonGuard.EnableFilter {
		go avalon.Filter(db, config.AvalonGuard.Filter)
	}
	go func() {
		feed = FetchFeed(config)
		ExportDB(feed)
	}()

	log.Println("Starting server...")
	r := gin.Default()
	r.Use(cors.Default())
	api := r.Group("/api/v1")
	// 获取文章列表
	api.GET("/feed", func(c *gin.Context) {
		c.JSON(200, feed)
	})
	// 发送评论
	api.POST("/comment/:article_id", func(c *gin.Context) {
		data, article_id := new(model.Comment), c.Param("article_id")
		if err := c.ShouldBindJSON(data); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if data.ArticleId, data.Status = article_id, "ok"; db.Create(data).Error != nil {
			c.JSON(500, gin.H{"error": "failed to create comment"})
			return
		}
		c.JSON(200, gin.H{"message": "Comment added"})
	})
	// 获取评论列表
	api.GET("/comment/:article_id", func(c *gin.Context) {
		article_id := c.Param("article_id")
		comments := new([]model.Comment)
		if db.Where("article_id = ? AND status IN ?", article_id, []string{"ok", "audit"}).Find(comments).Error != nil {
			c.JSON(500, gin.H{"error": "failed to get comments"})
			return
		}
		c.JSON(200, comments)
	})
	// 举报评论
	api.DELETE("/comment/:id", func(c *gin.Context) {
		id := c.Param("id")
		if db.Model(new(model.Comment)).Where("id = ?", id).Update("status", "audit").Error != nil {
			c.JSON(500, gin.H{"error": "failed to delete comment"})
			return
		}
		c.JSON(200, gin.H{"message": "Comment deleted"})
	})
	// TODO:申请恢复自己的评论
	// api.PUT("/comment/:id", func(c *gin.Context) {
	// 	id := c.Param("id")
	// 	data := new(model.Comment)
	// 	if err := c.ShouldBindJSON(data); err != nil {
	// 		c.JSON(400, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	if db.Model(&model.Comment{}).Where("id = ?", id).Updates(data).Error != nil {
	// 		c.JSON(500, gin.H{"error": "failed to update comment"})
	// 		return
	// 	}
	// 	c.JSON(200, gin.H{"message": "Comment updated"})
	// })
	admin := api.Group("/admin")
	admin.Use(lib.LoggerMiddleware(config.LogFile))
	admin.Use(lib.JWTMiddleware(func(c *gin.Context, token string) error {
		for _, t := range config.AdminToken {
			if t == token {
				return nil
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
		return fmt.Errorf("unauthorized")
	}))
	// 按状态获取评论列表
	admin.GET("/comment/:filter", func(ctx *gin.Context) {
		filter := ctx.Param("filter")
		comments := new([]model.Comment)
		if filter == "all" {
			if db.Find(comments).Error != nil {
				ctx.JSON(500, gin.H{"error": "failed to get comments"})
				return
			}
		} else {
			if db.Where("status = ?", filter).Find(comments).Error != nil {
				ctx.JSON(500, gin.H{"error": "failed to get comments"})
				return
			}
		}
		ctx.JSON(200, comments)
	})
	// 审核评论
	admin.POST("/comment/audit/:id", func(ctx *gin.Context) {
		id, data := ctx.Param("id"), new(struct {
			Status string `json:"status"`
		})
		if err := ctx.ShouldBindJSON(data); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if db.Model(new(model.Comment)).Where("id = ?", id).Update("status", data.Status).Error != nil {
			ctx.JSON(500, gin.H{"error": "failed to update comment"})
			return
		}
		ctx.JSON(200, gin.H{"message": "Comment updated"})
	})
	r.Static("/db/", "./db/")
	r.StaticFile("/db.json", "./db.json")
	r.StaticFile("/index.json", "./index.json")

	subFS, _ := fs.Sub(f, "frontend/dist")
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(subFS))))

	crontab := cron.New(cron.WithSeconds())
	if _, err := crontab.AddFunc("0 15 * * * *", func() {
		feed = FetchFeed(config)
		ExportDB(feed)
	}); err != nil {
		log.Fatal("Failed to start feed update daemon")
	}
	crontab.Start()
	log.Println("Feed update daemon started.")
	log.Fatal(r.Run(":8192"))
}

func FetchFeed(config *model.Config) *model.Feed {
	wg, mutex, feed := new(sync.WaitGroup), new(sync.Mutex), new(model.Feed)

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
			articles := make([]model.Article, len(res.Items))
			for i, item := range res.Items {
				articles[i] = model.Article{Title: item.Title, Time: item.PubDate, Content: item.Content, Url: item.Link}
			}
			mutex.Lock()
			feed.Author = append(feed.Author, model.Author{Name: res.Title, Email: res.Author.Email, Uri: res.Link, Description: res.Description, Article: articles})
			mutex.Unlock()
			log.Println("Fetched RSS:", url)
		}(url)
	}
	wg.Wait()
	feed.Update = time.Now()
	log.Println("Fetch RSS done.")
	return feed
}

func ExportDB(feed *model.Feed) {
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
			fileName := fmt.Sprintf("db/%s.txt", lib.GenerateShortLink(article.Url))
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
