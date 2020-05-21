package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/javicg/toggl-sync/config"
	"io"
	"log"
	"net/http"
	"time"
)

const workLogEntryCommentFooter = "Added automatically by toggl-sync"

type JiraApi struct {
	baseUrl string
	client  *http.Client
}

func NewJiraApi() (api *JiraApi) {
	api = &JiraApi{}
	api.baseUrl = config.GetJiraServerUrl() + "/rest/api/latest"
	api.client = &http.Client{}
	return
}

type WorkLogEntry struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

func (jira *JiraApi) LogWorkWithUserDescription(ticket string, userDescription string, timeSpent time.Duration) (err error) {
	entry := createWorkLogEntryWithUserDescription(userDescription, timeSpent)
	return jira.logEntry(ticket, entry)
}

func createWorkLogEntryWithUserDescription(userDescription string, timeSpent time.Duration) *WorkLogEntry {
	return &WorkLogEntry{
		Comment:          fmt.Sprintf("%s\n%s", userDescription, workLogEntryCommentFooter),
		TimeSpentSeconds: int(timeSpent.Seconds()),
	}
}

func (jira *JiraApi) LogWork(ticket string, timeSpent time.Duration) (err error) {
	entry := createWorkLogEntry(timeSpent)
	return jira.logEntry(ticket, entry)
}

func createWorkLogEntry(timeSpent time.Duration) *WorkLogEntry {
	return &WorkLogEntry{
		Comment:          workLogEntryCommentFooter,
		TimeSpentSeconds: int(timeSpent.Seconds()),
	}
}

func (jira *JiraApi) logEntry(ticket string, entry *WorkLogEntry) (err error) {
	entryJson, err := json.Marshal(entry)
	if err != nil {
		log.Fatalln("[LogWork] Marshalling of work entry failed! Error:", err)
	}

	resp, err := jira.postAuthenticated("/issue/"+ticket+"/worklog", bytes.NewBuffer(entryJson))
	if err != nil {
		return
	} else if resp.StatusCode != 201 {
		log.Fatalf("[LogWork] Request failed with status: %d", resp.StatusCode)
	}

	return resp.Body.Close()
}

func (jira *JiraApi) postAuthenticated(path string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", jira.baseUrl+path, body)
	if err != nil {
		return
	}

	req.SetBasicAuth(config.GetJiraUsername(), config.GetJiraPassword())

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return jira.client.Do(req)
}
