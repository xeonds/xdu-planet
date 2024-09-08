package avalon

import (
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"xyz.xeonds/xdu-planet/model"
)

// 自动将超时评论转换为屏蔽状态
func GraveTimer(db *gorm.DB) {
	for {
		time.Sleep(time.Minute)
		comments := new([]model.Comment)
		if db.Where("status = ?", "audit").Find(comments).Error != nil {
			log.Println("Failed to get comments")
			continue
		}
		for _, comment := range *comments {
			if time.Since(comment.CreatedAt) > time.Hour*24 {
				if db.Update("status", "block").Where("id = ?", comment.ID).Error != nil {
					log.Println("Failed to delete comment")
					continue
				}
				log.Printf("Comment %d deleted\n", comment.ID)
			}
		}
	}
}

func Filter(db *gorm.DB, filter []string) {
	for {
		time.Sleep(time.Minute)
		comments := new([]model.Comment)
		if db.Where("status = ?", "audit").Find(comments).Error != nil {
			log.Println("Failed to get comments")
			continue
		}
		for _, comment := range *comments {
			for _, word := range filter {
				// 如果评论中包含关键词
				if strings.Contains(comment.Content, word) {
					if db.Update("status", "block").Where("id = ?", comment.ID).Error != nil {
						log.Println("Failed to delete comment")
						continue
					}
					log.Printf("Comment %d deleted\n", comment.ID)
				}
			}
		}
	}
}
