package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/higuoxing/concourse-google-chat-alert-resource/concourse"
)

func main() {
	err := json.NewEncoder(os.Stdout).Encode(concourse.CheckResponse{})
	if err != nil {
		log.Fatalln(err)
	}
}
