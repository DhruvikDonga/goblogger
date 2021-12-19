package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Blogs struct {
	*gorm.Model
	UsersID         uint
	Title           string
	Slug            string
	Blogdescription string
	Content         string
	Banner          string    //one post can have many images
	CreatedAt       time.Time `gorm:"column:created_at"`
}
