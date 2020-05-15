package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

type JiraApi struct {
	baseUrl string
	client  *http.Client
}

func NewJiraApi() (api *JiraApi, err error) {
	api = &JiraApi{}
	api.baseUrl = viper.GetString("JIRA_SERVER_URL") + "/rest/api/latest"
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

	req.SetBasicAuth(viper.GetString("JIRA_USERNAME"), viper.GetString("JIRA_PASSWORD"))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return jira.client.Do(req)
}
