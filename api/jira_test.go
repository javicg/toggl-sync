package api

import (
	"encoding/json"
	"github.com/javicg/toggl-sync/config"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestJiraApi_LogWork(t *testing.T) {
	ticket := "EXAMPLE-1234"
	expectedEntry := WorkLogEntry{
		Comment:          "Added automatically by toggl-sync",
		TimeSpentSeconds: 61,
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint: "/issue/" + ticket + "/worklog",
			RequestValidator: func(r *http.Request) {
				validateBodyIsExactly(t, r, expectedEntry)
			},
			ResponseCode: http.StatusCreated,
		}).
		Create()
	defer server.Close()

	config.SetJiraServerUrl(server.URL)

	jiraApi := NewJiraApi()
	err := jiraApi.LogWork(ticket, time.Duration(60)*time.Second)
	if err != nil {
		t.Errorf("API call failed with an error: %s", err)
	}
}

func TestJiraApi_LogWork_ErrorWhenRequestFails(t *testing.T) {
	ticket := "EXAMPLE-1234"

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:     "/issue/" + ticket + "/worklog",
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.SetJiraServerUrl(server.URL)

	jiraApi := NewJiraApi()
	err := jiraApi.LogWork(ticket, time.Duration(60)*time.Second)
	if err == nil {
		t.Error("API errors should be returned to the client")
	}
}

func TestJiraApi_LogWorkWithUserDescription(t *testing.T) {
	ticket := "EXAMPLE-1234"
	expectedEntry := WorkLogEntry{
		Comment:          "Writing toggl-sync tests\nAdded automatically by toggl-sync",
		TimeSpentSeconds: 60,
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint: "/issue/" + ticket + "/worklog",
			RequestValidator: func(r *http.Request) {
				validateBodyIsExactly(t, r, expectedEntry)
			},
			ResponseCode: http.StatusCreated,
		}).
		Create()
	defer server.Close()

	config.SetJiraServerUrl(server.URL)

	jiraApi := NewJiraApi()
	err := jiraApi.LogWorkWithUserDescription(ticket, "Writing toggl-sync tests", time.Duration(60)*time.Second)
	if err != nil {
		t.Errorf("API call failed with an error: %s", err)
	}
}

func validateBodyIsExactly(t *testing.T, r *http.Request, expectedEntry WorkLogEntry) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Parsing request body failed with an error: %s", err)
	} else {
		var payload WorkLogEntry
		err := json.Unmarshal(bytes, &payload)
		if err != nil {
			t.Errorf("JSON unmarshalling failed: %s", err)
		}

		if payload != expectedEntry {
			t.Errorf("Unexpected payload: was [%#v] instead of [%#v]", payload, expectedEntry)
		}
	}
}
