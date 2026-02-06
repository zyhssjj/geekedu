package main

import (
	"gorm.io/gorm"
	"time"
)

// 用户表（users）
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;comment:'用户名';size:50"`
	Password string `gorm:"not null;comment:'密码（bcrypt加密）';size:100"`
	Role     int    `gorm:"not null;default:0;comment:'角色：0-学员，1-管理员'"`
}

// 课程表（courses）
type Course struct {
	gorm.Model
	Title        string  `gorm:"not null;comment:'课程标题';size:100"`
	Price        float64 `gorm:"not null;comment:'课程价格'"`
	Intro        string  `gorm:"comment:'课程简介';type:text"`
	CoverOssKey  string  `gorm:"not null;comment:'封面OSS Key';size:255"`
	CreateUserID uint    `gorm:"not null;comment:'创建者ID'"`
}

// 视频表（videos）
type Video struct {
	ID        uint      `gorm:"primarykey"`
	CourseID  uint      `gorm:"not null;index"`  // 外键关联课程
	Title     string    `gorm:"type:varchar(255);not null"` // 视频标题
	OssKey    string    `gorm:"type:varchar(500);not null"` // OSS存储key（必须非空）
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
// 订单表（orders）
type Order struct {
	gorm.Model
	UserID     uint `gorm:"not null;comment:'用户ID';index:idx_user_id"`
	CourseID   uint `gorm:"not null;comment:'课程ID';index:idx_course_id"`
	OrderStatus int `gorm:"not null;default:1;comment:'订单状态：1-已完成'"`
}

// 自定义表名（gorm默认复数形式，显式指定更清晰）
func (User) TableName() string { return "users" }
func (Course) TableName() string { return "courses" }
func (Video) TableName() string { return "videos" }
func (Order) TableName() string { return "orders" }