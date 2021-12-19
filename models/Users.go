package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Users struct {
	*gorm.Model
	Name      string
	Email     string
	Password  string
	Role      int
	Linkedin  string
	Facebook  string
	Instagram string
	Reddit    string
	Twitter   string
	Medium    string
	Youtube   string
	Personal  string
	Blog      []Blogs `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //one user can have many posts

	CreatedAt time.Time `gorm:"column:created_at"`
}
type Auth struct {
	Name          string
	Id            uint
	Authenticated bool
}
