package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type JiraApi struct {
	baseUrl string
	client  *http.Client
}

func NewJiraApi() (api *JiraApi, err error) {
	api = &JiraApi{}
	baseUrl, ok := os.LookupEnv("JIRA_BASE_URL")
	if !ok {
		err = errors.New(fmt.Sprintf("%s not specified!", "JIRA_BASE_URL"))
		return
	}
	api.baseUrl = baseUrl
	api.client = &http.Client{}
	return
}

type WorkLogEntry struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func (jira *JiraApi) LogWork(ticket string, timeSpent time.Duration) (err error) {
	entry := &WorkLogEntry{Comment: "Added automatically by toggl-sync", TimeSpentSeconds: int(timeSpent.Seconds())}
	entryJson, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("[LogWork] Marshalling of work entry failed! Error:", err)
		return
	}

	fmt.Println(string(entryJson))

	resp, err := jira.postAuthenticated("/issue/"+ticket+"/worklog", bytes.NewBuffer(entryJson))
	if err != nil {
		return
	} else if resp.StatusCode != 201 {
		err = errors.New(fmt.Sprintf("[LogWork] Request failed with status: %d", resp.StatusCode))
		return
	}

	return resp.Body.Close()
}

func (jira *JiraApi) postAuthenticated(path string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", jira.baseUrl+path, body)
	if err != nil {
		return
	}

	err = addBasicAuth(req, "JIRA_USERNAME", "JIRA_PASSWORD")
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return jira.client.Do(req)
}
