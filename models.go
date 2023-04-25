package main

import (
	"time"
)

type Strategy struct {
	Name    string    `json:"name"`
	Mq      string    `json:"mq"`
	Ex      string    `json:"ex"`
	Created time.Time `json:"created"`
}

type StrategyRequest struct {
	Id     string    `json:"id"`
	Name    string    `json:"name"`
	Ex      string    `json:"ex"`
	Created time.Time `json:"created"`
}
