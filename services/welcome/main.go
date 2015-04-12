package main

import (
	"time"

	"github.com/dinedal/message_droid/common"
)

type worker struct{}

func (_ *worker) GetServiceUpdate() string {
	const welcome = "Welcome!"

	return welcome
}

func main() {
	common.ServiceMainLoop(&worker{}, "welcome", 5*time.Second)
}
