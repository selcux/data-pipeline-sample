package model

import "time"

type ProductView struct {
	Event      string     `json:"event"`
	MessageId  string     `json:"messageid"`
	UserId     string     `json:"userid"`
	Properties Properties `json:"properties"`
	Context    Context    `json:"context"`
	Timestamp  time.Time  `json:"timestamp"`
}

type Context struct {
	Source string `json:"source"`
}

type Properties struct {
	ProductId string `json:"productid"`
}
