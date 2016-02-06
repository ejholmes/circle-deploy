package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ejholmes/hookshot/events"
)

type buildRequest struct {
	BuildParameters map[string]string `json:"build_parameters"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	var event events.Deployment

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("at=handle repo=%s ref=%s", event.Repository.FullName, event.Deployment.Ref))

	raw, err := json.Marshal(buildRequest{
		BuildParameters: map[string]string{
			"GITHUB_DEPLOYMENT": fmt.Sprintf("%d", event.Deployment.ID),
		},
	})
	if err != nil {
		panic(err)
	}
	token := r.FormValue("circle-token")

	url := fmt.Sprintf("https://circleci.com/api/v1/project/%s/tree/%s?circle-token=%s", event.Repository.FullName, event.Deployment.Ref, token)
	req, err := http.NewRequest("POST", url, bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode/100 != 2 {
		panic(fmt.Errorf("unexpected response: %v", resp.Status))
	}
}

func main() {
	http.Handle("/", http.HandlerFunc(handle))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
