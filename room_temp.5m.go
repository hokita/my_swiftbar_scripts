//usr/bin/env go run $0 $@; exit
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Device struct
type Device struct {
	NewestEvents struct {
		Te Event `json:"te"`
		Hu Event `json:"hu"`
	} `json:"newest_events"`
}

// Event struct
type Event struct {
	Val       float32   `json:"val"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	url   = "https://api.nature.global/1/devices"
	token = ""
)

func main() {
	res, err := getRemoDevices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var devices []Device
	if err := json.Unmarshal(body, &devices); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, d := range devices {
		fmt.Printf("%.1f ℃\n", d.NewestEvents.Te.Val)
		fmt.Printf("%.1f ％\n", d.NewestEvents.Hu.Val)
	}
	os.Exit(0)
}

func getRemoDevices() (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := new(http.Client)
	res, err := client.Do(req)

	return res, err
}
