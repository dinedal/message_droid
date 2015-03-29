package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/dinedal/message_droid/common"
)

var productionFlag = flag.Bool("production", false, "Use actual LED sign.")

var state struct {
	Services []common.ServiceUpdate // Queue of upcoming services. Index 0 is next service.

	sync.RWMutex
}

const refreshRateSeconds time.Duration = 10

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func updateLedSign(text string) {
	if *productionFlag {
		cmd := exec.Command("lowlevel.pl", "--speed", "3", "--effect", "scroll")
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

	var service common.ServiceUpdate
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

	for {
		var newText string

		state.RLock()
		if len(state.Services) > 0 {
			newText = state.Services[0].Text
			state.Services = state.Services[1:]
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
	flag.Parse()

	log.Println("Started.")

	// Set the working directory to the root of the package, so that its assets folder can be used.
	{
		bpkg, err := build.Import("github.com/dinedal/message_droid", "", build.ImportComment)
		if err != nil {
			log.Fatalln("Unable to find github.com/dinedal/message_droid package in your GOPATH, it's needed to load assets.")
		}

		err = os.Chdir(bpkg.Dir)
		if err != nil {
			log.Panicln("os.Chdir:", err)
		}
	}

	go background()

	http.HandleFunc("/update", update)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
