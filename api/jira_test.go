package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/javicg/toggl-sync/config"
	"github.com/stretchr/testify/assert"
)

func TestJiraApi_LogWork(t *testing.T) {
	ticket := "EXAMPLE-1234"
	expectedEntry := workLogEntry{
		Comment:          "Added automatically by toggl-sync",
		TimeSpentSeconds: 60,
	}

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:         "/issue/" + ticket + "/worklog",
			RequestValidator: validateBodyMatches(t, expectedEntry),
			ResponseCode:     http.StatusCreated,
		}).
		Create()
	defer server.Close()

	config.Set(config.JiraServerURL, server.URL)

	jiraAPI := NewJiraAPI()
	err := jiraAPI.LogWork(ticket, time.Duration(60)*time.Second)
	assert.Nil(t, err)
}

func TestJiraApi_LogWork_ErrorWhenRequestFails(t *testing.T) {
	ticket := "EXAMPLE-1234"

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/issue/" + ticket + "/worklog",
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.Set(config.JiraServerURL, server.URL)

	jiraAPI := NewJiraAPI()
	err := jiraAPI.LogWork(ticket, time.Duration(60)*time.Second)
	assert.NotNilf(t, err, "API errors should be returned to the client")
}

func TestJiraApi_LogWork_ErrorWhenRequestErrors(t *testing.T) {
	config.Set(config.JiraServerURL, "%#2")

	jiraAPI := NewJiraAPI()
	err := jiraAPI.LogWork("EXAMPLE-1234", time.Duration(60)*time.Second)
	assert.NotNil(t, err, "Request errors (e.g. misconfiguration) should be returned to the client")
}

func TestJiraApi_LogWorkWithUserDescription(t *testing.T) {
	ticket := "EXAMPLE-1234"
	expectedEntry := workLogEntry{
		Comment:          "Writing toggl-sync tests\nAdded automatically by toggl-sync",
		TimeSpentSeconds: 60,
	}

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:         "/issue/" + ticket + "/worklog",
			RequestValidator: validateBodyMatches(t, expectedEntry),
			ResponseCode:     http.StatusCreated,
		}).
		Create()
	defer server.Close()

	config.Set(config.JiraServerURL, server.URL)

	jiraAPI := NewJiraAPI()
	err := jiraAPI.LogWorkWithUserDescription(ticket, time.Duration(60)*time.Second, "Writing toggl-sync tests")
	assert.Nil(t, err)
}

func validateBodyMatches(t *testing.T, expectedBody workLogEntry) func(*http.Request) {
	return func(r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Parsing request body failed with an error: %s", err)
		} else {
			var body workLogEntry
			err := json.Unmarshal(bytes, &body)
			if err != nil {
				t.Errorf("JSON unmarshalling failed: %s", err)
			}

			if body != expectedBody {
				t.Errorf("Unexpected payload: was [%#v] instead of [%#v]", body, expectedBody)
			}
		}
	}
}
