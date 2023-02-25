package api

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/javicg/toggl-sync/config"
	"github.com/stretchr/testify/assert"
)

func TestTogglApi_GetMe(t *testing.T) {
	expectedMe := Me{
		Data: PersonalInfo{
			Email:    "tester@toggl-sync.com",
			Fullname: "TogglSync Tester",
		},
	}

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/me",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJSONString(expectedMe),
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	me, err := togglAPI.GetMe()
	assert.Nil(t, err)
	assert.Equal(t, expectedMe, *me)
}

func TestTogglApi_GetMe_ErrorWhenRequestFails(t *testing.T) {
	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/me",
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	_, err := togglAPI.GetMe()
	assert.NotNilf(t, err, "API errors should be returned to the client")
}

func TestTogglApi_GetMe_ErrorWhenResponseHasUnexpectedFormat(t *testing.T) {
	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/me",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJSONString("Bogus!"),
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	_, err := togglAPI.GetMe()
	assert.NotNilf(t, err, "JSON marshalling errors should be returned to the client")
}

func TestTogglApi_GetMe_ErrorWhenRequestErrors(t *testing.T) {
	config.Set(config.TogglServerURL, "%#2")

	togglAPI := NewTogglAPI()
	_, err := togglAPI.GetMe()
	assert.NotNil(t, err, "Request errors (e.g. misconfiguration) should be returned to the client")
}

func TestTogglApi_GetTimeEntries(t *testing.T) {
	expectedTimeEntries := []TimeEntry{
		{
			Id:          1,
			Duration:    120,
			Description: "Writing toggl-sync tests",
		},
		{
			Id:          2,
			Duration:    240,
			Description: "Increasing toggl-sync test coverage",
		},
	}

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/time_entries",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJSONString(expectedTimeEntries),
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	entries, err := togglAPI.GetTimeEntries(startDate, endDate)
	assert.Nil(t, err)
	assert.Equal(t, expectedTimeEntries, entries)
}

func TestTogglApi_GetTimeEntries_ErrorWhenRequestFails(t *testing.T) {
	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/time_entries",
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	_, err := togglAPI.GetTimeEntries(startDate, endDate)
	assert.NotNilf(t, err, "API errors should be returned to the client")
}

func TestTogglApi_GetTimeEntries_ErrorWhenResponseHasUnexpectedFormat(t *testing.T) {
	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/time_entries",
			ResponseCode: http.StatusOK,
			ResponseBody: AsJSONString("Bogus!"),
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	_, err := togglAPI.GetTimeEntries(startDate, endDate)
	assert.NotNilf(t, err, "JSON marshalling errors should be returned to the client")
}

func TestTogglApi_GetTimeEntries_ErrorWhenRequestErrors(t *testing.T) {
	config.Set(config.TogglServerURL, "%#2")

	togglAPI := NewTogglAPI()
	startDate, _ := time.Parse("2006-01-02", "2020-05-08")
	endDate, _ := time.Parse("2006-01-02", "2020-05-09")
	_, err := togglAPI.GetTimeEntries(startDate, endDate)
	assert.NotNil(t, err, "Request errors (e.g. misconfiguration) should be returned to the client")
}

func TestTogglApi_GetProjectById(t *testing.T) {
	projectId := 10
	expectedProject := Project{
		Data: ProjectData{
			Id:   projectId,
			Name: "Top Secret",
		},
	}

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/projects/" + strconv.Itoa(projectId),
			ResponseCode: http.StatusOK,
			ResponseBody: AsJSONString(expectedProject),
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	project, err := togglAPI.GetProjectById(projectId)
	assert.Nil(t, err)
	assert.Equal(t, expectedProject, *project)
}

func TestTogglApi_GetProjectById_ErrorWhenRequestFails(t *testing.T) {
	projectId := 10

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/projects/" + strconv.Itoa(projectId),
			ResponseCode: http.StatusBadGateway,
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	_, err := togglAPI.GetProjectById(projectId)
	assert.NotNilf(t, err, "API errors should be returned to the client")
}

func TestTogglApi_GetProjectById_ErrorWhenResponseHasUnexpectedFormat(t *testing.T) {
	projectId := 10

	server := NewHTTPServer().
		StubAPI(&Stubbing{
			Endpoint:     "/projects/" + strconv.Itoa(projectId),
			ResponseCode: http.StatusOK,
			ResponseBody: AsJSONString("Bogus!"),
		}).
		Create()
	defer server.Close()

	config.Set(config.TogglServerURL, server.URL)

	togglAPI := NewTogglAPI()
	_, err := togglAPI.GetProjectById(projectId)
	assert.NotNilf(t, err, "JSON marshalling errors should be returned to the client")
}

func TestTogglApi_GetProjectById_ErrorWhenRequestErrors(t *testing.T) {
	config.Set(config.TogglServerURL, "%#2")

	togglAPI := NewTogglAPI()
	_, err := togglAPI.GetProjectById(10)
	assert.NotNil(t, err, "Request errors (e.g. misconfiguration) should be returned to the client")
}
