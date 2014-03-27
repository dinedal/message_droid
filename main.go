package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

var state = struct {
	State map[string]string

	sync.RWMutex
}{State: make(map[string]string)}

func readAllString(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: Switch to json, etc.
// curl -i -X POST -d"your message here for now" http://10.0.0.223:8080/set

func set(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	defer state.RUnlock()

	if r.Method != "POST" {
		return
	}

	text := readAllString(r.Body)

	fmt.Println(text)

	cmd := exec.Command("/home/pi/muni-led-sign/client/lowlevel.pl", "--speed", "3", "--effect", "scroll")
	cmd.Stdin = strings.NewReader(text)
	out, err := cmd.CombinedOutput()
	//panicOnError(err)
	//fmt.Println(out)
	fmt.Println(string(out), err)
}

func list(w http.ResponseWriter, r *http.Request) {
	state.RLock()
	defer state.RUnlock()

	fmt.Fprintf(w, "We have %v connection(s).\n", len(state.State))
	fmt.Fprintf(w, "%#v", state.State)
}

func main() {
	fmt.Println("Started.")

	http.HandleFunc("/set", set)
	http.HandleFunc("/list", list)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
