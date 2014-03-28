package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var contentFunc = func() string {
	shuryearNow := 1970 + float64(time.Now().UnixNano())/(3600*24*3652422*100000)
	return fmt.Sprintf("%.7f", shuryearNow)
}

func main() {
	for {
		// TODO: Import type Service from MD lib, encode that.
		service := `{"service_id": "float_time", "text": "` + contentFunc() + `"}`

		_, err := http.Post("http://localhost:8080/update", "application/json", strings.NewReader(service))
		if err != nil {
			log.Println("http.Post:", err)
		}

		// TODO: Factor this logic to a common place.
		time.Sleep(5 * time.Second)
	}
}
