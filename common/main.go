package common

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ServiceUpdate struct {
	ServiceId string `json:"service_id"`
	Text      string `json:"text"`
}

type ServiceWorker interface {
	GetServiceUpdate() string
}

func ServiceMainLoop(worker ServiceWorker, serviceId string, updateInterval time.Duration) {
	for {
		started := time.Now()

		service := ServiceUpdate{
			ServiceId: serviceId,
			Text:      worker.GetServiceUpdate(),
		}

		body, err := json.Marshal(service)
		if err != nil {
			log.Println("json.Marshal:", err)
		}

		_, err = http.Post("http://localhost:8080/update", "application/json", bytes.NewReader(body))
		if err != nil {
			log.Println("http.Post:", err)
		}

		time.Sleep(updateInterval - time.Since(started))
	}
}
