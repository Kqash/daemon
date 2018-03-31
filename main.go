package main

import "flag"

var apiURL string
var payKeyStroke string

func main() {
	flag.StringVar(&apiURL, "api", "http://localhost:8000/", "api url to ask for the queue")
	flag.StringVar(&payKeyStroke, "keystroke", "f6", "Keystroke that controls credit adding")
	flag.Parse()
}
