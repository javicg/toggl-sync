package api

import (
	"encoding/json"
	"github.com/javicg/toggl-sync/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestJiraApi_LogWork(t *testing.T) {
	ticket := "EXAMPLE-1234"
	expectedEntry := workLogEntry{
		Comment:          "Added automatically by toggl-sync",
		TimeSpentSeconds: 60,
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:         "/issue/" + ticket + "/worklog",
			RequestValidator: validateBodyMatches(t, expectedEntry),
			ResponseCode:     http.StatusCreated,
		}).
		Create()
	defer server.Close()

	config.SetJiraServerUrl(server.URL)

	jiraApi := NewJiraApi()
	err := jiraApi.LogWork(ticket, time.Duration(60)*time.Second)
	assert.Nil(t, err)
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
	assert.NotNilf(t, err, "API errors should be returned to the client")
}

func TestJiraApi_LogWorkWithUserDescription(t *testing.T) {
	ticket := "EXAMPLE-1234"
	expectedEntry := workLogEntry{
		Comment:          "Writing toggl-sync tests\nAdded automatically by toggl-sync",
		TimeSpentSeconds: 60,
	}

	server := NewHttpServer().
		StubApi(&Stubbing{
			Endpoint:         "/issue/" + ticket + "/worklog",
			RequestValidator: validateBodyMatches(t, expectedEntry),
			ResponseCode:     http.StatusCreated,
		}).
		Create()
	defer server.Close()

	config.SetJiraServerUrl(server.URL)

	jiraApi := NewJiraApi()
	err := jiraApi.LogWorkWithUserDescription(ticket, time.Duration(60)*time.Second, "Writing toggl-sync tests")
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
