package models

import "gorm.io/gorm"

type Template struct {
	*gorm.Model
	TemplateName string `gorm:"primarykey;unique:true",json:"templateName"`
	Html         string `json:"html"`
}
