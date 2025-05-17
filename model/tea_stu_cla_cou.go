package model

import "gorm.io/gorm"

type Subject struct {
	gorm.Model
	Name       string
	Tags       []string       `gorm:"serializer:csv"`  // 课程标签
	Syllabus   []string       `gorm:"serializer:csv"`  // 课程大纲
	Properties map[string]any `gorm:"serializer:json"` // 课程属性
}
