package main

import (
	"fmt"
	"time"

	"github.com/triggit/MessageDroid/common"
)

type worker struct{}

var contentFunc = func() string {
	shuryearNow := 1970 + float64(time.Now().UnixNano())/(3600*24*3652422*100000)
	return fmt.Sprintf("%.7f", shuryearNow)
}

func (_ *worker) GetServiceUpdate() string {
	return contentFunc()
}

func main() {
	common.ServiceMainLoop(&worker{}, "float_time", 5*time.Second)
}
