package main

import (
	"time"

	"github.com/triggit/MessageDroid/common"
)

type worker struct{}

func (_ *worker) GetServiceUpdate() string {
	const welcome = "Welcome to Triggit!"

	return welcome
}

func main() {
	common.ServiceMainLoop(&worker{}, "welcome", 5*time.Second)
}
