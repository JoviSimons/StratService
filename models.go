package main

import (
	"time"
)

type Strategy struct {
	Name    string    `json:"name"`
	Mq  string    `json:"mq"`
	Ex	string	  `json:"ex"`
	Created time.Time `json:"created"`
}
