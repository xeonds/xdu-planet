package model

import (
	"time"

	"gorm.io/gorm"
)

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
type Comment struct {
	gorm.Model
	ArticleId string `json:"article_id"`
	Content   string `json:"content"`
	UserId    string `json:"user_id"`
	ReplyTo   string `json:"reply_to"`
	Status    string `json:"status"` // ok|block|delete|audit
}

type DatabaseConfig struct {
	Type     string // 数据库类型
	Host     string
	Port     string
	User     string
	Password string
	DB       string // 数据库名
	Migrate  bool
}

// config file
type Config struct {
	Version int
	DatabaseConfig
	AvalonGuard struct {
		EnableGraveTimer bool
		GraveTimeout     time.Duration
		EnableFilter     bool
		Filter           []string
	}
	LogFile string // 日志文件
	// 邮箱支持
	AdminToken []string // 管理员token
	Feeds      []string
}
