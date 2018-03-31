package drainer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

// Drainer contacts the queue and drains it
type Drainer struct {
	apiURL     string
	pollTime   int
	creditChan chan<- int
	stopChan   chan bool
}

// NewDrainer creates a new Drainer
func NewDrainer(apiURL string, pollTime int, creditChan chan int) (*Drainer, error) {
	stopChan := make(chan bool, 1)

	if pollTime <= 0 {
		pollTime = 10
	}
	drainer := Drainer{
		apiURL:     apiURL,
		pollTime:   pollTime,
		creditChan: creditChan,
		stopChan:   stopChan,
	}
	return &drainer, nil
}

// Drain continuously checks the queue after some timeout
func (d *Drainer) Drain() {
	for {
		select {
		case <-d.stopChan:
			fmt.Println("Received stop")
			return
		case <-time.After(time.Duration(d.pollTime) * time.Second):
			credits := d.CheckForCredits()
			if credits > 0 {
				d.creditChan <- credits
			}
		}
	}
}

// CheckForCredits gets credits from queue
func (d *Drainer) CheckForCredits() int {
	url := path.Join(d.apiURL, "qash")
	resp, err := http.Get(url)

	if err != nil {
		// Ok just assume that we have nothing
		return 0
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(body), &data)

	creditsIface, ok := data["credits"]
	if !ok {
		return 0
	}
	credits := creditsIface.(int)

	if credits <= 0 {
		return 0
	}
	go d.emptyCredits()
	return credits
}

func (d *Drainer) emptyCredits() {
	url := path.Join(d.apiURL, "qash")
	http.NewRequest("PATCH", url, nil)
	// Don't care about getting back response
}

// Stop stops the drainer
func (d *Drainer) Stop() {
	d.stopChan <- true
}
