package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	for {
		// TODO: Import type Service from MD lib, encode that.
		service := `{"service_id": "welcome", "text": "Welcome to Triggit!"}`

		_, err := http.Post("http://localhost:8080/update", "application/json", strings.NewReader(service))
		if err != nil {
			log.Println("http.Post:", err)
		}

		// TODO: Factor this logic to a common place.
		time.Sleep(5 * time.Second)
	}
}
