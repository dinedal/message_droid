package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var productionFlag = flag.Bool("production", false, "Use actual LED sign.")

type Service struct {
	ServiceId string `json:"service_id"`
	Text      string `json:"text"`
}

// TODO: Service obsoletion.
var state struct {
	Services []Service

	sync.RWMutex
}

const refreshRateSeconds time.Duration = 20

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func updateLedSign(text string) {
	if *productionFlag {
		cmd := exec.Command("/home/pi/muni-led-sign/client/lowlevel.pl", "--speed", "3", "--effect", "scroll")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err != nil {
			log.Println("updateLedSign: cmd.Run():", err)
		}
	} else {
		fmt.Printf("fake updateLedSign: %q\n", text)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	state.Lock()
	defer state.Unlock()

	for i, _ := range state.Services {
		if state.Services[i].ServiceId == service.ServiceId {
			state.Services[i] = service
			return
		}
	}
	state.Services = append(state.Services, service)
}

func background() {
	var currentText string
	var serviceIndex int

	for {
		var newText string

		state.RLock()
		if len(state.Services) > 0 {
			if serviceIndex >= len(state.Services) {
				serviceIndex = 0
			}
			newText = state.Services[serviceIndex].Text
			serviceIndex++
		}
		state.RUnlock()

		if newText != currentText {
			updateLedSign(newText)
			currentText = newText
		}

		time.Sleep(refreshRateSeconds * time.Second)
	}
}

func main() {
	log.Println("Started.")

	flag.Parse()

	go background()

	http.HandleFunc("/update", update)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
