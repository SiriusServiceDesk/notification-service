package models

type Email struct {
	Subject string   `json:"subject"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Data    string   `json:"data"`
}
