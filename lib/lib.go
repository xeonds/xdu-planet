package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"xyz.xeonds/xdu-planet/model"
)

// 加载配置文件
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

// 打开数据库连接
func NewDB(config *model.DatabaseConfig, migrator func(*gorm.DB) error) *gorm.DB {
	var db *gorm.DB
	var err error
	switch config.Type {
	case "mysql":
		dsn := config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DB), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if config.Migrate {
		if migrator == nil {
			log.Fatalf("Migrator is nil")
		}
		if err = migrator(db); err != nil {
			log.Fatalf("Failed to migrate tables: %v", err)
		}
	}
	return db
}

// 鉴权中间件
// 支持添加权限校验（返回error表示校验失败），以及上下文操作
func JWTMiddleware(authToken func(*gin.Context, string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if authToken != nil && authToken(c, c.GetHeader("Authorization")) != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

// 日志中间件，记录客户端IP，请求方法，请求路径，请求耗时，请求头的Aurhorization字段, 将日志保存在path/to/logFile.log中
func LoggerMiddleware(logFile string) gin.HandlerFunc {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Printf("%s %s %s %s %s\n", c.ClientIP(), c.Request.Method, c.Request.URL.Path, time.Since(start), c.GetHeader("Authorization"))
	}
}

func GenerateShortLink(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	hash := h.Sum(nil)
	shortLink := hex.EncodeToString(hash)[:24] // 取前8个字符作为短链接
	return shortLink
}
