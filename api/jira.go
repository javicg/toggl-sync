package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"io"
	"net/http"
	"time"
)

const workLogEntryCommentFooter = "Added automatically by toggl-sync"

type JiraApi interface {
	LogWork(ticket string, timeSpent time.Duration) error
	LogWorkWithUserDescription(ticket string, description string, timeSpent time.Duration) error
}

type JiraApiHttpClient struct {
	baseUrl string
	client  *http.Client
}

func NewJiraApi() JiraApi {
	api := &JiraApiHttpClient{}
	api.baseUrl = config.GetJiraServerUrl() + "/rest/api/latest"
	api.client = &http.Client{}
	return api
}

type WorkLogEntry struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func (jira *JiraApiHttpClient) LogWorkWithUserDescription(ticket string, userDescription string, timeSpent time.Duration) (err error) {
	entry := createWorkLogEntryWithUserDescription(userDescription, timeSpent)
	return jira.logEntry(ticket, entry)
}

func createWorkLogEntryWithUserDescription(userDescription string, timeSpent time.Duration) *WorkLogEntry {
	return &WorkLogEntry{
		Comment:          fmt.Sprintf("%s\n%s", userDescription, workLogEntryCommentFooter),
		TimeSpentSeconds: int(timeSpent.Seconds()),
	}
}

func (jira *JiraApiHttpClient) LogWork(ticket string, timeSpent time.Duration) (err error) {
	entry := createWorkLogEntry(timeSpent)
	return jira.logEntry(ticket, entry)
}

func createWorkLogEntry(timeSpent time.Duration) *WorkLogEntry {
	return &WorkLogEntry{
		Comment:          workLogEntryCommentFooter,
		TimeSpentSeconds: int(timeSpent.Seconds()),
	}
}

func (jira *JiraApiHttpClient) logEntry(ticket string, entry *WorkLogEntry) error {
	entryJson, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("[LogWork] Marshalling of work entry failed! Error: %s", err)
	}

	resp, err := jira.postAuthenticated("/issue/"+ticket+"/worklog", bytes.NewBuffer(entryJson))
	if err != nil {
		return err
	} else if resp.StatusCode != 201 {
		return fmt.Errorf("[LogWork] Request to log work for ticket [%s] failed with status [%d]", ticket, resp.StatusCode)
	}

	return resp.Body.Close()
}

func (jira *JiraApiHttpClient) postAuthenticated(path string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", jira.baseUrl+path, body)
	if err != nil {
		return
	}

	req.SetBasicAuth(config.GetJiraUsername(), config.GetJiraPassword())

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return jira.client.Do(req)
}
