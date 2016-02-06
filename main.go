package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ejholmes/hookshot"
	"github.com/ejholmes/hookshot/events"
)

type buildRequest struct {
	BuildParameters map[string]string `json:"build_parameters"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/ejholmes/circle-deploy#readme", 301)
}

func deployment(w http.ResponseWriter, r *http.Request) {
	var event events.Deployment

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("at=handle repo=%s ref=%s", event.Repository.FullName, event.Deployment.Ref))

	raw, err := json.Marshal(buildRequest{
		BuildParameters: map[string]string{
			"GITHUB_DEPLOYMENT":             fmt.Sprintf("%d", event.Deployment.ID),
			"GITHUB_DEPLOYMENT_ENVIRONMENT": fmt.Sprintf("%d", event.Deployment.Environment),
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode/100 != 2 {
		panic(fmt.Errorf("unexpected response: %v", resp.Status))
	}
}

func main() {
	r := hookshot.NewRouter()
	r.HandleFunc("ping", ping)
	r.HandleFunc("deployment", deployment)
	r.NotFoundHandler = http.HandlerFunc(ping)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
