package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type Message struct {
	*gorm.Model
	ID           string      `json:"ID"`
	Subject      string      `json:"subject"`
	From         string      `json:"from"`
	To           StringSlice `gorm:"type:json"`
	Data         JSONB       `gorm:"type:jsonb"`
	TemplateName string      `json:"templateName"`
	Status       string      `json:"status"`
	Type         string      `json:"type"`
}

type StringSlice []string

func (ss StringSlice) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

func (ss *StringSlice) Scan(src interface{}) error {
	switch data := src.(type) {
	case []byte:
		return json.Unmarshal(data, ss)
	case string:
		return json.Unmarshal([]byte(data), ss)
	}
	return nil
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(src interface{}) error {
	switch data := src.(type) {
	case []byte:
		return json.Unmarshal(data, src)
	case string:
		return json.Unmarshal([]byte(data), src)
	}
	return nil
}
